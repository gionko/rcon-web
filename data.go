package main

import (
	"bufio"
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"

	"github.com/madcitygg/rcon"
	"github.com/oschwald/geoip2-golang"
	"github.com/rumblefrog/go-a2s"
)

type User struct {
	Duration float32 `json:"duration"`
	ID       string  `json:"id"`
	IP       string  `json:"ip"`
	Name     string  `json:"name"`
	Ping     uint32  `json:"ping"`
	Score    uint32  `json:"score"`
	State    string  `json:"state"`

	City       string  `json:"city"`
	Country    string  `json:"country"`
	CountryISO string  `json:"country_iso"`
	Latitude   float64 `json:"latitude"`
	Longitude  float64 `json:"longitude"`
	TimeZone   string  `json:"timezone"`
}

func get_scores() ([]User, error) {
	var users []User

	// Get A2S player request response

	req, err := a2s.NewClient(fmt.Sprintf("%s:%d", config.ServerAddress, config.ServerPort))
	if err != nil {
		return nil, fmt.Errorf("A2S client error: %+v", err)
	}
	defer req.Close()

	resp, err := req.QueryPlayer()
	if err != nil {
		return nil, fmt.Errorf("A2S player query error: %+v", err)
	}

	// Convert response to User slice

	for _, p := range resp.Players {
		var u User
		u.Name     = p.Name
		u.Score    = p.Score
		u.Duration = p.Duration
		users = append(users, u)
	}

	return users, nil
}

func get_status() ([]string, error) {
	var status []string

	// Get RCON `status` command response

	req, err := rcon.Dial(fmt.Sprintf("%s:%d", config.ServerAddress, config.ServerPort))
	if err != nil {
		return nil, fmt.Errorf("RCON dial error: %+v", err)
	}
	defer req.Close()

	err = req.Authenticate(config.ServerPassword)
	if err != nil {
		return nil, fmt.Errorf("RCON authentication error: %+v", err)
	}

	resp, err := req.Execute("status")
	if err != nil {
		return nil, fmt.Errorf("RCON `status` command error: %+v", err)
	}

	// Convert response to array of strings

	scanner := bufio.NewScanner(strings.NewReader(resp.Body))
	for scanner.Scan() {
		status = append(status, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("RCON server response parsing error: %+v", err)
	}

	return status, nil
}

func get_users(status []string) ([]User, error) {
	var users []User

	// Get user scores

	scores, err := get_scores()
	if err != nil {
		return nil, err
	}


	geo, err := geoip2.Open("GeoLite2-City.mmdb")
	if err != nil {
		return nil, fmt.Errorf("Error opening GeoIP2 database: %+v", err)
	}
	defer geo.Close()

	// Extract user info from status

	for _, line := range status {
		// re := regexp.MustCompile("(?i).*?\"(.*?)\" (.*?) (.*?) (.*?) (.*?) (.*?) (.*?) (.*?):(.*?)$")
		re := regexp.MustCompile("(?i).*?\"(.*?)\" +(.*?) +(.*?) +(.*?) +(.*?) +(.*?) +(.*?) +(.*?):(.*?)$")
		match := re.FindStringSubmatch(line)

		// If match is successful, it will contain following data
		// 0: full match
		// 1: name
		// 2: steam id
		// 3: connected
		// 4: ping
		// 5: loss
		// 6: state
		// 7: rate
		// 8: ip
		// 9: port

		if match != nil {
			var user User
			user.ID = match[2]
			user.IP = match[8]
			user.Name = match[1]
			user.State = match[6]

			ping, err := strconv.ParseUint(match[4], 10, 32)
			if err != nil {
				log.Errorf("Could not extract user ping from RCON `status` response: (%s) %+v", line, err)
			} else {
				user.Ping = uint32(ping)
			}

			// Merge duration and score

			var del = -1
			for i, score := range scores {
				if user.Name == score.Name {
					user.Duration = score.Duration
					user.Score = score.Score
					del = i
					break
				}
			}

			// Delete user from scores if match was found

			if del >= 0 {
				scores = append(scores[:del], scores[del + 1:]...)
			}

			// Fill in data from GeoIP

			ip := net.ParseIP(user.IP)
			record, err := geo.City(ip)
			if err != nil {
				log.Errorf("GeoIP error: %+v", err)
			}
			user.City       = record.City.Names["en"]
			user.Country    = record.Country.Names["en"]
			user.CountryISO = record.Country.IsoCode
			user.Latitude   = record.Location.Latitude
			user.Longitude  = record.Location.Longitude
			user.TimeZone   = record.Location.TimeZone

			// Save the user

			users = append(users, user)
		}
	}

	return users, nil
}


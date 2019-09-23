package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/madcitygg/rcon"
	"github.com/oschwald/geoip2-golang"
	"github.com/rumblefrog/go-a2s"
)

type Map struct {
	Name       string  `json:"name"`
	Map        string  `json:"map"`
}

type Player struct {
	Duration   float32 `json:"duration"`
	ID         string  `json:"id"`
	IP         string  `json:"ip"`
	Name       string  `json:"name"`
	Ping       uint32  `json:"ping"`
	Score      uint32  `json:"score"`
	State      string  `json:"state"`

	City       string  `json:"city"`
	Country    string  `json:"country"`
	CountryISO string  `json:"country_iso"`
	Latitude   float64 `json:"latitude"`
	Longitude  float64 `json:"longitude"`
	TimeZone   string  `json:"timezone"`
}

type Steam struct {
	SteamID    string  `json:"steamid"`
	Name       string  `json:"personaname"`
	URL        string  `json:"profileurl"`
	Avatar     string  `json:"avatarfull"`
	Created    uint32  `json:"timecreated"`
	CountryISO string  `json:"loccountrycode"`
}

type PlayerSteam struct {
	Player
	Steam Steam `json:"steam"`
}

func get_maps(maps []string) ([]Map, error) {
	var list []Map

	// Extract map names from `maps` command reply

	for _, line := range maps {
		re := regexp.MustCompile("(?i).*?\\(fs\\) (.*?).bsp$")
		match := re.FindStringSubmatch(line)

		// If match is successful, it will contain following data
		// 0: full match
		// 1: map name

		if match != nil {
			var m Map
			m.Name = strings.Title(strings.Split(match[1], "_")[0])
			m.Map  = match[1]

			// Save the map

			list = append(list, m)
		}
	}

	// Done

	return list, nil
}

func get_players(status []string) ([]Player, error) {
	var players []Player

	// Get player scores

	scores, err := get_scores()
	if err != nil {
		return nil, err
	}

	geo, err := geoip2.Open(config.GeoIP2_DB)
	if err != nil {
		return nil, fmt.Errorf("Error opening GeoIP2 database: %+v", err)
	}
	defer geo.Close()

	// Extract player info from status

	for _, line := range status {
		re := regexp.MustCompile("(?i).*?\"(.*?)\" +(STEAM_.*?) +(.*?) +(.*?) +(.*?) +(.*?) +(.*?) +(.*?):(.*?)$")
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
			var player Player
			player.ID = match[2]
			player.IP = match[8]
			player.Name = match[1]
			player.State = match[6]

			ping, err := strconv.ParseUint(match[4], 10, 32)
			if err != nil {
				log.Errorf("Could not extract player ping from RCON `status` response: (%s) %+v", line, err)
			} else {
				player.Ping = uint32(ping)
			}

			// Merge duration and score

			var del = -1
			for i, score := range scores {
				if player.Name == score.Name {
					player.Duration = score.Duration
					player.Score = score.Score
					del = i
					break
				}
			}

			// Delete player from scores if match was found

			if del >= 0 {
				scores = append(scores[:del], scores[del + 1:]...)
			}

			// Fill in data from GeoIP

			ip := net.ParseIP(player.IP)
			record, err := geo.City(ip)
			if err != nil {
				log.Errorf("GeoIP error: %+v", err)
			}
			player.City       = record.City.Names["en"]
			player.Country    = record.Country.Names["en"]
			player.CountryISO = record.Country.IsoCode
			player.Latitude   = record.Location.Latitude
			player.Longitude  = record.Location.Longitude
			player.TimeZone   = record.Location.TimeZone

			// Save the player

			players = append(players, player)
		}
	}

	// Done

	return players, nil
}

func get_scores() ([]Player, error) {
	var players []Player

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

	// Convert response to Player slice

	for _, r := range resp.Players {
		var p Player
		p.Name     = r.Name
		p.Score    = r.Score
		p.Duration = r.Duration
		players = append(players, p)
	}

	// Done

	return players, nil
}

func get_steam(steamid64 string) (*Steam, error) {

	// Create new HTTP request

	req, err := http.NewRequest("GET", "https://api.steampowered.com/ISteamUser/GetPlayerSummaries/v2/", nil)
	if err != nil {
		return nil, fmt.Errorf("HTTP request error: %+v", err)
	}

	// Add request parameters

	q := req.URL.Query()
	q.Add("key", config.SteamKey)
	q.Add("steamids", steamid64)
	req.URL.RawQuery = q.Encode()

	// Execute request

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP execute error: %+v", err)
	}

	// Check response status

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Steam API response status = %s", resp.Status)
	}

	// Read response body

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("HTTP response read error: %+v", err)
	}

	// Decode response

	type Response struct {
		Players []Steam `json:"players"`
	}
	type Data struct {
		Response Response `json:"response"`
	}
	data := Data{}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, fmt.Errorf("HTTP response unmarshal error: %+v", err)
	}

	// Validate

	if len(data.Response.Players) != 1 {
		return nil, fmt.Errorf("Bad response array length: %d", len(data.Response.Players))
	}

	// Done

	steam := data.Response.Players[0]

	return &steam, nil
}

func rcon_command(command string, check string) ([]string, error) {
	var status []string

	// Get RCON command response

	req, err := rcon.Dial(fmt.Sprintf("%s:%d", config.ServerAddress, config.ServerPort))
	if err != nil {
		return nil, fmt.Errorf("RCON dial error: %+v", err)
	}
	defer req.Close()

	err = req.Authenticate(config.ServerPassword)
	if err != nil {
		return nil, fmt.Errorf("RCON authentication error: %+v", err)
	}

	resp, err := req.Execute(command)
	if err == io.EOF && len(check) == 0 {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("RCON `%s` command error: %+v", command, err)
	}

	// Convert response to array of strings

	scanner := bufio.NewScanner(strings.NewReader(resp.Body))
	for scanner.Scan() {
		status = append(status, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("RCON response parsing error: %+v", err)
	}

	// Check for empty response

	if len(status) == 0 {
		return nil, fmt.Errorf("RCON response is empty")
	}

	// Check command execution via regex on any line of response

	found := false
	re := regexp.MustCompile(check)

	for _, line := range status {
		match := re.FindStringSubmatch(line)
		if match != nil {
			found = true
		}
	}

	if !found {
		return nil, fmt.Errorf("RCON response check failed: %s", status)
	}

	// Done

	return status, nil
}

package main

type Map struct {
	Name           string `json:"name"`
	Map            string `json:"map"`
}

type Site struct {
	URL            string `json:"url"`
	Logo           string `json:"logo"`
	Title          string `json:"title"`
}

type Config struct {
	ApiPort        int    `json:"api_port"`
	GeoIP2_DB      string `json:"geoip2_db"`
	Site           Site   `json:"site"`

	GoogleKey      string `json:"google_key"`
	SteamKey       string `json:"steam_key"`

	ServerAddress  string `json:"server_address"`
	ServerPort     int    `json:"server_port"`
	ServerPassword string `json:"server_password"`

	Maps           []Map  `json:"maps"`
}

var config Config

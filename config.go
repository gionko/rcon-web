package main

type Site struct {
	URL         string `json:"url"`
	Static      string `json:"static"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type Config struct {
	ApiPort        int    `json:"api_port"`
	GeoIP2_DB      string `json:"geoip2_db"`
	Site           Site   `json:"site"`

	ServerAddress  string `json:"server_address"`
	ServerPort     int    `json:"server_port"`
	ServerPassword string `json:"server_password"`
}

var config Config

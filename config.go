package main

type Config struct {
	ApiPort        int    `json:"api_port"`

	ServerAddress  string `json:"server_address"`
	ServerPort     int    `json:"server_port"`
	ServerPassword string `json:"server_password"`
}

var config Config

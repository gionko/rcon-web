package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/user"

	"github.com/gin-gonic/gin"
	"github.com/rumblefrog/go-a2s"
)

const version = "1.0.0"

func init() {
	flag.Usage = usage
}

func usage() {
	fmt.Printf("Usage: %s [-config file] [-version]\n", os.Args[0])
	flag.PrintDefaults()
}

func a2s_test() {
	client, err := a2s.NewClient("46.183.163.222:27015")
	if err != nil {
		log.Infof("Error: %+v", err)
		os.Exit(1)
	}
	defer client.Close()

	info, err := client.QueryInfo()
	if err != nil {
		log.Infof("Error: %+v", err)
		os.Exit(1)
	}
	log.Infof("ServerInfo: %+v", info)
	log.Infof("\tExtendedServerInfo: %+v", info.ExtendedServerInfo)

	player, err := client.QueryPlayer()
	if err != nil {
		log.Infof("Error: %+v", err)
		os.Exit(1)
	}
	fmt.Printf("PlayerInfo: %+v", player)
	for i, p := range player.Players {
		log.Infof("\tPlayer %d: %+v", i, p)
	}

	rules, err := client.QueryRules()
	if err != nil {
		log.Infof("Error: %+v", err)
		os.Exit(1)
	}
	log.Infof("Rules: %+v", rules)
}

func main() {
	// TODO: remove me
	a2s_test()

	arg_config := new(string)

	// Parse command line arguments

	arg_version := flag.Bool(  "version", false, "Show version information")
	arg_config   = flag.String("config",  "",    "Config file")
	flag.Parse()

	// Set config filename if it was not provided

	if *arg_config == "" {
		// Try to extract config filename from RCON_CONF environment variable
		env, set := os.LookupEnv("RCON_CONF")
		if set {
			*arg_config = env
		} else {
			// No config argument & no RCON_CONF are set, use default one
			usr, err := user.Current()
			if err != nil {
				log.Errorf("Could not get current user: %+v", err)
				os.Exit(1)
			}
			*arg_config = fmt.Sprintf("%s/.rconrc", usr.HomeDir)
		}
	}

	// Show version info before parsing the configuration file

	if *arg_version {
		fmt.Printf("%s %s\n", os.Args[0], version)
		os.Exit(0)
	}

	// Read configuration file

	data, err := ioutil.ReadFile(*arg_config)
	if err != nil {
		log.Errorf("Could not read configuration file: %+v", err)
		os.Exit(1)
	}

		// Parse json data into Config structure

		err = json.Unmarshal(data, &config)
		if err != nil {
			log.Errorf("Error parsing configuration file: %+v", err)
			os.Exit(1)
		}

	// Suppress Gin debug output

	gin.SetMode(gin.ReleaseMode)

	// Create Gin instance

	router := gin.New()
	log.Debug("Gin instance created")

	// Set default middleware

	router.Use(gin.Recovery())
	log.Debug("Recovery middleware set")

	// Set logging middleware

	router.Use(LogMiddleware())
	log.Debug("Logging middleware set")

	// Create routing group, remove later

	group := router.Group("/v2")

	// Declare assets

	group.Static("/static", "static")
	group.StaticFile("/favicon.ico", "static/favicon.ico")
	group.StaticFile("/robots.txt", "static/robots.txt")
	log.Debug("Assets initialized")

	// Set routes

	group.GET("/users", RouteUsers)
	log.Debug("Routes set")

	// Start the server

	log.Info("Server started!")
	log.Infof("Port: %d", config.ApiPort)
	err = http.ListenAndServe(fmt.Sprintf(":%d", config.ApiPort), router)
	log.Fatalf("HTTP server error: %s", err)
}
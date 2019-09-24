package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/user"

	"github.com/gin-contrib/multitemplate"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

const version = "1.0.0"

func init() {
	flag.Usage = usage
}

func usage() {
	fmt.Printf("Usage: %s [-config file] [-version]\n", os.Args[0])
	flag.PrintDefaults()
}

func main() {
	// Parse command line arguments

	arg_version := flag.Bool(  "version", false, "Show version information")
	arg_config  := flag.String("config",  "",    "Config file")
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

	// Set session store
	store := cookie.NewStore([]byte("die, bitch"))
	router.Use(sessions.Sessions("rcon-web", store))

	// Set multi-template renderer

	router.HTMLRender = multi_render()

	// Create routing group, remove later

	group := router.Group("/")

	// Declare assets

	group.Static("/static", "static")
	group.StaticFile("/favicon.ico", "static/favicon.ico")
	group.StaticFile("/robots.txt", "static/robots.txt")
	log.Debug("Assets initialized")

	// Set API routes

	group.GET   ("/api/bots/:id",       RouteAPIBots)
	group.GET   ("/api/damage/:id",     RouteAPIDamage)
	group.GET   ("/api/difficulty/:id", RouteAPIDifficulty)
	group.POST  ("/api/login",          RouteAPILogin)
	group.GET   ("/api/logout",         RouteAPILogout)
	group.GET   ("/api/maps",           RouteAPIMaps)
	group.GET   ("/api/maps/:id",       RouteAPIMap)
	group.GET   ("/api/players",        RouteAPIPlayers)
	group.GET   ("/api/players/:id",    RouteAPIPlayer)
	group.PUT   ("/api/players/:id",    RouteAPIPlayersBan)
	group.DELETE("/api/players/:id",    RouteAPIPlayersKick)
	group.GET   ("/api/status",         RouteAPIStatus)
	log.Debug("API routes set")

	// Set frontend routes

	group.GET("/",            RouteFEIndex)
	group.GET("/bots",        RouteFEBots)
	group.GET("/login",       RouteFELogin)
	group.GET("/maps",        RouteFEMaps)
	group.GET("/players",     RouteFEPlayers)
	group.GET("/players/:id", RouteFEPlayer)
	log.Debug("Frontend routes set")

	// Start the server

	log.Info("Server started!")
	log.Infof("Port: %d", config.ApiPort)
	err = http.ListenAndServe(fmt.Sprintf(":%d", config.ApiPort), router)
	log.Fatalf("HTTP server error: %s", err)
}

func multi_render() multitemplate.Render {
	r := multitemplate.New()
	r.AddFromFiles("bots",    "templates/base.tmpl", "templates/bots.tmpl")
	r.AddFromFiles("index",   "templates/base.tmpl", "templates/index.tmpl")
	r.AddFromFiles("login",   "templates/base.tmpl", "templates/login.tmpl")
	r.AddFromFiles("maps",    "templates/base.tmpl", "templates/maps.tmpl")
	r.AddFromFiles("player",  "templates/base.tmpl", "templates/player.tmpl")
	r.AddFromFiles("players", "templates/base.tmpl", "templates/players.tmpl")
	return r
}

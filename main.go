package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jeromelesaux/ethereum-training/client"
	"github.com/jeromelesaux/ethereum-training/config"
	"github.com/jeromelesaux/ethereum-training/controller"
)

var (
	configFile            = flag.String("config", "", "Configuration file path.")
	displayHelp           = flag.Bool("help", false, "Display help message and quit.")
	configurationFilepath string
)

func main() {

	flag.Parse()
	if *displayHelp {
		help()
	}

	client.SafeNonceTx = &client.SafeNonce{}

	// laod main configuration
	if *configFile != "" {
		configurationFilepath = *configFile
		config.LoadConfigFile(configurationFilepath)
	} else {
		config.LoadConfig()
	}

	// authenticate on ethereum platforms
	client.Authenticate()

	// router creation
	router := gin.Default()

	// controller with routes definition
	controller := &controller.Controller{}
	router.LoadHTMLGlob("resources/*.html") // add static html files
	router.StaticFile("/logo-innovation-lab-v2.png", "resources/logo-innovation-lab-v2.png")
	// add certifications api
	router.POST("/anchor", controller.Anchoring)
	router.POST("/verify", controller.Verify)
	router.POST("/anchormultiple", controller.AnchorMultiple)
	router.POST("/verifymultiple", controller.VerifyMultiple)

	// add static html route
	root := router.Group("/")
	root.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	// start server at port 8080
	if err := router.Run(":8080"); err != nil {
		fmt.Fprintf(os.Stderr, "Can not start server error :%v\n", err)
		os.Exit(-1)
	}
}

func help() {
	flag.PrintDefaults()
	os.Exit(-1)
}

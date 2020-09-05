package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jeromelesaux/ethereum-training/client"
	"github.com/jeromelesaux/ethereum-training/config"
	"github.com/jeromelesaux/ethereum-training/controller"
	"github.com/jeromelesaux/ethereum-training/http/header"
	"github.com/jeromelesaux/ethereum-training/persistence"
	"github.com/jeromelesaux/ethereum-training/token"
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
	controller.LoadCredentials()
	if err := persistence.Initialise(); err != nil {
		fmt.Fprintf(os.Stderr, "Error while initialise the database with error (%v)\n", err)
		os.Exit(-1)
	}

	// authenticate on ethereum platforms
	client.Authenticate()

	// creation des cookies stores et token pour google oauth
	token, err := token.RandToken(64)
	if err != nil {
		log.Fatal("unable to generate random token: ", err)
	}
	store := sessions.NewCookieStore([]byte(token))
	store.Options(sessions.Options{
		Path:   "/",
		MaxAge: 86400 * 7,
	})

	// router creation
	router := gin.Default()

	// google oauth controller
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(header.NocacheHeaders(), sessions.Sessions("goquestsession", store))

	// controller with routes definition
	controller := &controller.Controller{}
	//router.LoadHTMLGlob("resources/*.html") // add static html files
	router.LoadHTMLGlob("resources/*.tmpl") // add static html files
	router.StaticFile("/logo-innovation-lab-v2.png", "resources/logo-innovation-lab-v2.png")
	router.StaticFile("/uploader.css", "resources/uploader.css")
	// add certifications api
	router.POST("/verify", controller.Verify)
	router.POST("/verifymultiple", controller.VerifyMultiple)
	router.GET("/login", controller.LoginHandler)
	router.GET("/auth", controller.AuthHandler)
	router.GET("/verification", controller.Verification)
	// add static html route
	root := router.Group("/")
	root.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", nil)
	})

	authorized := router.Group("/api")
	authorized.Use(controller.AuthorizeRequest(), header.NocacheHeaders())
	authorized.POST("/anchor", controller.Anchoring)
	authorized.POST("/anchormultiple", controller.AnchorMultiple)
	authorized.GET("/txhash", controller.GetFile)
	authorized.GET("/safebox", controller.Safebox)
	authorized.GET("/certification", controller.Certification)

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

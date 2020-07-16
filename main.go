package main

import (
	"fmt"
	"os"

	"github.com/jeromelesaux/ethereum-training/client"

	"github.com/jeromelesaux/ethereum-training/config"

	"github.com/gin-gonic/gin"
	"github.com/jeromelesaux/ethereum-training/controller"
)

func main() {

	client.SafeNonceTx = &client.SafeNonce{}

	// laod main configuration
	config.LoadConfig()

	// authenticate on ethereum platforms
	client.Authenticate()

	// router creation
	router := gin.Default()

	// controller with routes definition
	controller := &controller.Controller{}
	router.POST("/anchor", controller.Anchoring)
	router.POST("/verify", controller.Verify)
	router.POST("/anchormultiple", controller.AnchorMultiple)
	router.POST("/verifymultiple", controller.VerifyMultiple)
	// start server at port 8080
	if err := router.Run(":8080"); err != nil {
		fmt.Fprintf(os.Stderr, "Can not start server error :%v\n", err)
		os.Exit(-1)
	}
}

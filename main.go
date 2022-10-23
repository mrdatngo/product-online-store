package main

import (
	route "github.com/mrdatngo/gin-products/routes"
	util "github.com/mrdatngo/gin-products/utils"
	"log"
)

func main() {
	/**
	@description Setup Server
	*/
	router := route.SetupRouter()

	/**
	@description Run Server
	*/
	log.Fatal(router.Run(":" + util.GodotEnv("GO_PORT")))
}

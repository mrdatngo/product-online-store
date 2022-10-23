package config

import (
	"fmt"
	"os"

	util "github.com/mrdatngo/gin-products/utils"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Connection() *gorm.DB {
	databaseURI := ""

	if os.Getenv("GO_ENV") != "production" {
		databaseURI = util.GodotEnv("DATABASE_URI_DEV")
	} else {
		databaseURI = os.Getenv("DATABASE_URI_PROD")
	}

	fmt.Printf("databaseURI: %v", databaseURI)

	db, err := gorm.Open(mysql.Open(databaseURI), &gorm.Config{})

	if err != nil {
		defer logrus.Info("Connection to Database Failed")
		logrus.Fatal(err.Error())
	}

	if os.Getenv("GO_ENV") != "production" {
		logrus.Info("Connection to Database Successfully")
	}
	return db
}

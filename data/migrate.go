package main

import (
	"fmt"
	model "github.com/mrdatngo/gin-products/models"
	util "github.com/mrdatngo/gin-products/utils"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
	"strings"
)

func main() {
	databaseURI := make(chan string, 1)

	if os.Getenv("GO_ENV") != "production" {
		databaseURI <- util.GodotEnv("DATABASE_URI_DEV")
	} else {
		databaseURI <- os.Getenv("DATABASE_URI_PROD")
	}

	fmt.Println(databaseURI)

	db, err := gorm.Open(mysql.Open(<-databaseURI), &gorm.Config{})

	if err != nil {
		defer logrus.Info("Connection to Database Failed")
		logrus.Fatal(err.Error())
	}

	if os.Getenv("GO_ENV") != "production" {
		logrus.Info("Connection to Database Successfully")
	}
	AutoMigrate(db)
	InitData(db)
}

func AutoMigrate(db *gorm.DB) {
	err := db.AutoMigrate(
		&model.EntityBranch{},
		&model.EntityProduct{},
		&model.EntityOption{},
		&model.EntityOptionValue{},
		&model.EntitySkuValue{},
		&model.EntityProductSku{},
	)
	if err != nil {
		logrus.Fatal(err.Error())
	}
}

func InitData(db *gorm.DB) {
	err := loadSQLFile(db, "data/init-data.sql")
	if err != nil {
		logrus.Fatal(err)
	}
}

func loadSQLFile(db *gorm.DB, sqlFile string) error {
	file, err := os.ReadFile(sqlFile)
	if err != nil {
		return err
	}
	tx := db.Begin()
	defer func() {
		tx.Rollback()
	}()
	for _, q := range strings.Split(string(file), ";") {
		q := strings.TrimSpace(q)
		if q == "" {
			continue
		}
		tx.Exec(q)
	}
	tx.Commit()
	return nil
}

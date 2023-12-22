package db

import (
	"log"
	"os"

	"github.com/RuanScherer/journey-track-api/config"
	"github.com/RuanScherer/journey-track-api/domain/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

func GetConnection() *gorm.DB {
	if db == nil {
		db = connect()
	}
	return db
}

func connect() *gorm.DB {
	appConfig := config.GetAppConfig()

	gormConfig := &gorm.Config{}
	if appConfig.DbLogEnabled {
		gormConfig.Logger = logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{},
		)
	}

	db, err := gorm.Open(postgres.Open(appConfig.DbDsn), gormConfig)
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}

	err = db.AutoMigrate(&model.User{}, &model.Project{}, &model.ProjectInvite{}, &model.Event{})
	if err != nil {
		panic("failed to migrate database: " + err.Error())
	}

	return db
}

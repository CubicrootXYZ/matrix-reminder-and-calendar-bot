package database

import (
	"github.com/CubicrootXYZ/gormlogger"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/configuration"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlog "gorm.io/gorm/logger"
)

// Database holds all information for connecting to the database
type Database struct {
	config configuration.Database
	db     gorm.DB
	debug  bool
}

// Create creates a database object
func Create(config configuration.Database, debug bool) (*Database, error) {
	db := Database{
		config: config,
		debug:  debug,
	}

	err := db.connect()
	if err != nil {
		return nil, err
	}
	err = db.initialize()

	return &db, err
}

func (d *Database) connect() error {
	logger := gormlogger.NewLogger()

	db, err := gorm.Open(mysql.Open(d.config.Connection+"?parseTime=True"), &gorm.Config{
		Logger: logger,
	})
	if err != nil {
		return err
	}

	if !d.debug {
		db.Logger.LogMode(gormlog.Warn)
	}

	d.db = *db
	return nil
}

func (d *Database) initialize() error {
	err := d.db.AutoMigrate(&Reminder{})
	if err != nil {
		return err
	}

	err = d.db.AutoMigrate(&Channel{})
	if err != nil {
		return err
	}

	err = d.db.AutoMigrate(&Message{})
	if err != nil {
		return err
	}

	err = d.db.AutoMigrate(&Event{})
	if err != nil {
		return err
	}

	return nil
}

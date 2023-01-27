package database_test

import (
	"os"
	"testing"

	"github.com/CubicrootXYZ/gologger"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/database"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var logger gologger.Logger
var service database.Service
var gormDB *gorm.DB

func getConnection() string {
	host := os.Getenv("TEST_DB_HOST")
	if host == "" {
		host = "localhost"
	}

	return "root:mypass@tcp(" + host + ":3306)/remindme"
}

func getGormDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open(getConnection()+"?parseTime=True"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	return db
}

func getService(gormDB *gorm.DB) database.Service {
	service, err := database.New(gormDB)
	if err != nil {
		panic(err)
	}

	return service
}

func getLogger() gologger.Logger {
	return gologger.New(gologger.LogLevelDebug, 0)
}

func TestMain(m *testing.M) {
	logger = getLogger()
	gormDB = getGormDB()
	service = getService(gormDB)

	m.Run()
}

func testUser() *database.MatrixUser {
	return &database.MatrixUser{
		ID:    "@remindme:example.org",
		Rooms: []database.MatrixRoom{*testRoom()},
	}
}

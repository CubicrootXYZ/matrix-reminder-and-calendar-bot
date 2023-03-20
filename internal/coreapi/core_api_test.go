package coreapi_test

import (
	"net/http/httptest"
	"time"

	"github.com/CubicrootXYZ/gologger"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/api/middleware"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/coreapi"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
)

func testCoreAPI(ctrl *gomock.Controller) (*httptest.Server, *database.MockService) {
	logger := gologger.New(gologger.LogLevelDebug, 0)
	db := database.NewMockService(ctrl)

	api := coreapi.New(&coreapi.Config{
		Database:            db,
		DefaultAuthProvider: middleware.APIKeyAuth("123"),
	}, logger)

	r := gin.New()
	err := api.RegisterRoutes(r)
	if err != nil {
		panic(err)
	}
	server := httptest.NewServer(r)

	return server, db
}

func testDatabaseChannel() database.Channel {
	dailyReminder := uint(130)
	created, _ := time.Parse(time.RFC3339, "2006-01-02T15:04:05+07:00")

	c := database.Channel{
		Description:   "chan desc",
		DailyReminder: &dailyReminder,
		TimeZone:      "Europe/Berlin",
	}

	c.ID = 1
	c.CreatedAt = created
	return c
}

func testChannel() coreapi.Channel {
	dailyReminder := "02:10"
	tz := "Europe/Berlin"

	return coreapi.Channel{
		ID:            1,
		CreatedAt:     "2006-01-02T15:04:05+07:00",
		Description:   "chan desc",
		DailyReminder: &dailyReminder,
		TimeZone:      &tz,
	}
}

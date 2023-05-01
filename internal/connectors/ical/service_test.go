package ical_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/CubicrootXYZ/gologger"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/ical"
	icaldb "github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/ical/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/tj/assert"
	"gorm.io/gorm"
)

func testService(ctrl *gomock.Controller) (ical.Service, *icaldb.MockService, *database.MockService) {
	db := database.NewMockService(ctrl)
	icalDB := icaldb.NewMockService(ctrl)
	return ical.New(&ical.Config{
			Database: db,
			ICalDB:   icalDB,
		}, gologger.New(gologger.LogLevelDebug, 0)),
		icalDB,
		db
}

func TestService_InputRemoved(t *testing.T) {
	ctrl := gomock.NewController(t)
	service, icalDB, _ := testService(ctrl)

	icalDB.EXPECT().DeleteIcalInput(uint(1)).Return(nil)

	err := service.InputRemoved("ical", 1)
	require.NoError(t, err)
}

func TestService_InputRemovedWithNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	service, icalDB, _ := testService(ctrl)

	icalDB.EXPECT().DeleteIcalInput(uint(1)).Return(icaldb.ErrNotFound)

	err := service.InputRemoved("ical", 1)
	require.NoError(t, err)
}

func TestService_InputRemovedWithWrongType(t *testing.T) {
	ctrl := gomock.NewController(t)
	service, _, _ := testService(ctrl)

	err := service.InputRemoved("notical", 1)
	require.NoError(t, err)
}

func TestService_InputRemovedWithError(t *testing.T) {
	ctrl := gomock.NewController(t)
	service, icalDB, _ := testService(ctrl)

	icalDB.EXPECT().DeleteIcalInput(uint(1)).Return(errors.New("test"))

	err := service.InputRemoved("ical", 1)
	require.Error(t, err)
}

func TestService_OutputRemoved(t *testing.T) {
	ctrl := gomock.NewController(t)
	service, icalDB, _ := testService(ctrl)

	icalDB.EXPECT().DeleteIcalOutput(uint(1)).Return(nil)

	err := service.OutputRemoved("ical", 1)
	require.NoError(t, err)
}

func TestService_OutputRemovedWithNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	service, icalDB, _ := testService(ctrl)

	icalDB.EXPECT().DeleteIcalOutput(uint(1)).Return(icaldb.ErrNotFound)

	err := service.OutputRemoved("ical", 1)
	require.NoError(t, err)
}

func TestService_OutputRemovedWithWrongType(t *testing.T) {
	ctrl := gomock.NewController(t)
	service, _, _ := testService(ctrl)

	err := service.OutputRemoved("notical", 1)
	require.NoError(t, err)
}

func TestService_OutputRemovedWithError(t *testing.T) {
	ctrl := gomock.NewController(t)
	service, icalDB, _ := testService(ctrl)

	icalDB.EXPECT().DeleteIcalOutput(uint(1)).Return(errors.New("test"))

	err := service.OutputRemoved("ical", 1)
	require.Error(t, err)
}

func TestService_NewOutput(t *testing.T) {
	ctrl := gomock.NewController(t)
	service, icalDB, db := testService(ctrl)

	icalDB.EXPECT().NewIcalOutput(gomock.Any()).Return(&icaldb.IcalOutput{
		Model: gorm.Model{
			ID: 1,
		},
	}, nil)

	db.EXPECT().AddOutputToChannel(uint(2), &database.Output{
		ChannelID:  2,
		OutputType: "ical",
		OutputID:   1,
		Enabled:    true,
	}).Return(nil)

	output, err := service.NewOutput(2)
	require.NoError(t, err)

	assert.Equal(t, uint(1), output.ID)
}

func TestService_NewOutputWithAddError(t *testing.T) {
	ctrl := gomock.NewController(t)
	service, icalDB, db := testService(ctrl)

	icalDB.EXPECT().NewIcalOutput(gomock.Any()).Return(&icaldb.IcalOutput{
		Model: gorm.Model{
			ID: 1,
		},
	}, nil)

	db.EXPECT().AddOutputToChannel(uint(2), &database.Output{
		ChannelID:  2,
		OutputType: "ical",
		OutputID:   1,
		Enabled:    true,
	}).Return(errors.New("test"))

	_, err := service.NewOutput(2)
	require.Error(t, err)
}

func TestService_NewOutputWithNewIcalError(t *testing.T) {
	ctrl := gomock.NewController(t)
	service, icalDB, _ := testService(ctrl)

	icalDB.EXPECT().NewIcalOutput(gomock.Any()).Return(nil, errors.New("test"))

	_, err := service.NewOutput(2)
	require.Error(t, err)
}

func TestService_GetOutput(t *testing.T) {
	ctrl := gomock.NewController(t)
	service, icalDB, _ := testService(ctrl)

	icalDB.EXPECT().GetIcalOutputByID(uint(1)).Return(&icaldb.IcalOutput{
		Model: gorm.Model{
			ID: 1,
		},
	}, nil)

	output, err := service.GetOutput(1)
	require.NoError(t, err)

	assert.Equal(t, uint(1), output.ID)
}

func TestService_GetOutputWithNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	service, icalDB, _ := testService(ctrl)

	icalDB.EXPECT().GetIcalOutputByID(uint(1)).Return(nil, icaldb.ErrNotFound)

	_, err := service.GetOutput(1)
	require.ErrorIs(t, err, ical.ErrNotFound)
}

func TestService_GetOutputWithError(t *testing.T) {
	ctrl := gomock.NewController(t)
	service, icalDB, _ := testService(ctrl)

	icalDB.EXPECT().GetIcalOutputByID(uint(1)).Return(nil, errors.New("test"))

	_, err := service.GetOutput(1)
	require.Error(t, err)
}

func TestService_Fetcher(t *testing.T) {
	ctrl := gomock.NewController(t)
	service, icalDB, db := testService(ctrl)

	content, err := os.ReadFile("format/testdata/calendar1.ical")
	require.NoError(t, err)

	called := false
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(content)
		called = true
	}))

	f := false
	icalDB.EXPECT().ListIcalInputs(&icaldb.ListIcalInputsOpts{
		Disabled: &f,
	}).Return([]icaldb.IcalInput{
		{
			Model: gorm.Model{
				ID: 1,
			},
			URL: server.URL + "/",
		},
	}, nil)

	db.EXPECT().GetInputByType(uint(1), "ical").Return(&database.Input{
		Model: gorm.Model{
			ID: 2,
		},
		ChannelID: 3,
	}, nil)

	inputID := uint(2)
	// TODO ensure we get events
	db.EXPECT().NewEvents([]database.Event{
		{
			Time:      testTime().UTC(),
			Message:   "Event 1",
			Active:    true,
			Duration:  time.Minute * 5,
			ChannelID: 3,
			InputID:   &inputID,
		},
		{
			Time:      testTime().UTC(),
			Message:   "Event 2",
			Active:    true,
			Duration:  time.Minute * 5,
			ChannelID: 3,
			InputID:   &inputID,
		},
		{
			Time:      testTime().UTC(),
			Message:   "Event 3",
			Active:    true,
			Duration:  time.Minute * 5,
			ChannelID: 3,
			InputID:   &inputID,
		},
	}).Return(nil)
	icalDB.EXPECT().UpdateIcalInput(gomock.Any()).Return(nil, nil)

	go func() {
		require.NoError(t, service.Start())
	}()

	time.Sleep(time.Millisecond * 50)
	require.NoError(t, service.Stop())
	time.Sleep(time.Millisecond * 10)

	assert.True(t, called, "testserver was not called once")
}

func testTime() time.Time {
	t, _ := time.Parse(time.RFC3339, "2120-01-02T15:04:05+00:00")
	return t
}

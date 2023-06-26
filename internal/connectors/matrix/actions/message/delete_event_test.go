package message_test

import (
	"testing"
	"time"

	"github.com/CubicrootXYZ/gologger"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/ical"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/actions/message"
	matrixdb "github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/mautrixcl"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/messenger"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/tests"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestDeleteEventAction(t *testing.T) {
	action := &message.DeleteEventAction{}

	assert.NotEmpty(t, action.Name())

	title, desc, examples := action.GetDocu()
	assert.NotEmpty(t, title)
	assert.NotEmpty(t, desc)
	assert.NotEmpty(t, examples)

	assert.NotNil(t, action.Selector())
}

func TestDeleteEventAction_Selector(t *testing.T) {
	action := &message.DeleteEventAction{}
	r := action.Selector()

	_, _, examples := action.GetDocu()
	for _, example := range examples {
		assert.True(t, r.MatchString(example))
	}
}

func TestDeleteEventAction_HandleEvent(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	db := database.NewMockService(ctrl)
	matrixDB := matrixdb.NewMockService(ctrl)
	client := mautrixcl.NewMockClient(ctrl)
	msngr := messenger.NewMockMessenger(ctrl)
	icalBridge := ical.NewMockService(ctrl)

	action := &message.DeleteEventAction{}
	action.Configure(
		gologger.New(gologger.LogLevelDebug, 0),
		client,
		msngr,
		matrixDB,
		db,
		&matrix.BridgeServices{
			ICal: icalBridge,
		},
	)

	matrixDB.EXPECT().NewMessage(&matrixdb.MatrixMessage{
		ID:            "evt1",
		UserID:        toP("@user:example.com"),
		Body:          `delete 123`,
		BodyFormatted: `delete 123`,
		Type:          matrixdb.MessageTypeEventDelete,
		EventID:       toP(uint(123)),
		Incoming:      true,
		SendAt:        time.UnixMilli(tests.TestEvent().Event.Timestamp),
	},
	).Return(nil, nil)

	db.EXPECT().ListEvents(&database.ListEventsOpts{
		IDs:       []uint{123},
		ChannelID: toP(uint(68272)),
	}).Return([]database.Event{
		mockEvent(),
	}, nil)

	db.EXPECT().DeleteEvent(toP(mockEvent())).Return(nil)

	msngr.EXPECT().SendResponse(&messenger.Response{
		Message:                   "Deleted event \"\"",
		MessageFormatted:          "Deleted event \"\"",
		RespondToMessage:          "delete 123",
		RespondToMessageFormatted: "delete 123",
		RespondToUserID:           "@user:example.com",
		RespondToEventID:          "evt1",
		ChannelExternalIdentifier: "!room123",
	}).Return(&messenger.MessageResponse{
		ExternalIdentifier: "id1",
	}, nil)

	matrixDB.EXPECT().NewMessage(&matrixdb.MatrixMessage{
		ID:            "id1",
		UserID:        toP("@user:example.com"),
		Body:          `Deleted event ""`,
		BodyFormatted: `Deleted event ""`,
		Type:          matrixdb.MessageTypeEventDelete,
		EventID:       toP(uint(123)),
	},
	).Return(nil, nil)

	action.HandleEvent(tests.TestEvent(tests.WithBody("delete 123", "delete 123")))

	// Wait for async message sending.
	time.Sleep(time.Millisecond * 10)
}

func mockEvent() database.Event {
	return database.Event{
		Model: gorm.Model{
			ID: 123,
		},
	}
}

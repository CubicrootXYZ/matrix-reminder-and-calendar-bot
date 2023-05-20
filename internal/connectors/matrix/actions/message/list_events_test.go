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
	dbtests "github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database/tests"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestListEventsAction_Meta(t *testing.T) {
	action := &message.ListEventsAction{}
	assert.Greater(t, len(action.Name()), 2)

	title, expl, examples := action.GetDocu()
	assert.Greater(t, len(title), 2)
	assert.Greater(t, len(expl), 2)
	assert.Greater(t, len(examples), 0)
}

func TestListEventsAction_Selector(t *testing.T) {
	action := &message.ListEventsAction{}
	r := action.Selector()

	_, _, examples := action.GetDocu()

	for _, example := range examples {
		assert.True(t, r.MatchString(example))
	}
}

func TestListEventsAction_HandleEvent(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	db := database.NewMockService(ctrl)
	matrixDB := matrixdb.NewMockService(ctrl)
	client := mautrixcl.NewMockClient(ctrl)
	msngr := messenger.NewMockMessenger(ctrl)
	icalBridge := ical.NewMockService(ctrl)

	action := &message.ListEventsAction{}
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

	db.EXPECT().ListEvents(&database.ListEventsOpts{
		ChannelID: toP(uint(68272)),
	}).Return([]database.Event{
		dbtests.TestEvent(),
	}, nil)

	matrixDB.EXPECT().NewMessage(tests.TestMessage(
		tests.WithFromTestEvent(),
		tests.WithMessageType(matrixdb.MessageTypeEventList),
	)).Return(nil, nil)

	msngr.EXPECT().SendMessage(&messenger.Message{
		Body: `➡️ TEST EVENT
at 08:04 02.01.2006 (UTC) (ID: 2824) 
`,
		BodyHTML:                  `➡️ <b>test event</b><br>at 08:04 02.01.2006 (UTC) (ID: 2824) <br>`,
		ChannelExternalIdentifier: "!room123",
	}).Return(&messenger.MessageResponse{
		ExternalIdentifier: "!234",
	}, nil)

	matrixDB.EXPECT().NewMessage(&matrixdb.MatrixMessage{
		ID:     "!234",
		UserID: toP("@user:example.com"),
		Body: `➡️ TEST EVENT
at 08:04 02.01.2006 (UTC) (ID: 2824) 
`,
		BodyFormatted: `➡️ <b>test event</b><br>at 08:04 02.01.2006 (UTC) (ID: 2824) <br>`,
		Type:          matrixdb.MessageTypeEventList,
		RoomID:        0,
	})

	action.HandleEvent(tests.TestEvent())

	// Wait for async process.
	time.Sleep(time.Millisecond * 10)
}

func toP[T any](elem T) *T {
	return &elem
}

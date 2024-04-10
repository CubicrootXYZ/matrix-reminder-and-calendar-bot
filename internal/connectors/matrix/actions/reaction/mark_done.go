package reaction

import (
	"github.com/CubicrootXYZ/gologger"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix"
	matrixdb "github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/mautrixcl"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/messenger"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
)

// MarkDoneAction takes cafe of delete requests via reactions.
type MarkDoneAction struct {
	logger    gologger.Logger
	client    mautrixcl.Client
	messenger messenger.Messenger
	matrixDB  matrixdb.Service
	db        database.Service
}

// Configure is called on startup and sets all dependencies.
func (action *MarkDoneAction) Configure(logger gologger.Logger, client mautrixcl.Client, messenger messenger.Messenger, matrixDB matrixdb.Service, db database.Service, _ *matrix.BridgeServices) {
	action.logger = logger
	action.client = client
	action.matrixDB = matrixDB
	action.db = db
	action.messenger = messenger
}

// Name of the action.
func (action *MarkDoneAction) Name() string {
	return "Mark Event Done"
}

// GetDocu returns the documentation for the action.
func (action *MarkDoneAction) GetDocu() (title, explaination string, examples []string) {
	return "Mark Event Done",
		"React with a ✅ to mark the event as done.",
		[]string{"✅"}
}

// Selector defines on which reactions this action should be called.
func (action *MarkDoneAction) Selector() []string {
	return []string{"✅"}
}

// HandleEvent is where the reaction event and the related message get's send to if it matches the Selector.
func (action *MarkDoneAction) HandleEvent(event *matrix.ReactionEvent, reactionToMessage *matrixdb.MatrixMessage) {
	l := action.logger.WithFields(
		map[string]any{
			"reaction":        event.Content.RelatesTo.Key,
			"room":            reactionToMessage.RoomID,
			"related_message": reactionToMessage.ID,
			"user":            event.Event.Sender,
		},
	)
	if reactionToMessage.EventID == nil || reactionToMessage.Event == nil {
		l.Infof("skipping because message does not relate to any event")
		return
	}

	evt := reactionToMessage.Event
	evt.Active = false

	_, err := action.db.UpdateEvent(evt)
	if err != nil {
		l.Err(err)
		_ = action.messenger.SendMessageAsync(messenger.PlainTextMessage(
			"Whoopsie, can not update the event as requested.",
			event.Room.RoomID,
		))
		return
	}

	err = action.messenger.DeleteMessageAsync(&messenger.Delete{
		ExternalIdentifier:        reactionToMessage.ID,
		ChannelExternalIdentifier: reactionToMessage.Room.RoomID,
	})
	if err != nil {
		l.Err(err)
	}
}

package types

import (
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
)

// Messenger defines an interface for interacting with matrix messages
type Messenger interface {
	SendReplyToEvent(msg string, replyEvent *event.Event, channel *database.Channel, msgType database.MessageType) (resp *mautrix.RespSendEvent, err error)
	CreateChannel(userID string) (*mautrix.RespCreateRoom, error)
	SendFormattedMessage(msg, msgFormatted string, channel *database.Channel, msgType database.MessageType, relatedReminderID uint) (resp *mautrix.RespSendEvent, err error)
	DeleteMessage(messageID, roomID string) error
	SendNotice(msg, roomID string) (resp *mautrix.RespSendEvent, err error)
}

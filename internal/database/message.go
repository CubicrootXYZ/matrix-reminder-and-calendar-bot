package database

import (
	"time"

	"gorm.io/gorm"
	"maunium.net/go/mautrix/event"
)

// Message holds information about a single message
type Message struct {
	gorm.Model
	Body               string
	BodyHTML           string
	ReminderID         *uint
	Reminder           Reminder
	ResponseToMessage  string
	Type               MessageType
	ChannelID          uint
	Channel            Channel
	Timestamp          int64
	ExternalIdentifier string
}

// MessageType defines different types of messages
type MessageType string

// Message types differentiate the context of a message
const (
	// Reminder itself
	MessageTypeReminderRequest = MessageType("REMINDER_REQUEST")
	MessageTypeReminderSuccess = MessageType("REMINDER_SUCCESS")
	MessageTypeReminderFail    = MessageType("REMINDER_FAIL")
	MessageTypeReminder        = MessageType("REMINDER")
	// Arbitrary actions
	MessageTypeActions      = MessageType("ACTIONS")
	MessageTypeReminderList = MessageType("REMINDER_LIST")
	// Reminder edits
	MessageTypeReminderUpdate           = MessageType("REMINDER_UPDATE")
	MessageTypeReminderUpdateFail       = MessageType("REMINDER_UPDATE_FAIL")
	MessageTypeReminderUpdateSuccess    = MessageType("REMINDER_UPDATE_SUCCESS")
	MessageTypeReminderDelete           = MessageType("REMINDER_DELETE")
	MessageTypeReminderDeleteSuccess    = MessageType("REMINDER_DELETE_SUCCESS")
	MessageTypeReminderDeleteFail       = MessageType("REMINDER_DELETE_Fail")
	MessageTypeReminderRecurringRequest = MessageType("REMINDER_RECURRING_REQUEST")
	MessageTypeReminderRecurringSuccess = MessageType("REMINDER_RECURRING_SUCCESS")
	MessageTypeReminderRecurringFail    = MessageType("REMINDER_RECURRING_FAIL")
	// Settings
	MessageTypeTimezoneChangeRequest        = MessageType("TIMEZONE_CHANGE")
	MessageTypeTimezoneChangeRequestSuccess = MessageType("TIMEZONE_CHANGE_SUCCESS")
	MessageTypeTimezoneChangeRequestFail    = MessageType("TIMEZONE_CHANGE_FAIL")
	// Do not save!
	MessageTypeDoNotSave = MessageType("")
)

// AddMessageFromMatrix adds a message to the database
func (d *Database) AddMessageFromMatrix(id string, timestamp int64, content *event.MessageEventContent, reminder *Reminder, msgType MessageType, channel *Channel) (*Message, error) {
	relatesTo := ""
	if content.RelatesTo != nil {
		relatesTo = content.RelatesTo.EventID.String()
	}
	message := Message{
		ResponseToMessage:  relatesTo,
		Type:               msgType,
		ChannelID:          channel.ID,
		Channel:            *channel,
		Timestamp:          timestamp,
		ExternalIdentifier: id,
	}
	message.Model.CreatedAt = time.Now().UTC()

	if content != nil {
		message.Body = content.Body
		message.BodyHTML = content.FormattedBody
	}

	if reminder != nil {
		message.Reminder = *reminder
		message.ReminderID = &reminder.ID
	}

	err := d.db.Create(&message).Error

	return &message, err
}

// AddMessage adds a message to the database
func (d *Database) AddMessage(message *Message) (*Message, error) {

	err := d.db.Create(message).Error

	return message, err
}

// GetMessageByExternalID returns if found the message with the given external id
func (d *Database) GetMessageByExternalID(externalID string) (*Message, error) {
	message := &Message{}
	err := d.db.Preload("Reminder").First(&message, "external_identifier = ?", externalID).Error
	return message, err
}

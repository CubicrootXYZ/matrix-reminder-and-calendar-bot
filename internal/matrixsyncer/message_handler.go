package matrixsyncer

import (
	"fmt"
	"regexp"
	"strings"
	"time"
	"unicode"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/log"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/matrixmessenger"
	"github.com/tj/go-naturaldate"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
	"gorm.io/gorm"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/id"
)

// handles new messages
func (s *Syncer) handleMessages(source mautrix.EventSource, evt *event.Event) {
	log.Debug(fmt.Sprintf("New message: / Sender: %s / Room: / %s / Time: %d", evt.Sender, evt.RoomID, evt.Timestamp))

	// Do not answer our own and old messages
	if evt.Sender == id.UserID(s.botName) || evt.Timestamp/1000 <= time.Now().Unix()-60 {
		return
	}

	channel, err := s.daemon.Database.GetChannelByUserAndChannelIdentifier(evt.Sender.String(), evt.RoomID.String())

	content, ok := evt.Content.Parsed.(*event.MessageEventContent)
	if !ok {
		log.Warn("Event is not a message event. Can not handle it")
		return
	}

	// Unknown channel
	if err == gorm.ErrRecordNotFound || channel == nil {
		channel2, _ := s.daemon.Database.GetChannelByUserIdentifier(evt.Sender.String())
		// But we know the user
		if channel2 != nil {
			log.Info("User messaged us in a Channel we do not know")
			_, err := s.messenger.SendReplyToEvent("Hey, this is not our usual messaging channel ;)", evt, evt.RoomID.String())
			if err != nil {
				log.Warn(err.Error())
			}
		} else {
			log.Info("We do not know that user.")
		}
		return
	}

	// If it is a reply check if a reply action matches first
	if s.checkReplyActions(evt, channel, content) {
		return
	}

	// Check if a action matches
	if s.checkActions(evt, channel, content) {
		return
	}

	// Nothing left so it must be a reminder
	reminder, err := s.parseRemind(evt, channel)
	if err != nil {
		log.Warn(fmt.Sprintf("Failed parsing the Reminder with: %s", err.Error()))
		return
	}

	msg := fmt.Sprintf("Successfully added new reminder (ID: %d) for %s", reminder.ID, reminder.RemindTime.Format("15:04 02.01.2006"))

	log.Info(msg)
	_, err = s.messenger.SendReplyToEvent(msg, evt, evt.RoomID.String())
	if err != nil {
		log.Warn("Was not able to send success message to user")
	}
}

// checkActions checks if a message matches any special actions and performs them.
func (s *Syncer) checkActions(evt *event.Event, channel *database.Channel, content *event.MessageEventContent) (matched bool) {
	message := strings.ToLower(content.Body)

	// List action
	for _, action := range s.actions {
		log.Info("Checking for match with action " + action.Name)
		if matched, err := regexp.Match(action.Regex, []byte(message)); matched && err == nil {
			_ = action.Action(evt, channel)
			log.Info("Matched")
			return true
		}
	}

	return false
}

func (s *Syncer) handleReactions(source mautrix.EventSource, evt *event.Event) {
	log.Debug(fmt.Sprintf("New reaction: / Sender: %s / Room: / %s / Time: %d", evt.Sender, evt.RoomID, evt.Timestamp))

	// Do not answer our own and old messages
	if evt.Sender == id.UserID(s.botName) || evt.Timestamp/1000 <= time.Now().Unix()-60 {
		return
	}

	channel, err := s.daemon.Database.GetChannelByUserAndChannelIdentifier(evt.Sender.String(), evt.RoomID.String())
	if err != nil {
		log.Warn("Do not know that user and channel.")
	}

	content, ok := evt.Content.Parsed.(*event.ReactionEventContent)
	if !ok {
		log.Warn("Event is not a reaction event. Can not handle it.")
		return
	}

	if content.RelatesTo.EventID.String() == "" {
		log.Warn("Reaction with no realting message. Can not handle that.")
		return
	}

	message, err := s.daemon.Database.GetMessageByExternalID(content.RelatesTo.EventID.String())
	if err != nil {
		log.Info("Do not know the message related to the reaction.")
		return
	}

	for _, action := range s.reactionActions {
		log.Info("Checking for match with action " + action.Name)
		if action.Type != ReactionActionType(message.Type) && action.Type != ReactionActionTypeAll {
			continue
		}

		for _, key := range action.Keys {
			if content.RelatesTo.Key == key {
				err = action.Action(message, content, evt, channel)
				if err == nil {
					return
				}
			}
		}
	}

	log.Info("Nothing handled that reaction")

}

func (s *Syncer) checkReplyActions(evt *event.Event, channel *database.Channel, content *event.MessageEventContent) (matched bool) {
	if content.RelatesTo == nil {
		return false
	}
	if len(content.RelatesTo.EventID) < 2 {
		return false
	}

	message := strings.ToLower(content.Body)
	replyMessage, err := s.daemon.Database.GetMessageByExternalID(content.RelatesTo.EventID.String())
	if err != nil || replyMessage == nil {
		log.Info("Message replies to unknown message")
		return false
	}

	for _, action := range s.replyActions {
		if action.ReplyToType != "" && action.ReplyToType != replyMessage.Type {
			continue
		}

		if matched, err := regexp.Match(action.Regex, []byte(message)); matched && err == nil {
			_ = action.Action(evt, channel, replyMessage)
			log.Info("Matched")
			return true
		}
	}

	// Fallback change reminder date
	if replyMessage.ReminderID > 0 {
		baseTime := time.Now()

		// Clear body from characters the library can not handle
		t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
		strippedBody, _, _ := transform.String(t, matrixmessenger.StripReply(content.Body))
		log.Info(strippedBody)

		remindTime, err := naturaldate.Parse(strippedBody, baseTime, naturaldate.WithDirection(naturaldate.Future))
		if err != nil {
			log.Warn(err.Error())
			s.messenger.SendReplyToEvent("Sorry I was not able to get a time out of that message", evt, evt.RoomID.String())
			return true
		}

		reminder, err := s.daemon.Database.UpdateReminder(replyMessage.ReminderID, remindTime)
		if err != nil {
			log.Warn(err.Error())
			return true
		}

		s.messenger.SendReplyToEvent(fmt.Sprintf("I rescheduled your reminder \"%s\" to %s.", reminder.Message, reminder.RemindTime.Format("15:04 02.01.2006")), evt, evt.RoomID.String())
	}

	return true
}

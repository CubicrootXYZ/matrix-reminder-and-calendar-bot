package matrix

import (
	"errors"
	"time"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/format"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/messenger"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
)

// EventStateHandler handles state events from matrix.
func (service *service) EventStateHandler(source mautrix.EventSource, evt *event.Event) {
	logger := service.logger.WithFields(map[string]any{
		"sender":          evt.Sender,
		"room":            evt.RoomID,
		"event_timestamp": evt.Timestamp,
	})
	logger.Debugf("new state event")

	if service.crypto.enabled {
		service.crypto.olm.HandleMemberEvent(evt)
	}

	// Ignore old events or events from the bot itself
	if evt.Sender.String() == service.botname || evt.Timestamp/1000 < service.lastMessageFrom.Unix() {
		return
	}

	content, ok := evt.Content.Parsed.(*event.MemberEventContent)
	if !ok {
		logger.Infof("Event is not a member event. Can not handle it.")
		return
	}

	// Check if the event is known
	_, err := service.matrixDatabase.GetEventByID(evt.ID.String())
	if err == nil {
		return
	}

	switch content.Membership {
	case event.MembershipInvite, event.MembershipJoin:
		err := service.handleInvite(evt, content)
		if err != nil {
			logger.Errorf("Failed to handle membership invite with: " + err.Error())
		}
	case event.MembershipLeave, event.MembershipBan:
		err := service.handleLeave(evt, content)
		if err != nil {
			logger.Errorf("Failed to handle membership leave with: " + err.Error())
		}
	default:
		logger.Infof("No handling of this event as Membership %s is unknown.", content.Membership)
	}
}

func (service *service) handleInvite(evt *event.Event, content *event.MemberEventContent) error {
	declineInvites, err := service.maxUserReached()
	if err != nil {
		return err
	}

	if declineInvites {
		service.logger.Debugf(evt.Sender.String() + " ignored bot reached max users")
		return nil
	}

	user, err := service.matrixDatabase.GetUserByID(evt.Sender.String())
	if err != nil {
		if !errors.Is(err, database.ErrNotFound) {
			return err
		}
		user = nil
	}

	if user.Blocked {
		service.logger.Debugf("user '%s' is blocked - ignoring", evt.Sender.String())
		return nil
	}

	roomCreated := false
	room, err := service.matrixDatabase.GetRoomByID(evt.RoomID.String())
	if err != nil {
		if !errors.Is(err, database.ErrNotFound) {
			return err
		}
		room = nil
	}

	_, err = service.client.JoinRoom(evt.RoomID.String(), "", nil)
	if err != nil {
		service.logger.Errorf("Failed joining channel %s with: %s", evt.RoomID.String(), err.Error())
		return err
	}

	if user == nil {
		user, err = service.matrixDatabase.NewUser(&database.MatrixUser{
			ID: evt.Sender.String(),
		})
		if err != nil {
			return err
		}
	}

	if room == nil {
		roomCreated = true
		room, err = service.matrixDatabase.NewRoom(&database.MatrixRoom{
			RoomID: evt.RoomID.String(),
		})
		if err != nil {
			return err
		}
	}

	room.Users = append(room.Users, *user)
	room, err = service.matrixDatabase.UpdateRoom(room)
	if err != nil {
		return err
	}

	_, err = service.matrixDatabase.NewEvent(&database.MatrixEvent{
		ID:     evt.ID.String(),
		UserID: user.ID,
		RoomID: room.ID,
		Type:   string(content.Membership),
		SendAt: time.Unix(evt.Timestamp/1000, 0),
	})
	if err != nil {
		return err
	}

	if roomCreated {
		go service.sendWelcomeMessage(room, user)
	}

	return nil
}

func (service *service) sendWelcomeMessage(room *database.MatrixRoom, user *database.MatrixUser) {
	message, messageFormatted := getWelcomeMessage()

	resp, err := service.messenger.SendMessage(messenger.HTMLMessage(
		message,
		messageFormatted,
		room.RoomID,
	))
	if err != nil {
		service.logger.Infof("Failed to send message: " + err.Error())
		return
	}

	_, err = service.matrixDatabase.NewMessage(&database.MatrixMessage{
		ID:            resp.ExternalIdentifier,
		UserID:        user.ID,
		RoomID:        room.ID,
		Body:          message,
		BodyFormatted: messageFormatted,
		SendAt:        resp.Timestamp,
		Type:          database.MessageTypeWelcome,
		Incoming:      false,
	})
	if err != nil {
		service.logger.Infof("Failed saving message into database: " + err.Error())
	}
}

func (service *service) handleLeave(evt *event.Event, content *event.MemberEventContent) error {
	if evt.StateKey == nil {
		return nil
	}

	room, err := service.matrixDatabase.GetRoomByID(string(evt.RoomID))
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			return nil
		}
		return err
	}

	err = service.matrixDatabase.DeleteAllEventsFromRoom(room.ID)
	if err != nil {
		return err
	}

	err = service.matrixDatabase.DeleteAllMessagesFromRoom(room.ID)
	if err != nil {
		return err
	}

	err = service.matrixDatabase.DeleteRoom(room.ID)

	// TODO delete channel? At least if no other in/output is set

	return err
}

func getWelcomeMessage() (string, string) {
	msg := format.Formater{}
	msg.Title("Welcome to RemindMe")
	msg.TextLine("Hey, I am your personal reminder bot. Beep boop beep.")
	msg.Text("You want to now what I am capable of? Just text me ")
	msg.BoldLine("list all commands")
	msg.TextLine("First things you should do are setting your timezone and a daily reminder.")

	msg.SubTitle("Attribution")
	msg.TextLine("This bot is open for everyone and build with the help of voluntary software developers.")
	msg.Text("The source code can be found at ")
	msg.Link("GitHub", "https://github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot")
	msg.TextLine(". Star it if you like the bot, open issues or discussions with your findings.")

	return msg.Build()
}

func (service *service) maxUserReached() (bool, error) {
	if !service.config.AllowInvites {
		return true, nil
	}

	if service.config.RoomLimit > 0 {
		roomCount, err := service.matrixDatabase.GetRoomCount()
		if err != nil {
			return true, err
		}

		if service.config.RoomLimit > uint(roomCount) {
			return false, nil
		}

		return true, nil
	}

	return false, nil
}

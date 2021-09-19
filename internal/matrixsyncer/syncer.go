package matrixsyncer

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/configuration"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/eventdaemon"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/log"
)

// Syncer receives messages from a matrix channel
type Syncer struct {
	config          configuration.Matrix
	users           []string
	client          *mautrix.Client
	daemon          *eventdaemon.Daemon
	botName         string
	messenger       Messenger
	actions         []*Action         // Actions based on direct messages from the user
	reactionActions []*ReactionAction // Actions based on reactions by the user
	replyActions    []*ReplyAction    // Actions based on replies from the user on existing messages
}

// Create creates a new syncer
func Create(config configuration.Matrix, matrixUsers []string, messenger Messenger) *Syncer {
	syncer := &Syncer{
		config:    config,
		users:     matrixUsers,
		messenger: messenger,
	}

	// Add all actions
	syncer.actions = append(syncer.actions, syncer.getActionList())
	syncer.actions = append(syncer.actions, syncer.getActionCommands())
	syncer.actions = append(syncer.actions, syncer.getActionTimezone())
	syncer.actions = append(syncer.actions, syncer.getActionSetDailyReminder())
	syncer.actions = append(syncer.actions, syncer.getActionDeleteDailyReminder())

	syncer.reactionActions = append(syncer.reactionActions, syncer.getReactionActionDelete(ReactionActionTypeReminderRequest))
	syncer.reactionActions = append(syncer.reactionActions, syncer.getReactionsAddTime(ReactionActionTypeReminderRequest)...)
	syncer.reactionActions = append(syncer.reactionActions, syncer.getReactionActionDeleteDailyReminder(ReactionActionTypeDailyReminder))

	syncer.replyActions = append(syncer.replyActions, syncer.getReplyActionDelete(database.MessageTypesWithReminder))
	syncer.replyActions = append(syncer.replyActions, syncer.getReplyActionRecurring(database.MessageTypesWithReminder))

	return syncer
}

// Start starts the syncer
func (s *Syncer) Start(daemon *eventdaemon.Daemon) error {
	log.Info(fmt.Sprintf("Starting Matrix syncer for user %s on server %s", s.config.Username, s.config.Homeserver))

	s.daemon = daemon
	s.botName = fmt.Sprintf("@%s:%s", s.config.Username, strings.ReplaceAll(strings.ReplaceAll(s.config.Homeserver, "https://", ""), "http://", ""))

	// Log into matrix
	client, err := mautrix.NewClient(s.config.Homeserver, "", "")
	if err != nil {
		return err
	}

	s.client = client
	_, err = s.client.Login(&mautrix.ReqLogin{
		Type:             "m.login.password",
		Identifier:       mautrix.UserIdentifier{Type: mautrix.IdentifierTypeUser, User: s.config.Username},
		Password:         s.config.Password,
		StoreCredentials: true,
	})
	if err != nil {
		return err
	}
	log.Info("Logged in to matrix")

	// Make channel for each user we do not know and remove users we do not want
	channels := make([]*database.Channel, 0)

	for _, user := range s.users {
		channel, err := s.daemon.Database.GetChannelByUserIdentifier(user)
		if err == gorm.ErrRecordNotFound {
			channel, err = s.createChannel(user)
			if err != nil {
				return err
			}
		} else if err != nil {
			return err
		}

		channels = append(channels, channel)

		s.messenger.SendNotice("Sorry I was sleeping for a while. I am now ready for your requests!", channel.ChannelIdentifier)
	}

	err = s.daemon.Database.CleanChannels(channels)
	if err != nil {
		log.Warn("Can not clean channels list")
		panic(err)
	}

	// Get messages
	syncer := s.client.Syncer.(*mautrix.DefaultSyncer)
	syncer.OnEventType(event.EventMessage, s.handleMessages)
	syncer.OnEventType(event.EventReaction, s.handleReactions)
	return client.Sync()
}

// Stop stops the syncer
func (s *Syncer) Stop() {
	s.client.StopSync()
}

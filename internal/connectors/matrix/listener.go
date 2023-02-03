package matrix

import (
	"errors"

	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
)

func (service *service) startListener() error {
	syncer, ok := service.client.Syncer.(*mautrix.DefaultSyncer)
	if !ok {
		return errors.New("syncer of wrong type")
	}

	if service.crypto.enabled {
		syncer.OnSync(func(resp *mautrix.RespSync, since string) bool {
			service.crypto.olm.ProcessSyncResponse(resp, since)
			return true
		})
		syncer.OnEventType(event.EventEncrypted, service.MessageEventHandler)
		syncer.OnEventType(event.StateEncryption, func(_ mautrix.EventSource, event *event.Event) {
			service.crypto.stateStore.SetEncryptionEvent(event)
		})
	}

	syncer.OnEventType(event.EventMessage, service.MessageEventHandler)
	/* TODO syncer.OnEventType(event.EventReaction, reactionHandler.NewEvent)
	syncer.OnEventType(event.StateMember, stateMemberHandler.NewEvent)*/

	return service.client.Sync()
}

// setLastMessage so the handlers will know which messages can be ignored savely
func (service *service) setLastMessage() { // TODO use this!
	message, err := service.matrixDatabase.GetLastMessage()
	if err == nil {
		service.lastMessageFrom = message.SendAt
	}
}

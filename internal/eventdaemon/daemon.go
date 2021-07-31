package eventdaemon

import (
	"sync"
	"time"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/log"
)

// Daemon is the event daemon collecting events from a messenger
type Daemon struct {
	Database Database
	syncer   Syncer
}

// Command defines commands the daemon can handle
type Command int16

const (
	// CommandStop stops the daemon
	CommandStop = Command(1)
)

// Create returns a new event daemon
func Create(database Database, syncer Syncer) *Daemon {
	return &Daemon{
		Database: database,
		syncer:   syncer,
	}
}

// Start starts the daemon
func (d *Daemon) Start(wg *sync.WaitGroup) {
	for {
		log.Info("Starting syncer")
		err := d.syncer.Start(d)
		log.Warn("Syncer stopped.")
		if err != nil {
			log.Error("Syncer returned error: " + err.Error())
		}
		time.Sleep(time.Minute * 5)
	}
	//wg.Done()
}

// Stop stops the daemon
func (d *Daemon) Stop() {
	log.Info("Stopping syncer")
	d.syncer.Stop()
}

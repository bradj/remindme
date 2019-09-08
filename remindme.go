package remindme

import (
	"sync"
	"time"
)

type reminder struct {
	ID      int
	Author  string
	Body    string
	Channel string
	EndTime time.Time
}

type db struct {
	expiredReminders chan reminder

	mut       sync.Mutex
	Reminders map[int]reminder
	ID        int
}

// New creates a db
func New() db {
	return db{
		expiredReminders: make(chan reminder),
		Reminders:        make(map[int]reminder),
	}
}

func (d *db) Add(author string, body string, end time.Time) {
	d.mut.Lock()
	defer d.mut.Unlock()

	d.ID++

	d.Reminders[d.ID] = reminder{
		ID:      d.ID,
		Author:  author,
		Body:    body,
		EndTime: end,
	}
}

func (d *db) find(t time.Time) {
	d.mut.Lock()
	defer d.mut.Unlock()

	for _, rem := range d.Reminders {
		if rem.EndTime.After(t) {
			continue
		}

		// send reminder
		d.expiredReminders <- rem

		// remove reminder
		delete(d.Reminders, rem.ID)
	}
}

func (d *db) StartTicks() {
	ticker := time.NewTicker(time.Minute)

	for {
		t := <-ticker.C
		d.find(t)
	}
}

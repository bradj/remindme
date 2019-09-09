package remindme

import (
	"sync"
	"time"
)

// Reminder object
type Reminder struct {
	ID      int
	Author  string
	Body    string
	Channel string
	EndTime time.Time
}

// DB for reminders
type DB struct {
	ExpiredReminders chan Reminder

	mut       sync.RWMutex
	Reminders map[int]Reminder
	ID        int
}

// New creates a db
func New() *DB {
	return &DB{
		ExpiredReminders: make(chan Reminder),
		Reminders:        make(map[int]Reminder),
	}
}

// Add a new reminder
func (d *DB) Add(author string, body string, end time.Time) {
	d.mut.Lock()
	defer d.mut.Unlock()

	d.Reminders[d.ID] = Reminder{
		ID:      d.ID,
		Author:  author,
		Body:    body,
		EndTime: end,
	}

	d.ID++
}

// Remove a reminder
func (d *DB) Remove(r Reminder) {
	d.mut.Lock()
	defer d.mut.Unlock()

	delete(d.Reminders, r.ID)
}

// Count all reminders
func (d *DB) Count() int {
	d.mut.RLock()
	defer d.mut.RUnlock()

	return len(d.Reminders)
}

func (d *DB) findByTime(t time.Time) {
	d.mut.Lock()
	defer d.mut.Unlock()

	for _, rem := range d.Reminders {
		if rem.EndTime.After(t) {
			continue
		}

		// send reminder
		d.ExpiredReminders <- rem

		// remove reminder
		delete(d.Reminders, rem.ID)
	}
}

func (d *DB) findByKey(key int) *Reminder {
	d.mut.RLock()
	defer d.mut.RUnlock()

	rem, ok := d.Reminders[key]

	if !ok {
		return nil
	}

	return &rem
}

func (d *DB) findByAuthor(author string) []Reminder {
	d.mut.RLock()
	defer d.mut.RUnlock()

	reminders := make([]Reminder, 0)

	for _, rem := range d.Reminders {
		if rem.Author != author {
			continue
		}

		reminders = append(reminders, rem)
	}

	return reminders
}

// WaitForReminders starts the reminder timer and emits expired reminders
// on DB.expiredReminders
func (d *DB) WaitForReminders() {
	ticker := time.NewTicker(time.Minute)

	for {
		t := <-ticker.C
		d.findByTime(t)
	}
}

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
	Network string
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
func (d *DB) Add(reminder Reminder) {
	d.mut.Lock()
	defer d.mut.Unlock()

	reminder.ID = d.ID
	d.Reminders[d.ID] = reminder

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

func (d *DB) findByTime(t time.Time) []Reminder {
	d.mut.Lock()
	defer d.mut.Unlock()

	reminders := make([]Reminder, 0)

	for _, rem := range d.Reminders {
		if rem.EndTime.After(t) {
			continue
		}

		reminders = append(reminders, rem)
	}

	return reminders
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

func (d *DB) expireReminders(reminders []Reminder) {
	for _, rem := range reminders {
		// send reminder
		d.ExpiredReminders <- rem

		// remove reminder
		d.Remove(rem)
	}
}

// WaitForReminders starts the reminder timer and emits expired reminders
// on DB.expiredReminders
func (d *DB) WaitForReminders() {
	ticker := time.NewTicker(time.Second * 10)

	for t := range ticker.C {
		d.expireReminders(d.findByTime(t))
	}
}

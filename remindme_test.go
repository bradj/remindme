package remindme

import (
	"testing"
	"time"
)

var body = "do remindme with brad"
var author = "aaron"
var reminder = Reminder{
	Author:  author,
	Body:    body,
	EndTime: time.Now().Add(-time.Minute),
}

func TestFindByTime(t *testing.T) {
	db := New()

	db.Add(reminder)

	if db.Count() != 1 {
		t.Error("Reminders map is not of length 1")
	}

	reminders := db.findByTime(time.Now())

	if len(reminders) != 1 {
		t.Errorf("Reminders should be of length 1 instead it is %d", len(reminders))
	}

	rem := reminders[0]

	if a := rem.Author; a != author {
		t.Errorf("Reminder author should be %v but was %v instead", author, a)
	}

	if b := rem.Body; b != body {
		t.Error("Body mismatch in findByTime")
	}
}

func TestFindByAuthor(t *testing.T) {
	db := New()

	db.Add(reminder)

	if db.Count() != 1 {
		t.Error("Reminders map is not of length 1")
	}

	reminders := db.findByAuthor(author)

	if count := len(reminders); count != 1 {
		t.Errorf("Reminders should be of length 1 but is %d instead", count)
	}

	if a := reminders[0].Author; a != author {
		t.Errorf("Reminder author should be %v but instead is %v", author, a)
	}
}

func TestFindByKey(t *testing.T) {
	db := New()
	const key = 0

	db.Add(reminder)

	if db.Count() != 1 {
		t.Error("Reminders map is not of length 1")
	}

	rem := db.findByKey(key)

	if rem == nil {
		t.Errorf("Reminder was not found for key %d", key)
	}

	if a := rem.Author; a != author {
		t.Errorf("Reminder Author was supposed to be %v but was %v instead", author, a)
	}
}

func TestExpiredReminders(t *testing.T) {
	db := New()

	db.Add(reminder)

	if db.Count() != 1 {
		t.Error("Reminders map is not of length 1")
	}

	go db.expireReminders(db.findByAuthor(reminder.Author))

	rem := <-db.ExpiredReminders

	if rem.Author != reminder.Author {
		t.Errorf("Reminder Author was supposed to be %s but was %s instead", reminder.Author, rem.Author)
	}

	if rem.Body != reminder.Body {
		t.Errorf("Reminder Body was supposed to be %s but was %s instead", reminder.Author, rem.Author)
	}
}

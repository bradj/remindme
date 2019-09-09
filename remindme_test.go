package remindme

import (
	"testing"
	"time"
)

func TestFindByTime(t *testing.T) {
	db := New()
	author := "aaron"

	db.Add(author, "do remindme with brad",
		time.Now().Add(-time.Minute))

	if db.Count() != 1 {
		t.Error("Reminders map is not of length 1")
	}

	go db.findByTime(time.Now())

	rem := <-db.ExpiredReminders

	if a := rem.Author; a != author {
		t.Errorf("Reminder author should be %v but was %v instead", author, a)
	}

	if db.Count() != 0 {
		t.Error("Reminders map is not of length 0")
	}
}

func TestFindByAuthor(t *testing.T) {
	db := New()
	author := "aaron"

	db.Add(author, "do remindme with brad",
		time.Now().Add(-time.Minute))

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

	const author = "aaron"
	const key = 0

	db.Add(author, "do remindme with brad",
		time.Now().Add(-time.Minute))

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

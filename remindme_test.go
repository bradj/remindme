package remindme

import (
	"testing"
	"time"
)

func testReminder() Reminder {
	body := "do remindme with brad"
	author := "aaron"
	reminder := Reminder{
		Author:  author,
		Body:    body,
		EndTime: time.Now().Add(time.Second),
	}

	return reminder
}

func TestFindByTime(t *testing.T) {
	t.Parallel()
	db := New()

	testRem := testReminder()
	db.Add(testRem)

	if db.Count() != 1 {
		t.Error("Reminders map is not of length 1")
	}

	reminders := db.findByTime(time.Now().Add(time.Second * 2))

	if len(reminders) != 1 {
		t.Errorf("Reminders should be of length 1 instead it is %d", len(reminders))
	}

	rem := reminders[0]

	if a := rem.Author; a != testRem.Author {
		t.Errorf("Reminder author should be %v but was %v instead", testRem.Author, a)
	}

	if b := rem.Body; b != testRem.Body {
		t.Error("Body mismatch in findByTime")
	}
}

func TestFindByAuthor(t *testing.T) {
	t.Parallel()
	db := New()

	testRem := testReminder()
	db.Add(testRem)

	if db.Count() != 1 {
		t.Error("Reminders map is not of length 1")
	}

	reminders := db.findByAuthor(testRem.Author)

	if count := len(reminders); count != 1 {
		t.Errorf("Reminders should be of length 1 but is %d instead", count)
	}

	if a := reminders[0].Author; a != testRem.Author {
		t.Errorf("Reminder author should be %v but instead is %v", testRem.Author, a)
	}
}

func TestFindByKey(t *testing.T) {
	t.Parallel()
	db := New()
	const key = 0

	testRem := testReminder()
	db.Add(testRem)

	if db.Count() != 1 {
		t.Error("Reminders map is not of length 1")
	}

	rem := db.findByKey(key)

	if rem == nil {
		t.Errorf("Reminder was not found for key %d", key)
	}

	if a := rem.Author; a != testRem.Author {
		t.Errorf("Reminder Author was supposed to be %v but was %v instead", testRem.Author, a)
	}
}

func TestExpiredReminders(t *testing.T) {
	t.Parallel()
	db := New()

	testRem := testReminder()
	db.Add(testRem)

	if db.Count() != 1 {
		t.Error("Reminders map is not of length 1")
	}

	go db.expireReminders(db.findByAuthor(testRem.Author))

	rem := <-db.ExpiredReminders

	if rem.Author != testRem.Author {
		t.Errorf("Reminder Author was supposed to be %s but was %s instead", testRem.Author, rem.Author)
	}

	if rem.Body != testRem.Body {
		t.Errorf("Reminder Body was supposed to be %s but was %s instead", testRem.Author, rem.Author)
	}
}

func TestAddWithPastTime(t *testing.T) {
	t.Parallel()
	db := New()

	errorMessage := "End time must be after now"

	rem := Reminder{
		Author:  "aaron",
		Body:    "fix the bot, brad",
		EndTime: time.Now().Add(-time.Minute),
	}

	err := db.Add(rem)

	if err != errPastEndTime {
		t.Errorf("Error message should be \"%s\" but was \"%s\"", errorMessage, err.Error())
	}
}

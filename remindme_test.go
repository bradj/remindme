package remindme

import (
	"testing"
	"time"
)

func TestFind(t *testing.T) {
	db := New()

	db.Add("aaron", "do remindme with brad",
		time.Now().Add(-time.Minute))

	if len(db.Reminders) != 1 {
		t.Error("Reminders map is not of length 1")
	}

	go db.find(time.Now())

	<-db.expiredReminders

	db.mut.Lock()
	defer db.mut.Unlock()

	if len(db.Reminders) != 0 {
		t.Error("Reminders map is not of length 0")
	}
}

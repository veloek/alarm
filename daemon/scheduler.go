package daemon

import (
	"fmt"
	"os"
	"sync"
	"time"
)

func newScheduler(u <-chan []*Alarm) *scheduler {
	s := &scheduler{updates: u}
	go s.listenForUpdates()
	return s
}

type scheduler struct {
	alarms  []*Alarm
	updates <-chan []*Alarm
	mutex   sync.Mutex
}

func (s *scheduler) listenForUpdates() {
	for u := range s.updates {
		s.mutex.Lock()
		s.alarms = u
		s.mutex.Unlock()
	}
}

func (s *scheduler) start() {
	for {
		s.tick()
		time.Sleep(900 * time.Millisecond)
	}
}

func (s *scheduler) tick() {
	now := time.Now()
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for _, a := range s.alarms {
		h := int(a.Time.Hours)
		m := int(a.Time.Minutes)
		s := int(a.Time.Seconds)

		if a.Time.Format == Time_AM && h == 12 {
			h = 0
		}

		if a.Time.Format == Time_PM && h < 12 {
			h += 12
		}

		if s == now.Second() && m == now.Minute() && (h == now.Hour() || a.Recurrence == Alarm_HOURLY) {
			triggerAlarm(a)
		}
	}
}

func triggerAlarm(a *Alarm) {
	if a.Recurrence == Alarm_NO_RECURRENCE {
		err := removeAlarm(int(a.Id))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error removing one-time alarm: %v\n", err)
		}
	}

	if err := playSound(); err != nil {
		fmt.Fprintf(os.Stderr, "Error playing sound: %v\n", err)
	}
}

func removeAlarm(alarmID int) error {
	repo, err := dbConnect()
	if err != nil {
		return err
	}
	return repo.RemoveAlarm(alarmID)
}

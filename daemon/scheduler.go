package daemon

import (
	"sync"
	"time"
)

func newScheduler(u <-chan []*Alarm, cb func(a *Alarm)) *scheduler {
	s := &scheduler{updates: u, handler: cb}
	go s.listenForUpdates()
	return s
}

type scheduler struct {
	alarms  []*Alarm
	updates <-chan []*Alarm
	mutex   sync.Mutex
	handler func(a *Alarm)
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
		hr := int(a.Time.Hours)
		min := int(a.Time.Minutes)
		sec := int(a.Time.Seconds)

		if a.Time.Format == Time_AM && hr == 12 {
			hr = 0
		}

		if a.Time.Format == Time_PM && hr < 12 {
			hr += 12
		}

		if sec == now.Second() && min == now.Minute() && (hr == now.Hour() || a.Recurrence == Alarm_HOURLY) {
			s.handler(a)
		}
	}
}

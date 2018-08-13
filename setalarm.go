package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/veloek/alarm/daemon"
)

func setAlarm(t string, r string) {
	a, err := parseAlarm(t, r)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing alarm %s%s: %v\n", t, r, err)
		printUsageAndExit()
	}

	c := newClient()
	_, err = c.SetAlarm(context.Background(), &daemon.SetAlarmRequest{
		Alarm: a})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error setting alarm: %v\n", err)
		os.Exit(2)
	}
}

var alarmRegex = regexp.MustCompile("^(\\d+):(\\d+)(?::(\\d+))?(AM|PM)?$")

func parseAlarm(t string, r string) (al *daemon.Alarm, err error) {
	alarmRecurrence, err := parseRecurrence(r)
	if err != nil {
		return
	}

	m := alarmRegex.FindStringSubmatch(t)
	if m == nil {
		err = fmt.Errorf("invalid time: %s", t)
		return
	}

	hours, _ := strconv.Atoi(m[1])
	minutes, _ := strconv.Atoi(m[2])
	var seconds int
	if m[3] != "" {
		seconds, _ = strconv.Atoi(m[3])
	}
	var format daemon.Time_Format
	switch m[4] {
	case "AM":
		format = daemon.Time_AM
	case "PM":
		format = daemon.Time_PM
	default:
		format = daemon.Time_TWENTYFOUR
	}

	if format == daemon.Time_AM || format == daemon.Time_PM {
		if hours > 12 {
			err = errors.New("invalid time: hours more than 12")
			return
		}
	} else if hours > 23 {
		err = errors.New("invalid time: hours more than 23")
		return
	}

	if minutes > 59 {
		err = errors.New("invalid time: minutes more than 59")
		return
	}

	if seconds > 59 {
		err = errors.New("invalid time: seconds more then 59")
		return
	}

	alarmTime := &daemon.Time{
		Hours:   int32(hours),
		Minutes: int32(minutes),
		Seconds: int32(seconds),
		Format:  format,
	}

	al = &daemon.Alarm{
		Time:       alarmTime,
		Recurrence: alarmRecurrence,
	}
	return
}

func parseRecurrence(r string) (rec daemon.Alarm_Recurrence, err error) {
	switch strings.ToLower(r) {
	case "hourly":
		rec = daemon.Alarm_HOURLY
	case "daily":
		rec = daemon.Alarm_DAILY
	case "":
		rec = daemon.Alarm_NO_RECURRENCE
	default:
		err = fmt.Errorf("invalid recurrence: %s", r)
	}
	return
}

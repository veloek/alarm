package daemon

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type alarmRepo struct {
	db *sql.DB
}

func (r *alarmRepo) SetAlarm(a *Alarm) error {
	stmt, err := r.db.Prepare(`INSERT INTO alarm (time, time_format, recurrence) VALUES (?, ?, ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		fmt.Sprintf("%02d:%02d:%02d", a.Time.Hours, a.Time.Minutes, a.Time.Seconds),
		Time_Format_name[int32(a.Time.Format)],
		Alarm_Recurrence_name[int32(a.Recurrence)])
	return err
}

func (r *alarmRepo) GetAlarms() ([]*Alarm, error) {
	rows, err := r.db.Query(`SELECT alarm_id, time, time_format, recurrence FROM alarm`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	alarms := make([]*Alarm, 0)

	var id int
	var time string
	var format string
	var recurrence string
	for rows.Next() {
		err = rows.Scan(&id, &time, &format, &recurrence)
		if err != nil {
			return nil, err
		}

		t := parseTime(time, format)
		alarms = append(alarms, &Alarm{
			Id:         int32(id),
			Time:       t,
			Recurrence: Alarm_Recurrence(Alarm_Recurrence_value[recurrence]),
		})
	}
	return alarms, nil
}

func parseTime(t, f string) *Time {
	split := strings.Split(t, ":")
	h, _ := strconv.Atoi(split[0])
	m, _ := strconv.Atoi(split[1])
	s, _ := strconv.Atoi(split[2])
	return &Time{
		Hours:   int32(h),
		Minutes: int32(m),
		Seconds: int32(s),
		Format:  Time_Format(Time_Format_value[f]),
	}
}

func (r *alarmRepo) RemoveAlarm(id int) error {
	_, err := r.db.Exec(`DELETE FROM alarm WHERE alarm_id = ?`, id)
	return err

}

func dbConnect() (*alarmRepo, error) {
	dbFile := filepath.Join(userHomeDir(), ".alarms")
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return nil, err
	}

	err = createDatabase(db)
	if err != nil {
		db.Close()
		return nil, err
	}

	return &alarmRepo{db}, nil
}

func createDatabase(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS alarm (
		alarm_id INTEGER PRIMARY KEY AUTOINCREMENT,
		time VARCHAR(8),
		time_format VARCHAR(10),
		recurrence VARCHAR2(20)
	)`)
	return err
}

// userHomeDir grabbed from:
// https://stackoverflow.com/questions/7922270/obtain-users-home-directory
func userHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	} else if runtime.GOOS == "linux" {
		home := os.Getenv("XDG_CONFIG_HOME")
		if home != "" {
			return home
		}
	}
	return os.Getenv("HOME")
}

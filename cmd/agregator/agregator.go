package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"time"
)

const (
	Last        Minutes = 0
	OneMinute           = 1
	FiveMinutes         = 5
	OneHour             = 60
	OneDay              = 1440
	OneWeek             = 10080
	OneMonth            = 302400 //30 days
)

type Minutes int

type Average struct {
	Usd float32
	Eur float32
}

var (
	Debug   = true
	UsdCode = "USDRUB"
	EurCode = "EURRUB"
)

func main() {
	doEvery(time.Minute, doAgregate)
}

func toLog(v ...interface{}) {
	if Debug {
		log.Println(v...)
	}
}

func doEvery(d time.Duration, f func(time.Time)) {
	for x := range time.Tick(d) {
		f(x)
	}
}
func doAgregate(t time.Time) {

	now := t.UTC()
	previousMinute := now.Add(time.Duration(-1) * time.Minute)

	db, err := sql.Open("sqlite3", "go_test.sqlite")

	if err != nil {

		toLog("Connect DB Error => " + err.Error())

	} else {

		saveAverageToHistory(db, previousMinute, OneMinute)

		if previousMinute.Minute()%5 == 0 {

			saveAverageToHistory(db, previousMinute, FiveMinutes)
		}

		if previousMinute.Minute() == 0 {

			saveAverageToHistory(db, previousMinute, OneHour)

			if now.Hour() == 0 {

				saveAverageToHistory(db, previousMinute, OneDay)

				if int(now.Weekday()) == 1 {

					saveAverageToHistory(db, previousMinute, OneWeek)
				}

				if now.Day() == 1 {

					saveAverageToHistory(db, previousMinute, OneMonth)

				}

			}

		}

		saveAverage(db, previousMinute)

		err = db.Close()

		if err != nil {

			toLog("Close DB Error=> " + err.Error())

		}
	}
}

func getAverage(db *sql.DB, dateTime time.Time, interval Minutes) (average Average) {

	var (
		sqlTable string
	)

	switch interval {
	case OneDay:
		sqlTable = "average_per_hour"
	case OneWeek:
		sqlTable = "average_per_day"
	case OneMonth:
		sqlTable = "average_per_day"
	default:
		sqlTable = "raw_data"
	}

	query := fmt.Sprintf(`SELECT AVG(usd) AS usd, AVG(eur) AS eur FROM "%s" WHERE `, sqlTable)

	if interval == Last {
		query = query + "1 ORDER BY utc LIMIT 1"
	} else {
		dateString := dateTime.Format("2006-01-02 15:04:00")

		query = query + fmt.Sprintf(
			`utc BETWEEN datetime('%s', '-%d minutes') AND  datetime('%s')`,
			dateString, interval, dateString)
	}

	rows, err := db.Query(query)
	if err != nil {
		toLog("Select DB Error => " + err.Error())
	} else {

		for rows.Next() {

			if err := rows.Scan(&average.Usd, &average.Eur); err != nil {
				toLog("DB Select Result Scan Error => " + err.Error())
			}

		}
		if err := rows.Err(); err != nil {
			toLog("DB Select Rows Error => " + err.Error())
		}

		err = rows.Close()
		if err != nil {
			toLog("DB Close Rows Error => " + err.Error())
		}

	}
	return
}

func saveAverageToHistory(db *sql.DB, time time.Time, interval Minutes) {

	var (
		sqlTable string
	)

	switch interval {
	case OneMinute:
		sqlTable = "average_per_minute"
	case FiveMinutes:
		sqlTable = "average_per_five_minute"
	case OneHour:
		sqlTable = "average_per_hour"
	case OneDay:
		sqlTable = "average_per_day"
	case OneWeek:
		sqlTable = "average_per_week"
	case OneMonth:
		sqlTable = "average_per_month"
	default:
		toLog("[func saveAverageToHistory] => unknown interval.")
		return
	}

	averagePerInterval := getAverage(db, time, interval)

	stmt, err := db.Prepare("INSERT INTO " + sqlTable + "(usd, eur, utc) values(?,?,?)")
	if err != nil {
		toLog("Insert Prepare DB Error => " + err.Error())
	} else {

		_, err := stmt.Exec(
			averagePerInterval.Usd,
			averagePerInterval.Eur,
			time.Format("2006-01-02 15:04:00"),
		)
		if err != nil {
			toLog("Insert Exec DB Error => " + err.Error())
		}

		err = stmt.Close()
		if err != nil {
			toLog("Stmt Close DB Error => " + err.Error())
		}
	}
}

func saveAverage(db *sql.DB, time time.Time) {

	lastValue := getAverage(db, time, Last)
	averagePerDay := getAverage(db, time, OneDay)
	averagePerWeek := getAverage(db, time, OneWeek)
	averagePerMonth := getAverage(db, time, OneMonth)

	stmt, err := db.Prepare("INSERT INTO average_for_interval(last, day, week, month, utc, code) values(?,?,?,?,?,?)")
	if err != nil {
		toLog("Insert Prepare DB Error => " + err.Error())
	} else {

		dateString := time.Format("2006-01-02 15:04:00")
		var err error

		_, err = stmt.Exec(lastValue.Usd, averagePerDay.Usd, averagePerWeek.Usd, averagePerMonth.Usd, dateString, UsdCode)
		if err != nil {
			toLog("Insert Exec DB Error => " + err.Error())
		}

		_, err = stmt.Exec(lastValue.Eur, averagePerDay.Eur, averagePerWeek.Eur, averagePerMonth.Eur, dateString, EurCode)
		if err != nil {
			toLog("Insert Exec DB Error => " + err.Error())
		}

		err = stmt.Close()
		if err != nil {
			toLog("Stmt Close DB Error => " + err.Error())
		}
	}

}

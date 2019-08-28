package handlers

import (
	"TestTask/models"
	"TestTask/renderings"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo"

	"net/http"

	"time"
)

func GetHistory(ctx echo.Context) error {

	var (
		code     = checkCode(ctx.Param("code"))
		tableRow = getCurrencyRowName(code)
		sqlTable = GetTableNameByInterval(ctx.QueryParam("interval"))
		resp     = NewResponse()
		db       = ctx.Get(models.DBContextKey).(*sql.DB)
		item     = renderings.HistoryItem{}
		dt       *time.Time
	)

	dateFrom, errFrom := checkFrom(ctx.QueryParam("from"))
	dateTo, errTo := checkTo(ctx.QueryParam("to"))

	if errFrom != nil && errTo != nil {
		resp.Message = "Invalid datetime in query parameter 'from' or 'to'"
	} else {

		query := `SELECT printf("%.4f", ` + tableRow + ` ) AS rate, utc AS time FROM "` + sqlTable + `" 
				WHERE utc BETWEEN datetime(?) AND  datetime(?)`

		rows, err := db.Query(query, dateFrom, dateTo)
		if err != nil {
			resp.Code = 500
			resp.Message = "Internal Server Error"
		} else {
			for rows.Next() {
				if err = rows.Scan(&item.Rate, &dt); err != nil {
					resp.Code = 500
					resp.Message = "Internal Server Error"
				} else {
					item.Time = dt.Format("2006-01-02 15:04")
					resp.Payload = append(resp.Payload, item)
				}
			}

			if err := rows.Err(); err != nil {
				resp.Code = 500
				resp.Message = "Internal Server Error"
			}

			if err := rows.Close(); err != nil {
				resp.Code = 500
				resp.Message = "Internal Server Error"
			}

			if len(resp.Payload) > 0 {
				resp.Message = "Success"
				resp.Code = http.StatusOK
			} else {
				resp.Code = http.StatusNotFound
				resp.Message = "Page Not Found"
			}

		}
	}

	return ctx.JSON(resp.Code, resp)
}

func GetTableNameByInterval(interval string) string {
	switch interval {
	case "1m":
		return "average_per_minute"
	case "5m":
		return "average_per_five_minute"
	case "1h":
		return "average_per_hour"
	case "1d":
		return "average_per_day"
	default:
		return "average_per_day"
	}
}

func NewResponse() renderings.HistoryResponse {
	resp := renderings.HistoryResponse{}
	resp.Payload = []renderings.HistoryItem{}
	resp.Message = "Bad Request"
	resp.Code = http.StatusBadRequest
	return resp
}

func checkCode(code string) string {
	if len(code) == 0 {
		return "usdrub"
	} else {
		return code
	}
}

func checkFrom(From string) (DateFrom string, err error) {

	if len(From) == 0 {
		From = time.Now().UTC().Add(time.Duration(-24) * time.Hour).Format("2006-01-02 15:04")
	}

	dateTime, err := time.Parse("2006-01-02 15:04", From)
	if err != nil {
		return DateFrom, err
	} else {
		DateFrom = dateTime.Format("2006-01-02 15:04:00")
	}
	return DateFrom, nil
}

func checkTo(To string) (DateTo string, err error) {

	if len(To) == 0 {
		To = time.Now().UTC().Format("2006-01-02 15:04")
	}

	dateTime, err := time.Parse("2006-01-02 15:04", To)
	if err != nil {
		return DateTo, err
	} else {
		DateTo = dateTime.Format("2006-01-02 15:04:00")
	}
	return DateTo, nil
}

func getCurrencyRowName(code string) string {

	switch code {
	case "usdrub":
		return "usd"
	case "eurrub":
		return "eur"
	default:
		return "usd"
	}
}

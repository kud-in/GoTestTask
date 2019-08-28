package handlers

import (
	"TestTask/models"
	"TestTask/renderings"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo"
	"net/http"
	"strings"
)

func GetStatus(ctx echo.Context) error {

	resp := renderings.StatusResponse{}
	resp.Payload = []renderings.IntervalAverage{}

	resp.Message = "Bad Request"
	resp.Code = http.StatusBadRequest

	var code = ctx.Param("code")
	if len(code) == 0 {
		code = "usdrub"
	}

	db := ctx.Get(models.DBContextKey).(*sql.DB)

	query := `SELECT 
					printf("%.4f", last ) AS last, 
					printf("%.4f", day ) AS day, 
					printf("%.4f", week ) AS week, 
					printf("%.4f", month ) AS month
				FROM "average_for_interval" 
				WHERE code = upper(?) ORDER BY utc LIMIT 1`

	row := db.QueryRow(query, code)

	var average = renderings.IntervalAverage{}

	average.Code = strings.ToUpper(code)

	err := row.Scan(&average.Last, &average.Day, &average.Week, &average.Month)
	switch err {
	case sql.ErrNoRows:
		resp.Code = http.StatusNotFound
		resp.Message = "Page Not Found"

	case nil:
		resp.Message = "Success"
		resp.Code = http.StatusOK
		resp.Payload = append(resp.Payload, average)
	default:
		panic(err)
	}

	return ctx.JSON(resp.Code, resp)
}

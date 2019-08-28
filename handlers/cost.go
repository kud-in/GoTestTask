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

func GetCost(ctx echo.Context) error {

	var(
		code = checkCode(ctx.Param("code"))
		db = ctx.Get(models.DBContextKey).(*sql.DB)
		tableRow = getCurrencyRowName(code)
		resp = NewCostResponse()
		cost = renderings.Cost{}
	)



	datetime, err := checkTo(ctx.QueryParam("datetime"))

	if err != nil {
		resp.Message = "Invalid query parameter 'datetime'"
	} else {

		query := `SELECT printf("%.4f", ` + tableRow + ` ) AS rate FROM "average_per_minute" WHERE utc = ? LIMIT 1`

		row := db.QueryRow(query, datetime)

		err := row.Scan(&cost.Rate)

		switch err {
		case sql.ErrNoRows:
			resp.Code = http.StatusNotFound
			resp.Message = "Page Not Found"
		case nil:
			resp.Message = "Success"
			resp.Code = http.StatusOK
			cost.Code = strings.ToUpper(code)
			resp.Payload = append(resp.Payload, cost)
		default:
			panic(err)
		}
	}

	return ctx.JSON(resp.Code, resp)
}

func NewCostResponse() renderings.CostResponse {
	resp := renderings.CostResponse{}
	resp.Payload = []renderings.Cost{}
	resp.Message = "Bad Request"
	resp.Code = http.StatusBadRequest
	return resp
}

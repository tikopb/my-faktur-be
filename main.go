package main

import (
	"bemyfaktur/internal/database"
	"fmt"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	db := database.GetDb()

	fmt.Sprintln(db)

	//after all set, push the start
	e.Logger.Fatal(e.Start((":4000")))

}

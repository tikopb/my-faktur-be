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

	e.Logger.Fatal(e.Start((":3000")))

}

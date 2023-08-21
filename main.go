package main

import (
	"bemyfaktur/internal/database"

	paRepo "bemyfaktur/internal/repository/partner"

	paUsecase "bemyfaktur/internal/usecase/partner"

	"bemyfaktur/internal/delivery/rest"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	db := database.GetDb()

	partnerRepo := paRepo.GetRepository(db)
	parnerUsecase := paUsecase.GetUsecase(partnerRepo)

	h := rest.NewHandler(parnerUsecase)

	rest.LoadRoutes(e, h)
	//after all set, push the start
	e.Logger.Fatal(e.Start((":4000")))

}

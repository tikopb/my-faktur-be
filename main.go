package main

import (
	"bemyfaktur/internal/database"

	"bemyfaktur/internal/delivery/rest"

	"bemyfaktur/internal/usecase"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	db := database.GetDb()

	container := usecase.NewContainer(db)
	h := rest.NewHandler(container.PartnerUsecase, container.ProductUsecase, container.InvoiceUsecase, container.PaymentUsecase, container.PgUtil, container.AuthUsecase, container.Middleware, db)

	rest.LoadRoutes(e, h)

	//after all set, push the start
	e.Logger.Fatal(e.Start((":4000")))

}

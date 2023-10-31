package main

import (
	"bemyfaktur/internal/database"

	"bemyfaktur/internal/delivery/logger"
	"bemyfaktur/internal/delivery/rest"

	"bemyfaktur/internal/usecase"

	"github.com/labstack/echo/v4"
)

func main() {
	logger.Init()
	e := echo.New()
	rest.LoadMiddlewares(e)

	db := database.GetDb()

	container := usecase.NewContainer(db)
	h := rest.NewHandler(container.PartnerUsecase, container.ProductUsecase, container.InvoiceUsecase, container.PaymentUsecase, container.PgUtil, container.AuthUsecase, container.Middleware, db)

	rest.LoadRoutes(e, h)

	e.Logger.Fatal(e.Start((":4000")))

}

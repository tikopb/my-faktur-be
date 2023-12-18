package main

import (
	"bemyfaktur/internal/database"
	"bemyfaktur/internal/database/seeders"

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

	seeders.DBSeed(db)

	e.Logger.Fatal(e.Start((":4000")))

}

// func (server *Server) initCommands(config AppConfig, dbConfig DBConfig) {
// 	server.initializeDB(dbConfig)

// 	cmdApp := cli.NewApp()
// 	cmdApp.Commands = []cli.Command{
// 		{
// 			Name: "db:migrate",
// 			Action: func(c *cli.Context) error {
// 				server.dbMigrate()
// 				return nil
// 			},
// 		},
// 		{
// 			Name: "db:seed",
// 			Action: func(c *cli.Context) error {
// 				err := seeders.DBSeed(server.DB)
// 				if err != nil {
// 					log.Fatal(err)
// 				}

// 				return nil
// 			},
// 		},
// 	}

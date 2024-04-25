package main

import (
	"bemyfaktur/internal/database"
	"bemyfaktur/internal/database/seeders"
	"flag"
	"log"
	"os"

	"bemyfaktur/internal/delivery/logger"
	"bemyfaktur/internal/delivery/rest"

	"bemyfaktur/internal/usecase"

	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"github.com/urfave/cli"
)

func main() {
	logger.Init()
	e := echo.New()
	rest.LoadMiddlewares(e)

	db := database.GetDb(false)

	container := usecase.NewContainer(db)
	h := rest.NewHandler(container.PartnerUsecase, container.ProductUsecase, container.InvoiceUsecase, container.PaymentUsecase, container.FileserviceUsecase, container.PgUtil, container.AuthUsecase, container.Middleware, db)

	rest.LoadRoutes(e, h)

	flag.Parse()
	arg := flag.Arg(0)

	appPort := GetEnvVariabel("be_port")
	if appPort == "" {
		e.Logger.Fatal("be_port config is required on env")
	}

	if arg != "" {
		initCommands()
	} else {
		e.Logger.Fatal(e.Start((":" + appPort)))
	}
}

func initCommands() {
	cmdApp := cli.NewApp()
	db := database.GetDb(true)

	cmdApp.Commands = []cli.Command{
		{
			Name: "db:seed",
			Action: func(c *cli.Context) error {
				seeders.DBSeed(db)
				return nil
			},
		},
		{
			Name: "db:migrate",
			Action: func(c *cli.Context) error {
				seeders.MigrateDb(db)
				return nil
			},
		},
	}

	err := cmdApp.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func GetEnvVariabel(envName string) string {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		panic("env of " + envName + "not found")
	}

	value := viper.GetString(envName)
	return value
}

package rest

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func LoadRoutes(e *echo.Echo, handler *handler) {
	//userAuth
	userGroup := e.Group("/user")
	userGroup.POST("/register", handler.RegisterUser)
	userGroup.POST("/login", handler.Login)
	userGroup.GET("/refresh", handler.RefreshSession)
	userGroup.GET("/logout", handler.LogOut)

	//partner
	partnerGroup := e.Group("/partner")
	partnerGroup.GET("s/", handler.IndexPartner, handler.middleware.CheckAuth)
	partnerGroup.GET("/:id", handler.GetPartner, handler.middleware.CheckAuth)
	partnerGroup.GET("/patial/", handler.PartialPartner, handler.middleware.CheckAuth)
	partnerGroup.POST("/", handler.CreatePartner, handler.middleware.CheckAuth)
	partnerGroup.PUT("/:id", handler.UpdatedPartner, handler.middleware.CheckAuth)
	partnerGroup.DELETE("/:id", handler.DeletePartner, handler.middleware.CheckAuth)

	//product
	productGroup := e.Group("/product")
	productGroup.GET("s/", handler.IndexProduct, handler.middleware.CheckAuth)
	productGroup.GET("/:id", handler.GetProduct, handler.middleware.CheckAuth)
	productGroup.GET("/patial/", handler.PartialProduct, handler.middleware.CheckAuth)
	productGroup.POST("/", handler.CreateProduct, handler.middleware.CheckAuth)
	productGroup.PUT("/:id", handler.UpdatedProduct, handler.middleware.CheckAuth)
	productGroup.DELETE("/:id", handler.DeleteProduct, handler.middleware.CheckAuth)

	//group invoice
	invoiceGroup := e.Group("/invoice")
	invoiceGroup.GET("s/", handler.IndexInvoice, handler.middleware.CheckAuth)
	invoiceGroup.GET("/:id", handler.GetInvoice, handler.middleware.CheckAuth)
	invoiceGroup.GET("/patial/", handler.Partialnvoice, handler.middleware.CheckAuth)
	invoiceGroup.POST("", handler.CreateInvoice, handler.middleware.CheckAuth)
	invoiceGroup.POST("/v2", handler.CreateInvoiceV2, handler.middleware.CheckAuth)
	invoiceGroup.PUT("/:id", handler.UpdateInvoice, handler.middleware.CheckAuth)
	invoiceGroup.DELETE("/:id", handler.DeleteInvoice, handler.middleware.CheckAuth)

	//group invoice
	invoiceGroupLine := e.Group("/invoiceline")
	invoiceGroupLine.GET("s/:id", handler.IndexInvoiceLine, handler.middleware.CheckAuth)
	invoiceGroupLine.GET("/:id", handler.GetInvoiceLine, handler.middleware.CheckAuth)
	invoiceGroupLine.POST("", handler.CreateInvoiceLine, handler.middleware.CheckAuth)
	invoiceGroupLine.PUT("/:id", handler.UpdatedInvoiceLine, handler.middleware.CheckAuth)
	invoiceGroupLine.DELETE("/:id", handler.DeleteInvoiceLine, handler.middleware.CheckAuth)

	//payment group
	paymentGroup := e.Group("/payment")
	paymentGroup.GET("s/", handler.IndexPayment, handler.middleware.CheckAuth)
	paymentGroup.GET("/:id", handler.Getpayment, handler.middleware.CheckAuth)
	paymentGroup.POST("", handler.CreatePayment, handler.middleware.CheckAuth)
	paymentGroup.POST("/v2", handler.CreatePaymentV2, handler.middleware.CheckAuth)
	paymentGroup.PUT("/:id", handler.UpdatePayment, handler.middleware.CheckAuth)
	paymentGroup.DELETE("/:id", handler.DeletePayment, handler.middleware.CheckAuth)

	//group invoice
	paymentGroupLine := e.Group("/paymentline")
	paymentGroupLine.GET("s/:id", handler.IndexPaymentLine, handler.middleware.CheckAuth)
	paymentGroupLine.GET("/:id", handler.GetPaymentLine, handler.middleware.CheckAuth)
	paymentGroupLine.POST("", handler.CreatePaymentLine, handler.middleware.CheckAuth)
	paymentGroupLine.PUT("/:id", handler.UpdatePaymentLine, handler.middleware.CheckAuth)
	paymentGroupLine.DELETE("/:id", handler.DeletePaymentLine, handler.middleware.CheckAuth)

	//fileservice
	e.GET("/v1/filestore", handler.GetTheFileBaseUrl)

	//cors test
	e.GET("/cors", handler.CorsTest)
}

func LoadMiddlewares(e *echo.Echo) {
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		AllowCredentials: true,
	}))
}

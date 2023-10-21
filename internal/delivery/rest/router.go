package rest

import (
	"github.com/labstack/echo/v4"
)

func LoadRoutes(e *echo.Echo, handler *handler) {
	authMiddleware := GetAuthMiddleware(handler.authUsecase)

	//userAuth
	userGroup := e.Group("/user")
	userGroup.POST("/register", handler.RegisterUser)
	userGroup.POST("/login", handler.Login)
	userGroup.POST("/refresh", handler.RefreshSession)
	userGroup.POST("/logout", handler.LogOut)

	//partner
	partnerGroup := e.Group("/partner")
	partnerGroup.GET("", handler.IndexPartner, authMiddleware.CheckAuth)
	partnerGroup.GET("/:id", handler.GetPartner, authMiddleware.CheckAuth)
	partnerGroup.POST("", handler.CreatePartner, authMiddleware.CheckAuth)
	partnerGroup.PUT("/:id", handler.UpdatedPartner, authMiddleware.CheckAuth)
	partnerGroup.DELETE("/:id", handler.DeletePartner, authMiddleware.CheckAuth)

	//group
	productGroup := e.Group("/product")
	productGroup.GET("", handler.IndexProduct, authMiddleware.CheckAuth)
	productGroup.GET("/:id", handler.GetProduct, authMiddleware.CheckAuth)
	productGroup.POST("", handler.CreateProduct, authMiddleware.CheckAuth)
	productGroup.PUT("/:id", handler.UpdatedProduct, authMiddleware.CheckAuth)
	productGroup.DELETE("/:id", handler.DeleteProduct, authMiddleware.CheckAuth)

	//group invoice
	invoiceGroup := e.Group("/invoice")
	invoiceGroup.GET("", handler.IndexInvoice, authMiddleware.CheckAuth)
	invoiceGroup.GET("/:id", handler.GetInvoice, authMiddleware.CheckAuth)
	invoiceGroup.POST("", handler.CreateInvoice, authMiddleware.CheckAuth)
	invoiceGroup.PUT("/:id", handler.UpdateInvoice, authMiddleware.CheckAuth)
	invoiceGroup.DELETE("/:id", handler.DeleteInvoice, authMiddleware.CheckAuth)

	//group invoice
	invoiceGroupLine := e.Group("/invoiceline")
	invoiceGroupLine.GET("", handler.IndexInvoiceLine, authMiddleware.CheckAuth)
	invoiceGroupLine.GET("/:id", handler.GetInvoiceLine, authMiddleware.CheckAuth)
	invoiceGroupLine.POST("", handler.CreateInvoiceLine, authMiddleware.CheckAuth)
	invoiceGroupLine.PUT("/:id", handler.UpdatedInvoiceLine, authMiddleware.CheckAuth)
	invoiceGroupLine.DELETE("/:id", handler.DeleteInvoiceLine, authMiddleware.CheckAuth)

	//payment group
	paymentGroup := e.Group("/payment")
	paymentGroup.GET("", handler.IndexPayment, authMiddleware.CheckAuth)
	paymentGroup.GET("/:id", handler.Getpayment, authMiddleware.CheckAuth)
	paymentGroup.POST("", handler.CreatePayment, authMiddleware.CheckAuth)
	paymentGroup.PUT("/:id", handler.UpdatePayment, authMiddleware.CheckAuth)
	paymentGroup.DELETE("/:id", handler.DeletePayment, authMiddleware.CheckAuth)

	//group invoice
	paymentGroupLine := e.Group("/paymentline")
	paymentGroupLine.GET("", handler.IndexPaymentLine, authMiddleware.CheckAuth)
	paymentGroupLine.GET("/:id", handler.GetPaymentLine, authMiddleware.CheckAuth)
	paymentGroupLine.POST("", handler.CreatePaymentLine, authMiddleware.CheckAuth)
	paymentGroupLine.PUT("/:id", handler.UpdatePaymentLine, authMiddleware.CheckAuth)
	paymentGroupLine.DELETE("/:id", handler.DeletePaymentLine, authMiddleware.CheckAuth)

}

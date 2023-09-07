package rest

import "github.com/labstack/echo/v4"

func LoadRoutes(e *echo.Echo, handler *handler) {

	//partner
	partnerGroup := e.Group("/partner")
	partnerGroup.GET("", handler.IndexPartner)
	partnerGroup.GET("/:id", handler.GetPartner)
	partnerGroup.POST("", handler.CreatePartner)
	partnerGroup.PUT("/:id", handler.UpdatedPartner)
	partnerGroup.DELETE("/:id", handler.DeletePartner)

	//group
	productGroup := e.Group("/product")
	productGroup.GET("", handler.IndexProduct)
	productGroup.GET("/:id", handler.GetProduct)
	productGroup.POST("", handler.CreateProduct)
	productGroup.PUT("/:id", handler.UpdatedProduct)
	productGroup.DELETE("/:id", handler.DeleteProduct)

	//group invoice
	invoiceGroup := e.Group("/invoice")
	invoiceGroup.GET("", handler.IndexInvoice)
	invoiceGroup.GET("/:id", handler.GetInvoice)
	invoiceGroup.POST("", handler.CreateInvoice)
	invoiceGroup.PUT("/:id", handler.UpdateInvoice)
	invoiceGroup.DELETE("/:id", handler.DeleteInvoice)

	//group invoice
	invoiceGroupLine := e.Group("/invoiceline")
	invoiceGroupLine.GET("", handler.IndexInvoiceLine)
	invoiceGroupLine.GET("/:id", handler.GetInvoiceLine)
	invoiceGroupLine.POST("", handler.CreateInvoiceLine)
	invoiceGroupLine.PUT("/:id", handler.UpdatedInvoiceLine)
	invoiceGroupLine.DELETE("/:id", handler.DeleteInvoiceLine)

	//payment group
	paymentGroup := e.Group("/payment")
	paymentGroup.GET("", handler.IndexPayment)
	paymentGroup.GET("/:id", handler.Getpayment)
	paymentGroup.POST("", handler.CreatePayment)
	paymentGroup.PUT("/:id", handler.UpdatePayment)
	paymentGroup.DELETE("/:id", handler.DeletePayment)

	//group invoice
	paymentGroupLine := e.Group("/paymentline")
	paymentGroupLine.GET("", handler.IndexPaymentLine)
	paymentGroupLine.GET("/:id", handler.GetPaymentLine)
	paymentGroupLine.POST("", handler.CreatePaymentLine)
	paymentGroupLine.PUT("/:id", handler.UpdatePaymentLine)
	paymentGroupLine.DELETE("/:id", handler.DeletePaymentLine)
	userGroup := e.Group("/user")
	userGroup.GET("/:id", handler.Getuser)

}

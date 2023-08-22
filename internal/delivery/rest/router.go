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

}

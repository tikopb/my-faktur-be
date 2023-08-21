package rest

import "github.com/labstack/echo/v4"

func LoadRoutes(e *echo.Echo, handler *handler) {

	partnerGroup := e.Group("/partner")
	partnerGroup.GET("", handler.IndexPartner)
	partnerGroup.GET("/:id", handler.GetPartner)
	partnerGroup.POST("", handler.CreatePartner)
	partnerGroup.PUT("/:id", handler.UpdatedPartner)
	partnerGroup.DELETE("/:id", handler.DeletePartner)

}

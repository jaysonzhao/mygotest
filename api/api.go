package api

import (
	"github.com/jaysonzhao/gotest/handlers"

	"github.com/labstack/echo"
)

func MainGroup(e *echo.Echo) {
	// Route / to handler function
	e.GET("/health-check", handlers.HealthCheck)

	e.GET("/cats/:data", handlers.GetCats)
	e.GET("/pods", handlers.GetPods)
	e.POST("/cats", handlers.AddCat)

}

func AdminGroup(g *echo.Group) {
	g.GET("/main", handlers.MainAdmin)
}

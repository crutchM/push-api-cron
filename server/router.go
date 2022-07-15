package server

import "github.com/gin-gonic/gin"

type Router struct {
}

func NewRouter() *Router {
	return &Router{}
}

func (r *Router) InitRoutes() *gin.Engine {
	router := gin.New()
	groups := router.Group("/groups")
	{
		groups.GET("/:app_id")
		groups.POST("/")
		groups.POST("/start:id")
		groups.POST("/stop:id")
	}
	device := router.Group("/device")
	{
		device.POST("/createToken")
		device.POST("/updateToken")

	}
	return router
}

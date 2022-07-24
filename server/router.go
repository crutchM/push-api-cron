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
		groups.GET("/:app_id")   // получить данные о группе
		groups.POST("/")         //создать
		groups.POST("/start:id") // запустить отправку пушей
		groups.POST("/stop:id")  // остановить отправку пушей
	}
	device := router.Group("/device")
	{
		device.POST("/createToken") // запись токена устройства в базу
		device.POST("/updateToken") // получение токена устройства в базу

	}
	return router
}

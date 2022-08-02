package server

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"push-api-cron/core/models"
	"push-api-cron/core/models/device"
	"push-api-cron/data/service"
)

type Router struct {
	service service.Service
	ch      chan struct{}
}

func NewRouter(service2 service.Service) *Router {
	return &Router{
		service: service2,
		ch:      make(chan struct{}),
	}
}

func (r *Router) InitRoutes() *gin.Engine {
	router := gin.New()
	groups := router.Group("/groups")
	{
		groups.POST("/", r.CreateGroup) //создать
		groups.POST("/start", r.Start)  // запустить отправку пушей
		groups.POST("/stop", r.Stop)    // остановить отправку пушей
	}
	device := router.Group("/device")
	{
		device.POST("/add", r.AddDevice) // запись токена устройства в базу

	}
	return router
}

type input struct {
	Group    int             `json:"group_id"`
	Messages models.Messages `json:"messages"`
	Interval int             `json:"interval"`
}

func (r *Router) Stop(c *gin.Context) {
	r.service.Stop(r.ch)
}
func (r *Router) Start(c *gin.Context) {
	var inp input
	if err := c.BindJSON(&inp); err != nil {
		logrus.Info()
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}
	if err := r.service.Start(r.ch, inp.Group, inp.Messages, inp.Interval); err != nil {
		logrus.Info()
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"status": "success",
	})
	return
}

func (r *Router) AddDevice(c *gin.Context) {
	var inp device.Device
	if err := c.BindJSON(&inp); err != nil {
		logrus.Info()
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}
	if err := r.service.AddDevice(inp); err != nil {
		logrus.Info()
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"status": "success",
	})
}

func (r *Router) CreateGroup(c *gin.Context) {
	var inp models.InputGroup

	if err := c.BindJSON(&inp); err != nil {
		logrus.Info(err.Error())
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}
	res, err := r.service.CreateGroup(inp)
	if err != nil {
		logrus.Info(err.Error())
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"new_group": res,
	})
}

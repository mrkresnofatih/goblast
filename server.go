package goblast

import (
	"sync"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type IServer interface {
	Initialize()
	AddController(controller IController)
}

// Description: a server object that renders all endpoints and handlers along with additional configs
type ApplicationServer struct {
	RunState    *sync.WaitGroup
	Router      *echo.Echo
	Port        string
	Controllers []IController
	CorsConfig  ApplicationServerCorsConfiguration
}

// Description: simple cors config class/struct
type ApplicationServerCorsConfiguration struct {
	AllowMethods     []string
	AllowHeaders     []string
	AllowCredentials bool
	AllowOrigins     []string
}

// Description: call this method after all controllers have been added to ApplicationServer.Controllers[]
func (a *ApplicationServer) Initialize() {
	a.RunState.Add(1)
	go func() {
		a.Router = echo.New()
		a.Router.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowMethods:     a.CorsConfig.AllowHeaders,
			AllowHeaders:     a.CorsConfig.AllowHeaders,
			AllowCredentials: a.CorsConfig.AllowCredentials,
			AllowOrigins:     a.CorsConfig.AllowOrigins,
		}))
		for _, controller := range a.Controllers {
			controller.Register(a.Router)
		}
		a.Router.Logger.Fatal(a.Router.Start(":" + a.Port))
		a.RunState.Done()
	}()
}

// Description: to add one-by-one controllers to the ApplicationServer
func (a *ApplicationServer) AddController(controller IController) {
	a.Controllers = append(a.Controllers, controller)
}

type IController interface {
	Register(echo *echo.Echo)
}

type IEndpoint interface {
	GetHandler() echo.HandlerFunc
	GetPath() string
	Register(group *echo.Group)
}

type IRouter interface {
	Build()
	AddEndpoint(endpoint IEndpoint)
}

// Description: a controller object, for registering all endpoints under one specific path prefix
type ControllerRouter struct {
	MainRouter *echo.Echo
	router     *echo.Group
	PathPrefix string
	Endpoints  []IEndpoint
}

// Description: call this method when all endpoints have been registered to the controller
func (c *ControllerRouter) Build() {
	c.router = c.MainRouter.Group(c.PathPrefix)
	for _, endpoint := range c.Endpoints {
		endpoint.Register(c.router)
	}
}

// Description: for adding one-by-one endpoints to controller
func (c *ControllerRouter) AddEndpoint(endpoint IEndpoint) {
	c.Endpoints = append(c.Endpoints, endpoint)
}

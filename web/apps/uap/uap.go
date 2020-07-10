package uap

import (
	"github.com/gin-gonic/gin"
	"github.com/oceanho/gw/web/apps/uap/api"
)

func init() {
}

type App struct {
}

func New() *App {
	return &App{}
}

func (u App) Name() string {
	return "oceanho.uap"
}

func (u App) BaseRouter() string {
	return "uap"
}

func (u App) Register(router *gin.RouterGroup) {
	router.GET("tenant/<uid:int64>", api.CreateTenant)
}

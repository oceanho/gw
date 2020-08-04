package stor

import (
	"github.com/oceanho/gw"
	"github.com/oceanho/gw/web/apps/stor/api"
)

func init() {
}

type App struct {
}

func New() App {
	return App{}
}

func (a App) Name() string {
	return "gw.stor"
}

func (a App) Router() string {
	return "stor"
}

func (a App) Register(router *gw.RouterGroup) {
	router.GET("object/create", api.CreateObject)
	router.POST("object/modify", api.ModifyObject)
}

func (a App) Migrate(store gw.Store) {
	// db := store.GetDbStore()
}

func (a App) Use(opt *gw.ServerOption) {
}
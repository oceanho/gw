package Decorator

import (
	"github.com/oceanho/gw"
	"github.com/oceanho/gw/contrib/apps/uap/Const"
)

var (
	resOpCheckFs map[string]gw.DecoratorHandler
)

func init() {
	initial()
	initialRoleDecorators()
	initialUserDecorators()
}

func initial() {
	resOpCheckFs = make(map[string]gw.DecoratorHandler)
}

func generalCheck(ctx *gw.Context) (r bool, status int, err error, payload interface{}) {
	author := ctx.User()
	if author.IsUser() {
		return false, 403, gw.ErrPermissionDenied, "Non-user has no this permission"
	}
	return true, 0, nil, nil
}
func initialRoleDecorators() {
	resOpCheckFs[Const.RoleDecorator.CreationPermKey()] = func(ctx *gw.Context) (status int, err error, payload interface{}) {
		if ok, s, e, p := generalCheck(ctx); !ok {
			return s, e, p
		}
		return 0, nil, nil
	}
	resOpCheckFs[Const.RoleDecorator.ModificationPermKey()] = func(ctx *gw.Context) (status int, err error, payload interface{}) {
		if ok, s, e, p := generalCheck(ctx); !ok {
			return s, e, p
		}
		return 0, nil, nil
	}
	resOpCheckFs[Const.RoleDecorator.ModificationPermKey()] = func(ctx *gw.Context) (status int, err error, payload interface{}) {
		if ok, s, e, p := generalCheck(ctx); !ok {
			return s, e, p
		}
		return 0, nil, nil
	}
	resOpCheckFs[Const.RoleDecorator.ReadAllPermKey()] = func(ctx *gw.Context) (status int, err error, payload interface{}) {
		if ok, s, e, p := generalCheck(ctx); !ok {
			return s, e, p
		}
		return 0, nil, nil
	}
	resOpCheckFs[Const.RoleDecorator.ReadDetailPermKey()] = func(ctx *gw.Context) (status int, err error, payload interface{}) {
		if ok, s, e, p := generalCheck(ctx); !ok {
			return s, e, p
		}
		return 0, nil, nil
	}
}

func initialUserDecorators() {
	resOpCheckFs[Const.UserDecorator.CreationPermKey()] = func(ctx *gw.Context) (status int, err error, payload interface{}) {
		if ok, s, e, p := generalCheck(ctx); !ok {
			return s, e, p
		}
		return 0, nil, nil
	}
	resOpCheckFs[Const.UserDecorator.ModificationPermKey()] = func(ctx *gw.Context) (status int, err error, payload interface{}) {
		if ok, s, e, p := generalCheck(ctx); !ok {
			return s, e, p
		}
		return 0, nil, nil
	}
	resOpCheckFs[Const.UserDecorator.ModificationPermKey()] = func(ctx *gw.Context) (status int, err error, payload interface{}) {
		if ok, s, e, p := generalCheck(ctx); !ok {
			return s, e, p
		}
		return 0, nil, nil
	}
	resOpCheckFs[Const.UserDecorator.ReadAllPermKey()] = func(ctx *gw.Context) (status int, err error, payload interface{}) {
		if ok, s, e, p := generalCheck(ctx); !ok {
			return s, e, p
		}
		return 0, nil, nil
	}
	resOpCheckFs[Const.UserDecorator.ReadDetailPermKey()] = func(ctx *gw.Context) (status int, err error, payload interface{}) {
		if ok, s, e, p := generalCheck(ctx); !ok {
			return s, e, p
		}
		return 0, nil, nil
	}
}

func CheckPermissionDecorator(permKey string) *gw.Decorator {
	if f, ok := resOpCheckFs[permKey]; ok {
		return &gw.Decorator{
			Catalog: Const.DecoratorCatalog,
			Before:  f,
		}
	}
	return nil
}

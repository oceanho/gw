package gwImpls

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/oceanho/gw"
	"github.com/oceanho/gw/contrib/apps/uap/conf"
	"github.com/oceanho/gw/contrib/apps/uap/dbModel"
	"github.com/oceanho/gw/libs/gwjsoner"
	"github.com/oceanho/gw/logger"
	"gorm.io/gorm"
	"time"
)

type UserManager struct {
	gw.ServerState
	cachePrefix      string
	cacheStoreName   string
	backendStoreName string
	cacheExpiration  time.Duration
	permPagerExpr    gw.PagerExpr
}

func (u UserManager) Backend() *gorm.DB {
	return u.Store().GetDbStoreByName(u.backendStoreName)
}

func (u UserManager) Cache() *redis.Client {
	return u.Store().GetCacheStoreByName(u.cacheStoreName)
}

func (u UserManager) User() *gorm.DB {
	return u.Backend().Model(dbModel.User{})
}

func (u UserManager) GetFromCache(passport string) gw.User {
	var bytes, err = u.Cache().Get(context.Background(), u.CacheKey(passport)).Bytes()
	if err != nil {
		logger.Error("operation cache fail, err: %v", err)
		return gw.EmptyUser
	}
	var user gw.User
	if gwjsoner.Unmarshal(bytes, &user) != nil {
		return gw.EmptyUser
	}
	return user
}

func (u UserManager) CacheKey(passport string) string {
	return fmt.Sprintf("%s.%s", u.cachePrefix, passport)
}

func (u UserManager) SaveToCache(user gw.User) {
}

//
// APIs
//
func (u UserManager) Create(user *gw.User) error {
	store := u.Store()
	db := store.GetDbStore()
	var model dbModel.User
	err := db.First(&model, "tenant_id=? and passport=?", user.TenantId, user.Passport).Error
	if err != nil && err.Error() != "record not found" {
		return err
	}
	if model.ID > 0 {
		return gw.ErrorUserHasExists
	}

	// default
	model.IsUser = false
	model.IsAdmin = false
	model.IsTenancy = false

	// passport
	model.Passport = user.Passport
	model.TenantId = user.TenantId
	model.Secret = user.Secret

	switch user.UserType {
	case gw.Administrator:
		model.IsAdmin = true
		break
	case gw.Tenancy:
		model.IsTenancy = true
		break
	case gw.NonUser:
		model.IsUser = true
		break
	}
	tx := store.GetDbStore().Begin()
	err = tx.Create(&model).Error
	if err != nil {
		return err
	}
	err = tx.Commit().Error
	user.ID = model.ID
	return err
}

func (u UserManager) Modify(user gw.User) error {
	panic("implement me")
}

func (u UserManager) Delete(tenantId, userId uint64) error {
	panic("implement me")
}

func (u UserManager) QueryByUser(tenantId uint64, passport, password string) (gw.User, error) {
	var user gw.User
	var err = u.User().Find(&user, "tenant_id=? and passport=? and secret=?", tenantId, passport, password).Error
	if err != nil {
		return gw.EmptyUser, err
	}
	if user.IsEmpty() {
		return user, gw.ErrorUserNotFound
	}
	_, perms, err := u.PermissionManager().QueryByUser(user.TenantId, user.ID, gw.DefaultPageExpr)
	if err != nil {
		return user, err
	}
	user.Permissions = make([]gw.Permission, len(perms))
	copy(user.Permissions, perms)
	return user, nil
}

func (u UserManager) QueryByAKS(tenantId uint64, accessKey, accessSecret string) (gw.User, error) {
	panic("implement me")
}

func (u UserManager) Query(tenantId, userId uint64) (gw.User, error) {
	panic("implement me")
}

func (u UserManager) QueryList(tenantId uint64, expr gw.PagerExpr, total int64, out []gw.User) error {
	panic("implement me")
}

func DefaultUserManager(state gw.ServerState) UserManager {
	var cnf = conf.GetUAP(state.ApplicationConfig())
	return UserManager{
		ServerState:      state,
		cacheStoreName:   cnf.User.Backend.Name,
		backendStoreName: cnf.User.Cache.Name,
		cachePrefix:      cnf.User.Cache.Prefix,
		cacheExpiration:  time.Hour * time.Duration(cnf.User.Cache.ExpirationHours), // One week.
	}
}

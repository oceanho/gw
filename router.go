package gw

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/oceanho/gw/conf"
	"github.com/oceanho/gw/logger"
	"net/http"
	"path"
	"reflect"
	"runtime"
	"sync"
	"time"
)

const (
	gwRouterInfoKey = "gw-router"
	gwDbContextKey  = "gw-context"
)

// Handler defines a http handler for gw framework.
type Handler func(ctx *Context)

// Context represents a gw Context object, it's extension from gin.Context.
type Context struct {
	*gin.Context
	requestId  string
	user       User
	startTime  time.Time
	logger     Logger
	server     *HostServer
	state      ContextState
	queries    map[string][]string
	params     map[string]interface{}
	bindModels map[string]interface{}
}

// ServerState represents a Server state context object.
type ContextState struct {
	store  IStore
	state  ServerState
	values map[string]interface{}
}

const (
	IStoreName               = "github.com/oceanho/gw.IStore"
	IPasswordSignerName      = "github.com/oceanho/gw.IPasswordSigner"
	IPermissionManagerName   = "github.com/oceanho/gw.IPermissionManager"
	IPermissionCheckerName   = "github.com/oceanho/gw.IPermissionChecker"
	ISessionStateManagerName = "github.com/oceanho/gw.ISessionStateManager"
	IdentifierGeneratorName  = "github.com/oceanho/gw.IdentifierGenerator"
	IAuthManagerName         = "github.com/oceanho/gw.IAuthManager"
	IUserManagerName         = "github.com/oceanho/gw.IUserManager"
	ContextStateName         = "github.com/oceanho/gw.ContextState"
	ServerStateName          = "github.com/oceanho/gw.ServerState"
	HostServerName           = "github.com/oceanho/gw.HostServer"
)

var nullReflectValue = reflect.ValueOf(nil)

func (state ContextState) objectTypers() map[string]ObjectTyper {
	var typers = make(map[string]ObjectTyper)
	typers[IStoreName] = newObjectTyper(IStoreName, state.store)
	typers[IPermissionManagerName] = newObjectTyper(IPermissionManagerName, state.state.PermissionManager())
	typers[IAuthManagerName] = newObjectTyper(IAuthManagerName, state.state.AuthManager())
	typers[IPermissionCheckerName] = newObjectTyper(IPermissionCheckerName, state.state.PermissionChecker())
	typers[ISessionStateManagerName] = newObjectTyper(ISessionStateManagerName, state.state.SessionStateManager())
	typers[IPasswordSignerName] = newObjectTyper(IPasswordSignerName, state.state.PasswordSigner())
	typers[IUserManagerName] = newObjectTyper(IUserManagerName, state.state.UserManager())
	typers[IdentifierGeneratorName] = newObjectTyper(IdentifierGeneratorName, state.state.IDGenerator())
	typers[ContextStateName] = newObjectTyper(ContextStateName, state)
	typers[HostServerName] = newObjectTyper(HostServerName, state.state.s)
	typers[ServerStateName] = newObjectTyper(HostServerName, state.state)
	return typers
}

func (c *Context) RequestId() string {
	return c.requestId
}

func (c Context) User() User {
	return c.user
}

func (c Context) Store() IStore {
	return c.state.store
}

func (c Context) CryptoHash() ICryptoHash {
	return c.state.state.CryptoHash()
}

func (c Context) CryptoProtect() ICryptoProtect {
	return c.state.state.CryptoProtect()
}

func (c Context) PasswordSigner() IPasswordSigner {
	return c.state.state.PasswordSigner()
}

func (c Context) AuthManager() IAuthManager {
	return c.state.state.AuthManager()
}

func (c Context) SessionStateManager() ISessionStateManager {
	return c.state.state.SessionStateManager()
}

func (c Context) PermissionManager() IPermissionManager {
	return c.state.state.PermissionManager()
}

func (c Context) PermissionChecker() IPermissionChecker {
	return c.PermissionManager().Checker()
}

func (c Context) UserManager() IUserManager {
	return c.state.state.UserManager()
}

func (c Context) IDGenerator() IdentifierGenerator {
	return c.state.state.IDGenerator()
}

func (c Context) Resolve(typerName string) interface{} {
	return c.state.state.s.DIProvider.ResolveWithState(c.state, typerName)
}

func (c Context) ResolveByTyper(typer reflect.Type) interface{} {
	return c.state.state.s.DIProvider.ResolveByTyperWithState(c.state, typer)
}

func (c Context) RespBodyBuildFunc() RespBodyBuildFunc {
	return c.state.state.RespBodyBuildFunc()
}

// Router represents a gw's Router info.
type Router struct {
	locker        sync.Mutex
	prefix        string
	server        *gin.Engine
	router        *gin.RouterGroup
	currentRouter *gin.RouterGroup
	routerInfos   []RouterInfo
}

// RouterGroup represents a gw's Group Router info.
type RouterGroup struct {
	*Router
}

// RouterInfo represents a gw Router Info.
type RouterInfo struct {
	Method            string
	UrlPath           string
	Handler           Handler
	Decorators        []Decorator
	Permissions       []Permission
	handlerActionName string
	beforeDecorators  []Decorator
	afterDecorators   []Decorator
}

func (r *RouterInfo) String() string {
	// https://books.studygolang.com/gobyexample/string-formatting
	// Padding left.
	before, after := splitDecorators(r.Decorators...)
	return fmt.Sprintf("%-8s%s -> %s "+
		"[decorators(%d before, %d after)]", r.Method, r.UrlPath, r.handlerActionName, len(before), len(after))
}

func (router *Router) createRouter(method, relativePath string, handler Handler, handlerActionName string, decorators ...Decorator) {
	urlPath := path.Join(router.currentRouter.BasePath(), relativePath)
	//var _decorators []Decorator
	//_decorators = append(_decorators, NewStoreDbSetupDecorator(func(ctx Context, db *gorm.DB) *gorm.DB {
	//	return db.Set(gwDbUserInfoKey, ctx.User())
	//}))
	//_decorators = append(_decorators, decorators...)
	routerInfo := RouterInfo{
		Method:     method,
		UrlPath:    urlPath,
		Handler:    handler,
		Decorators: decorators,
	}
	beforeDecorators, afterDecorators := splitDecorators(decorators...)
	pds := FilterDecorator(func(d Decorator) bool {
		return d.Catalog == permissionDecoratorCatalog
	}, decorators...)
	var perms []Permission
	for _, p := range pds {
		if pd, ok := p.MetaData.([]Permission); ok {
			perms = append(routerInfo.Permissions, pd...)
		}
	}
	if handlerActionName == "" {
		handlerActionName = getHandlerFullName(handler)
	}
	routerInfo.Permissions = perms
	routerInfo.afterDecorators = afterDecorators
	routerInfo.beforeDecorators = beforeDecorators
	routerInfo.handlerActionName = handlerActionName
	if method == "any" {
		router.currentRouter.Any(relativePath, func(c *gin.Context) {
			c.Set(gwRouterInfoKey, routerInfo)
			handle(c)
		})
		return
	}
	if method == "group" {
		router.currentRouter.Group(relativePath, func(c *gin.Context) {
			c.Set(gwRouterInfoKey, routerInfo)
			handle(c)
		})
		return
	}
	// gin router
	router.currentRouter.Handle(method, relativePath, func(c *gin.Context) {
		c.Set(gwRouterInfoKey, routerInfo)
		handle(c)
	})
	router.locker.Lock()
	defer router.locker.Unlock()
	router.routerInfos = append(router.routerInfos, routerInfo)
}

func getHandlerFullName(handler Handler) string {
	var val = runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name()
	return fmt.Sprintf("%s(ctx *Context)", val)
}

// Gw framework's APIs.
func (c Context) GetHostServer() HostServer {
	return *c.server
}

// Query returns a string from queries.
func (c Context) Query(key string) string {
	val := c.Context.Query(key)
	if val == "" {
		queries := c.QueryArray(key)
		if len(queries) > 0 {
			val = queries[0]
		}
	}
	return val
}

// QueryArray returns a string array from queries.
func (c Context) QueryArray(key string) []string {
	return c.queries[key]
}

// StartTime returns the Context start *time.Time
func (c Context) StartTime() *time.Time {
	return &c.startTime
}

// GET register a http Get router of handler.
func (router *Router) GET(relativePath string, handler Handler, decorators ...Decorator) {
	router.createRouter(http.MethodGet, relativePath, handler, "", decorators...)
}

// POST register a http POST router of handler.
func (router *Router) POST(relativePath string, handler Handler, decorators ...Decorator) {
	router.createRouter(http.MethodPost, relativePath, handler, "", decorators...)
}

// PUT register a http PUT router of handler.
func (router *Router) PUT(relativePath string, handler Handler, decorators ...Decorator) {
	router.createRouter(http.MethodPut, relativePath, handler, "", decorators...)
}

// HEAD register a http HEAD router of handler.
func (router *Router) HEAD(relativePath string, handler Handler, decorators ...Decorator) {
	router.createRouter(http.MethodHead, relativePath, handler, "", decorators...)
}

// DELETE register a http DELETE router of handler.
func (router *Router) DELETE(relativePath string, handler Handler, decorators ...Decorator) {
	router.createRouter(http.MethodDelete, relativePath, handler, "", decorators...)
}

// OPTIONS register a http OPTIONS router of handler.
func (router *Router) OPTIONS(relativePath string, handler Handler, decorators ...Decorator) {
	router.createRouter(http.MethodOptions, relativePath, handler, "", decorators...)
}

// PATCH register a http PATCH router of handler.
func (router *Router) PATCH(relativePath string, handler Handler, decorators ...Decorator) {
	router.createRouter(http.MethodPatch, relativePath, handler, "", decorators...)
}

// Any register a all HTTP methods router of handler.
func (router *Router) Any(relativePath string, handler Handler, decorators ...Decorator) {
	router.createRouter("any", relativePath, handler, "", decorators...)
}

// Any register a HTTP Connect router of handler.
func (router *Router) Connect(relativePath string, handler Handler, decorators ...Decorator) {
	router.createRouter(http.MethodConnect, relativePath, handler, "", decorators...)
}

// Trace register a HTTP Trace router of handler.
func (router *Router) Trace(relativePath string, handler Handler, decorators ...Decorator) {
	router.createRouter(http.MethodTrace, relativePath, handler, "", decorators...)
}

// Use register a gin Middleware of handler.
func (router *Router) Use(middleware ...gin.HandlerFunc) {
	router.currentRouter.Use(middleware...)
}

// Handlers returns the current router of gin.HandlersChain
func (router *Router) Handlers() gin.HandlersChain {
	return router.currentRouter.Handlers
}

// Group returns a new route group.
func (router *Router) Group(relativePath string, handler Handler, decorators ...Decorator) *RouterGroup {
	rg := &RouterGroup{
		router,
	}
	if handler == nil {
		rg.currentRouter = rg.router.Group(relativePath)
	} else {
		router.createRouter("group", relativePath, handler, "", decorators...)
	}
	rg.prefix = path.Join(router.prefix, relativePath)
	return rg
}

// Config returns a snapshot of the current Context's conf.Config object.
func (c *Context) Config() conf.ApplicationConfig {
	return config(c.Context)
}

// Bind define a Api that can be bind data to out object by gin.Context's Bind(...) APIs.
// It's auto response 400, invalid request parameter to client if bind fail.
// returns a error message for c.Bind(...).
func (c *Context) Bind(out interface{}) error {
	if err := c.Context.Bind(out); err != nil {
		c.JSON400Msg(400, fmt.Sprintf("invalid request parameters, details: \n%v", err))
		return err
	}
	return nil
}

// BindQuery define a Api that can be bind data to out object by gin.Context's Bind(...) APIs.
// It's auto response 400, invalid request parameter to client if bind fail.
// returns a error message for c.BindQuery(...).
func (c *Context) BindQuery(out interface{}) error {
	if err := c.Context.BindQuery(out); err != nil {
		c.JSON400Msg(400, fmt.Sprintf("invalid request parameters, details: \n%v", err))
		return err
	}
	return nil
}

func handle(c *gin.Context) {
	var router, ok = c.MustGet(gwRouterInfoKey).(RouterInfo)
	if !ok {
		logger.Error("invalid handler, can not be got RouterInfo.")
		return
	}
	var status int
	var err error
	var shouldStop bool = false
	var payload interface{}
	// action before Decorators
	var s = *hostServer(c)
	var requestID = getRequestID(s, c)
	var ctx = makeCtx(c, requestID)
	for _, d := range router.beforeDecorators {
		status, err, payload = d.Before(ctx)
		if err != nil || status != 0 {
			shouldStop = true
			break
		}
	}
	if shouldStop {
		if payload == "" {
			payload = "caller decorator fail."
		}
		status, body := parseErrToRespBody(s, requestID, payload, err)
		c.JSON(status, body)
		return
	}

	// process Action handler.
	router.Handler(ctx)

	// action after Decorators
	l := len(router.afterDecorators)
	if l < 1 {
		return
	}
	shouldStop = false
	for i := l - 1; i >= 0; i-- {
		status, err, payload = router.afterDecorators[i].After(ctx)
		if err != nil || status != 0 {
			shouldStop = true
			break
		}
	}
	if shouldStop {
		if payload == "" {
			payload = "caller decorator fail."
		}
		status, body := parseErrToRespBody(s, requestID, payload, err)
		c.JSON(status, body)
	}
}

func parseErrToRespBody(s HostServer, requestID string, msgBody interface{}, err error) (int, interface{}) {
	var status = http.StatusBadRequest
	if err == ErrPermissionDenied {
		status = http.StatusForbidden
	} else if err == ErrInternalServerError {
		status = http.StatusInternalServerError
	} else if err == ErrUnauthorized {
		status = http.StatusUnauthorized
	}
	return status, s.RespBodyBuildFunc(status, requestID, err.Error(), msgBody)
}

func makeCtx(c *gin.Context, requestID string) *Context {
	s := hostServer(c)
	user := getUser(c)
	serverState := ServerState{s: s}
	ctx := &Context{
		server:    s,
		Context:   c,
		user:      user,
		requestId: requestID,
		startTime: time.Now(),
		logger:    getLogger(c),
		state: ContextState{
			state: serverState,
		},
		queries:    make(map[string][]string),
		params:     make(map[string]interface{}),
		bindModels: make(map[string]interface{}),
	}
	var dbSetups []StoreDbSetupHandler
	dbSetups = append(dbSetups, s.storeDbSetupHandler)
	var cacheSetups []StoreCacheSetupHandler
	cacheSetups = append(cacheSetups, s.storeCacheSetupHandler)
	store := &backendWrapper{
		user:                    user,
		ctx:                     ctx,
		storeDbSetupHandlers:    dbSetups,
		storeCacheSetupHandlers: cacheSetups,
		store:                   serverState.Store(),
	}
	// FIXME(Ocean): has circular reference?
	ctx.state.store = store
	return ctx
}

func splitDecorators(decorators ...Decorator) (before, after []Decorator) {
	for _, d := range decorators {
		if d.After != nil {
			after = append(after, d)
		}
		if d.Before != nil {
			before = append(before, d)
		}
	}
	return before, after
}

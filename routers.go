package swagger

//type Router struct {
//	g      gin.IRouter
//	path   string
//	method string
//}
//
//func (g *Router) Path() string {
//	return g.path
//}
//
//func (g *Router) Use(middleware ...gin.HandlerFunc) *Router {
//	g.g.Use(middleware...)
//	return g
//}
//func (g *Router) Group(pth string, handlers ...gin.HandlerFunc) *Router {
//	gg := g.g.Group(pth)
//	gg.Use(handlers...)
//
//	return &Router{
//		g:    gg,
//		path: path.Join(g.path, pth),
//	}
//}
//
//func (g *Router) GET(pth string, handlers ...gin.HandlerFunc) *Router {
//	return g.Handle("GET", pth, handlers...)
//}
//func (g *Router) POST(pth string, handlers ...gin.HandlerFunc) *Router {
//	return g.Handle("POST", pth, handlers...)
//}
//
//func (g *Router) Handle(method, pth string, handlers ...gin.HandlerFunc) *Router {
//	g.g.Handle(method, pth, handlers...)
//	return &Router{
//		g:      g.g,
//		method: method,
//		path:   path.Join(g.path, pth),
//	}
//}

//func NewRouter(r gin.IRouter) *Router {
//	return &Router{
//		g: r,
//	}
//}
//
//func RegisterAPI[Req, Resp any](r *ApiGroup, router *Router, method, path string, handler Handler[Req, Resp], opts ...OptFunc) {
//	rsc := generateSchema(reflect.ValueOf(new(Req)), "")
//	a := &Api{
//		Request:  new(Req),
//		Response: new(Resp),
//
//		Method:         method,
//		RequestSchema:  rsc,
//		ResponseSchema: generateSchema(reflect.ValueOf(new(Resp)), ""),
//	}
//	for _, opt := range opts {
//		opt(a)
//	}
//	rsc.Description = a.Description
//
//	route := router.Handle(method, path, r.schemaHandler(a), WrapHandler[Req, Resp](handler, a.ErrHandler))
//	a.Route = route.Path()
//	r.apis = append(r.apis, a)
//}

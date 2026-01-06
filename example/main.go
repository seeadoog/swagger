package main

import (
	"github.com/gin-gonic/gin"
	"github.com/seeadoog/swagger"
)

type Child struct {
	Name  string  `json:"name" desc:"name of user"`
	Age   int     `json:"age"`
	Point float32 `json:"point" example:"1.1" binding:"min=0,max=5" default:"1.4" desc:"point of user"`
}
type createUserRequest struct {
	Username string         `json:"username" required:"true" example:"user01" binding:"required" desc:"user，username for creating a new user,user，username for creating a new user,user，username for creating a new user,user，username for creating a new user,user，username for creating a new user,user，username for creating a new user,user，username for creating a new user" `
	Password string         `json:"password" required:"true" example:"12345678" desc:"password for creating a new user"`
	Class    *string        `json:"class" required:"true" enum:"1,2,3,4,5,6" example:"userClass" default:"2"  desc:"class for creating a new user"`
	Children []*Child       `json:"children" required:"true" desc:"children for creating a new user"`
	Tags     map[string]any `json:"tags" required:"true" desc:"tags for creating a new user" example:"{}" location:"header"`
	Email    string         `json:"email" binding:"required,email" example:"user01@example.com" desc:"email of user" location:"query"`
}

func (c *createUserRequest) Validate(ctx swagger.ValidateCtx) swagger.ValidateFuncs {
	return swagger.ValidateFuncs{
		ctx.NotEmpty("username", c.Username),
		ctx.NotEmpty("password", c.Password),
		ctx.MinLength("password", c.Password, 8),
		ctx.StrIn("class", swagger.PtrVal(c.Class), "1", "2", "3"),
	}
}

type getUserRequest struct {
	Username string `json:"username" required:"true" example:"user01" location:"path,username"`
	Class    string `json:"class" required:"true" example:"class1" location:"query,class"`
	Level    int    `json:"level" required:"true" example:"1" location:"query,level"`
}
type Resp[T any] struct {
	Data T `json:"data"`
}

func main() {
	apiGroup := swagger.NewAPIGroup()
	ginEngine := gin.New()
	router := ginEngine

	userGroup := router.Group("/api/users")

	type createUserResponse struct {
		Code    int                `json:"code" desc:"response code , 0 indicates success"`
		Message string             `json:"err_msg" example:"" desc:"response error message"`
		Data    *createUserRequest `json:"data"`
	}
	swagger.RegisterAPI(apiGroup, userGroup, "POST", "", func(ctx *gin.Context, req *createUserRequest) *createUserResponse {
		return &createUserResponse{Data: req}
	},
		swagger.WithTitle("create_user"), swagger.WithDescription("create_user"))

	swagger.RegisterAPI(apiGroup, userGroup, "GET", ":username", func(ctx *gin.Context, req *getUserRequest) *createUserResponse { return &createUserResponse{} },
		swagger.WithTitle("get_user"), swagger.WithDescription("get_user"), swagger.WithUnExported())

	swagger.RegisterAPIWithDoc(apiGroup, userGroup, "POST", "/df", HandlerCreateUser, "create_user", "create user")

	a := apiHandler{}

	hd := &handlers{}

	//swagger.RegisterApiTemplate(apiGroup, userGroup, hd.ApiGetUser())

	apiGroup.RegisterAllApi(ginEngine, hd, func(funcName string) bool {
		return true
	})
	swagger.RegisterAPIWithDoc(apiGroup, userGroup, "GET", "/dics/:doc_name", a.HandlerDocumentHtml, "get doc", "get doc")
	ginEngine.GET("/apidoc.html", apiGroup.HandlerDocumentHtml())
	ginEngine.GET("/apidoc.md", apiGroup.HandlerDocumentMd())
	ginEngine.GET("/apischema", apiGroup.HandlerAllApiSchemas())
	ginEngine.Run(":8902")
}

func HandlerCreateUser(ctx *gin.Context, req *struct {
	Name  string `json:"name" binding:"required" example:"user01" desc:"name of user"`
	Age   int    `json:"age" binding:"required,max=100,min=18" example:"18" desc:"age of user"`
	Class int    `json:"class" binding:"required,max=20,min=1" example:"13" desc:"class for creating a new user,class"`
}) *Resp[any] {

	return &Resp[any]{Data: req}
}

type apiHandler struct {
}

type documentResponse struct {
	Req any
}

func (a *apiHandler) HandlerDocumentHtml(ctx *gin.Context, req *struct {
	DocName string `location:"path,doc_name" example:"doc01" desc:"doc name"`
	Ttl     string `location:"header,x-ttl" example:"100s" desc:"req ttl"`
}) *documentResponse {
	return &documentResponse{
		Req: req,
	}
}

type handlers struct {
}

func (h *handlers) ApiGetUser() any {
	type request struct {
		Username string `json:"username" example:"user01" location:"path,username"`
	}
	type response struct {
		Data *request `json:"data"`
	}

	return swagger.ApiTemplate[request, response]{
		Title:       "get user",
		Description: "get user",
		Method:      "GET",
		Path:        "/sapi/users/:username",
		Handler: func(ctx *gin.Context, req *request) *response {

			return &response{
				Data: req,
			}
		},
	}
}

func (h *handlers) CreateGetUser() any {
	type request struct {
		Username string `json:"username" example:"user01" location:"path,username"`
	}
	type response struct {
		Data *request `json:"data"`
	}

	return swagger.ApiTemplate[request, response]{
		Title:       "get user",
		Description: "get user",
		Method:      "POST",
		Path:        "/sapi/users",
		Handler: func(ctx *gin.Context, req *request) *response {

			return &response{
				Data: req,
			}
		},
	}
}

type request struct {
	Username *string `json:"username" example:"user01" location:"json,username" desc:"this use username"`
	Xttl     *int64  `location:"header,x-ttl" example:"100" desc:"req ttl"`
	IsUsed   *bool   `location:"query,is-used" example:"true" desc:"req is-used"`
}

type response struct {
	Data *request `json:"data" desc:"data"`
}

func (h *handlers) CreateGetUser3() swagger.ApiTemplate[request, response] {

	return swagger.ApiTemplate[request, response]{
		Title:       "get user",
		Description: "get user",
		Method:      "POST",
		Path:        "/bapi/users",
		Handler: func(ctx *gin.Context, req *request) *response {

			return &response{
				Data: req,
			}
		},
	}
}

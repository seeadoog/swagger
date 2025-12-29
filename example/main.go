package main

import (
	"github.com/gin-gonic/gin"
	"github.com/seeadoog/swagger"
)

type Child struct {
	Name string `json:"name" desc:"name of user"`
	Age  int    `json:"age"`
}
type createUserRequest struct {
	Username string         `json:"username" required:"true" example:"user01" binding:"required" desc:"user，username for creating a new user,user，username for creating a new user,user，username for creating a new user,user，username for creating a new user,user，username for creating a new user,user，username for creating a new user,user，username for creating a new user" `
	Password string         `json:"password" required:"true" example:"12345678" desc:"password for creating a new user"`
	Class    *string        `json:"class" required:"true" enum:"1,2,3,4,5,6" example:"userClass" default:"2"  desc:"class for creating a new user"`
	Children []*Child       `json:"children" required:"true" desc:"children for creating a new user"`
	Tags     map[string]any `json:"tags" required:"true" desc:"tags for creating a new user" example:"{}"`
	Email    string         `json:"email" binding:"required,email" example:"user01@example.com" desc:"email of user"`
}

func (c *createUserRequest) Validate(ctx swagger.ValidateCtx) swagger.ValidateFuncs {
	return swagger.ValidateFuncs{
		ctx.NotEmpty("username", c.Username),
		ctx.NotEmpty("password", c.Password),
		ctx.MinLength("password", c.Password, 8),
		ctx.StrIn("class", swagger.PtrVal(c.Class), "1", "2", "3"),
	}
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
		return &createUserResponse{
			Data: req,
		}
	}, swagger.WithTitle("create_user"), swagger.WithDescription("create_user"))

	type getUserRequest struct {
		Username string `json:"username" required:"true" example:"user01" location:"path,username"`
		Class    string `json:"class" required:"true" example:"class1" location:"query,class"`
		Level    int    `json:"level" required:"true" example:"1" location:"query,level"`
	}
	swagger.RegisterAPI(apiGroup, userGroup, "GET", ":username", func(ctx *gin.Context, req *getUserRequest) *createUserResponse {
		return &createUserResponse{}
	}, swagger.WithTitle("get_user"), swagger.WithDescription("get_user"))

	ginEngine.GET("/apidoc.html", apiGroup.HandlerDocumentHtml())
	ginEngine.GET("/apidoc.md", apiGroup.HandlerDocumentMd())

	ginEngine.Run(":8902")
}

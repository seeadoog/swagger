package swagger

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"reflect"
	"strconv"
	"testing"
)

type Class struct {
	Name string `json:"name" example:"6"`
	No   int    `json:"class" example:"1" default:"0"`
}
type request struct {
	Name  string  `json:"name" desc:"username" example:"xiaoming" required:"true"`
	Age   int     `json:"age" example:"18"`
	Ns    string  `json:"ns"  location:"path,ns" default:"nsx"`
	Key   string  `json:"key" location:"path,key" default:"xxxxx" example:"test_key"`
	File  *string `json:"file" location:"query,file"`
	Def   string  `json:"def" location:"query,def" default:"defs"`
	Mod   string  `json:"mod" location:"json,mod" enum:"read,write"`
	Bdf   *string `json:"bdf" location:"json,bdf" default:"bdfs"`
	Child struct {
		SName *string `json:"name" default:"xiaoming"`
	} `json:"child"`
	Class []Class  `json:"class"`
	Dirs  []string `json:"dirs" example:"dir,dir2"`
	Ids   []int    `json:"ids" example:"1,2,3"`
	Nes   string   `json:"nes" example:"-"`
	Path  string   `json:"path" location:"path,path" example:"abc/ttx"`
}

func HandlerReq(ctx *gin.Context, req *request) *response {
	return &response{
		Type:    "ok",
		Message: "success",
		Data:    req,
	}
}

func TestName(t *testing.T) {

	gine := gin.New()

	apiGroup := NewAPIGroup()
	gp := gine.Group("/api/users").Group("/dd/:ddd")
	RegisterAPI(apiGroup, gp, "GET", "hello", func(ctx *gin.Context, req *request) *response {

		return &response{}
	}, WithDescription("hello"), WithTitle(" set hello with id"))
	RegisterAPI(apiGroup, gine, "POST", "/hello/:key", HandlerReq, WithTitle("hello"), WithDescription("hello world"))
	RegisterAPI(apiGroup, gine, "POST", "/hello_33/:key/*path", HandlerReq, WithTitle("hello33"), WithDescription("hello world"))

	for i := 0; i < 50; i++ {
		RegisterAPI(apiGroup, gine, "POST", "/helloN"+strconv.Itoa(i), HandlerReq, WithTitle("hello"+strconv.Itoa(i)), WithDescription("hello world"))
	}
	gine.GET("/apidoc.md", apiGroup.HandlerDocumentMd())
	gine.GET("/apidoc.html", apiGroup.HandlerDocumentHtml())

	gine.Run(":8901")
}

func printJSON(b any) {
	bs, _ := json.MarshalIndent(b, "", "    ")
	fmt.Println(string(bs))
}
func printJSON2(b any) {
	bs, _ := json.Marshal(b)
	fmt.Println(string(bs))
}

func TestToken(t *testing.T) {
	fmt.Println(parsePathParams("/:name/*path/ase/:ns/:pp"))
}

type response struct {
	Type    string   `json:"type" example:"ok" binding:"required,email"`
	Message string   `json:"message" example:"success"`
	Data    *request `json:"data"`
	Age     int      `json:"age" example:"18" binding:"oneof=1 2"`
}

func TestVad(t *testing.T) {
	v := validator.New()
	v.SetTagName("binding")
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		return fld.Tag.Get("json")
	})
	err := v.Struct(&response{
		Type: "ok",
		Age:  3,
	})

	for _, e := range err.(validator.ValidationErrors) {

		fmt.Println(FormatValidatorError(e))
	}

}

package swagger

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"reflect"
	"strings"
)

type ApiTemplate[Req, Res any] struct {
	Title       string
	Description string
	Method      string
	Path        string
	Handler     Handler[Req, Res]
}

func RegisterApiTemplate[Req, Res any](ag *ApiGroup, gg BasicRouter, a ApiTemplate[Req, Res], opt ...OptFunc) {
	opt = append(opt, WithTitle(a.Title), WithDescription(a.Description))
	RegisterAPI(ag, gg, a.Method, a.Path, a.Handler, opt...)
}

func (a *ApiGroup) RegisterAllApi(gg BasicRouter, apiStruct any, isapi func(funcName string) bool, opts ...OptFunc) {

	v := reflect.ValueOf(apiStruct)
	//if v.Kind() == reflect.Ptr {
	//	a.RegisterAllApi(gg, v.Elem().Interface(), isapi, opts...)
	//	return
	//}
	num := 0
	for i := 0; i < v.NumMethod(); i++ {

		m := v.Method(i)
		mt := v.Type().Method(i)

		if isapi(mt.Name) {
			num++
			out := m.Call(nil)
			if len(out) != 1 {
				panic(fmt.Sprintf("%s.%s is not vaid api func, output args num should be 1", v.Type(), mt.Name))
			}

			aa := out[0]
			if out[0].Kind() == reflect.Interface {
				aa = aa.Elem()
			}

			if aa.Kind() != reflect.Struct {
				panic(fmt.Sprintf("%s.%s is not vaid api func,return type should be ApiTemplate", v.Type(), mt.Name))
			}

			if !strings.Contains(aa.Type().String(), "ApiTemplate") {
				panic(fmt.Sprintf("api template is not ApiTemplate"))
			}
			hd := aa.FieldByName("Handler")

			hdt := hd.Type()

			o := append([]OptFunc{
				WithTitle(aa.FieldByName("Title").String()),
				WithDescription(aa.FieldByName("Description").String()),
			}, opts...)
			inType := hdt.In(1).Elem()
			outType := hdt.Out(0).Elem()

			var errH ErrHandler

			aaa := a.RegisterGin(gg, reflect.New(inType).Interface(), reflect.New(outType).Interface(),
				aa.FieldByName("Method").String(), aa.FieldByName("Path").String(), func(ctx *gin.Context) {
					req := reflect.New(inType)

					err := bindRequest(a, ctx, req.Interface())
					if err != nil {
						if ctx.IsAborted() {
							return
						}
						errH(ctx, err)
						return
					}
					out := hd.Call([]reflect.Value{reflect.ValueOf(ctx), req})
					res := out[0].Interface()
					if ctx.IsAborted() {
						return
					}

					abortWithStatusJson(ctx, 200, res)
				}, o...)

			errH = aaa.ErrHandler
			if errH == nil {
				errH = defaultErrHandler
			}
		}
	}

	if num == 0 {
		panic(fmt.Sprintf("%s type has no handler to register", v.Type().String()))
	}

}

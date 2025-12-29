package swagger

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"strings"
)

type Handler[Req, Resp any] func(ctx *gin.Context, req *Req) *Resp

type ErrHandler func(ctx *gin.Context, err error)

func WrapHandler[Req, Resp any](a *ApiGroup, hd Handler[Req, Resp], errHandler ErrHandler) gin.HandlerFunc {
	if errHandler == nil {
		errHandler = func(ctx *gin.Context, err error) {
			ss := []string{}
			es, ok := err.(validator.ValidationErrors)
			if ok {
				for _, e := range es {
					ss = append(ss, FormatValidatorError(e))
				}
			}

			errmsg := ""
			if len(ss) > 0 {
				errmsg = strings.Join(ss, ",")
			} else {
				errmsg = err.Error()
			}
			ctx.AbortWithStatusJSON(400, gin.H{"error": errmsg})
		}
	}

	return func(ctx *gin.Context) {
		req := new(Req)

		err := bindRequest(a, ctx, req)
		if err != nil {
			if ctx.IsAborted() {
				return
			}
			errHandler(ctx, err)
			return
		}
		res := hd(ctx, req)
		if ctx.IsAborted() {
			return
		}
		ctx.AbortWithStatusJSON(200, res)
	}
}

type HandlerChain[Req, Resp any] struct {
	handlers []Handler[Req, Resp]
	idx      int
}

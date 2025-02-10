package binding

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Handle[IN any, OUT any](handler func(ctx context.Context, in IN) (OUT, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		var errorRenderFn = DefaultErrorRenderFunc()
		var in IN
		if err := c.BindJSON(in); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, errorRenderFn(err))
			return
		}

		out, err := handler(c, in)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, errorRenderFn(err))
		}
		c.JSON(http.StatusOK, out)
	}
}

type ErrorRenderFunc func(err error) any

var defaultErrorRender ErrorRenderFunc

func init() {
	defaultErrorRender = func(err error) any {
		return gin.H{"error": err.Error()}
	}
}

func SetDefaultErrorRenderer(f ErrorRenderFunc) {
	defaultErrorRender = f
}

func DefaultErrorRenderFunc() ErrorRenderFunc {
	return defaultErrorRender
}

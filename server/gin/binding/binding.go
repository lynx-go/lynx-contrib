package binding

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Handle[IN any, OUT any](handler func(ctx context.Context, in IN) (OUT, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		var errorRenderFn = DefaultErrorRenderFunc()
		var in IN
		if err := c.ShouldBind(in); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, errorRenderFn(err))
			return
		}

		out, err := handler(c, in)
		if err != nil {
			switch e := err.(type) {
			case *Error:
				if e.Status != nil {
					c.AbortWithStatusJSON(*e.Status, errorRenderFn(err))
				} else {
					c.AbortWithStatusJSON(http.StatusInternalServerError, errorRenderFn(err))
				}
			default:
				c.AbortWithStatusJSON(http.StatusInternalServerError, errorRenderFn(err))
			}
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

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Reason  string `json:"reason"`
	Status  *int   `json:"status"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("[Error %d] %s - %s", e.Code, e.Reason, e.Message)
}

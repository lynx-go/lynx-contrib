package handler

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Option struct {
	DefaultErrorStatus      int
	DefaultBadRequestStatus int
}

var defaultOption = DefaultOption()

func SetOption(o Option) {
	defaultOption = o
}

func DefaultOption() Option {
	return Option{
		DefaultErrorStatus:      http.StatusOK,
		DefaultBadRequestStatus: http.StatusOK,
	}
}

func Wrap[IN any, OUT any](handler func(ctx context.Context, in IN) (OUT, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		var errorRenderFn = DefaultErrorRenderFunc()
		var in IN
		if err := c.ShouldBind(in); err != nil {
			c.AbortWithStatusJSON(defaultOption.DefaultBadRequestStatus, errorRenderFn(err))
			return
		}

		out, err := handler(c, in)
		if err != nil {
			if e, ok := err.(WithStatus); ok {
				c.JSON(e.Status(), errorRenderFn(err))
				return
			} else {
				c.JSON(defaultOption.DefaultErrorStatus, errorRenderFn(err))
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
	Status  *int   `json:"-"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("[Error %d] %s - %s", e.Code, e.Reason, e.Message)
}

type WithStatus interface {
	Status() int
}

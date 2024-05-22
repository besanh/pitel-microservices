package v1

import (
	"github.com/gin-gonic/gin"
)

type Test struct{}

func NewTest(engine *gin.Engine) {
	handler := &Test{}
	Group := engine.Group("bss-message/v1/test")
	{
		Group.GET("", handler.Test)
	}
}

func (handler *Test) Test(c *gin.Context) {

}

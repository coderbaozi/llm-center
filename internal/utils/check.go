package utils

import (
	"github.com/cloudwego/hertz/pkg/app"
)

func CheckRequestParam(c *app.RequestContext, param interface{}) error {
	if err := c.BindAndValidate(&param); err != nil {
		return err
	}
	return nil
}

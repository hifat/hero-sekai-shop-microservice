package request

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type (
	httpContextService interface {
		Bind(data any) error
	}

	httpContext struct {
		Context   echo.Context
		validator *validator.Validate
	}
)

func NewHttpContext(ctx echo.Context) httpContextService {
	return &httpContext{
		Context:   ctx,
		validator: validator.New(),
	}
}

func (c *httpContext) Bind(data any) error {
	if err := c.Context.Bind(data); err != nil {
		return err
	}

	if err := c.validator.Struct(data); err != nil {
		return err
	}

	return nil
}

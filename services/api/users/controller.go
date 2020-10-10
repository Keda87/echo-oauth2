package users

import (
	"github.com/Keda87/echo-oauth2/models"
	"github.com/labstack/echo/v4"
	"net/http"
)

type userController struct {
	userService UserServiceInterface
}

func NewController(userService UserServiceInterface) *userController {
	return &userController{
		userService,
	}
}

func (ctrl *userController) HandleRegisterUser(c echo.Context) error {
	body := new(models.UserPayload)
	if err := c.Bind(body); err != nil {
		return err
	}

	if err := c.Validate(body); err != nil {
		return err
	}

	result, err := ctrl.userService.Register(c.Request().Context(), body)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, models.ResponseData{
		Data: result,
	})
}

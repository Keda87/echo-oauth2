package authentications

import (
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/labstack/echo/v4"
)

type authController struct {
	srv *server.Server
}

func NewController(server *server.Server) *authController {
	return &authController{
		srv: server,
	}
}

func (ctrl *authController) HandleObtainToken(c echo.Context) error {
	return ctrl.srv.HandleTokenRequest(c.Response().Writer, c.Request())
}

package internalapi

import (
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/commonapi"
	"github.com/labstack/echo"
)

func (o *InternalAPI) GetTokenVerify(c echo.Context) error {
	return commonapi.GetTokenVerify(c)
}

package commonapi

import (
	"net/http"

	"github.com/ONBUFF-IP-TOKEN/baseapp/auth/inno"
	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/config"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/auth"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/context"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/resultcode"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/model"
	"github.com/labstack/echo"
)

func PostWebAccountLogin(c echo.Context, accountWeb *context.AccountWeb) error {
	resp := new(base.BaseResponse)
	resp.Success()

	// 1. 소셜 정보 검증
	userID, email, err := auth.GetIAuth().SocialAuths[accountWeb.SocialType].VerifySocialKey(accountWeb.SocialKey)
	if err != nil || len(userID) == 0 || len(email) == 0 {
		log.Errorf("%v", err)
		resp.SetReturn(resultcode.Result_Auth_VerifySocial_Key)
		return c.JSON(http.StatusOK, resp)
	}

	payload := &context.Payload{
		LoginType: context.WebAccountLogin,
		InnoUID: inno.AESEncrypt(inno.MakeInnoID(userID, email),
			[]byte(config.GetInstance().Secret.Key),
			[]byte(config.GetInstance().Secret.Iv)),
	}

	reqAccountWeb := &context.ReqAccountWeb{
		InnoUID:    payload.InnoUID,
		SocialID:   userID,
		SocialType: accountWeb.SocialType,
	}

	// 2. 웹 로그인/가입
	resAccountWeb, err := model.GetDB().AuthWebAccounts(reqAccountWeb)
	if err != nil {
		log.Errorf("%v", err)
		resp.SetReturn(resultcode.Result_DBError)
		return c.JSON(http.StatusOK, resp)
	}
	payload.AUID = resAccountWeb.AUID
	resAccountWeb.InnoUID = payload.InnoUID

	// 3. Access, Refresh 토큰 생성
	if jwtInfoValue, err := auth.GetIAuth().MakeToken(payload); err != nil {
		log.Errorf("%v", err)
		resp.SetReturn(resultcode.Result_Auth_MakeTokenError)
		return c.JSON(http.StatusOK, resp)
	} else {
		resAccountWeb.JwtInfo = *jwtInfoValue
	}

	resp.Value = *resAccountWeb

	return c.JSON(http.StatusOK, resp)
}

func DelWebAccountLogout(c echo.Context) error {
	ctx := base.GetContext(c).(*context.InnoAuthContext)
	resp := new(base.BaseResponse)
	resp.Success()

	// Check if the token has expired
	if _, err := auth.GetIAuth().GetJwtInfo(ctx.Payload.LoginType, ctx.Payload.Uuid); err != nil {
		resp.SetReturn(resultcode.Result_Auth_ExpiredJwt)
	} else {
		// Delete the uuid in Redis.
		if err := auth.GetIAuth().DeleteUuidRedis(ctx.Payload.LoginType, ctx.Payload.Uuid); err != nil {
			resp.SetReturn(resultcode.Result_RedisError)
		}
	}

	return c.JSON(http.StatusOK, resp)
}

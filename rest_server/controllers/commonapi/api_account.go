package commonapi

import (
	"net/http"

	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/context"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/resultcode"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/model"
	"github.com/labstack/echo"
)

func PostAccountLogin(c echo.Context, reqAuthAccountApp *context.ReqAuthAccountApplication) error {
	resp := new(base.BaseResponse)
	resp.Success()

	// 1. 인증 프로시저 호출 (신규 유저, 기존 유저를 체크)
	ctx := base.GetContext(c).(*context.InnoAuthContext)
	if respAuth, err := model.GetDB().AccountAuthApplication(reqAuthAccountApp, ctx.Payload); err != nil {
		// 에러
		log.Error(err)
	} else {
		reqAuthAccountApp.Account.AUID = respAuth.AUID

		// 2. 신규/기존 유저에 따른 분기 처리
		// 신규 유저
		if respAuth.IsJoined == 1 {
			// 1. token-manager에 새 지갑 주소 생성 요청
			reqNewWallet := new(context.ReqNewWallet)
			reqNewWallet.Symbol = "ETH"
			reqNewWallet.NickName = reqAuthAccountApp.Account.SocialID
			if tokenInfo, err := GetTokenAddressNew(reqNewWallet); err != nil {
				resp.SetReturn(resultcode.Result_TokenManager_AddressNew)
			} else {
				resp.Value = tokenInfo
			}

			// 2. [DB] 지갑 생성 프로시저 호출

		} else { // 기존 유저
			// 1. token-manager 호출X -> point-manager에 포인트 수량 정보 요청
		}

		// 3. 액세스/리프레시 토큰 생성

		// 4. 같이 데이터를 담아서 액세스토큰과 함께 게임서버로 전달해줌.
	}

	return c.JSON(http.StatusOK, resp)
}

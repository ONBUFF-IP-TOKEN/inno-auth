package commonapi

import (
	"net/http"

	"github.com/ONBUFF-IP-TOKEN/baseapp/auth/inno"
	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/config"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/auth"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/commonapi/inner"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/context"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/resultcode"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/model"
	"github.com/labstack/echo"
)

// Web 계정 로그인/가입
func PostWebAccountLogin(c echo.Context, accountWeb *context.AccountWeb) error {
	resp := new(base.BaseResponse)
	resp.Success()
	conf := config.GetInstance()

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
			[]byte(conf.Secret.Key),
			[]byte(conf.Secret.Iv)),
	}

	reqAccountWeb := &context.ReqAccountWeb{
		InnoUID:    payload.InnoUID,
		SocialID:   userID,
		SocialType: accountWeb.SocialType,
	}

	// 2. 웹 로그인/가입
	resAccountWeb, err := model.GetDB().AuthAccounts(reqAccountWeb)
	if err != nil {
		log.Errorf("%v", err)
		resp.SetReturn(resultcode.Result_DBError)
		return c.JSON(http.StatusOK, resp)
	}
	payload.AUID = resAccountWeb.AUID
	resAccountWeb.InnoUID = payload.InnoUID

	// 3. ONIT 지갑이 없는 유저는 지갑을 생성
	if !resAccountWeb.ExistsMainWallet {
		// 3-1. [token-manager] ETH 지갑 생성
		var baseCoinList []context.CoinInfo
		for i, value := range conf.BaseCoin.Symbol {
			baseCoinList = append(baseCoinList, context.CoinInfo{
				CoinID:     conf.BaseCoin.ID[i],
				CoinSymbol: value,
			})
		}
		walletInfo, err := inner.TokenAddressNew(baseCoinList, payload.InnoUID)
		if err != nil {
			log.Errorf("%v", err)
			resp.SetReturn(resultcode.Result_Api_Get_Token_Address_New)
			return c.JSON(http.StatusOK, resp)
		}

		// 3-2. [DB] ETH 지갑 생성 프로시저 호출
		if err := model.GetDB().AddAccountBaseCoins(resAccountWeb.AUID, walletInfo); err != nil {
			log.Errorf("%v", err)
			resp.SetReturn(resultcode.Result_Procedure_Add_Base_Account_Coins)
			return c.JSON(http.StatusOK, resp)
		}

		// 3-3. [DB] ONIT 사용자 코인 등록
		if err := model.GetDB().AddAccountCoins(resAccountWeb.AUID, conf.OnitCoin.ID); err != nil {
			log.Errorf("%v", err)
			resp.SetReturn(resultcode.Result_Procedure_Add_Account_Coins)
			return c.JSON(http.StatusOK, resp)
		}
	}

	// 4. Access, Refresh 토큰 생성

	// 4-1. 기존에 발급된 토큰이 있는지 확인
	if oldJwtInfo, err := auth.GetIAuth().GetJwtInfoByInnoUID(payload.LoginType, context.AccessT, payload.InnoUID); err != nil || oldJwtInfo == nil {
		// 4-2. 기존에 발급된 토큰이 없다면 토큰을 발급한다. (Redis 확인)
		if jwtInfoValue, err := auth.GetIAuth().MakeWebToken(payload); err != nil {
			log.Errorf("%v", err)
			resp.SetReturn(resultcode.Result_Auth_MakeTokenError)
			return c.JSON(http.StatusOK, resp)
		} else {
			// 4-3. 새로 발급된 토큰으로 응답
			resAccountWeb.JwtInfo = *jwtInfoValue
		}
	} else {
		// 4-2. 기존 발급된 토큰으로 응답
		resAccountWeb.JwtInfo = *oldJwtInfo
	}

	resp.Value = *resAccountWeb

	return c.JSON(http.StatusOK, resp)
}

// Web 계정 로그아웃
func DelWebAccountLogout(c echo.Context) error {
	ctx := base.GetContext(c).(*context.InnoAuthContext)
	resp := new(base.BaseResponse)
	resp.Success()

	// Check if the token has expired
	if _, err := auth.GetIAuth().GetJwtInfoByInnoUID(ctx.Payload.LoginType, context.AccessT, ctx.Payload.InnoUID); err != nil {
		resp.SetReturn(resultcode.Result_Auth_ExpiredJwt)
	} else {
		// Delete the innoUID in Redis.
		if err := auth.GetIAuth().DeleteInnoUIDRedis(ctx.Payload.LoginType, context.AccessT, ctx.Payload.InnoUID); err != nil {
			resp.SetReturn(resultcode.Result_RedisError)
		}
	}

	return c.JSON(http.StatusOK, resp)
}

// Web 계정 로그인 정보 확인
func PostWebAccountInfo(c echo.Context, reqAccountInfo *context.ReqAccountInfo) error {
	resp := new(base.BaseResponse)
	resp.Success()

	if jwtInfo, err := auth.GetIAuth().GetJwtInfoByInnoUID(context.WebAccountLogin, context.AccessT, reqAccountInfo.InnoUID); err != nil {
		resp.SetReturn(resultcode.Result_Auth_Invalid_InnoUID)
	} else {
		if _, atClaims, err := auth.GetIAuth().VerifyAccessToken(jwtInfo.AccessToken); err != nil {
			resp.SetReturn(resultcode.Result_Auth_InvalidJwt)
		} else {
			resWebAccountInfo := &context.ResWebAccountInfo{
				JwtInfo: *jwtInfo,
				InnoUID: reqAccountInfo.InnoUID,
				AUID:    int64(atClaims["au_id"].(float64)),
			}
			resp.Value = *resWebAccountInfo
		}
	}

	return c.JSON(http.StatusOK, resp)
}

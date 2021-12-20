package context

import (
	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/baseutil/datetime"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/config"
)

type LoginType int

const (
	NoneLogin LoginType = iota
	CpLogin
	AppLogin
	AppAccountLogin
	WebAccountLogin
)

var LoginTypeText = map[LoginType]string{
	CpLogin:         "CP",
	AppLogin:        "APP",
	AppAccountLogin: "APPACCOUNT",
	WebAccountLogin: "WEBACCOUNT",
}

type Payload struct {
	CompanyID int       `json:"company_id,omitempty"`
	AppID     int       `json:"app_id,omitempty"`
	LoginType LoginType `json:"login_type,omitempty"`
	Uuid      string    `json:"uuid,omitempty"`
	IsEnabled bool      `json:"is_enabled,omitempty"`
}

// InnoAuthServerContext API의 Request Context
type InnoAuthContext struct {
	*base.BaseContext
	Payload *Payload
	JwtInfo *JwtInfo
}

// NewInnoAuthServerContext 새로운 InnoAuthServer Context 생성
func NewInnoAuthServerContext(baseCtx *base.BaseContext) interface{} {
	if baseCtx == nil {
		return nil
	}

	ctx := new(InnoAuthContext)
	ctx.BaseContext = baseCtx

	return ctx
}

// AppendRequestParameter BaseContext 이미 정의되어 있는 ReqeustParameters 배열에 등록
func AppendRequestParameter() {
}

func (o *InnoAuthContext) SetAuthContext(payload *Payload) {
	o.Payload = payload
}

func MakeDt(data *int64) {
	*data = datetime.GetTS2MilliSec()
}

func GetTokenExpiryperiod(loginType LoginType) (int64, int64) {
	confAuth := config.GetInstance().Auth
	switch loginType {
	case AppLogin:
		return confAuth.AppAccessTokenExpiryPeriod, confAuth.AppRefreshTokenExpiryPeriod
	case WebAccountLogin:
		return confAuth.WebAccessTokenExpiryPeriod, confAuth.WebRefreshTokenExpiryPeriod
	}
	return 0, 0
}

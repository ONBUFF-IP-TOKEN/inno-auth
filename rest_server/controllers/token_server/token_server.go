package token_server

import (
	"github.com/ONBUFF-IP-TOKEN/baseInnoClient/context"
	"github.com/ONBUFF-IP-TOKEN/baseInnoClient/token_manager"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/config"
)

var gTokenServer *token_manager.TokenServer

func GetInstance() *token_manager.TokenServer {
	return gTokenServer
}

func InitTokenManager(conf *config.ServerConfig) error {
	tokenServerInfo := &context.ServerInfo{
		HostInfo: context.HostInfo{
			IntHostUri: conf.TokenManager.InternalUri,
			ExtHostUri: conf.TokenManager.ExternalUri,
			IntVer:     conf.TokenManager.InternalVersion,
			ExtVer:     conf.TokenManager.ExternalVersion,
		},
	}

	gTokenServer = token_manager.NewTokenManagerServerInfo(tokenServerInfo)
	return nil
}

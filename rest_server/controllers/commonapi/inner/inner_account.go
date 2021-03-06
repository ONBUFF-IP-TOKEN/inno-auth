package inner

import (
	"github.com/ONBUFF-IP-TOKEN/baseapp/auth/inno"
	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/config"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/context"
)

func TokenAddressNew(coinList []context.CoinInfo, nickName string) ([]context.WalletInfo, error) {
	var addressList []context.WalletInfo

	for _, coin := range coinList {
		reqAddressNew := &context.ReqAddressNew{
			Symbol:   coin.CoinSymbol,
			NickName: nickName,
		}
		if resp, err := GetTokenAddressNew(reqAddressNew); err != nil {
			log.Errorf("%v", err)
			return nil, err
		} else {
			respAddressNew := &context.WalletInfo{
				CoinID:     coin.CoinID,
				Symbol:     coin.CoinSymbol,
				Address:    resp.Address,
				PrivateKey: resp.PrivateKey,
			}
			addressList = append(addressList, *respAddressNew)
		}
	}
	return addressList, nil
}

func DecryptInnoUID(innoUID string) string {
	return inno.AESDecrypt(innoUID, []byte(config.GetInstance().Secret.Key), []byte(config.GetInstance().Secret.Iv))
}

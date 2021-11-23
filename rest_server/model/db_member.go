package model

import (
	"fmt"

	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/context"
)

func (o *DB) InsertApp(appInfo *context.AppInfo) error {
	sqlQuery := fmt.Sprintf("INSERT INTO onbuff_inno.dbo.auth_app(app_name, cp_idx, login_id, login_pwd, create_dt) output inserted.idx "+
		"VALUES('%v', %v, '%v', '%v', %v)",
		appInfo.AppName, appInfo.CpIdx, appInfo.Account.LoginId, appInfo.Account.LoginPwd, appInfo.CreateDt)

	var lastInsertId int64
	err := o.Mssql.QueryRow(sqlQuery, &lastInsertId)

	if err != nil {
		log.Error(err)
		return err
	}

	log.Debug("InsertApp idx:", lastInsertId)

	return nil
}

func (o *DB) SelectApp(appInfo *context.AppInfo) (*context.AppInfo, error) {
	sqlQuery := fmt.Sprintf("SELECT * FROM onbuff_inno.dbo.auth_app WHERE app_name='%v'", appInfo.AppName)
	rows, err := o.Mssql.Query(sqlQuery)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	defer rows.Close()

	app := new(context.AppInfo)

	for rows.Next() {
		if err := rows.Scan(&app.Idx, &app.AppName, &app.CpIdx, &app.Account.LoginId, &app.Account.LoginPwd, &app.CreateDt); err != nil {
			log.Error(err)
			return nil, err
		}
	}
	return app, err
}

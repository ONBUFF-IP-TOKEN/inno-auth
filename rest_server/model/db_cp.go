package model

import (
	"fmt"

	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/context"
)

func (o *DB) InsertCP(company *context.Company) error {
	sqlQuery := fmt.Sprintf("INSERT INTO onbuff_inno.dbo.auth_cp(cp_name, create_dt) output inserted.idx "+
		"VALUES('%v', %v)",
		company.CompanyName, 0)

	var lastInsertId int64
	err := o.Mssql.QueryRow(sqlQuery, &lastInsertId)

	if err != nil {
		log.Error(err)
		return err
	}

	log.Debug("InsertCP idx:", lastInsertId)

	return nil
}

func (o *DB) SelectGetCpInfoByIdx(idx int) (*context.Company, error) {
	sqlQuery := fmt.Sprintf("SELECT * FROM onbuff_inno.dbo.auth_cp WHERE idx=%v", idx)
	rows, err := o.Mssql.Query(sqlQuery)

	if err != nil {
		log.Error(err)
		return nil, err
	}
	defer rows.Close()

	cp := context.NewCompany()

	var createDt int
	for rows.Next() {
		if err := rows.Scan(&cp.CompanyID, &cp.CompanyName, &createDt); err != nil {
			log.Error(err)
			return nil, err
		}
	}
	return cp, err
}

func (o *DB) SelectGetCpInfoByCpName(cpName string) (*context.Company, error) {
	sqlQuery := fmt.Sprintf("SELECT * FROM onbuff_inno.dbo.auth_cp WHERE cp_name='%v'", cpName)
	rows, err := o.Mssql.Query(sqlQuery)

	if err != nil {
		log.Error(err)
		return nil, err
	}
	defer rows.Close()

	cp := context.NewCompany()
	var create_dt int64

	for rows.Next() {
		if err := rows.Scan(&cp.CompanyID, &cp.CompanyName, &create_dt); err != nil {
			log.Error(err)
			return nil, err
		}
	}

	return cp, err
}

func (o *DB) DeleteCP(company *context.Company) error {
	sqlQuery := fmt.Sprintf("DELETE FROM onbuff_inno.dbo.auth_cp WHERE cp_name='%v'", company.CompanyName)
	rows, err := o.Mssql.Query(sqlQuery)
	if err != nil {
		log.Error(err)
		return err
	}
	defer rows.Close()

	return nil
}

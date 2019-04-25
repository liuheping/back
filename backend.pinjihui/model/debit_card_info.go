package model

import "pinjihui.com/backend.pinjihui/util"

type DebitCardInfo struct {
	CardHolder string
	BankName   string
	CardNumber string
	ProvinceID int32
	CityID     int32
	Branch     string
}

func NewDebitCardInfo(row string) (*DebitCardInfo, error) {
	info := &DebitCardInfo{}
	if err := util.Row2Struct(&row, info); err != nil {
		return nil, err
	}
	return info, nil
}

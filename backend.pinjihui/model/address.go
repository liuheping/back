package model

import "pinjihui.com/backend.pinjihui/util"

type Address struct {
	// ProvinceId int32
	// CityId     int32
	AreaId     int32
	RegionName *string
	Address    string
}

func NewAddress(row string) (*Address, error) {
	info := &Address{}
	if err := util.Row2Struct(&row, info); err != nil {
		return nil, err
	}
	return info, nil
}

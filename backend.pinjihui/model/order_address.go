package model

import "pinjihui.com/backend.pinjihui/util"

type OrderAddress struct {
	Consignee  string
	Zipcode    *string
	Mobile     string
	AreaID     int32   //`db:"area_id"`
	RegionName *string //`db:"region_name"`
	Address    string  //`db:"address"`
}

func NewOrderAddress(row string) (*OrderAddress, error) {
	info := &OrderAddress{}
	if err := util.Row2Struct(&row, info); err != nil {
		return nil, err
	}
	return info, nil
}

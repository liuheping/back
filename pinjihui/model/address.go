package model

import "pinjihui.com/pinjihui/util"

type Address struct {
    // 省份ID
    ProvinceId  int32    `row:"-"`
    // 城市ID
    CityId      int32    `row:"-"`
    // 区县ID
    AreaId      int32    `db:"area_id"`
    // 地区全名, 如 四川 成都市 高新区
    RegionName  *string  `db:"region_name"`
    // 区县后的详细地址, 如 xx街xx号
    Address     string   `valid:"required"`
}

func NewAddress(addrRow *string) (*Address, error) {
    addr := &Address{}
    if err := util.Row2Struct(addrRow, addr); err != nil {
        return nil, err
    }
    return addr, nil
}

type Location struct {
    Lat float64
    Lng float64
} 

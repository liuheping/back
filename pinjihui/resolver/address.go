
package resolver

import (
    "pinjihui.com/pinjihui/model"
    rp "pinjihui.com/pinjihui/repository"
)

type addressResolver struct {
    m *model.Address
}

func (r *addressResolver) ProvinceId() int32 {
    if r.m.ProvinceId == 0 {
        if err := rp.L("address").(*rp.AddressRepository).FillAddr(r.m); err != nil {
            panic(err)
        }
    }
    return r.m.ProvinceId
}

func (r *addressResolver) CityId() int32 {
    if r.m.CityId == 0 {
        if err := rp.L("address").(*rp.AddressRepository).FillAddr(r.m); err != nil {
            panic(err)
        }
    }
    return r.m.CityId
}

func (r *addressResolver) AreaId() int32 {
    return r.m.AreaId
}

func (r *addressResolver) RegionName() *string {
    return r.m.RegionName
}

func (r *addressResolver) Address() string {
    return r.m.Address
}

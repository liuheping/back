package resolver

import "pinjihui.com/backend.pinjihui/model"

type addressResolver struct {
	m *model.Address
}

// func (r *addressResolver) ProvinceId() int32 {
// 	return r.m.ProvinceId
// }

// func (r *addressResolver) CityId() int32 {
// 	return r.m.CityId
// }

func (r *addressResolver) AreaId() int32 {
	return r.m.AreaId
}

func (r *addressResolver) RegionName() *string {
	return r.m.RegionName
}

func (r *addressResolver) Address() string {
	return r.m.Address
}

package model

type ShippingMethod struct {
    ID                string
    Name              string
    EnableForPlatform bool `db:"enable_for_platform"`
}

var ShippingMethods = []*ShippingMethod{
    {"getBySelf", "上门取货(自提)", false},
    {"bigPackage", "大件物流", true},
}
var CodeMap = map[string]int32{"getBySelf": 0, "bigPackage": 1}

package model

type ShippingAddress struct {
    ID        string
    UserId    string   `db:"user_id"`
    Mobile    string   `valid:"numeric,required,length(11|11)~mobile invalid"`
    Consignee string   `valid:"required"`
    Address   *Address `db:"*address"`
    Zipcode   *string
    IsDefault bool     `db:"is_default"`
    //CreatedAt string   `db:"created_at"`
    //UpdatedAt string   `db:"updated_at"`
}

type ShippingAddressDB struct {
    Address string
    ShippingAddress
}

type OrderAddress struct {
    Consignee  string
    Zipcode    *string
    Mobile     string
    AreaID     int32   `db:"area_id"`
    RegionName *string `db:"region_name"`
    Address    string  `db:"address"`
}

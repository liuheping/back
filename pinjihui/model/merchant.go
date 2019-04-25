package model

type Merchant struct {
    ID              string   `db:"user_id"`
    CompanyName     string   `db:"company_name"`
    CompanyAddress  *Address
    DeliveryAddress *Address
    CompanyImage    *string  `db:"company_image"`
    Logo            *string
    Telephone       *string
    Type            string   `fi:"-"`
    Distance        *float64 `fi:"-"`
    SalesVolume     *int     `db:"sales_volume" fi:"-"`
    CommentQuantity *int     `db:"comment_quantity" fi:"-"`
    Lat             *float64
    Lng             *float64
}

type MerchantWithStock struct {
    MerchantDB
    ProductId string  `db:"product_id"`
    Stock     int32
    Price     float64 `db:"retail_price"`
    Distance  *float64
}

type MerchantDB struct {
    CompanyAddress  *string `db:"company_address"`
    DeliveryAddress *string `db:"delivery_address"`
    Merchant
}

var Platform = &Merchant{
    CompanyName: "成都拼机惠",
}

const ALLY = "ally"
const PROVIDER = "provider"

func (m *MerchantDB) ParseAddrs() error {
    var err error
    if (m.CompanyAddress != nil) {
        if m.Merchant.CompanyAddress, err = NewAddress(m.CompanyAddress); err != nil {
            return err
        }
    }
    if m.DeliveryAddress != nil {

        if m.Merchant.DeliveryAddress, err = NewAddress(m.DeliveryAddress); err != nil {
            return err
        }
    }
    return nil
}

func (m *Merchant) IsPlatform() bool {
    return m.Type == PROVIDER
}

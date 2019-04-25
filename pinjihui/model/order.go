package model

import (
    "pinjihui.com/pinjihui/util"
    "strconv"
    "sync"
)

type Order struct {
    ID     string
    UserID string `db:"user_id"`
    // 订单状态
    Status string
    // 子订单
    children []*Order `fi:"-"`
    // 订单收货地址
    Address *OrderAddress
    // 订单商家
    MerchantID *string `db:"merchant_id"`
    // 商品
    products []*OrderProduct `fi:"-"`
    // 配送方法
    ShippingID   *string `db:"shipping_id"`
    ShippingName *string `db:"shipping_name"`
    // 支付方式
    PayCode     string  `db:"pay_code"`
    PayName     string  `db:"pay_name"`
    InvPayee    *string `db:"inv_payee"`
    InvType     *string `db:"inv_type"`
    InvTaxpayer *string `db:"inv_taxpayer_id"`
    // 无货处理方法,暂不用
    //howOos *string
    // 订单金额
    Amount float64
    // 已支付金额
    MoneyPaid float64 `db:"money_paid"`
    // 应付款金额, 订单金额减去优惠金额
    OrderAmount float64 `db:"order_amount"`
    // 下单时间
    CreatedAt string `db:"created_at"`
    // 支付时间
    PayTime *string `db:"pay_time"`
    // 发货时间
    ShippingTime *string `db:"shipping_time"`
    UsedCoupon   *string `db:"used_coupon"`
    PostScript   *string `db:"postscript"`
    ParentID     *string `db:"parent_id"`
    AddressRow   string  `db:"address"`
    //优惠金额
    OfferAmount float64 `db:"offer_amount"`
    //供货商总成本, 即商品批发价总和
    ProviderIncome *float64 `db:"provider_income"`
    //加盟商提成
    AllyIncome *float64 `db:"ally_income"`

    muxForProducts sync.Mutex `fi:"-"`
    muxForChildren sync.Mutex `fi:"-"`
}

func (o *Order) ParseAddress() error {
    o.Address = &OrderAddress{}
    return util.Row2Struct(&o.AddressRow, o.Address)
}

type OrderProduct struct {
    ID            string
    OrderId       string   `db:"order_id"`
    ProductId     string   `db:"product_id"`
    ProductName   string   `db:"product_name"`
    ProductNumber int32    `db:"product_number"`
    ProductPrice  float64  `db:"product_price"`
    BatchPrice    float64  `db:"batch_price"`
    SecondPrice   float64  `db:"second_price"`
    ProductImage  string   `db:"product_image"`
    Spec1Name     *string  `db:"spec_1_name"`
    Spec2Name     *string  `db:"spec_2_name"`
    Spec1         *string  `db:"spec_1"`
    Spec2         *string  `db:"spec_2"`
    ShippingFee   *float64 `db:"shipping_fee"`
    CommentID     *string  `db:"comment_id" fi:"-"`
}

type WxPayParams struct {
    AppID      string
    TimeStamp  string
    NonceStr   string
    PackageStr string
    SignType   string
    PaySign    string
}

type ShippingInfo struct {
    ID             string
    Company        *string
    DeliveryNumber *string `db:"delivery_number"`
    Images         *string
}

const (
    OS_UNPAID    = "unpaid"
    OS_PAID      = "paid"
    OS_CANCELLED = "cancelled"
    OS_SHIPPED   = "shipped"
    OS_FINISH    = "finish"
)

func GetFeeWithFenUnit(v float64) string {
    return strconv.FormatInt(int64(v*100), 10)
}

func (o *Order) GetChildren(f func(string) ([]*Order, error)) ([]*Order, error) {
    o.muxForChildren.Lock()
    defer o.muxForChildren.Unlock()
    var err error
    if o.children == nil {
        o.children, err = f(o.ID)
    }
    return o.children, err
}

func (o *Order) GetProducts(f func(string) ([]*OrderProduct, error)) ([]*OrderProduct, error) {
    o.muxForProducts.Lock()
    defer o.muxForProducts.Unlock()
    var err error
    if o.products == nil {
        o.products, err = f(o.ID)
    }
    return o.products, err
}

func (s *ShippingInfo) ImageArr() []string {
    return util.ParseArray(s.Images)
}

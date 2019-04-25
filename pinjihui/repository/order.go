package repository

import (
    "github.com/jmoiron/sqlx"
    "github.com/op/go-logging"
    "golang.org/x/net/context"
    "pinjihui.com/pinjihui/model"
    "github.com/rs/xid"
    gc "pinjihui.com/pinjihui/context"
    "errors"
    "pinjihui.com/pinjihui/util"
    "fmt"
    "github.com/yearnfar/wxpay"
    "database/sql"
    //"gopkg.in/go-with/wxpay.v1"
    "strconv"
    "time"
    "qiniupkg.com/x/log.v7"
    "pinjihui.com/pinjihui/loader"
)

type OrderRepository struct {
    BaseRepository
}

func NewOrderRepository(db *sqlx.DB, log *logging.Logger) *OrderRepository {
    return &OrderRepository{BaseRepository{db: db, log: log}}
}

const (
    balancePaymentName = "余额支付"
    balancePaymentCode = "Balance"
)

func (o *OrderRepository) Create(ctx context.Context, input *model.OrderInput) ([]*model.Order, error) {
    gc.CheckAuth(ctx)

    if len(input.ItemIDs) == 0 {
        return nil, gc.ErrInvalidParam
    }
    parentOrder := model.Order{}
    parentOrder.ID = xid.New().String()
    parentOrder.UserID = *gc.CurrentUser(ctx)
    //设置订单状态
    parentOrder.Status = model.OS_UNPAID
    //支付方式
    parentOrder.PayCode = "WechartPay"
    parentOrder.PayName = "微信"
    //发票信息
    if input.InvOpt != nil {
        parentOrder.InvType = &input.InvType
        parentOrder.InvPayee = input.InvPayee
        parentOrder.InvTaxpayer = input.InvTaxpayer
    }

    //查出所有购物车项
    items, err := L("cart").(*CartRepository).FindByIDs(ctx, &input.ItemIDs)
    if err != nil {
        return nil, err
    }
    if len(items) != len(input.ItemIDs) {
        return nil, errors.New("invalid item ids")
    }
    //计算商品总价格
    parentOrder.Amount = calculateAmount(items)
    //是否使用优惠券
    if input.UsedCoupon != nil && 0 != len(*input.UsedCoupon) {
        coupons, err := L("coupon").(*CouponRepository).FindAvailableByIDs(ctx, *input.UsedCoupon)
        if err != nil {
            return nil, err
        }
        if len(coupons) != len(*input.UsedCoupon) {
            return nil, errors.New("invalid coupon ids")
        }
        //检查优惠券是否可用
        for _, coupon := range coupons {
            if coupon.LimitAmount != nil && parentOrder.Amount < *coupon.LimitAmount {
                return nil, fmt.Errorf("coupon is not available, limit order amount: %.2f", parentOrder.Amount)
            }
            parentOrder.MoneyPaid += coupon.Value
            parentOrder.OfferAmount += coupon.Value
        }
        parentOrder.UsedCoupon = util.EncodeArrayForPG(input.UsedCoupon)
    }
    parentOrder.OrderAmount = parentOrder.Amount - parentOrder.MoneyPaid
    //优惠金额已足够支付订单, 直接成为已支付状态
    if parentOrder.OrderAmount <= 0 {
        parentOrder.OrderAmount = 0
        parentOrder.Status = model.OS_PAID
        payTime := time.Now().Format(time.RFC3339)
        parentOrder.PayTime = &payTime

        //支付方式
        parentOrder.PayCode = balancePaymentCode
        parentOrder.PayName = balancePaymentName
    }
    //处理地址
    adds, err := L("address").(*AddressRepository).FindByID(ctx, input.ShippingAddrID)
    if err == gc.ErrNoRecord {
        return nil, errors.New("invalid address")
    }
    parentOrder.Address = GenerateOrderAddress(adds)
    //分组商品
    merchantMapItems := make(map[string][]*model.CartItem)
    for _, v := range items {
        merchantMapItems[v.MerchantID] = append(merchantMapItems[v.MerchantID], v)
    }
    //对每个商家都创建一个子订单
    if len(input.SmOpt) != len(merchantMapItems) {
        return nil, gc.ErrInvalidParam
    }
    //生成所有订单
    orders, err := o.GenerateOrders(ctx, merchantMapItems, &parentOrder, input.SmOpt)
    if err != nil {
        return nil, err
    }
    //保存订单
    tx := o.db.MustBegin()
    for _, order := range orders {
        sqls := util.NewSQLBuilder(order).InsertSQLBuild([]string{"CreatedAt", "AddressRow"})
        _, err = tx.NamedExec(sqls, order)
        if err != nil {
            tx.Rollback()
            return nil, err
        }
        if order.MerchantID != nil {
            err := o.saveOrderProduct(tx, ctx, order, merchantMapItems)
            if err != nil {
                o.log.Errorf("saveOrderProduct failed, error: %v", err)
                tx.Rollback()
                return nil, err
            }
        }
    }
    //设置优惠券为已用状态
    if input.UsedCoupon != nil && len(*input.UsedCoupon) != 0 {
        err := o.SetUsed(tx, *input.UsedCoupon)
        if err != nil {
            tx.Rollback()
            return nil, err
        }
    }
    //从购物车中删除商品
    err = o.DeleteCartItems(ctx, tx, items)
    if err != nil {
        return nil, err
    }
    if err = tx.Commit(); err != nil {
        tx.Rollback()
        return nil, err
    }
    return orders, nil
}
func (o *OrderRepository) SetUsed(tx *sqlx.Tx, ids [] string) error {
    couponSQL := `UPDATE user_coupons SET used=true WHERE id in (?)`
    query, args, err := sqlx.In(couponSQL, ids)
    if err != nil {
        return err
    }
    query = o.db.Rebind(query)
    res, err := tx.Exec(query, args...)
    if err != nil {
        return err
    }
    aff, err := res.RowsAffected()
    if err != nil {
        return err
    }
    if aff != int64(len(ids)) {
        o.log.Errorf("Set coupon as used failed, RowsAffected:%d not eq ids len. ids: %v", aff, ids)
        return errors.New("set coupon as used fail")
    }
    return nil
}

func (o *OrderRepository) DeleteCartItems(ctx context.Context, tx *sqlx.Tx, items []*model.CartItem) error {
    for _, v := range items {
        sqls := "DELETE FROM carts WHERE product_id=$1 AND merchant_id=$2 AND user_id=$3"
        _, err := tx.Exec(sqls, v.ProductID, v.MerchantID, *gc.CurrentUser(ctx))
        if err != nil {
            tx.Rollback()
            return err
        }
    }
    return nil
}

func (o *OrderRepository) GenerateOrders(ctx context.Context, merchantMapItems map[string][]*model.CartItem, parentOrder *model.Order, SmOpt []*model.ShppingMethodOption) ([]*model.Order, error) {
    //记录商家id是否重复
    uniqueMap := make(map[string]bool)
    //是否是单一订单
    isSingleOrder := len(merchantMapItems) == 1
    orders := make([]*model.Order, 0)
    orders = append(orders, parentOrder)
    for _, shippingMethodOption := range SmOpt {
        itemsInOneMT, ok := merchantMapItems[shippingMethodOption.MerchantID]
        if !ok {
            return nil, gc.ErrInvalidParam
        }
        if _, ok := uniqueMap[shippingMethodOption.MerchantID]; ok {
            return nil, errors.New("duplicated merchant id")
        }
        uniqueMap[shippingMethodOption.MerchantID] = true

        //处理加盟商订单
        var newOrder *model.Order
        if isSingleOrder {
            newOrder = parentOrder
        } else {
            newOrder = new(model.Order)
            *newOrder = *parentOrder
            newOrder.ID = xid.New().String()
        }
        //设置商家/消息/和配送方式
        newOrder.MerchantID = &shippingMethodOption.MerchantID
        newOrder.PostScript = shippingMethodOption.Message
        newOrder.ShippingID = &shippingMethodOption.ShippingMethod
        newOrder.ShippingName = &model.ShippingMethods[model.CodeMap[shippingMethodOption.ShippingMethod]].Name
        if !isSingleOrder {
            //计算子订单总价
            newOrder.Amount = calculateAmount(itemsInOneMT)
            //分摊已支付和优惠
            newOrder.MoneyPaid = newOrder.Amount / parentOrder.Amount * parentOrder.MoneyPaid
            newOrder.OfferAmount = newOrder.Amount / parentOrder.Amount * parentOrder.OfferAmount
            newOrder.OrderAmount = newOrder.Amount - newOrder.MoneyPaid
            newOrder.ParentID = &parentOrder.ID
            orders = append(orders, newOrder)
        } else {
            //如果是单一订单,设置父订单为nil
            newOrder.ParentID = nil
        }
    }
    return orders, nil
}

func (o *OrderRepository) saveOrderProduct(tx *sqlx.Tx, ctx context.Context, order *model.Order, merchantMapItems map[string][]*model.CartItem) error {
    //构建订单商品
    items := merchantMapItems[*order.MerchantID]
    order.ProviderIncome = new(float64)
    order.AllyIncome = new(float64)
    for _, item := range items {
        product, err := loader.LoadProduct(ctx, item.ProductID, item.MerchantID)
        if err != nil {
            return err
        }
        if product.IsSale == false {
            return errors.New("product isn't on sale")
        }
        images, err := L("product").(*ProductRepository).FindImagesById(item.ProductID)
        if err != nil {
            return err
        }
        var ProductImage string
        if len(images) > 0 {
            ProductImage = images[0].BigImage
        }
        op := &model.OrderProduct{
            ID:            xid.New().String(),
            OrderId:       order.ID,
            ProductId:     item.ProductID,
            ProductName:   product.Name,
            ProductNumber: item.ProductCount,
            ProductPrice:  item.Price,
            BatchPrice:    product.BatchPrice,
            SecondPrice:   product.SecondPrice,
            ProductImage:  ProductImage,
            Spec1Name:     product.Spec1Name,
            Spec2Name:     product.Spec2Name,
            Spec1:         product.Spec1,
            Spec2:         product.Spec2,
            ShippingFee:   product.ShippingFee,
        }
        *order.ProviderIncome += op.BatchPrice * float64(op.ProductNumber)
        *order.AllyIncome += float64(op.ProductNumber) * (op.ProductPrice - op.SecondPrice)
        sqls := util.NewSQLBuilder(op).InsertSQLBuild(nil)
        _, err = tx.NamedExec(sqls, op)
        if err != nil {
            return err
        }
        //如果是秒杀产品，减数量
        err = o.updateSpikeCount(tx, op.ProductId, *order.MerchantID, int(op.ProductNumber))
        if err != nil {
            log.Errorf("update Spike count failed: %v", err)
            return err
        }
    }
    //更新订单商户收入
    updateIncome := fmt.Sprintf(`UPDATE orders SET provider_income=:provider_income,ally_income=(SELECT CASE WHEN type = 'ally' THEN :ally_income ELSE GREATEST(%.2f, 0) END FROM users WHERE users.id=orders.merchant_id) WHERE id=:id`, *order.AllyIncome)
    _, err := tx.NamedExec(updateIncome, order)
    return err
}

func GenerateOrderAddress(sa *model.ShippingAddress) *model.OrderAddress {
    return &model.OrderAddress{sa.Consignee, sa.Zipcode, sa.Mobile, sa.Address.AreaId, sa.Address.RegionName, sa.Address.Address}
}

func calculateAmount(items []*model.CartItem) float64 {
    a := 0.0
    for _, v := range items {
        a += v.Price * float64(v.ProductCount)
    }
    return a
}

func (o *OrderRepository) FindByID(ctx context.Context, id string) (*model.Order, error) {
    gc.CheckAuth(ctx)
    order := model.Order{}
    sqls := util.NewSQLBuilder(&order).WhereRow("id=$1 AND user_id=$2").BuildQuery()
    row := o.db.QueryRowx(sqls, id, *gc.CurrentUser(ctx))
    err := row.StructScan(&order)
    if err == sql.ErrNoRows {
        return nil, gc.ErrNoRecord
    }
    if err != nil {
        return nil, err
    }
    return &order, nil
}

func (o *OrderRepository) SelectByID(id string) (*model.Order, error) {
    order := model.Order{}
    sqls := util.NewSQLBuilder(&order).WhereRow("id=$1").BuildQuery()
    row := o.db.QueryRowx(sqls, id)
    err := row.StructScan(&order)
    if err == sql.ErrNoRows {
        return nil, gc.ErrNoRecord
    }
    if err != nil {
        return nil, err
    }
    return &order, nil
}

func (o *OrderRepository) List(ctx context.Context, status *string, first int, after *string) ([]*model.Order, error) {
    gc.CheckAuth(ctx)
    orders := make([]*model.Order, 0)
    builder := util.NewSQLBuilder(orders).WhereRow("parent_id IS NULL AND user_id=$1")
    if status != nil {
        builder.Where("status", "=", *status)
    }
    if after != nil {
        builder.WhereRow(fmt.Sprintf("created_at<(SELECT created_at FROM orders WHERE id = '%s')", *after))
    }
    sqls := builder.OrderBy("created_at", "DESC").Limit(first, nil).BuildQuery()
    if err := o.db.Select(&orders, sqls, *gc.CurrentUser(ctx)); err != nil {
        return nil, err
    }
    return orders, nil
}

func (o *OrderRepository) FindChildren(parentID string) ([]*model.Order, error) {
    orders := make([]*model.Order, 0)
    sqls := util.NewSQLBuilder(orders).WhereRow("parent_id=$1").BuildQuery()
    err := o.db.Select(&orders, sqls, parentID)
    return orders, err
}

func (o *OrderRepository) Count(ctx context.Context, status *string) (int, error) {
    sqls := `SELECT COUNT(*) FROM orders WHERE parent_id is null AND user_id=$1`
    if status != nil {
        sqls += fmt.Sprintf(" AND status='%s'", *status)
    }
    var count int
    err := o.db.Get(&count, sqls, *gc.CurrentUser(ctx))
    if err != nil {
        return 0, err
    }
    return count, nil
}

func (o *OrderRepository) FindOrderProducts(orderID string) ([]*model.OrderProduct, error) {
    sqls := `SELECT op.*, c2.id AS comment_id FROM order_products op LEFT JOIN comments c2 ON op.order_id = c2.order_id and op.product_id = c2.product_id WHERE op.order_id=$1`
    items := make([]*model.OrderProduct, 0)
    err := o.db.Select(&items, sqls, orderID)
    return items, err
}

func (o *OrderRepository) FindOrderProduct(orderID, productID string) (*model.OrderProduct, error) {
    p := model.OrderProduct{}
    b := util.NewSQLBuilder(&p).WhereRowWithHolderPlace("product_id=$1 AND order_id=$2", productID, orderID)
    err := o.db.Get(&p, b.BuildQuery(), b.Args...)
    return &p, err
}

func (o *OrderRepository) FindShippingInfo(ctx context.Context, orderID string) (*model.ShippingInfo, error) {
    gc.CheckAuth(ctx)
    info := model.ShippingInfo{}
    sqls := util.NewSQLBuilder(&info).Table("shipping_info").WhereRow("order_id=$1").BuildQuery()
    row := o.db.QueryRowx(sqls, orderID)
    err := row.StructScan(&info)
    if err == sql.ErrNoRows {
        return nil, gc.ErrNoRecord
    }
    if err != nil {
        return nil, err
    }
    return &info, nil
}

const (
    JSAPI    = "JSAPI"
    signType = "MD5"
)

func (o *OrderRepository) WechartPayPrepare(ctx context.Context, order *model.Order) (*model.WxPayParams, error) {
    gc.CheckAuth(ctx)
    appID := ctx.Value("config").(*gc.Config).WechartAppID
    secret := ctx.Value("config").(*gc.Config).WechartMchKey
    mchID := ctx.Value("config").(*gc.Config).WechartMchid
    payConfig := wxpay.WxPayConfig{
        AppId:     appID,
        AppSecret: secret,
        MchId:     mchID,
        NotifyUrl: ctx.Value("config").(*gc.Config).WechartNotifyUrl,
        TradeType: JSAPI, // 支持 JSAPI，NATIVE，APP
    }

    // 统一下单
    openid := gc.User(ctx).Openid
    if openid == "" {
        userID := *gc.CurrentUser(ctx)
        wxUser, err := L("user").(*UserRepository).FindWechartByUserID(userID)
        if err == gc.ErrNoRecord {
            return nil, fmt.Errorf("can't find openid for user: %s", userID)
        }
        if err != nil {
            return nil, err
        }
        openid = wxUser.OpenId
    }
    params := map[string]string{
        "out_trade_no":     order.ID,                                   // 商户订单号
        "body":             "成都拼机惠-机械零配件",                              // 商品描述
        "total_fee":        model.GetFeeWithFenUnit(order.OrderAmount), // 100分 = 1元
        "openid":           openid,
        "spbill_create_ip": util.GetString(ctx.Value("requester_ip").(*string), ""),
    }
    ret, err := wxpay.UnifiedOrder(payConfig, params)

    if err != nil {
        o.log.Errorf("wxpay return %v", err)
    }

    if ret.ResultCode != "SUCCESS" {
        o.log.Errorf("wxpay UnifiedOrder return err: %+v", ret)
        return nil, fmt.Errorf("wxpay request failed: %s", ret.ErrCodeDes)
    }
    wxPayParams := map[string]string{
        "appId":     appID,
        "timeStamp": strconv.FormatInt(time.Now().Unix(), 10),
        "nonceStr":  util.GetNonceStr(32),
        "package":   fmt.Sprintf("prepay_id=%s", ret.PrepayId),
        "signType":  signType,
    }
    wxPayParams["paySign"] = util.WXMakeSign(wxPayParams, secret)
    return &model.WxPayParams{
        AppID:      wxPayParams["appId"],
        TimeStamp:  wxPayParams["timeStamp"],
        NonceStr:   wxPayParams["nonceStr"],
        PackageStr: wxPayParams["package"],
        SignType:   wxPayParams["signType"],
        PaySign:    wxPayParams["paySign"],
    }, nil
}

// 更新支付状态
func (o *OrderRepository) updateOrderPayStatus(tx *sqlx.Tx, order *model.Order) error {
    order.Status = model.OS_PAID
    query := `UPDATE orders SET status=:status, money_paid=money_paid+order_amount,pay_time=CURRENT_TIMESTAMP WHERE id=:id OR parent_id=:id`
    result, err := tx.NamedExec(query, order)
    if err != nil {
        o.log.Errorf("update order status failed when call NamedExec: %v", err)
        return err
    }
    aff, err := result.RowsAffected()
    if err != nil {
        o.log.Errorf("update order status failed when call RowsAffected(): %v", err)
        return err
    }
    if aff != 1 {
        log.Errorf("update order status failed, RowsAffected is not 1, but: %d, order id: %s", aff, order.ID)
        return errors.New("RowsAffected wrong")
    }
    log.Info("update order status success!")
    return nil
}
func (o *OrderRepository) OnOrderPaid(ctx context.Context, order *model.Order) (error) {
    tx := o.db.MustBegin()
    //更新订单状态, 已支付金额,付款时间
    err := o.updateOrderPayStatus(tx, order)
    if err != nil {
        return err
    }
    log.Info("updateOrderPayStatus Success!")
    //增加产品销量/减库存
    err = o.updateRelatedDataWhenOrderPaid(tx, ctx, order)
    if err != nil {
        tx.Rollback()
        return err
    }
    log.Info("updateRelatedDataWhenOrderPaid Success!")
    err = tx.Commit()
    if err != nil {
        log.Error("commit failed: %v", err)
    }
    log.Info("OnOrderPaid Success!")
    return err
}

func (o *OrderRepository) updateProductStatusWhenPaid(tx *sqlx.Tx, order *model.Order) error {
    orderProducts, err := o.FindOrderProducts(order.ID)
    if err != nil {
        o.log.Errorf("update order status failed when call FindOrderProducts(): %v", err)
        return err
    }
    if len(orderProducts) == 0 {
        o.log.Errorf("order has no product, that should not happen")
        return errors.New("order has no product")
    }
    //增加产品销量
    updateProductSV := `UPDATE products SET sales_volume = sales_volume + $1 WHERE id=$2`
    //增加商户产品销量，减少库存
    updatePMSS := `UPDATE rel_merchants_products SET stock = stock - $1, sales_volume = sales_volume+$1 WHERE product_id=$2 AND merchant_id=$3`
    for _, p := range orderProducts {
        _, err := tx.Exec(updateProductSV, p.ProductNumber, p.ProductId)
        if err != nil {
            log.Errorf("update products sales_volume failed: %v", err)
            return err
        }
        _, err = tx.Exec(updatePMSS, p.ProductNumber, p.ProductId, order.MerchantID)
        if err != nil {
            log.Errorf("update products stock failed: %v", err)
            return err
        }
    }
    return nil
}

func (o *OrderRepository) updateMerchantBalance(tx *sqlx.Tx, order *model.Order) error {
    //获取商户类型
    merchant, err := L("merchant").(*MerchantRepository).FindByID(*order.MerchantID)

    if err != nil {
        return err
    }
    income := &merchantIncome{id: *order.MerchantID, inoutType: "order"}
    incomes := []*merchantIncome{income}

    if merchant.Type == model.ALLY {
        //加盟商直接收入商品的总额
        income.income = order.Amount
    } else {
        income.income = *order.ProviderIncome
        //加盟商提成
        if *order.AllyIncome != 0 {
            queryAlly := `SELECT user_id FROM merchant_profiles join users u on merchant_profiles.user_id = u.id WHERE u.type='ally' AND (company_address).area_id=$1`
            order.ParseAddress()
            var allyId string
            err = tx.Get(&allyId, queryAlly, order.Address.AreaID)
            if err == nil {
                income = new(merchantIncome)
                income.id = allyId
                income.income = *order.AllyIncome
                income.inoutType = "commission"
                incomes = append(incomes, income)
            } else if err != sql.ErrNoRows {
                log.Errorf("get merchant by area_id failed, error: %v", err)
                return err
            }
        }
    }
    for _, income := range incomes {
        //增加商户余额
        updateBalance := `UPDATE merchant_profiles SET balance=balance+$1 WHERE user_id=$2`
        _, err = tx.Exec(updateBalance, income.income, income.id)
        if err != nil {
            log.Errorf("update merchant balance failed, error: %v", err)
            return err
        }
        //记录流水
        insertLog := `INSERT INTO merchant_balance_logs (merchant_id,inout,inout_type,"references") VALUES($1,$2,$3,$4)`
        _, err = tx.Exec(insertLog, income.id, income.income, income.inoutType, order.ID)
        if err != nil {
            log.Errorf("save merchant balance log failed, error: %v", err)
            return err
        }
    }
    return nil
}

func (o *OrderRepository) updateRelatedDataWhenOrderPaid(tx *sqlx.Tx, ctx context.Context, pOrder *model.Order) error {
    //找到子订单
    orders, err := o.FindChildren(pOrder.ID)
    if err != nil {
        return err
    }
    if len(orders) == 0 {
        orders = append(orders, pOrder)
    }
    for _, order := range orders {
        err = o.updateProductStatusWhenPaid(tx, order)
        if err != nil {
            return err
        }
        //todo 用户确认收货才能更新商户余额.
        err = o.updateMerchantBalance(tx, order)
        if err != nil {
            return err
        }
    }
    return nil
}

type merchantIncome struct {
    id        string
    income    float64
    inoutType string
}

func (o *OrderRepository) updateSpikeCount(tx *sqlx.Tx, productId, merchantId string, productQuantity int) error {
    canSpike, err := L("spike").(*SpikeRepository).CanSpike(tx, productId, merchantId, productQuantity)
    if err != nil {
        return err
    }
    if canSpike {
        updateSpikeCount := `UPDATE spikes SET total_count = total_count - $1 WHERE product_id=$2 AND merchant_id=$3 AND start_at<current_timestamp AND expired_at>current_timestamp AND total_count >= $1`
        _, err = tx.Exec(updateSpikeCount, productQuantity, productId, merchantId)
    }
    return err
}

func (o *OrderRepository) UpdateAllyStock(tx *sqlx.Tx, user *model.User, products []*model.OrderProduct) error {
    for _, p := range products {
        query := `INSERT INTO rel_merchants_products AS r (product_id, merchant_id, stock, retail_price, origin_price)
VALUES ($1, $2, $3, $4, $4 * (SELECT 1+value::float FROM configs WHERE code = 'retail_price_ratio' LIMIT 1))
ON CONFLICT (merchant_id,product_id) DO UPDATE SET stock=r.stock+EXCLUDED.stock`
        _, err := tx.Exec(query, p.ProductId, user.ID, p.ProductNumber, p.BatchPrice)
        if err != nil {
            return err
        }
    }
    return nil
}

func (o *OrderRepository) MyOrderProducts(ctx context.Context, first int, offset int, hasCommented bool) ([]*model.OrderProduct, error) {
    gc.CheckAuth(ctx)
    orderProducts := make([]*model.OrderProduct, 0)
    notOrNone := util.If(hasCommented, "NOT", "").(string)
    sqls := `SELECT op.*, c2.id AS comment_id FROM order_products op
  join orders o on op.order_id = o.id
  LEFT JOIN comments c2 on op.order_id = c2.order_id and op.product_id = c2.product_id
where c2.id IS ` + notOrNone + ` NULL and o.user_id = $1 AND o.status = $2 ORDER BY o.created_at DESC LIMIT $3 OFFSET $4`
    if err := o.db.Select(&orderProducts, sqls, *gc.CurrentUser(ctx), model.OS_FINISH, first, offset); err != nil {
        return nil, err
    }
    return orderProducts, nil
}

func (o *OrderRepository) MyOrderProductsCount(ctx context.Context, hasCommented bool) (int, error) {
    notOrNone := util.If(hasCommented, "NOT", "").(string)
    sqls := `SELECT count(*) FROM order_products op
  join orders o on op.order_id = o.id
  LEFT JOIN comments c2 on op.order_id = c2.order_id and op.product_id = c2.product_id
where c2.id IS ` + notOrNone + ` NULL and o.user_id = $1 AND o.status = $2`
    return o.GetIntFormDB(sqls, *gc.CurrentUser(ctx), model.OS_FINISH)
}

func (o *OrderRepository) Cancel(ctx context.Context, id string) error {
    tx := o.db.MustBegin()
    err := o.UserUpdateOrderStatus(tx, ctx, id, model.OS_UNPAID, model.OS_CANCELLED)
    if err == nil {
        //复原秒杀商品
        query := `update spikes
set total_count = total_count + order_products.product_number
from
  orders, order_products
where orders.id = order_products.order_id and orders.merchant_id = spikes.merchant_id and
      order_products.product_id = spikes.product_id
      and orders.id = $1 AND
      spikes.start_at < current_timestamp and expired_at > current_timestamp`
        _, err = tx.Exec(query, id)
        if err != nil {
            tx.Rollback()
            return err
        }
        query = `UPDATE user_coupons SET used=FALSE WHERE id IN (SELECT UNNEST(used_coupon) FROM orders WHERE id=$1)`
        _, err = tx.Exec(query, id)
        if err != nil {
            tx.Rollback()
            return err
        }
    }
    tx.Commit()
    return err
}

func (o *OrderRepository) ConfirmReceipt(ctx context.Context, id string) error {
    tx := o.db.MustBegin()
    sqls := `UPDATE orders SET status=$1 WHERE id=$2 AND status=$3 AND user_id=$4`
    result, err := tx.Exec(sqls, model.OS_FINISH, id, model.OS_SHIPPED, *gc.CurrentUser(ctx))
    err = ExceReturn(result, err)
    if err != nil {
        log.Errorf("ConfirmReceipt failed, error: %v", err)
        return err
    }
    //获取买家类型
    user, err := L("user").(*UserRepository).FindByID(*gc.CurrentUser(ctx))
    if err != nil {
        tx.Rollback()
        log.Errorf("ConfirmReceipt failed, error: %v", err)
        return err
    }
    //加盟商进货
    if user.Type == model.ALLY {
        orderProducts, err := o.FindOrderProducts(id)
        if err != nil {
            tx.Rollback()
            log.Errorf("ConfirmReceipt failed, error: %v", err)
            return err
        }
        if len(orderProducts) == 0 {
            //不该出现此情况，可以确认收货的订单只能是有商品的订单
            log.Errorf("invalid order, no product!")
            tx.Rollback()
            return errors.New("invalid order")
        }
        err = o.UpdateAllyStock(tx, user, orderProducts)
        if err != nil {
            log.Errorf("UpdateAllyStock failed, error: %v", err)
            tx.Rollback()
            return err
        }
    }
    err = tx.Commit()
    return err
}

func (o *OrderRepository) UserUpdateOrderStatus(tx *sqlx.Tx, ctx context.Context, id, from, to string) error {
    gc.CheckAuth(ctx)
    sqls := `UPDATE orders SET status=$1 WHERE id=$2 AND status=$3 AND user_id=$4`
    result, err := tx.Exec(sqls, to, id, from, *gc.CurrentUser(ctx))
    return ExceReturn(result, err)
}

func ExceReturn(result sql.Result, err error) error {
    if err != nil {
        return err
    }
    n, err := result.RowsAffected()
    if err != nil {
        return err
    }
    if n < 1 {
        return errors.New("no row affected")
    }
    return nil
}

func (o *OrderRepository) UpdateOrderStatusBySchedule() {
    //取消订单时将秒杀商品数量复原, 并且归还优惠券
    query := `update spikes
set total_count = total_count + order_products.product_number
from
  orders, order_products
where orders.id = order_products.order_id and orders.merchant_id = spikes.merchant_id and
      order_products.product_id = spikes.product_id
      and orders.status = 'unpaid' AND
      orders.created_at < (current_timestamp - interval '1 day') and
      spikes.start_at < current_timestamp and expired_at > current_timestamp;
      UPDATE user_coupons SET used=FALSE WHERE id IN (SELECT UNNEST(used_coupon) FROM orders WHERE orders.status = 'unpaid' AND
      orders.created_at < (current_timestamp - interval '1 day'));`
    query += fmt.Sprintf(`UPDATE orders SET status='%s' WHERE status='%s' AND created_at<(current_timestamp-interval '1 day');`, model.OS_CANCELLED, model.OS_UNPAID)
    query += fmt.Sprintf(`UPDATE orders SET status='%s' WHERE status='%s' AND created_at<(current_timestamp-interval '7 days') AND user_id IN (SELECT id FROM users WHERE type='consumer');`, model.OS_FINISH, model.OS_SHIPPED)
    _, err := o.db.Exec(query)
    if err != nil {
        fmt.Println(err)
    }
}

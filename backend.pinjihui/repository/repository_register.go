package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/op/go-logging"
)

func RegisterAll(db *sqlx.DB, log *logging.Logger) {
	Admin.Register("address", NewAddressRepository(db, log))
	Admin.Register("region", NewRegionRepository(db, log))
	Admin.Register("brand", NewBrandRepository(db, log))
	Admin.Register("category", NewCategoryRepository(db, log))
	Admin.Register("attribute", NewAttributeRepository(db, log))
	Admin.Register("attributeset", NewAttributeSetRepository(db, log))
	Admin.Register("stock", NewStockRepository(db, log))
	Admin.Register("product", NewProductRepository(db, log))
	Admin.Register("comment", NewCommentRepository(db, log))
	Admin.Register("payment", NewPaymentRepository(db, log))
	Admin.Register("config", NewConfigRepository(db, log))
	Admin.Register("wechartprofile", NewWechartProfileRepository(db, log))
	Admin.Register("coupon", NewCouponRepository(db, log))
	Admin.Register("productinorder", NewProductInOrderRepository(db, log))
	Admin.Register("order", NewOrderRepository(db, log))
	Admin.Register("cashrequest", NewCashRequestRepository(db, log))
	Admin.Register("spike", NewSpikeRepository(db, log))
	Admin.Register("public", NewPublicRepository(db, log))
	Admin.Register("operationlog", NewOperationLogRepository(db, log))
	Admin.Register("ad", NewAdRepository(db, log))
	Admin.Register("brandseries", NewBrandSeriesRepository(db, log))
	Admin.Register("waiter", NewWaiterRepository(db, log))
	Admin.Register("shippinginfo", NewShippingInfoRepository(db, log))
}

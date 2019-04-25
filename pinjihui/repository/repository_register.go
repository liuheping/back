package repository

import (
    "github.com/jmoiron/sqlx"
    "github.com/op/go-logging"
)

func RegisterAll(db *sqlx.DB, log *logging.Logger) {
    Admin.Register("user", NewUserRepository(db, log))
    Admin.Register("address", NewAddressRepository(db, log))
    Admin.Register("region", NewRegionRepository(db, log))
    Admin.Register("brand", NewBrandRepository(db, log))
    Admin.Register("category", NewCategoryRepository(db, log))
    Admin.Register("product", NewProductRepository(db, log))
    Admin.Register("merchant", NewMerchantRepository(db, log))
    Admin.Register("cart", NewCartRepository(db, log))
    Admin.Register("order", NewOrderRepository(db, log))
    Admin.Register("attr", NewAttributeRepository(db, log))
    Admin.Register("shipping", NewShippingMethodRepository(db, log))
    Admin.Register("coupon", NewCouponRepository(db, log))
    Admin.Register("comment", NewCommentRepository(db, log))
    Admin.Register("spike", NewSpikeRepository(db, log))
    Admin.Register("favorite", NewFavoriteRepository(db, log))
    Admin.Register("ad", NewADRepository(db, log))
}

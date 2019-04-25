package model

type MerchantProfile struct {
	UserId             string   `db:"user_id"`
	RepName            *string  `db:"rep_name"`
	SocialId           *string  `db:"social_id"`
	CompanyName        *string  `db:"company_name"`
	LicenseImage       *string  `db:"license_image"`
	CompanyImage       *string  `db:"company_image"`
	CompanyAddress     *Address `db:"company_address"`
	DeliveryAddress    *Address `db:"delivery_address"`
	CompanyAddressRow  string   `db:"company_address"`
	DeliveryAddressRow string   `db:"delivery_address"`
	Created_at         string
	Updated_at         string
	Lat                *string
	Lng                *string
	Balance            float64
	Logo               *string
	Telephone          *string
	// Waiters            *string
}

type MerchantProfileARR struct {
	MerchantProfile
	CompanyImage *[]string
	// Waiters      *[]string
}

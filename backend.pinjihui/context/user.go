package context

import "golang.org/x/net/context"

func CurrentUser(ctx context.Context) *string {
	return ctx.Value("user_id").(*string)
}

func QueryFileds(ctx context.Context) []string {
	return ctx.Value("query_filed").([]string)
}

func CheckAuth(ctx context.Context) {
	if !ctx.Value("is_authorized").(bool) {
		panic(CredentialsError)
	}
}

func WhereMerchant(ctx context.Context) string {
	var str string
	str = ` AND merchant_id = '` + *CurrentUser(ctx) + `'`
	return str
}

func WhereMerchantNULL() string {
	var str string
	str = ` AND merchant_id IS NULL`
	return str
}

func WhereMerchantOR(ctx context.Context) string {
	var str string
	str = ` AND (merchant_id = '` + *CurrentUser(ctx) + `' OR merchant_id IS NULL)`
	return str
}

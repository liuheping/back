package resolver

type userProfileResolver struct {
	result interface{}
}

func (r *userProfileResolver) ToMerchantProfile() (*merchantProfileResolver, bool) {
	res, ok := r.result.(*merchantProfileResolver)
	return res, ok
}

func (r *userProfileResolver) ToCustomerProfile() (*customerProfileResolver, bool) {
	res, ok := r.result.(*customerProfileResolver)
	return res, ok
}

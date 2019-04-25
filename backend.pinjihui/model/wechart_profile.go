package model

type WechartProfile struct {
	Openid      string
	Session_key string
	Nick_name   string
	Gender      int32
	Language    *string
	City        *string
	Province    *string
	Country     *string
	Avatar_url  *string
	User_id     string
}

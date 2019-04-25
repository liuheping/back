package model

type ProductInOrder struct {
	ID             string
	Order_id       string
	Product_id     string
	Product_image  string
	Product_number int32
	Product_price  float64
	Product_name   string
	Spec_1_name    *string
	Spec_2_name    *string
	Spec_1         *string
	Spec_2         *string
	Shipping_fee   *float64
	Batch_price    float64
	Second_price   float64
	Agent_id       *string
}

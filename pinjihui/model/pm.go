package model

type PM struct {
    ProductID  string
    MerchantID string
}

func (pm *PM) String() string {
    return pm.MerchantID + "," + pm.ProductID
}

func (pm *PM) Raw() interface{} {
    return pm
}

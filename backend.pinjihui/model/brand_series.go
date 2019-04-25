package model

type BrandSeries struct {
	ID            string
	Brand_id      string
	Series        string
	Image         *string
	Machine_types string
	Sort_order    *int32
}

type BrandSeriesARR struct {
	BrandSeries
	Machine_types []string
}

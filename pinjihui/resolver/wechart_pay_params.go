package resolver

import "pinjihui.com/pinjihui/model"

type wechartPayParamsResolver struct {
    m *model.WxPayParams
}

func (w *wechartPayParamsResolver) TimeStamp() string {
    return w.m.TimeStamp
}

func (w *wechartPayParamsResolver) NonceStr() string {
    return w.m.NonceStr
}
func (w *wechartPayParamsResolver) Package() string {
    return w.m.PackageStr
}
func (w *wechartPayParamsResolver) SignType() string {
    return w.m.SignType
}
func (w *wechartPayParamsResolver) PaySign() string {
    return w.m.PaySign
}

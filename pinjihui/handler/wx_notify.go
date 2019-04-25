package handler

import (
    "net/http"
    "github.com/yearnfar/wxpay"
    gc "pinjihui.com/pinjihui/context"
    "io/ioutil"
    "encoding/xml"
    "github.com/op/go-logging"
    "pinjihui.com/pinjihui/repository"
    "pinjihui.com/pinjihui/model"
    "fmt"
)

type notifyResponse struct {
    XMLName    xml.Name `xml:"xml"`
    ReturnCode string   `xml:"return_code"`
    ReturnMsg  string   `xml:"return_msg"`
}

const (
    FAIL    = "FAIL"
    SUCCESS = "SUCCESS"
)

func WxNotify() http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/xml")
        secret := r.Context().Value("config").(*gc.Config).WechartMchKey
        config := wxpay.WxPayConfig{
            AppSecret: secret,
        }
        body, err := ioutil.ReadAll(r.Body)
        log := r.Context().Value("log").(*logging.Logger)
        if err != nil {
            log.Errorf("wx notify err: read body failed")
            writeFailToWX(w, err)
            return
        }
        log.Info("wx notify content: %s", body)
        resp, err := wxpay.Notify(config, body)
        if err != nil {
            log.Errorf("wx notify err: %v", err)
            writeFailToWX(w, err)
            return
        }
        orderRep := repository.L("order").(*repository.OrderRepository)
        order, err := orderRep.SelectByID(resp.OutTradeNo)
        if err != nil {
            log.Errorf("wx notify err: db err when try to find order: %s, error info: %v", resp.OutTradeNo, err)
            writeFailToWX(w, err)
            return
        }
        //如果订单不是未支付状态, 从业务逻辑上看说明此次支付已经通知过了.
        if order.Status != model.OS_UNPAID {
            //直接返回成功
            log.Warning("wx notify warning: notify a not-unpaid order")
            writeSuccessToWX(w)
            return
        }
        if model.GetFeeWithFenUnit(order.OrderAmount) != resp.TotalFee {
            err := fmt.Errorf("order amount: %.2f yuan != %s fen total fee", order.OrderAmount, resp.TotalFee)
            log.Errorf("wx notify err: order amount not match response total fee: %v", err)
            writeFailToWX(w, err)
            return
        }
        err = orderRep.OnOrderPaid(r.Context(), order)
        if err != nil {
            log.Errorf("OnOrderPaid failed: %v", err)
            writeFailToWX(w, err)
            return
        }
        writeSuccessToWX(w)
        log.Info("notify success for order: %s", order.ID)
    })
}

func writeFailToWX(w http.ResponseWriter, err error) {
    response := &notifyResponse{
        ReturnCode: FAIL,
        ReturnMsg:  err.Error(),
    }
    writeToWXResponse(w, response, http.StatusInternalServerError)
}

func writeSuccessToWX(w http.ResponseWriter) {
    response := &notifyResponse{
        ReturnCode: SUCCESS,
        ReturnMsg:  "OK",
    }
    writeToWXResponse(w, response, http.StatusOK)
}

func writeToWXResponse(w http.ResponseWriter, response interface{}, code int) {
    xmlResponse, _ := xml.Marshal(response)
    w.WriteHeader(code)
    w.Write(xmlResponse)
}

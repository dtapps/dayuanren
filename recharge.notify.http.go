package dayuanren

import (
	"go.dtapp.net/gojson"
	"net/http"
)

type ResponseRechargeNotifyHttp struct {
	Userid       string `json:"userid"`        // 商户ID
	OrderNumber  string `json:"order_number"`  // true
	OutTradeNum  string `json:"out_trade_num"` // 商户订单号
	Otime        string `json:"otime"`         // 成功/失败时间，10位时间戳
	State        string `json:"state"`         // 充值状态；-1取消， 0充值中， 1充值成功 ，2充值失败，3部分成功（-1,2做失败处理；1做成功处理；3做部分成功处理）
	Mobile       string `json:"mobile"`        // 充值手机号
	Remark       string `json:"remark"`        // 备注信息
	ChargeAmount string `json:"charge_amount"` // 充值成功面额
	Voucher      string `json:"voucher"`       // 凭证
	ChargeKami   string `json:"charge_kami"`   // 卡密/流水号
	Sign         string `json:"sign"`          // 签名字符串，用于验签,以保证回调可靠性。
}

// RechargeNotifyHttp 充值结果通知-异步通知
// https://www.kancloud.cn/boyanyun/boyanyun_huafei/3097255
func (c *Client) RechargeNotifyHttp(w http.ResponseWriter, r *http.Request) (validateJson ResponseRechargeNotifyHttp, err error) {
	if r.Method == http.MethodPost {
		err = gojson.NewDecoder(r.Body).Decode(&validateJson)
	} else if r.Method == http.MethodGet {
		validateJson.Userid = r.URL.Query().Get("userid")
		validateJson.OrderNumber = r.URL.Query().Get("order_number")
		validateJson.OutTradeNum = r.URL.Query().Get("out_trade_num")
		validateJson.Otime = r.URL.Query().Get("otime")
		validateJson.State = r.URL.Query().Get("state")
		validateJson.Mobile = r.URL.Query().Get("mobile")
		validateJson.Remark = r.URL.Query().Get("remark")
		validateJson.ChargeAmount = r.URL.Query().Get("charge_amount")
		validateJson.Voucher = r.URL.Query().Get("voucher")
		validateJson.ChargeKami = r.URL.Query().Get("charge_kami")
		validateJson.Sign = r.URL.Query().Get("sign")
	}
	return validateJson, err
}

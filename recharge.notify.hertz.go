package dayuanren

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
)

type ResponseRechargeNotifyHertz struct {
	Userid       int64   `form:"userid" json:"userid" query:"userid"`                      // 商户ID
	OrderNumber  string  `form:"order_number" json:"order_number" query:"order_number"`    // true
	OutTradeNum  string  `form:"out_trade_num" json:"out_trade_num" query:"out_trade_num"` // 商户订单号
	Otime        int64   `form:"otime" json:"otime" query:"otime"`                         // 成功/失败时间，10位时间戳
	State        int64   `form:"state" json:"state" query:"state"`                         // 充值状态；-1取消， 0充值中， 1充值成功 ，2充值失败，3部分成功（-1,2做失败处理；1做成功处理；3做部分成功处理）
	Mobile       string  `form:"mobile" json:"mobile" query:"mobile"`                      // 充值手机号
	Remark       string  `form:"remark" json:"remark" query:"remark"`                      // 备注信息
	ChargeAmount float64 `form:"charge_amount" json:"charge_amount" query:"charge_amount"` // 充值成功面额
	Voucher      string  `form:"voucher" json:"voucher" query:"voucher"`                   // 凭证
	ChargeKami   string  `form:"charge_kami" json:"charge_kami" query:"charge_kami"`       // 卡密/流水号
	Sign         string  `form:"sign" json:"sign" query:"sign"`                            // 签名字符串，用于验签,以保证回调可靠性。
}

// RechargeNotifyHertz 充值结果通知-异步通知
// https://www.kancloud.cn/boyanyun/boyanyun_huafei/3097255
func (c *Client) RechargeNotifyHertz(ctx context.Context, h *app.RequestContext) (validateJson ResponseRechargeNotifyHertz, err error) {
	err = h.BindAndValidate(&validateJson)
	return validateJson, err
}

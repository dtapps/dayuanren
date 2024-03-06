package dayuanren

import (
	"context"
	"github.com/gin-gonic/gin"
)

type ResponseRechargeNotifyGin struct {
	Userid       int64   `form:"userid" json:"userid" uri:"userid" binding:"omitempty"`                      // 商户ID
	OrderNumber  string  `form:"order_number" json:"order_number" uri:"order_number" binding:"omitempty"`    // true
	OutTradeNum  string  `form:"out_trade_num" json:"out_trade_num" uri:"out_trade_num" binding:"omitempty"` // 商户订单号
	Otime        int64   `form:"otime" json:"otime" uri:"otime" binding:"omitempty"`                         // 成功/失败时间，10位时间戳
	State        int64   `form:"state" json:"state" uri:"state" binding:"omitempty"`                         // 充值状态；-1取消， 0充值中， 1充值成功 ，2充值失败，3部分成功（-1,2做失败处理；1做成功处理；3做部分成功处理）
	Mobile       string  `form:"mobile" json:"mobile" uri:"mobile" binding:"omitempty"`                      // 充值手机号
	Remark       string  `form:"remark" json:"remark" uri:"remark" binding:"omitempty"`                      // 备注信息
	ChargeAmount float64 `form:"charge_amount" json:"charge_amount" uri:"charge_amount" binding:"omitempty"` // 充值成功面额
	Voucher      string  `form:"voucher" json:"voucher" uri:"voucher" binding:"omitempty"`                   // 凭证
	ChargeKami   string  `form:"charge_kami" json:"charge_kami" uri:"charge_kami" binding:"omitempty"`       // 卡密/流水号
	Sign         string  `form:"sign" json:"sign" uri:"sign" binding:"omitempty"`                            // 签名字符串，用于验签,以保证回调可靠性。
}

// RechargeNotifyGin 充值结果通知-异步通知
// https://www.kancloud.cn/boyanyun/boyanyun_huafei/3097255
func (c *Client) RechargeNotifyGin(ctx context.Context, g *gin.Context) (validateJson ResponseRechargeNotifyGin, err error) {
	err = g.ShouldBind(&validateJson)
	return validateJson, err
}

package promotion

import (
	"encoding/json"
	"errors"
	"strconv"

	"github.com/XiBao/jos/api"
	"github.com/XiBao/jos/sdk"
	"github.com/XiBao/jos/sdk/request/union/promotion"
	"github.com/daviddengcn/ljson"
)

type UnionPromotionCodeRequest struct {
	api.BaseRequest
	MaterialId string `json:"materialId"`           // 推广物料
	SiteId     string `json:"siteId"`               // 站点ID是指在联盟后台的推广管理中的网站Id、APPID（1、通用转链接口禁止使用社交媒体id入参；2、订单来源，即投放链接的网址或应用必须与传入的网站ID/AppID备案一致，否则订单会判“无效-来源与备案网址不符”）
	PositionId uint64 `json:"positionId,omitempty"` // 推广位id
	SubUnionId string `json:"subUnionId,omitempty"` // 子联盟ID（需要联系运营开通权限才能拿到数据）
	Ext1       string `json:"ext1,omitempty"`       // 推客生成推广链接时传入的扩展字段（查看订单对应字段信息，需要联系运营开放白名单才能看到）
	Pid        string `json:"pid,omitempty"`        // 联盟子站长身份标识，格式：子站长ID_子站长网站ID_子站长推广位ID
	ChainType  uint   `json:"chainType,omitempty"`  // 转链类型，1：长链， 2 ：短链 ，3： 长链+短链，默认短链
	CouponUrl  string `json:"couponUrl,omitempty"`  // 优惠券领取链接，在使用优惠券、商品二合一功能时入参，且materialId须为商品详情页链接
}

type UnionPromotionCodeResponse struct {
	ErrorResp *api.ErrorResponnse             `json:"error_response,omitempty"`
	Data      *UnionPromotionCodeResponseData `json:"jd_union_open_promotion_common_get_response,omitempty"`
}

type UnionPromotionCodeResponseData struct {
	Result string `json:"result,omitempty"`
}

type UnionPromotioncodeResult struct {
	Code    int             `json:"code,omitempty"`
	Message string          `json:"message,omitempty"`
	Data    json.RawMessage `json:"data,omitempty"`
}

type PromotionCodeResp struct {
	ClickURL string `json:"clickURL,omitempty"`
	ShortURL string `json:"shortURL,omitempty"`
}

// 获取通用推广链接
func UnionPromotionCodeGet(req *UnionPromotionCodeRequest) (string, error) {
	client := sdk.NewClient(req.AnApiKey.Key, req.AnApiKey.Secret)
	client.Debug = req.Debug
	r := promotion.NewUnionPromotionCodeRequest()
	codeReq := &promotion.PromotionCodeReq{
		MaterialId: req.MaterialId,
		SiteId:     req.SiteId,
		PositionId: req.PositionId,
		SubUnionId: req.SubUnionId,
		Ext1:       req.Ext1,
		Pid:        req.Pid,
		CouponUrl:  req.CouponUrl,
	}
	r.SetPromotionCodeReq(codeReq)

	result, err := client.Execute(r.Request, req.Session)
	if err != nil {
		return "", err
	}
	var response UnionPromotionCodeResponse
	err = ljson.Unmarshal(result, &response)
	if err != nil {
		return "", err
	}

	if response.Data == nil {
		return "", errors.New("no data")
	}
	var ret UnionPromotioncodeResult
	err = ljson.Unmarshal([]byte(response.Data.Result), &ret)
	if err != nil {
		return "", err
	}

	if ret.Code != 200 {
		return "", &api.ErrorResponnse{Code: strconv.FormatInt(int64(ret.Code), 10), ZhDesc: ret.Message}
	}

	var codeResp PromotionCodeResp
	err = json.Unmarshal(ret.Data, &codeResp)
	if err != nil {
		return "", err
	}
	return codeResp.ClickURL, nil
}

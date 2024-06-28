package service

import (
	"github.com/gin-gonic/gin"
	"github.com/pipiguanli/douyincloud_mock/consts"
	Err "github.com/pipiguanli/douyincloud_mock/errors"
	"github.com/pipiguanli/douyincloud_mock/utils"
	"log"
	"time"
)

type WebhookQaExtra struct {
	QaPath           *string `json:"webhook_qa_path"`
	WebhookSignature *string `json:"webhookheader_x_douyin_signature"`
	WebhookMsgId     *string `json:"webhookheader_msg_id"`
}

func WebhookCallback(ctx *gin.Context) {
	reqPath := ctx.FullPath()

	// 请求头
	webhookSignature := ctx.Request.Header.Get(consts.WebhookHeader_X_Douyin_Signature)
	webhookMsgId := ctx.Request.Header.Get(consts.WebhookHeader_Msg_Id)
	if err := utils.CheckHeaders(ctx); err != nil {
		TemplateFailure(ctx, Err.NewQaError(Err.InvalidParamErr, err.Error()))
		return
	}
	if len(utils.GetHeaderByName(ctx, consts.Header_StressTag)) > 0 {
		// sleep 随机 100ms ~ 1000ms（0.1s ~ 0.5s）
		num := utils.GenerateRandInt(100, 500)
		time.Sleep(time.Duration(num) * time.Millisecond)
	}
	qaExtra := &WebhookQaExtra{
		QaPath:           &reqPath,
		WebhookSignature: &webhookSignature,
		WebhookMsgId:     &webhookMsgId,
	}

	// 请求体
	var commonReq WebhookCallbackReq
	err := ctx.Bind(&commonReq)
	if err != nil {
		TemplateFailure(ctx, Err.NewQaError(Err.ParamsResolveErr))
		return
	}
	log.Printf("[QA] request=%+v", utils.ToJsonString(&commonReq))

	switch commonReq.Event {
	case "verify_webhook":
		type ContentVerifyWebhook struct {
			Challenge string `json:"challenge"`
		}
		type ReqVerifyWebhook struct {
			Event      string               `json:"event"`
			ClientKey  string               `json:"client_key"`
			FromUserId string               `json:"from_user_id"`
			ToUserId   string               `json:"to_user_id"`
			Content    ContentVerifyWebhook `json:"content"`
		}
		type RespVerifyWebhook struct {
			Challenge string          `json:"challenge"`
			QaExtra   *WebhookQaExtra `json:"qa_extra"`
		}
		var req ReqVerifyWebhook
		err := ctx.Bind(&req)
		if err != nil {
			TemplateFailure(ctx, Err.NewQaError(Err.ParamsResolveErr))
			return
		}
		resp := &RespVerifyWebhook{
			Challenge: req.Content.Challenge,
			QaExtra:   qaExtra,
		}
		httpStatusCode := 200
		ctx.JSON(httpStatusCode, resp)
		log.Printf("[QA] response=%+v, httpStatusCode=%+v", utils.ToJsonString(resp), httpStatusCode)
	}
}

type WebhookCallbackReq struct {
	Event      string      `json:"event"`
	ClientKey  string      `json:"client_key"`
	FromUserId string      `json:"from_user_id"`
	ToUserId   string      `json:"to_user_id"`
	Content    interface{} `json:"content"`
}

type WebhookCallbackResp struct {
	ErrNo   int         `json:"err_no"`
	ErrTips string      `json:"err_tips"`
	Content interface{} `json:"content"`
	QaExtra *QaExtra    `json:"qa_extra"`
}

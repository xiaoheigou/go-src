package controller

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"strconv"
	"sync"
	"yuudidi.com/pkg/models"
	"yuudidi.com/pkg/protocol/response"
	"yuudidi.com/pkg/protocol/response/err-code"
	"yuudidi.com/pkg/service"
	"yuudidi.com/pkg/utils"
)

// use default options
var upgrader = websocket.Upgrader{
	ReadBufferSize:   1024,
	WriteBufferSize:  1024,
	//HandshakeTimeout: 5 * time.Second,
	// 取消ws跨域校验
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var ACKMsg = models.Msg{ACK:models.ACK}

var clients = new(sync.Map)

var engine = service.NewOrderFulfillmentEngine(nil)

func HandleWs(context *gin.Context) {

	var connIdentify string
	merchantId := context.Query("merchantId")
	h5 := context.Query("h5")

	if merchantId != "" {
		connIdentify = merchantId
	} else if h5 != "" {
		connIdentify = h5
	} else {
		context.JSON(400, "bad request")
	}

	c, err := upgrader.Upgrade(context.Writer, context.Request, nil)
	if err != nil {
		utils.Log.Infof("upgrade:", err)
		return
	}
	utils.Log.Debugf("connIdentify: %s", connIdentify)


	var msg models.Msg

	clients.Store(connIdentify, c)

	sessionKey := utils.UniqueMerchantOnlineKey()
	if merchantId != "" {
		utils.SetCacheSetMember(sessionKey,merchantId)
	}
	defaultCloseHandler := c.CloseHandler()
	c.SetCloseHandler(func(code int, text string) error {
		result := defaultCloseHandler(code, text)
		utils.Log.Debugf("Disconnected from server ", result)
		if merchantId != "" {
			utils.DelCacheSetMember(sessionKey, merchantId)
		}
		clients.Delete(connIdentify)
		return result
	})
	defer c.Close()
	ACKMsg.Data = make([]interface{},0)
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			utils.Log.Debugf("read:%v", err)
			break
		}
		utils.Log.Debugf("message: %s", message)
		err = json.Unmarshal(message, &msg)
		if err == nil {
			utils.Log.Debugf("recv: %v", msg)
			engine.UpdateFulfillment(msg)
			ACKMsg.MsgType = msg.MsgType
			c.WriteJSON(ACKMsg)
			switch msg.MsgType {
			case models.StartOrder:
				//开始接单
			case models.StopOrder:
				//停止接单
			}
		}
	}
}

// @Summary 通知消息
// @Tags WebSocket
// @Description websocket通知
// @Accept  json
// @Produce  json
// @Param body body models.Msg true "输入参数"
// @Success 200 {object} response.UpdateDistributorsRet "成功（status为success）失败（status为fail）都会返回200"
// @Router /w/users/{uid} [put]
func Notify(c *gin.Context) {
	var param models.Msg
	var ret response.EntityResponse
	if err := c.ShouldBind(&param); err != nil {
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.RequestParamErr.Data()
		c.JSON(200, ret)
	}

	value, err := json.Marshal(param)
	if err != nil {
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.RequestParamErr.Data()
		c.JSON(200, ret)
	}

	// send message to merchant
	for _, merchantId := range param.MerchantId {
		temp := strconv.FormatInt(merchantId,10)
		if conn, ok := clients.Load(temp); ok {
			c := conn.(*websocket.Conn)
			err := c.WriteMessage(websocket.TextMessage, value)
			if err != nil {
				utils.Log.Errorf("client.WriteJSON error: %v", err)
				clients.Delete(temp)
			}
		}
	}

	// send message to h5
	for _, h5 := range param.H5 {
		if conn, ok := clients.Load(h5); ok {
			c := conn.(*websocket.Conn)
			err := c.WriteMessage(websocket.TextMessage, value)
			if err != nil {
				utils.Log.Errorf("client.WriteJSON error: %v", err)
				clients.Delete(h5)
			}
		}
	}
	ret.Status = response.StatusSucc
	c.JSON(200, ret)
}

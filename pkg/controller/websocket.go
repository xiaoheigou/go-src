package controller

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"strconv"
	"sync"
	"time"
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
	HandshakeTimeout: 5 * time.Second,
	// 取消ws跨域校验
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

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
	utils.Log.Debugf("connIdentify: %v", connIdentify)
	defer c.Close()

	var msg models.Msg

	clients.Store(connIdentify, c)

	sessionAutoKey := utils.UniqueMerchantOnlineAutoKey()
	sessionKey := utils.UniqueMerchantOnlineKey()
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			utils.Log.Infof("read:", err)
			//socket 链接断掉，如果是merchant 将session删掉
			if merchantId != "" {
				clients.Delete(connIdentify)
				utils.DelCacheSetMember(sessionKey, merchantId)
				utils.DelCacheSetMember(sessionAutoKey, merchantId)
			}
			break
		}

		err = json.Unmarshal(message, &msg)
		if err == nil {
			utils.Log.Debugf("recv: %v", msg)
			engine.UpdateFulfillment(msg)
		}
	}
}

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

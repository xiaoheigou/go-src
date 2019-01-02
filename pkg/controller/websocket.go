package controller

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/typa01/go-utils"
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

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = 40 * time.Second

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var ACKMsg = models.Msg{ACK: models.ACK}

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

	var id int
	if merchantId != "" {
		temp, err := strconv.ParseInt(merchantId, 10, 64)
		if err != nil {
			context.JSON(400, "bad request")
		} else {
			service.SetOnlineMerchant(int(temp))
		}
	}

	c, err := upgrader.Upgrade(context.Writer, context.Request, nil)
	if err != nil {
		utils.Log.Infof("upgrade:", err)
		return
	}
	utils.Log.Debugf("connIdentify: %s", connIdentify)

	var msg models.Msg

	var orderTofulfill service.OrderToFulfill
	clients.Store(connIdentify, c)

	defaultCloseHandler := c.CloseHandler()
	c.SetCloseHandler(func(code int, text string) error {
		result := defaultCloseHandler(code, text)
		utils.Log.Debugf("Disconnected from server: %s", connIdentify)
		if merchantId != "" {
			service.DelOnlineMerchant(id)
		}
		clients.Delete(connIdentify)
		return result
	})
	ACKMsg.Data = make([]interface{}, 0)
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Close()
	}()
	//定时发送ping消息
	go func() {
		for {
			select {
			case <-ticker.C:
				utils.Log.Debugf("send ping message,connidentify:%s", connIdentify)
				if err := c.WriteMessage(websocket.PingMessage, nil); err != nil {
					utils.Log.Errorf("send PingMessage is error;error:%v", err)
					return
				}
			}
		}
	}()
	//处理返回的pong消息
	c.SetPongHandler(func(string) error {
		utils.Log.Debugf("receive pong message:%s", connIdentify)
		c.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			utils.Log.Debugf("read websocket  connIdentify:%s is error :%v", connIdentify, err)
			if merchantId != "" {
				service.DelOnlineMerchant(id)
			}
			clients.Delete(connIdentify)
			break
		}
		utils.Log.Debugf("message: %s", message)
		err = json.Unmarshal(message, &msg)
		if err == nil {
			if msg.MsgType == models.Accept {
				data := msg.Data
				if len(data) > 0 && merchantId != "" {
					if id, err := strconv.ParseInt(merchantId, 10, 64); err == nil {
						if b, err := json.Marshal(data[0]); err == nil {
							if err := json.Unmarshal(b, &orderTofulfill); err == nil {
								utils.Log.Debugf("accept msg,%v", orderTofulfill)
								engine.AcceptOrder(orderTofulfill, id)
							}
						}
					}
				}
			} else {
				engine.UpdateFulfillment(msg)
			}
			ACKMsg.MsgType = msg.MsgType
			ACKMsg.MsgId = tsgutils.GUID()
			if err := c.WriteJSON(ACKMsg); err != nil {
				utils.Log.Errorf("can't send ACKMsg,error:%v", err)
			}
			//switch msg.MsgType {
			//case models.StartOrder:
			//	//开始接单
			//case models.StopOrder:
			//	//停止接单
			//}
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
// @Router /notify [post]
func Notify(c *gin.Context) {
	var param models.Msg
	var ret response.EntityResponse
	if err := c.ShouldBind(&param); err != nil {
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.RequestParamErr.Data()
		c.JSON(200, ret)
	}
	param.MsgId = tsgutils.GUID()
	value, err := json.Marshal(param)
	if err != nil {
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.RequestParamErr.Data()
		c.JSON(200, ret)
	}
	utils.Log.Debugf("notify message:%s", value)

	// send message to merchant
	for _, merchantId := range param.MerchantId {
		temp := strconv.FormatInt(merchantId, 10)
		if conn, ok := clients.Load(temp); ok {
			c := conn.(*websocket.Conn)
			err := c.WriteMessage(websocket.TextMessage, value)
			if err != nil {
				utils.Log.Errorf("client.WriteJSON merchantId:%s error: %v ", temp, err)
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
				utils.Log.Errorf("client.WriteJSON h5:%s error: %v", h5, err)
				clients.Delete(h5)
			}
		}
	}
	ret.Status = response.StatusSucc
	c.JSON(200, ret)
}

func init() {
	//服务重启删掉redis里面的key
	utils.RedisClient.Del(utils.UniqueMerchantOnlineAutoKey(), utils.UniqueMerchantOnlineKey(), utils.UniqueMerchantOnlineAcceptKey())
}

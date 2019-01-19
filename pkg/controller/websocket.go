package controller

import (
	"encoding/json"
	"errors"
	jwt_lib "github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/typa01/go-utils"
	"github.com/zzh20/timewheel"
	"net/http"
	"strconv"
	"sync"
	"time"
	"yuudidi.com/pkg/models"
	"yuudidi.com/pkg/protocol/response"
	"yuudidi.com/pkg/protocol/response/err_code"
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

var (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = utils.Config.GetInt("websocket.timeout.pong")

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = utils.Config.GetInt("websocket.timeout.ping")

	// Maximum message size allowed from peer.
	maxMessageSize = 512
	//ACK message
	ACKMsg = models.Msg{ACK: models.ACK}
	//client maps
	clients = new(sync.Map)
	//fulfillment engine
	engine = service.NewOrderFulfillmentEngine(nil)
	//ping time wheel
	pingWheel *timewheel.TimeWheel
	//logout time wheel timeout
	logoutTimeout = utils.Config.GetInt("websocket.timeout.logout")
	//logout wheel
	logoutWheel *timewheel.TimeWheel
)

func HandleWs(context *gin.Context) {

	var connIdentify string
	var id int
	merchantId := context.Query("merchantId")
	h5 := context.Query("h5")

	if merchantId != "" {
		connIdentify = merchantId
		if utils.Config.GetString("appauth.skipauth") != "true" {
			// 当appauth.skipauth不为true时，才认证token
			if !tokenVerify(context, merchantId) {
				utils.Log.Warnf("merchant [%v] auth fail", merchantId)
				return
			}
		}
		temp, err := strconv.ParseInt(merchantId, 10, 64)
		if err != nil {
			context.JSON(400, "bad request")
		} else {
			service.SetOnlineMerchant(int(temp))
			id = int(temp)
		}
	} else if h5 != "" {
		connIdentify = h5
		//判断订单是否存在
		if utils.DB.First(&models.Order{}, "order_number = ?", h5).RecordNotFound() {
			context.JSON(400, "bad request")
			return
		}
	} else {
		context.JSON(400, "bad request")
		return
	}

	c, err := upgrader.Upgrade(context.Writer, context.Request, nil)
	if err != nil {
		utils.Log.Infof("upgrade:", err)
		return
	}
	utils.Log.Debugf("connIdentify: %s", connIdentify)

	var orderToFulfill service.OrderToFulfill
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
	defer c.Close()
	//添加ping时间轮
	pingWheel.Add(connIdentify)
	c.SetPingHandler(func(string) error {
		utils.Log.Debugf("receive ping message:%s", connIdentify)
		if _, ok := clients.Load(connIdentify); ok {
			utils.Log.Debugf("websocket conn is exist :%s", connIdentify)
			if err := c.WriteMessage(websocket.PongMessage, nil); err != nil {
				utils.Log.Errorf("reply PongMessage is error;error:%v", err)
				clients.Delete(connIdentify)
				return err
			}
		} else {
			utils.Log.Debugf("websocket conn is not exist :%s", connIdentify)
		}
		return nil
	})
	//处理返回的pong消息
	c.SetPongHandler(func(string) error {
		utils.Log.Debugf("receive pong message:%s", connIdentify)
		pingWheel.Add(connIdentify)
		logoutWheel.Remove(connIdentify)
		c.SetReadDeadline(time.Now().Add(time.Duration(pongWait) * time.Second))
		return nil
	})

	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			utils.Log.Debugf("read websocket connIdentify:%s is error :%v", connIdentify, err)
			if merchantId != "" {
				service.DelOnlineMerchant(id)
			}
			clients.Delete(connIdentify)
			break
		}
		utils.Log.Debugf("receive message: %s", message)
		var msg models.Msg
		err = json.Unmarshal(message, &msg)
		if err == nil {
			ACKMsg.MsgType = msg.MsgType
			ACKMsg.MsgId = tsgutils.GUID()
			if h5 != "" && msg.MsgType == models.NotifyPaid {
				//TODO 不要从订单里面进行查询
				order := models.Order{}
				if utils.DB.First(&order, "order_number = ?", connIdentify).RecordNotFound() {
					utils.Log.Debugf("websocket not found order,")
				} else {
					distributor := models.Distributor{}
					if utils.DB.First(&distributor, "id = ? ", order.DistributorId).RecordNotFound() {
						utils.Log.Debugf("websocket not found order,")
					} else {
						data := models.OrderData{
							PageUrl:       distributor.PageUrl,
							OrderNumber:   connIdentify,
							DistributorId: distributor.Id,
						}
						ACKMsg.Data = append(ACKMsg.Data, data)
					}
				}
			} else {
				ACKMsg.Data = make([]interface{}, 0)
			}
			if err := c.WriteJSON(ACKMsg); err != nil {
				utils.Log.Errorf("can't send ACKMsg,error:%v", err)
			}
			if msg.MsgType == models.Accept {
				data := msg.Data
				if len(data) > 0 && merchantId != "" {
					if id, err := strconv.ParseInt(merchantId, 10, 64); err == nil {
						if b, err := json.Marshal(data[0]); err == nil {
							if err := json.Unmarshal(b, &orderToFulfill); err == nil {
								utils.Log.Debugf("accept msg,%v", orderToFulfill)
								engine.AcceptOrder(orderToFulfill, id)
							}
						}
					}
				}
			} else {
				engine.UpdateFulfillment(msg)
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
		return
	}
	param.MsgId = tsgutils.GUID()
	value, err := json.Marshal(param)
	if err != nil {
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.RequestParamErr.Data()
		c.JSON(200, ret)
		return
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

func ping(data interface{}) {
	connIdentify := data.(string)
	logoutWheel.Add(connIdentify)
	if conn, ok := clients.Load(connIdentify); ok {
		c := conn.(*websocket.Conn)
		utils.Log.Debugf("send ping message,connidentify:%s", connIdentify)
		if err := c.WriteMessage(websocket.PingMessage, nil); err != nil {
			utils.Log.Errorf("send PingMessage is error;error:%v", err)
			clients.Delete(connIdentify)
			return
		}
	}
}

func logout(data interface{}) {
	connIdentify := data.(string)

	tx := utils.DB.Begin()
	var merchant models.Merchant
	if tx.Set("gorm:query_option", "FOR UPDATE").First(&merchant, "id = ?", connIdentify).RecordNotFound() {
		utils.Log.Debugf("websocket logout not found merchant,id:%s", connIdentify)
	}
	if err := tx.Model(&models.Preferences{}).Where("id = ?", merchant.PreferencesId).
		Updates(map[string]interface{}{"in_work": 0, "auto_accept": 0, "auto_confirm": 0}).Error; err != nil {
		utils.Log.Errorf("websocket logout update merchant in_work and auto_accept and auto_confirm is 0 failed,merchantId:%s", connIdentify)
		logoutWheel.Add(connIdentify)
	}
	if err := tx.Commit().Error; err != nil {
		utils.Log.Errorf("websocket logout commit is failed,merchantId:%s", connIdentify)
		logoutWheel.Add(connIdentify)
	}

}

func tokenVerify(context *gin.Context, merchantId string) bool {
	secret := utils.Config.GetString("appauth.authkey")
	token, err := request.ParseFromRequest(context.Request, request.OAuth2Extractor, func(token *jwt_lib.Token) (interface{}, error) {
		b := ([]byte(secret))
		return b, nil
	})
	if err != nil {
		utils.Log.Errorf("Authorization fail [%v]", err)
		context.AbortWithError(401, err)
		return false
	}

	if claims, ok := token.Claims.(jwt_lib.MapClaims); ok && token.Valid {
		tokenUid := claims["uid"]
		if tokenUid != merchantId {
			utils.Log.Errorf("jwt can only access resource belong to uid [%v], but you want to access resource belong to uid [%s]", tokenUid, merchantId)
			context.AbortWithError(401, errors.New("Authorization fail"))
			return false
		}
	} else {
		utils.Log.Errorln("Parse jwt error")
		context.AbortWithError(401, errors.New("Parse jwt error"))
		return false
	}
	return true
}

func InitWheel() {
	//启动ping时间轮
	utils.Log.Debugf("ping wheel period,%d", pingPeriod)
	pingWheel = timewheel.New(1*time.Second, pingPeriod, ping)
	pingWheel.Start()

	//启动
	utils.Log.Debugf("init logout wheel,period:%d", logoutTimeout)
	logoutWheel = timewheel.New(1*time.Second, logoutTimeout, logout)
	logoutWheel.Start()
}

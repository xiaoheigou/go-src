package service

import (
	"bytes"
	"encoding/json"

	"yuudidi.com/pkg/models"
	"yuudidi.com/pkg/utils"
)

var serviceURL = utils.Config.GetString("websocket.url") + ":" + utils.Config.GetString("websocket.port") + "/notify"

type notification struct {
	MsgType   models.MsgType `json:"msg_type"`
	Merchants []int64        `json:"merchant_id"`
	H5        []string       `json:"h5"`
	Timeout   uint           `json:"timeout"`
	Data      interface{}    `json:"data"`
}

// NotifyThroughWebSocketTrigger - send notification message through websocket server
func NotifyThroughWebSocketTrigger(msgType models.MsgType, merchants *[]int64, h5 *[]string, timeout uint, data interface{}) error {
	body := notification{
		MsgType:   msgType,
		Merchants: *merchants,
		H5:        *h5,
		Timeout:   timeout,
		Data:      data,
	}
	var bodyBytes []byte
	var err error
	if bodyBytes, err = json.Marshal(body); err != nil {
		utils.Log.Errorf("Unable to marshal notification message: %v", err)
		return err
	}
	if _, err := utils.HTTPPost(serviceURL, "application/json", bytes.NewBuffer(bodyBytes)); err != nil {
		utils.Log.Errorf("Error occured in sending notification through websocket: %v", err)
	}
	return nil
}

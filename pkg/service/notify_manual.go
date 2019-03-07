package service

import (
	"fmt"
	"github.com/pkg/errors"
	"yuudidi.com/pkg/models"
	"yuudidi.com/pkg/utils"
)

func ManualNotify(orderNumber string) error {
	order := models.Order{}
	if utils.DB.First(&order, "order_number = ?", orderNumber).RecordNotFound() {
		utils.Log.Errorf("Unable to find order %s", orderNumber)
		return errors.New(fmt.Sprintf("can not find order %s", orderNumber))
	}
	notify := Order2Notify(order)
	resp, err := PostNotifyToServer(order, notify)
	if err == nil && resp.Status == SUCCESS {
		return nil
	} else {
		utils.Log.Warnf("call PostNotifyToServer fail for order %s", orderNumber)
		return errors.New("call PostNotifyToServer fail")
	}
}

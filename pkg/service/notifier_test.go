package service

import (
	"testing"

	"yuudidi.com/pkg/models"
)

func TestNotifyThroughWebSocket(t *testing.T) {
	NotifyThroughWebSocketTrigger(
		models.SendOrder,
		&[]int64{1, 2, 3},
		&[]string{"ORD12342"},
		600,
		[]string{"anything"})
}

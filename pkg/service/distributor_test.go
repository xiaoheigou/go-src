package service

import (
	"fmt"
	"testing"
	"yuudidi.com/pkg/models"
	"yuudidi.com/pkg/utils"
)

func TestGetDistributorByAPIkey(t *testing.T) {
	distributor := models.Distributor{
		ApiKey: "test1",
		Name:   "test1",
	}
	if err := utils.DB.Create(&distributor).Error; err != nil {
		fmt.Printf("create distributor is error,error:%v",err)
		t.Fail()
	}

	if result,err := GetDistributorByAPIKey("test1");err != nil {
		fmt.Printf("get distributor by apikey error,error:%v",err)
		t.Fail()
	} else if result.Id != distributor.Id {
		fmt.Printf("distributor not match error,error:%v",err)
		t.Fail()
	}
	utils.DB.Unscoped().Delete(&distributor,"id = ?" ,distributor.Id)
}

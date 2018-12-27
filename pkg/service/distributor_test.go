package service

import (
	"fmt"
	"testing"
	"yuudidi.com/pkg/models"
	"yuudidi.com/pkg/protocol/response"
	"yuudidi.com/pkg/utils"
)

func TestGetDistributorByAPIkey(t *testing.T) {
	distributor := models.Distributor{
		ApiKey: "test1",
		Name:   "test1",
	}
	if err := utils.DB.Create(&distributor).Error; err != nil {
		fmt.Printf("create distributor is error,error:%v", err)
		t.Fail()
	}

	if result, err := GetDistributorByAPIKey("test1"); err != nil {
		fmt.Printf("get distributor by apikey error,error:%v", err)
		t.Fail()
	} else if result.Id != distributor.Id {
		fmt.Printf("distributor not match error,error:%v", err)
		t.Fail()
	}
	utils.DB.Unscoped().Delete(&distributor, "id = ?", distributor.Id)
}

func TestCreateDistributor(t *testing.T) {
	args := response.CreateDistributorsArgs{
		Password: "123456",
		ApiKey:   "test1",
		Name:     "test1",
		Username: "test",
	}

	if result := CreateDistributor(args); result.Status == response.StatusSucc {
		var user models.User
		if err := utils.DB.First(&user, "username = ?", args.Username).Error; err != nil {
			fmt.Printf("create distributor is failed")
			t.Fail()
		}
		if distributor, err := GetDistributorByAPIKey(args.ApiKey); distributor.Id <= 0 || err != nil {
			fmt.Printf("create distributor is failed")
			t.Fail()
		}
	} else {
		fmt.Printf("create distributor is failed")
		t.Fail()
	}
	utils.DB.Unscoped().Delete(&models.User{}, "username = ?", args.Username)
	utils.DB.Unscoped().Delete(&models.Distributor{}, "api_key = ?", args.ApiKey)
}

func TestUpdateDistributor(t *testing.T) {
	args := response.CreateDistributorsArgs{
		Password: "123456",
		ApiKey:   "test1",
		Name:     "test1",
		Username: "test",
	}

	if result := CreateDistributor(args); result.Status == response.StatusSucc {
		var user models.User
		if err := utils.DB.First(&user, "username = ?", args.Username).Error; err != nil {
			fmt.Printf("create distributor is failed")
			t.Fail()
		}
		if distributor, err := GetDistributorByAPIKey(args.ApiKey); distributor.Id <= 0 || err != nil {
			fmt.Printf("create distributor is failed")
			t.Fail()
		}
	} else {
		fmt.Printf("create distributor is failed")
		t.Fail()
	}
	utils.DB.Unscoped().Delete(&models.User{}, "username = ?", args.Username)
	utils.DB.Unscoped().Delete(&models.Distributor{}, "api_key = ?", args.ApiKey)
}

package service

import (
	"fmt"
	"testing"
	"yuudidi.com/pkg/protocol/response"
)

func TestCreateUser_Admin(t *testing.T) {
	user := response.UserArgs{
		Username: "admin",
		Password: "admin",
		Role:     0,
		Phone:    "13112345678",
		Email:    "admin@123.com",
		Address:  "123",
	}
	result := CreateUser(user, nil)
	if result.Status == response.StatusFail {
		t.Fail()
	}
	if result.Status == response.StatusSucc {
		t.Log("create admin user is success")
	}
}

func TestCreateUser_Distributor(t *testing.T) {
	user := response.UserArgs{
		Username: "distributor",
		Password: "distributor",
		Role:     2,
		Phone:    "13112345678",
		Email:    "admin@123.com",
		Address:  "123",
	}
	result := CreateUser(user, nil)
	if result.Status == response.StatusFail {
		t.Fail()
	}
	if result.Status == response.StatusSucc {
		t.Log("create admin user is success")
	}
	args := response.CreateDistributorsArgs{
		//Password: "123456",
		ApiKey:    "123456",
		ApiSecret: "654321",
		Name:      "test1",
		//Username: "test",
	}

	if result := CreateDistributor(args); result.Status == response.StatusSucc {
		if distributor, err := GetDistributorByAPIKey(args.ApiKey); distributor.Id <= 0 || err != nil {
			fmt.Printf("create distributor is failed")
			t.Fail()
		}
	} else {
		fmt.Printf("create distributor is failed")
		t.Fail()
	}
}

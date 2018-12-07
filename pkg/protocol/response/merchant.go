package response

import "yuudidi.com/pkg/models"

type GetProfileRet struct {
	CommonRet
	Entity struct {
		// 用户id
		Uid      int `json:"uid" example:123`
		// 用户昵称
		NickName string `json:"nickname" example:"老王"`
		// 平台币的符号
		AssetSymbol string `json:"asset_symbol" example:"BTUSD"`
		// 当前承兑商所有的平台币余额（包含被冻结的平台币）
		AssetTotal  string `json:"asset_total" example:"2000"`
		// 当前承兑商被冻结的平台币数量
		AssetFrozen  string `json:"asset_frozen" example:"100"`
	}
}

type GetAuditStatusRet struct {
	CommonRet
	Entity struct {
		// 用户id
		Uid      int `json:"uid" example:123`
		// user_status可以为0/1/2，分别表示“正常/待审核/冻结”
		UserStatus int `json:"user_status" example:0`
		// 客服联系信息
		ContactPhone string `json:"contact_phone" example:"13812341234"`
		// 额外的信息
		ExtraMessage string `json:"extra_message" example:"您由于xx原因，未通过审核"`
	}
}

type SetNickNameArg struct {
 	// 想设置的新昵称
	NickName string `json:"nickname" example:"王老板"`
}

type SetNickNameRet struct {
	CommonRet
	Entity struct {
		// 用户id
		Uid int `json:"uid" example:123`
	}
}

type SetWorkModeArg struct {
	// 用户id
	Uid int `json:"uid" example:123`
	// 是否接单(1:开启，0:关闭)
	Accept int `json:"accept" example:1`
	// 是否自动接单(1:开启，0:关闭)
	Auto int `json:"auto" example:1`
}

type SetWorkModeRet struct {
	CommonRet
	Entity struct {
		// 用户id
		Uid int `json:"uid" example:123`
	}
}

type GetWorkModeRet struct {
	CommonRet
	Entity struct {
		// 用户id
		Uid int `json:"uid" example:123`
		// 是否接单(1:开启，0:关闭)
		Accept int `json:"accept" example:1`
		// 是否自动接单(1:开启，0:关闭)
		Auto int `json:"auto" example:1`
	}
}

type SetIdentifyArg struct {
	// 用户id
	Uid    int    `json:"uid" example:123`
	Phone  string `json:"phone" example:13012341234`
	Email  string `json:"email" example:"xxx@xxx.com"`
	IdCard int    `json:"idcard" example:11088888888888888`
}

type SetIdentifyRet struct {
	CommonRet
	Entity struct {
		// 用户id
		Uid int `json:"uid" example:123`
	}
}

type GetIdentifyRet struct {
	CommonRet
	Entity struct {
		// 用户id
		Uid    int    `json:"uid" example:123`
		Phone  string `json:"phone" example:13012341234`
		Email  string `json:"email" example:"xxx@xxx.com"`
		IdCard string `json:"idcard" example:"11088888888888888"`
	}
}

type MerchantRet struct {
	CommonRet

	Entity struct {
		Data []models.Merchant `json:"data"`
	}
}

type RechargeArgs struct {
	Currency string `json:"currency" binding:"required" example:"BTUSD"`
	Count    string `json:"count" binding:"required" example:"200"`
}

type RechargeRet struct {
	CommonRet

	Entity struct {
		Balance string `json:"balance"`
	}
}

type ApproveArgs struct {
	//审核操作 1：通过 0：不通过
	Operation    int    `json:"operation" binding:"required" example:"1"`
	ContactPhone string `json:"currency" binding:"required" example:"BTUSD"`
	ExtraMessage string `json:"count" binding:"required" example:"200"`
}

type ApproveRet struct {
	CommonRet

	Entity struct {
		//
		Uid    int `json:"uid" example:"1"`
		Status int `json:"status"`
	}
}

type FreezeArgs struct {
	//冻结操作 1：冻结 0：解冻
	Operation    int    `json:"operation" binding:"required" example:"1"`
	ContactPhone string `json:"currency" binding:"required" example:"BTUSD"`
	ExtraMessage string `json:"count" binding:"required" example:"200"`
}

type FreezeRet struct {
	CommonRet
	Entity struct {
		Uid    int `json:"uid" example:"1"`
		Status int `json:"status"`
	}
}

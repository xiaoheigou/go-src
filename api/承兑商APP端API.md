

# 承兑商APP端API 设计文档
* [1.登陆](#登陆)
* [2.注册](#注册)
* [3.忘记密码](#忘记密码)
* [4.审核状态](#审核状态)
* [5.抢单](#抢单)
* [6.支付](#支付)
* [7.订单](#订单)
* [8.个人设置](#个人设置)
	* [1.个人设置信息](#个人设置信息) 
	* [2.推送设置](#推送设置) 
	* [3.收款账户信息](#收款账户信息) 
	* [4.身份认证](#身份认证) 
* 	[9.退出](#退出)
*  [10.session验证](#session验证)

API问题
1. 版本号要返回什么？
2. 订单界面需要的所有进行中的订单是通过websocket实时推送，还是在返回的订单列表数组数据在每一个数据里面添加一个字段，还是在entities外添加一个字段？
3. 身份认证暂时先不用做，以后再说？
****
所有的请求的header应该带有session的信息，这个还没有考虑好，所以请求里面还没有写。
所有支付流程需要输入密码字段
版本更新单独借口，放在session验证接口里面
抢单界面需要角标websocket实时推送
订单界面接口加字段返回所有进行中的数据

****

## 登陆
  
```
url：/login
method：POST
request body：
{
  "account":"131",
  "password":"123456"，
  "version": ""
}
response
{
  "status": "success",
  "entities": [
    {
      "user_status": 0,
      "nickname": "123"
    }
  ]
}
状态字段值
uid：主键
user_status: 0//用户状态值，
nickname: 昵称
```
## 注册
	
```
url: /register
method: POST
request body
{
  "phone": "13112345678",
  "email": "123@163.com",
  "password": "123456qwe"
}

response body
成功返回
{
  "status": "success",
  "entities": []
}
失败返回
{
	errCode:"10001",
	errMsg:"phone number is already exist",
}
status: 状态返回 
errCode:错误码
errMsg: 错误信息描述
```
## 忘记密码

忘记密码使用应该是邮箱或者手机号，发送随机数给对应的手机或者邮箱
	
**获取验证码**
```
url: verification/codes
method: post
request body
{
  "account": "13112345678"
}
account: 可以为邮箱或者手机号

response:
{
  "status": "success",
  "entities": []
}
status: 验证码发送成功状态,success,failed
```
**确认输入验证码是否正确**
```
url: verification/verify/codes
method: post
request body
{
  "account": "13112345678",
  "password": "123123",
  "code": "123456"
}

response:
{
  "status": "success",
  "entities": []
}

status：验证码是否正确，success，failed
```


## 审核状态

**获取审核状态**
```
url: /audit?uid=123
method: get
request param
uid : 登陆用户的id

response body:
{
  "audit_status": "unaudited",
  "auditer_phone": "12312312313",
  "audit_message": "123123"
}

```
## 抢单
1. 抢单列表
	方案：
		使用websocket推送
		使用自动接单，订单不会推送到此界面

推送消息格式
```
{
  "order_id": "123123123123",
  "user_id": "123",
  "user_nickname": "sk",
  "quantity": "123",
  "price": "6.5",
  "rate": "0",
  "order_status": "unpaid",
  "merchant_order_type": "2",
  "pay_type": "1",
  "name": "",
  "bank_account": "",
  "bank": "",
  "bank_branch": "",
  "gr_code": "",
  "remain_time": 900,
  "total":123,
  "create_timestamp": 15144444444
}

order_id: 订单id
user_id: 用户id
user_nickname: 用户昵称
quantity：订单数量
price: 订单价格
order_status:订单状态
merchant_order_type：承兑商订单类型 （1:buy 2:sell）
承兑商订单类型为1时，下面的字段会有值
pay_type: 支付类型 1:微信 2:支付宝 3:银行卡
name:账户名(支付宝,微信)
bank_account: 银行卡账号
bank：所属银行
bank_branch: 银行分行
gr_code: 二维码地址
```

**抢单**
使用websocket抢单,抢单发送消息格式
```
消息格式
{
  "uid": "1",
  "order_id": "2",
  "operation": 1
}
uid: 用户id
order_id: 订单id
operation: 订单操作 //1:抢单 
```


**取消订单**
需要通知对方
```
url: orders/{order_id}/cancel
method: put

path value:
order_id: 订单id

request body:
{
  "password": "123456"
}


订单取消成功
{
  "status": "success",
  "entities": []
}
```
**订单确认收款**
需要通知对方已经确认收款
```
url: orders/{order_id}/verify
method: put

path value:
order_id: 订单id

request body:
{
  "password": "123456"
}

订单确认收款成功
{
  "status": "success",
  "entities": []
}
```


## 支付

**去付款**
二次输入密码确认,订单信息里面会有支付的信息
```
url: orders/{order_id}/pay
method: post

path value:
order_id: 订单id

request body:
{
  "password": "123456"
}

订单确认收款成功
{
  "status": "success",
  "entities": []
}
```

**确认付款成功**
通知另一方已经付款
```
url: orders/{order_id}/paid/verify
method: post

path value:
order_id: 订单id

request body:
{
  "password": "123456"
}

订单确认付款成功
{
  "status": "success",
  "entities": []
}
```

## 订单
**订单列表数据获取**
可能需要支持字段搜索
```
url: /orders/{uid}
method: get

request param:

page: 0 //从0开始
size: 10 //默认为10
order_type: 0,1,2,3//0:进行中（所有没有确认收款的订单）1:买入 2:卖出 3:全部
下面的param用于搜索
order_id: 订单id
timestamp: 时间 (可能是订单的创建时间)

response body:
{
  "status": "success",
  "entities": [
    {
      "order_id": "123",
      "order_status": "unpaid",
      "user_id": "123",
      "user_nickname": "sk",
      "quantity": "123",
      "price": "6.5",
      "rate": "",
      "merchant_order_type": "2",
      "pay_type": "1",
      "name": "",
      "bank_account": "",
      "bank": "",
      "bank_branch": "",
      "gr_code": "",
      "remain_time": 900,
      "create_timestamp": 15144444444,
      "grap_timestamp": 15115151151,
      "total":123,
      "accomplish_timestamp": ""
    }
  ]
}

order_id: 订单id
user_id: 用户id
user_nickname: 用户昵称
quantity：订单数量
price: 订单价格
order_status:订单状态
merchant_order_type：承兑商订单类型 （1:buy 2:sell）
grap_timestamp: 抢单时间
total:总共进行的订单
accomplish_timestamp: 订单完成时间
```

## 个人设置

### 个人设置信息
**获取个人设置信息**
```
url： infos/{uid}
method: get

path value:
uid: 用户id

response body:
{
  "status": "success",
  "entities": [
    {
      "nickname": "sk",
      "available_count": "0",
      "frozen_count": "12345",
      "total_count": "12345"
    }
  ]
}
```
**修改昵称**
```
url: infos/{uid}
method: put

path value:
uid: 用户id

request body:
{
  "nickname": "sk"
}

修改成功返回
response body:
{
	"status":"success",
	"entities": []
}
```

### 推送设置
**获取推送设置信息**
```
url: settings/{uid}
method: get

path value:
uid:用户id

//自动接单开关
response body:
{
  "status": "success",
  "entities": [
    {
      "order": 1,
      "automatic_order": 1
    }
  ]
}
order: 是否接单(1:开启，0:关闭)
automatic_order:是否自动接单(1:开启，0:关闭)
```
**修改推送设置**
```
url: otc/settings/{uid}
method: put

path value:
uid:用户id

request body:
{
  "order": 1,
  "automatic_order": 1
}

修改成功返回
response body:
{
  "status": "success",
  "entities": []
}
```
### 收款账户信息
**获取收款账户信息**
```
url: otc/settings/payment-account/{uid}
method: get

path value:
uid: 用户id

response body:
{
  "status": "success",
  "entitys": [
    {
      "id": 1,
      "uid": 1,
      "pay_type": 1,
      "name": "",
      "bank_account": "6666",
      "bank": "招商银行",
      "bankBranch": "上海分行",
      "gr_code": ""
    },
    {
      "id": 2,
      "uid": 1,
      "pay_type": 2,
      "name": "sky",
      "bank_account": "",
      "bank": "",
      "bankBranch": "",
      "gr_code": "http://code.img"
    },
    {
      "id": 3,
      "uid": 1,
      "pay_type": 3,
      "name": "sky",
      "bank_account": "",
      "bank": "",
      "bankBranch": "",
      "gr_code": "http://code.img"
    }
  ]
}
id: 主键，无实际意义
pay_type: 支付类型 1:银行卡 2:微信 3:支付宝
name:账户名(支付宝,微信)
bank_account: 银行卡账号
bank：所属银行
bankBranch: 银行分行
gr_code: 二维码地址
```


**添加收款账户信息**

```
url: otc/settings/payment-account/{uid}
method: post

path value:
uid: 用户id

request body:
{
  "pay_type": 3,
  "name": "sky",
  "bank_account": "",
  "bank": "",
  "bankBranch": "",
  "gr_code": "http://code.img"
}
pay_type: 支付类型 1:银行卡 2:微信 3:支付宝
name:账户名(支付宝,微信)
bank_account: 银行卡账号
bank：所属银行
bankBranch: 银行分行
gr_code: 二维码地址

添加成功返回
response body:
{
  "status": "success",
  "entities": []
}
```
**修改收款账户信息**
需要先确定是否有正在进行的订单，如果有不允许修改
```
url: otc/settings/payment-account/{uid}
method: put

path value:
uid: 收款账号信息主键

request body:
{
  "id": 1,
  "pay_type": 3,
  "name": "sky",
  "bank_account": "",
  "bank": "",
  "bankBranch": "",
  "gr_code": "http://code.img"
}

id: 收款账号信息主键
pay_type: 支付类型 1:银行卡 2:微信 3:支付宝
name:账户名(支付宝,微信)
bank_account: 银行卡账号
bank：所属银行
bankBranch: 银行分行
gr_code: 二维码地址

修改成功返回
response body:
{
  "status": "success",
  "entities": []
}
```
**删除收款账户信息**
需要先确定是否有正在进行的订单，如果有不允许删除
```
url: otc/settings/payment-account/{uid}
method: delete

path value:
uid: 用户id 

request body:
{
  "id": 1
}
id: 收款账号信息主键

删除成功返回
response body:
{
  "status": "success",
  "entities": []
}
```
### 身份认证
**获取身份认证信息**
```
url: otc/settings/identity/{uid}
method: get

response body:
{
  "phone": "13112345678",
  "email": "123@163.com",
  "idcard": 1234123123123123
}
phone: 手机号
email: 邮箱
idcard：身份证号
```

**添加身份认证信息**
```
url: otc/settings/identity/{uid}
method: post

path value:
uid: 用户id

request body:
{
  "phone": "13112345678",
  "email": "123@163.com",
  "idcard": 1234123123123123
}
phone: 手机号
email: 邮箱
idcard：身份证号
```

**修改身份认证信息**
修改手机号或者邮箱的认证信息？不确定是否需要
```
url: otc/settings/identity
method: put

```

## 退出

```
url: otc/logout
method: post

response body:
{
  "status": "success",
  "entities": []
}
```

## session验证

```
url: otc/verification/sessions/{uid}
method: get

response body:
{
  "status": "success",
  "entities": [
    {
      "version":"",
      "expire": "0"
    }
  ]
}
expire: 是否过期，1：过期 0：未过期
version: 版本号
```

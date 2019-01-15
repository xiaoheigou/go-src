# 下载和安装依赖：
mkdir -p ~/go/src/
cd ~/go/src/
git clone https://YOUR_GIT_ACCOUNT@github.com/AAAChain/YuuPay_core-service.git yuudidi.com # 请替换账号名字
cd yuudidi.com
dep init -v     # 安装依赖
dep ensure -update  # 更新依赖


# 创建数据库
CREATE DATABASE `otc` CHARACTER SET utf8mb4;

# 为金融滴滴平台创建assets记录（设定distributor_id为1的记录为金融滴滴平台），用户提现订单中金融滴滴平台赚的钱会放到这行记录中
# 没有在程序中自动增加assets记录的原因是：避免并发导致创建两条或多条记录
insert into distributors (id, name, created_at, updated_at) values (1, "金融滴滴平台", now(), now());
insert into assets (distributor_id, currency_crypto, created_at, updated_at) values (1, "BTUSD", now(), now());


# 为避免死锁问题。在一个事务中，请按下面顺序去修改表：
orders
fulfillment_events
fulfillment_logs
assets （同一个表内的顺序：先平台、再币商、最后金融滴滴平台自己）
asset_histories
payment_infos

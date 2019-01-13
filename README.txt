# 下载和安装依赖：
mkdir -p ~/go/src/
cd ~/go/src/
git clone https://YOUR_GIT_ACCOUNT@github.com/AAAChain/YuuPay_core-service.git yuudidi.com # 请替换账号名字
cd yuudidi.com
dep init -v     # 安装依赖
dep ensure -update  # 更新依赖


# 创建数据库
CREATE DATABASE `otc` CHARACTER SET utf8mb4;

# 为金融滴滴平台创建assets记录，用户提现订单中金融滴滴平台赚的钱会放到这行记录中
# 没有在程序中自动增加这条记录的原因是：避免并发导致创建两条或多条记录
insert into assets (distributor_id, currency_crypto, created_at, updated_at) values (1, "BTUSD", now(), now());

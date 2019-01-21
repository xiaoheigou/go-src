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
insert into assets (merchant_id, distributor_id, currency_crypto, quantity, qty_frozen, created_at, updated_at) values (0, 1, "BTUSD", 0, 0, now(), now());

# 初始化管理后台管理员和平台商
INSERT INTO `users` VALUES (1, 'admin', 0xEE8B96D019A045B545C22AFB473930EA8DE80EB2168D1B054C2A4C0D90107E3D, 0x9B15DD242C9797949D60D7960161FC1007FD096A17BC0E8EF063295C44AD6B900423BD6CEA504246B1E262D15859FFC2B4253A6ADD5B8904D7DF01DFF1A76F03, 'Argon2', '13112345678', 'admin@123.com', '123', 0, 0, '2018-12-16 21:09:12', '2018-12-16 21:09:12', NULL);
INSERT INTO `users` VALUES (2, 'distributor', 0x97533E322F25E5134A1B29D5FEA24AEE5272AB34A5113120DBF72EA0D129FA41, 0x3B7D9B1887D5616CF3C8654F4E543A60B4347EFCDF1AC32505EF91250E0736860F1DD6574C9FDF85B827E6DA876A6E108FAECAC462CF5D3D97867C9F4CE892D3, 'Argon2', '13112345678', '', '', 0, 2, '2018-12-17 11:12:36', '2018-12-17 11:12:36', NULL);

# 为避免死锁问题。在一个事务中，请按下面顺序去修改表：
orders
fulfillment_events
fulfillment_logs
assets （同一个表内的顺序：先平台、再币商、最后金融滴滴平台自己）
asset_histories
payment_infos

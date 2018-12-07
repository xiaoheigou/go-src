# 下载和安装依赖：
mkdir -p ~/go/src/
cd ~/go/src/
git clone https://YOUR_GIT_ACCOUNT@github.com/AAAChain/YuuPay_core-service.git yuudidi.com # 请替换账号名字
cd yuudidi.com
dep init -v     # 安装依赖
dep ensure -update  # 更新依赖


# 启动app服务器：
go run cmd/app-server/main.go

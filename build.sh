proj_dir=${HOME}/go/src/yuudidi.com
cd $proj_dir
git pull origin dev

go build -o ${HOME}/dist/web cmd/server/web/main.go
go build -o ${HOME}/dist/app cmd/server/app/main.go
go build -o ${HOME}/dist/ticket cmd/server/ticket/main.go
go build -o ${HOME}/dist/background cmd/server/background/main.go
go build -o ${HOME}/dist/websocket cmd/server/websocket/main.go

cp configs/config.yml ${HOME}/dist/config.yml

cd ${HOME}

zip -r ${HOME}/dist.zip dist/*


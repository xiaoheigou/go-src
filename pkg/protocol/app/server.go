package app

import (
	"os"
	"path"
	"runtime"

	"github.com/gin-gonic/gin"
	"yuudidi.com/pkg/protocol/route"
	"yuudidi.com/pkg/utils"
)

func RunServer(port string) error {
	defer utils.DB.Close()
	defer utils.Log.OSFile.Close()
	r := gin.Default()

	// Set a lower memory limit for multipart forms (default is 32 MiB)
	r.MaxMultipartMemory = 1 << 20 // 1 MiB

	//store := cookie.NewStore([]byte("secret"))
	//r.Use(sessions.Sessions("session", store))

	route.AppServer(r)

	_, fileName, _, _ := runtime.Caller(0)
	rootPath := path.Join(fileName, "../../../../configs/")
	err := os.Chdir(rootPath)
	if err != nil {
		panic(err)
	}
	r.Run(":" + port)
	return nil
}

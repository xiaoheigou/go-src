package ticket
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

	route.WebServer(r)

	_, fileName, _, _ := runtime.Caller(0)
	rootPath := path.Join(fileName, "../../../../configs/")
	err := os.Chdir(rootPath)
	if err != nil {
		panic(err)
	}
	r.Run(":" + port)
	return nil
}


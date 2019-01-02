package swagger

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"os"
	"path"
	"runtime"
	"yuudidi.com/pkg/protocol/route"
	"yuudidi.com/pkg/utils"

	"github.com/swaggo/gin-swagger"              // gin-swagger middleware
	"github.com/swaggo/gin-swagger/swaggerFiles" // swagger embed files

	_ "yuudidi.com/docs" // docs is generated by Swag CLI, you have to import it.
)

func RunServer(port string) error {
	defer utils.DB.Close()
	defer utils.Log.OSFile.Close()
	r := gin.Default()

	// Set a lower memory limit for multipart forms (default is 32 MiB)
	r.MaxMultipartMemory = 1 << 20 // 1 MiB


	//store := cookie.NewStore([]byte("secret"))
	//r.Use(sessions.Sessions("session", store))

	printRequestBody := true
	if printRequestBody {
		r.Use(RequestLogger())
	}

	printResponseBody := true
	if printResponseBody {
		r.Use(ginBodyLogMiddleware)
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	route.AppServer(r)
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

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		buf, _ := ioutil.ReadAll(c.Request.Body)
		rdr1 := ioutil.NopCloser(bytes.NewBuffer(buf))
		rdr2 := ioutil.NopCloser(bytes.NewBuffer(buf)) //We have to create a new Buffer, because rdr1 will be read.

		body := readBody(rdr1)

		fmt.Println("====Request body begin==== [" + c.Request.Method + "] url: " + c.Request.URL.Path)
		if c.Request.Header.Get("Content-Type") == "multipart/form-data" {
			// Ignore form-data (probably binary data)
		} else if len(body) > 0 {
			fmt.Println(body) // Print request body
		}
		fmt.Println("====Request  body  end====")

		c.Request.Body = rdr2
		c.Next()
	}
}

func readBody(reader io.Reader) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(reader)

	s := buf.String()
	return s
}


// https://stackoverflow.com/questions/38501325/how-to-log-response-body-in-gin
type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func ginBodyLogMiddleware(c *gin.Context) {
	blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
	c.Writer = blw
	c.Next()
	// statusCode := c.Writer.Status()
	fmt.Println("====Response body begin==== [" + c.Request.Method + "] url: " + c.Request.URL.Path)
	if len(blw.body.String()) > 0 {
		fmt.Println(blw.body.String()) // Print response body
	}
	fmt.Println("====Response  body  end====")
}

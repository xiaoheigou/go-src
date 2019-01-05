package main

import (
	"fmt"
	"io"
	"net/http"
)

func main() {

	//这里的url只是个例子，具体的url及端口号需要根据平台商提供的serverUrl来确定
	http.HandleFunc("/create-order", func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/create-order" {
			http.Error(w, "404 not found.", http.StatusNotFound)
			return
		}

		if err := req.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		//从body里取status的值，以判断订单的状态
		status := req.FormValue("status")
		fmt.Println("Success")

		io.WriteString(w, "Success")

		//status=1,订单创建成功，此时需要返回Success，并告知服务端创建订单成功，当status=4,等待确认付款，返回Success，并告诉服务端等待确认付款。
		if status == "1" {
			io.WriteString(w, "Success")
			fmt.Println("create order success!")
			return
		} else if status == "4" {
			io.WriteString(w, "Success")
			fmt.Println("wait to confirm!")
			return
		}

	})

	//server.crt ,server.key是平台商服务需要提供的相应证书，client也需要使用统一套证书，只有当证书验证通过才能正常访问服务,8085是服务的监听端口（例子）
	if e := http.ListenAndServeTLS(":8085", "server.crt", "server.key", nil); e != nil {
	}

}

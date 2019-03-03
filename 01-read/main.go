package main

import (
	"fmt"
	"io"
	"strings"
)

func main() {
	// 1. readFrom
	// data, err := ReadFrom(os.Stdin, 11)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(string(data))

	//2.readat

	reader := strings.NewReader("Gohahahfdsfgs dd")

	p := make([]byte, 6)
	n, err := reader.ReadAt(p, 2)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%s,%d====if n<len(p),return err\n", string(p), n)

	//2.go standard
	// reader := strings.NewReader("Go语言中文网")
	// p := make([]byte, 6)
	// n, err := reader.ReadAt(p, 2)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Printf("%s, %d\n", p, n)

	//2.compare to read !!!
	reader2 := strings.NewReader("asdf")
	p2 := make([]byte, 6)
	n2, err2 := reader2.Read(p2)
	if err2 != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%s,%d====readat2\n", string(p2), n2)
	n3, err3 := reader2.Read(p2)
	if err3 != nil {
		fmt.Println(err3, n3, "====read3,this is less stricker")
		return
	}

}

// 多态的实现
func ReadFrom(read io.Reader, num int) ([]byte, error) {
	p := make([]byte, num)
	n, err := read.Read(p)
	if n > 0 {
		return p[:n], nil
	}
	return p, err

}

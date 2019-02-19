package timewheel

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func TestAddString(t *testing.T) {
	tw := New(1*time.Second, 1, "timewheel", func(obj interface{}) {
		fmt.Printf("%+v\n", obj)
	})
	tw.Start()
	go tw.Add("timewheel string")
	time.Sleep(2 * time.Second)
}

type object struct {
	A string
	B int
}

func (o object) MarshalBinary() (data []byte, err error) {
	return json.Marshal(o)
}

func (o object) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, o)
}

func TestAddObject(t *testing.T) {
	tw := New(1*time.Second, 1, "timewheel", func(obj interface{}) {
		fmt.Printf("%+v\n", obj)
	})
	tw.Start()
	go tw.Add(object{A: "timewheel", B: 10})
	time.Sleep(2 * time.Second)
}

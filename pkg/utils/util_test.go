package utils

import (
	"fmt"
	"reflect"
	"testing"
)

func TestInterSetInt64(t *testing.T) {
	list1 := []int64{1, 2, 3}
	list2 := []int64{3, 3, 2, 4}

	expect := []int64{2, 3}

	if !reflect.DeepEqual(InterSetInt64(list1, list2), expect) {
		t.Fail()
	}
}

func TestRemoveElement(t *testing.T) {
	t1 := []int64{0,1,2,3}

	t2 := RemoveElement(t1,1)

	fmt.Printf("%v",t2)
}
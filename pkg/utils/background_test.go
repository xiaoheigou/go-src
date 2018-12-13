package utils

import (
	"testing"
)

func TestBackgroupJob(t *testing.T) {
	SetSettings()
	//add job
	for i := 0; i < 100; i++ {
		var prio Priority
		if i%2 == 0 {
			prio = NormalPriority
		} else {
			prio = HighPriority
		}
		AddBackgroundJob("test", prio, []interface{}{"hello", "world"})
	}
}

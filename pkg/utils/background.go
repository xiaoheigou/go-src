package utils

import (
	"fmt"
	"strconv"

	"github.com/benmanns/goworker"
)

//Priority - to defint priority of the task
type Priority string

const (
	//HIGH - higest priority
	HIGH Priority = "high" //weight=4
	//NORMAL - normal priority
	NORMAL Priority = "normal" //weight=2
	//LOW - lowest priority
	LOW Priority = "low" //weight=1
)

func init() {
	//get non-string configuration then convert them
	var skipTLS bool
	var connections, concurrency int64
	var interval float64
	if skipTLS, err = strconv.ParseBool(Config.GetString("background.store.skiptlsverify")); err != nil {
		Log.Errorf("Wrong configuration: background.store.skiptlsverify, should be boolean. Set to default false.")
		skipTLS = false
	}
	if connections, err = strconv.ParseInt(Config.GetString("background.connections"), 10, 0); err != nil {
		Log.Errorf("Wrong configuration: background.connections, should be int. Set to default 10.")
		connections = 10
	}
	if concurrency, err = strconv.ParseInt(Config.GetString("background.concurrency"), 10, 0); err != nil {
		Log.Errorf("Wrong configuration: background.concurrency, should be int. Set to default 2.")
		concurrency = 2
	}
	if interval, err = strconv.ParseFloat(Config.GetString("background.interval"), 32); err != nil {
		Log.Errorf("Wrong configuration: background.interval, should be int. Set to default 5.0.")
		interval = 5.0
	}
	settings := goworker.WorkerSettings{
		URI:            Config.GetString("background.store.uri"),
		TLSCertPath:    Config.GetString("background.store.tlscertpath"),
		SkipTLSVerify:  skipTLS,
		Connections:    int(connections),
		Queues:         []string{string(HIGH) + "=4", string(NORMAL) + "=2", string(LOW) + "=1"},
		UseNumber:      true,
		ExitOnComplete: false,
		Concurrency:    int(concurrency),
		Namespace:      Config.GetString("background.namespace"),
		IntervalFloat:  interval,
	}
	goworker.SetSettings(settings)
	//launch background worker engine
	if err = goworker.Work(); err != nil {
		fmt.Printf("Can't launch background worker engine: %v", err)
	}
}

//RegisterWorkerFunc - register job consumer/worker function.
func RegisterWorkerFunc(jobName string, workerFunc func(string, ...interface{}) error) {
	goworker.Register(jobName, workerFunc)
}

//AddBackgroundJob - job producer to push running job at background
func AddBackgroundJob(jobName string, prio Priority, params ...interface{}) {
	goworker.Enqueue(&goworker.Job{
		Queue: string(prio),
		Payload: goworker.Payload{
			Class: jobName,
			Args:  params,
		},
	})
}

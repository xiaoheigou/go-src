package utils

import (
	"strconv"

	"github.com/benmanns/goworker"
)

//Priority - to defint priority of the task
type Priority string

const (
	//HighPriority - highest priority
	HighPriority Priority = "high" //weight=4
	//NormalPriority - normal priority
	NormalPriority Priority = "normal" //weight=2
	//LowPriority - lowest priority
	LowPriority Priority = "low" //weight=1
)

var engineInited = false

//SetSettings - set background job engine configuration
func SetSettings() {
	//get non-string configuration then convert them
	var skipTLS bool
	var connections, concurrency int64
	var interval float64
	var err error
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
		QueuesString:   string(HighPriority) + "=4," + string(NormalPriority) + "=2," + string(LowPriority) + "=1",
		UseNumber:      true,
		ExitOnComplete: false,
		Concurrency:    int(concurrency),
		Namespace:      Config.GetString("background.namespace"),
		IntervalFloat:  interval,
	}
	goworker.SetSettings(settings)
}

//RegisterWorkerFunc - register job consumer/worker function.
func RegisterWorkerFunc(jobName string, workerFunc func(string, ...interface{}) error) {
	goworker.Register(jobName, workerFunc)
}

//AddBackgroundJob - job producer to push running job at background
func AddBackgroundJob(jobName string, prio Priority, params ...interface{}) {
	if err := goworker.Enqueue(&goworker.Job{
		Queue: string(prio),
		Payload: goworker.Payload{
			Class: jobName,
			Args:  params,
		},
	}); err != nil {
		Log.Warnf("Failed to add background job: %v", err)
	}
}

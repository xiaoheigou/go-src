package background

import (
	"github.com/benmanns/goworker"
	"yuudidi.com/pkg/utils"
)

//LaunchBackgroundEngine - launch background engine
func LaunchBackgroundEngine() error {
	utils.SetSettings()
	//launch background worker engine
	if err := goworker.Work(); err != nil {
		utils.Log.Errorf("Can't launch background worker engine: %v", err)
		return err
	}
	return nil
}
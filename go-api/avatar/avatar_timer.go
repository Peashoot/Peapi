package avatar

import (
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/peashoot/peapi/common"
	"github.com/peashoot/peapi/config"
)

var c common.Mutex

// ClearOverdueFiles 清除过期文件
func ClearOverdueFiles() {
	if !c.TryLock() {
		return
	}
	defer c.Unlock()
	fileInfoList, err := ioutil.ReadDir(config.Config.AvatarConfig.AvatarFileLocalPath)
	if err != nil {
		log.Print("can't read path:", config.Config.AvatarConfig.AvatarFileLocalPath)
		return
	}
	limitTime := time.Now().Add(time.Duration(-config.Config.AvatarConfig.AvatarSaveDuration) * time.Minute)
	for _, fileInfo := range fileInfoList {
		if fileInfo.ModTime().Before(limitTime) {
			if err := os.Remove(fileInfo.Name()); err != nil {
				log.Print("failure to delete overdue file:", fileInfo.Name(), "err is:", err.Error())
			}
			log.Print("delete overdue file:", fileInfo.Name())
		}
	}
}

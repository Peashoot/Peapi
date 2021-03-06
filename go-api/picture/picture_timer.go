package picture

import (
	"io/ioutil"
	"log"
	"os"
	"path"
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
	fileInfoList, err := ioutil.ReadDir(config.Config.PictureOperateConfig.PictureGenerateFolder)
	if err != nil {
		log.Print("can't read path:", config.Config.PictureOperateConfig.PictureGenerateFolder)
		return
	}
	limitTime := time.Now().Add(time.Duration(-config.Config.PictureOperateConfig.PictureSaveDuration) * time.Minute)
	for _, fileInfo := range fileInfoList {
		if fileInfo.ModTime().Before(limitTime) {
			if err := os.Remove(path.Join(config.Config.PictureOperateConfig.PictureGenerateFolder, fileInfo.Name())); err != nil {
				log.Print("failure to delete overdue file:", fileInfo.Name(), "err is:", err.Error())
			}
			log.Print("delete overdue file:", fileInfo.Name())
		}
	}
}

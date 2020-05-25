package common

import (
	"fmt"
	"myProject/videoCli/makerCli"
	"myProject/videoMaker/account"
	"myTool/file"
	"os"
	"strings"
	"sync"
)

var MakerEngine *VideoMakerEngine
var isDealing = false
var lock = sync.Mutex{}
type VideoMakerEngine struct {
	MakerCli *makerCli.Engine
	Account  *account.Account
}

func NewMakerEngine(cli *makerCli.Engine, acc *account.Account) *VideoMakerEngine  {

	MakerEngine = &VideoMakerEngine{
		MakerCli: cli,
		Account:  acc,
	}

	return MakerEngine
}

func (v *VideoMakerEngine) ClearRemoteCache() {

	if v.MakerCli.AliYunOss != nil && v.MakerCli.AliYunOss.Expiration > 0 && v.MakerCli.AliYunOss.Check() == nil {
		v.MakerCli.AliYunOss.RemoveOldObject(v.MakerCli.AliYunOss.Expiration)
	}
}

func (v *VideoMakerEngine)ClearTemp()  {

	files,_ := file.GetAllFiles(v.MakerCli.ProjectDir)
	for _, f := range files {
		if strings.HasSuffix(f,".DS_Store") {
			os.Remove(f)
		}
	}
}


func (v *VideoMakerEngine)DoMaker() {

	lock.Lock()
	defer lock.Unlock()

	files, err := file.GetCurrentFiles(v.MakerCli.WorkDir)
	_, dirs, err := file.GetCurrentFilesAndDirs(v.MakerCli.WorkDir)
	if err != nil {
		fmt.Println("获取文件失败", err)
		return
	}
	if len(files) == 0 {
		fmt.Println("文件为空，请将文件放到video目录")
		return
	}
	if isDealing {
		fmt.Println("正在处理中，请稍后")
		return
	}

	isDealing = true

	if len(dirs) > 0 {
		for _,d := range dirs {
			files, err := file.GetCurrentFiles(d)
			if err != nil || len(files) == 0 {
				continue
			}
			v.MakerCli.DoComposite(files)
		}

	}

	v.MakerCli.DoComposite(files)


	isDealing = false

}
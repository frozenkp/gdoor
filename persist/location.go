package persist

import(
  "os"
  "os/exec"
  "io/ioutil"
  "path/filepath"

  "../info"
  "../config"
  "../debug"
)

func moveToTarget(i *info.Info){
  targetDir := filepath.Join(i.GetHomePath(), config.TargetDir)

  if err := os.MkdirAll(targetDir, 0755); err != nil {
    debug.Println(err)
    return
  }

  bin, err := ioutil.ReadFile(i.GetCurPath())
  if err != nil {
    debug.Println(err)
    return
  }

  if err := ioutil.WriteFile(filepath.Join(targetDir, config.TargetName), bin, 0755); err != nil {
    debug.Println(err)
    return
  }
}

func execTarget(i *info.Info){
  cmd := exec.Command(filepath.Join(i.GetHomePath(), config.TargetDir, config.TargetName))
  if err := cmd.Start(); err != nil {
    debug.Println(err)
  }
}

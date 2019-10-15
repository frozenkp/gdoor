package info

import(
  "os"
  "path/filepath"
  "strings"

  "../debug"
)

type Info struct{
  curDir    string
  curName   string
  homePath  string
}

func Init()(*Info, error){
  i := &Info{}

  absPath, err := filepath.Abs(os.Args[0])
  if err != nil{
    debug.Println(err)
    return nil, err
  }

  pathSlice := strings.Split(absPath, "/")
  i.curName = pathSlice[len(pathSlice)-1]
  i.homePath = os.Getenv("HOME")
  i.curDir = filepath.Dir(absPath)

  return i, nil
}

func (i *Info) GetCurDir()string{
  return i.curDir
}

func (i *Info) GetCurName()string{
  return i.curName
}

func (i *Info) GetHomePath()string{
  return i.homePath
}

func (i *Info) GetCurPath()string{
  return filepath.Join(i.curDir, i.curName)
}

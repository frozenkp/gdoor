package debug

import(
  "log"

  "../config"
)

func Printf(format string, v ...interface{}){
  if(config.DEBUG){
    log.Printf(format, v...)
  }
}

func Println(v interface{}){
  if(config.DEBUG){
    log.Println(v)
  }
}

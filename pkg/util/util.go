package util

import (
	"time"

	"github.com/astaxie/beego/logs"
)

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	logs.Info("%s took %s", name, elapsed)
}

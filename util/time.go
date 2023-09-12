package util

import (
	"fmt"
	"strings"
	"time"
)

func GetTimeNowMetaFormat() string {
	t := time.Now()
	datetimeString := t.Format(time.RFC3339)
	splitted := strings.Split(datetimeString, "T")
	date := splitted[0]
	time := splitted[1][0:8]
	timeDiff := splitted[1][8:]
	timeDiff = strings.Replace(timeDiff, ":", "", 1)

	return fmt.Sprintf("%s %s %s", date, time, timeDiff)
}

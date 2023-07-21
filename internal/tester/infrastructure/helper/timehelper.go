package helper

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

func ParsePhpMicroTime(withMicro string) time.Time {
	floatTime, err := strconv.ParseFloat(withMicro, 64)
	if err != nil {
		log.Fatalln(err)
	}
	microTime := floatTime * 1000000
	return time.UnixMicro(int64(microTime))

	//////////////
	//parts := strings.Split(micro, ".")
	//
	//secsStr := parts[0]
	//microSecsStr := parts[1]
	//
	//secs, err := strconv.ParseInt(secsStr, 10, 64)
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//
	//microSecs, err := strconv.ParseInt(microSecsStr, 10, 64)
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//
	//return time.Unix(secs, microSecs*1000)
}

func ParsePDur(pDur string) time.Duration {
	_, timeStr, found := strings.Cut(pDur, "p;dur=")
	if !found {
		log.Fatalln("pDur time is not found in string" + fmt.Sprintf(" %s", pDur))
	}

	return ParseMillisecondsDur(timeStr)
}

func ParseMillisecondsDur(milli string) time.Duration {
	t, err := time.ParseDuration(milli + "ms")
	if err != nil {
		log.Fatalln(err)
	}
	return t
}

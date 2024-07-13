package ycsb

import (
	"log"
	"os"
	"strconv"

	"github.com/provok1ng/sp/client"
	"github.com/provok1ng/sp/curp"
	"github.com/provok1ng/sp/swift"
	"github.com/provok1ng/sp/dlog"
)

type SwiftClient interface {
	Connect() error
	Disconnect()
	Read(int64) []byte
	Scan(int64, int64) []byte
	Write(int64, []byte)
}

func NewswiftClient(protocol, collocated, maddr string, mport int, fast, leaderless bool, args string) SwiftClient {
	var sc SwiftClient

	l := newLogger("out")
	c := client.NewClientLog(collocated, maddr, mport, fast, leaderless, false, l)
	b := client.NewBufferClient(c, 1, 1024, 0, 5, 0)

	if err := b.Connect(); err != nil {
		log.Fatal(err)
	}

	switch protocol {
	case "client":
		sc = b
		b.WaitReplies(b.ClosestId)
	case "SwiftPaxos":
		sc = swift.NewClient(b, 1)
	case "curp":
		num, err := strconv.Atoi(args)
		if err != nil {
			sc=b
			return sc
		}
		sc = curp.NewClient(b, 0, 0, num)
	}

	return sc
}

func newLogger(logPath string) *dlog.Logger {
    logF := os.Stdout
    if logPath != "" {
        var err error
        logF, err = os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY, 0644)
        if err != nil {
            log.Fatal("Can't open log file:", logPath)
        }
        // 使用单个变量接收logF.Name()的返回值
        fileName := logF.Name()
        // 假设dlog.New需要的第二个参数是一个布尔值
        boolFlag := true
        // 调用dlog.New，传递文件名和布尔标志
        return dlog.New(fileName, boolFlag)
    }
    // 如果logPath为空，使用标准输出作为日志文件
    return dlog.New("default.log", false)
}

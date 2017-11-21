package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/anchnet/rds-agent/cron"
	"github.com/anchnet/rds-agent/funcs"
	"github.com/anchnet/rds-agent/g"
	"github.com/anchnet/rds-agent/http"
)

func main() {
	cfg := flag.String("c", "cfg.json", "configuration file")
	version := flag.Bool("v", false, "show version")
	check := flag.Bool("check", false, "check collector")
	flag.Parse()

	if *version {
		fmt.Println(g.VERSION)
		os.Exit(0)
	}

	if *check {
		funcs.CheckCollector()
		os.Exit(0)
	}

	g.ParseConfig(*cfg)
	//init seelog
	g.InitSeeLog()

	g.InitRootDir()
	g.InitLocalIps()
	g.InitRpcClients()
	funcs.BuildMappers()
	cron.Collect()

	go http.Start()

	select {}

}

//
//import (
//	"github.com/denverdino/aliyungo/rds"
//)
//
//func main() {
//	DescribeDBInstancePerformance()
//}
//
//func DescribeDBInstancePerformance() {
//	AccessKeyId := "LTAIHEAN6tTQ5BAD"
//	AccessKeySecret := "mQxK7zSMSJT0t6Dw3E89xmFGH1Brty"
//	client := rds.NewClient(AccessKeyId, AccessKeySecret)
//	args := &rds.DescribeDBInstancePerformanceArgs{
//		DBInstanceId :"rm-m5eq2938h3h1rwpv5",
//		StartTime    :"2017-06-02T17:00Z",
//		EndTime      :"2017-06-02T18:00Z",
//	}
//	args.Setkey("MySQL_NetworkTraffic")
//	resp, err := client.DescribeDBInstancePerformance(args)
//	if err != nil {
//		println("Failed to describe rds regions %v", err)
//	}
//
//	PerformanceKey := resp.PerformanceKeys.PerformanceKey
//	println("all PerformanceKey %++v.", PerformanceKey)
//
//}

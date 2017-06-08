package funcs

import (
	"github.com/open-falcon/common/model"
	"time"
	"github.com/51idc/rds-agent/g"
	"github.com/denverdino/aliyungo/rds"
	"strings"
	"log"
)

type alerror struct {
	RequestId string
	HostId    string
	Code      string
	Message   string
}

func RDSMetrics() (L []*model.MetricValue) {
	db_type := g.Config().DBType
	// push version
	if g.Config().SmartAPI.Enabled {
		describeDBInstanceAttribute, err := DescribeDBInstanceAttribute()
		if err != nil {
			log.Println("GET Instance Version Interface Error ：", err.Error())
		} else {
			dbInstanceAttribute := describeDBInstanceAttribute.Items.DBInstanceAttribute
			if len(dbInstanceAttribute) != 1 {
				log.Println("GET Instance Version No Data Error .")
			} else {
				version := ""
				version += string(dbInstanceAttribute[0].Engine)
				version += " "
				version += dbInstanceAttribute[0].EngineVersion
				smartAPI_Push(version)
			}
		}
	}
	var metric_list map[string]bool
	if (db_type == "rds_mysql") {
		metric_list = g.Config().MySQLMetric
	} else if (db_type == "rds_sqlserver") {
		metric_list = g.Config().SQLServerMetric
	}
	if len(metric_list) < 1 {
		log.Println("Metric Error.Check config")
		return
	}
	if len(g.Config().AccessKeyId) < 1 && len(g.Config().AccessKeySecret) < 1 {
		log.Println("AccessKeyId or AccessKeySecret Error.Check config")
		return
	}

	var metric_str string
	for k, v := range metric_list {
		if v {
			metric_str += k + ","
		}
	}

	dbInstancePerformanceResponse, err := DescribeDBInstancePerformance(metric_str)
	if err != nil {
		log.Println("GET Instance Monitor Info  Interface Error ：", err.Error())
		return
	}
	if len(dbInstancePerformanceResponse.PerformanceKeys.PerformanceKey) < 1 {
		log.Println("GET Instance Monitor Info NO  Data Error .")
	}
	for _, performanceKey := range dbInstancePerformanceResponse.PerformanceKeys.PerformanceKey {
		if len(performanceKey.Key) > 0 && metric_list[performanceKey.Key] && len(performanceKey.Values.PerformanceValue) > 0 {
			if len(performanceKey.ValueFormat) > 0 && strings.Contains(performanceKey.ValueFormat, "&") {
				performanceValue := performanceKey.Values.PerformanceValue[len(performanceKey.Values.PerformanceValue) - 1].Value
				for i, valueFormat := range strings.Split(performanceKey.ValueFormat, "&") {
					L = append(L, GaugeValue(performanceKey.Key + "_" + strings.ToUpper(valueFormat), performanceValue[i]))
				}
			} else {
				L = append(L, GaugeValue(performanceKey.Key, performanceKey.Values.PerformanceValue[len(performanceKey.Values.PerformanceValue) - 1].Value))
			}
		} else {
			log.Println("GET Instance Monitor Info  Data Info Error .")
		}
	}

	return L
}

func DescribeDBInstancePerformance(metric_str string) (rds.DescribeDBInstancePerformanceResponse, error) {
	client := ALNewClient()
	start_time, end_time := time2rfc()
	args := &rds.DescribeDBInstancePerformanceArgs{
		DBInstanceId :g.Config().DBInstanceId,
		StartTime    :start_time,
		EndTime      :end_time,
	}
	args.Setkey(metric_str)
	resp, err := client.DescribeDBInstancePerformance(args)
	return resp, err
}

func DescribeDBInstanceAttribute() (*rds.DescribeDBInstanceAttributeResponse, error) {
	client := ALNewClient()
	arrtArgs := rds.DescribeDBInstancesArgs{
		DBInstanceId: g.Config().DBInstanceId,
	}
	resp, err := client.DescribeDBInstanceAttribute(&arrtArgs)
	return resp, err
}

func time2rfc() (string, string) {
	now := time.Now().UTC()
	cus := "2006-01-02T15:04Z"
	end_time_str := now.Format(cus)
	m, _ := time.ParseDuration("-5m")
	start_time := now.Add(m)
	return start_time.Format(cus), end_time_str
}
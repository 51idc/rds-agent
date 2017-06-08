package funcs

import (
	"github.com/51idc/rds-agent/g"
	"github.com/denverdino/aliyungo/rds"
)

var alClient *rds.Client

func ALNewClient() *rds.Client {
	if alClient == nil {
		AccessKeyId := g.Config().AccessKeyId
		AccessKeySecret := g.Config().AccessKeySecret
		alClient = rds.NewClient(AccessKeyId, AccessKeySecret)
	}
	return alClient
}

var alDebugClient *rds.Client

func ALNewClientForDebug() *rds.Client {
	if alDebugClient == nil {
		AccessKeyId := g.Config().AccessKeyId
		AccessKeySecret := g.Config().AccessKeySecret
		alDebugClient = rds.NewClient(AccessKeyId, AccessKeySecret)
		alDebugClient.SetDebug(true)
	}
	return alDebugClient

}
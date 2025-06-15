package collector

import (
	"testing"

	model2 "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ces/v1/model"
	"github.com/stretchr/testify/assert"

	"github.com/huaweicloud/cloudeye-exporter/logs"
)

func TestGeipInfo_GetResourceInfo(t *testing.T) {
	CloudConf.Global.RmsRetryTimes = 1
	conf.AccessKey = "test_ak"
	conf.SecretKey = "test_sk"
	patches := getPatches()
	defer patches.Reset()
	logs.InitLog("")

	patches.ApplyFuncReturn(listResources, mockRmsResource(), nil)
	patches.ApplyFuncReturn(listAllMetrics, []model2.MetricInfoList{
		{
			MetricName: "xxxxx",
			Dimensions: []model2.MetricsDimension{
				{
					Name:  "test-dim",
					Value: "0001-0001-000001",
				},
			},
			Namespace: "SYS.GEIP",
		},
	}, nil)
	var geipInfo GeipInfo
	labels, metrics := geipInfo.GetResourceInfo()
	assert.Equal(t, 1, len(labels))
	assert.Equal(t, 1, len(metrics))
}

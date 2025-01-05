package collector

import (
	"testing"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ces/v1/model"
	"github.com/stretchr/testify/assert"

	"github.com/huaweicloud/cloudeye-exporter/logs"
)

func TestBmsGetResourceInfo(t *testing.T) {
	instances := []EcsInstancesInfo{
		{ResourceBaseInfo: ResourceBaseInfo{ID: "0001-0001-000000001", Name: "host01", EpId: "0"}},
	}
	metricConf = map[string]MetricConf{
		"SYS.BMS": {
			Resource: "rms",
			DimMetricName: map[string][]string{
				"instance_id": {"cpu_utils"},
			},
		},
	}

	patches := getPatches()
	logs.InitLog("")
	patches.ApplyFuncReturn(getAllServerFromRMS, instances, nil)
	patches.ApplyFuncReturn(getIPFromEcsInfo, "")
	defer patches.Reset()

	var bmsGetter BMSInfo
	labels, metrics := bmsGetter.GetResourceInfo()
	assert.Equal(t, 1, len(labels))
	assert.Equal(t, 1, len(metrics))

	var servicesGetter SERVICEBMSInfo
	patches.ApplyFuncReturn(listAllMetrics, []model.MetricInfoList{
		{
			Dimensions: []model.MetricsDimension{
				{
					Name:  "instance_id",
					Value: "111111",
				},
			},
			Namespace:  "SERVICE.BMS",
			MetricName: "cpu_usage",
		},
	}, nil)

	serviceLabel, _ := servicesGetter.GetResourceInfo()
	assert.Equal(t, 1, len(serviceLabel))
}

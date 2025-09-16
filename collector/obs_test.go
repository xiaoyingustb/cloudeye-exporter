package collector

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ces/v1/model"
)

func TestOBSInfo_GetResourceInfo(t *testing.T) {
	metricConf = map[string]MetricConf{
		"SYS.OBS": {
			DimMetricName: map[string][]string{
				"bucket_name": {"download_traffic"},
			},
		},
	}
	patches := gomonkey.NewPatches()
	defer patches.Reset()
	patches.ApplyFuncReturn(listAllMetrics, []model.MetricInfoList{
		{
			Dimensions: []model.MetricsDimension{
				{
					Name:  "bucket_name",
					Value: "test_name",
				},
			},
			Namespace:  "SYS.OBS",
			MetricName: "download_traffic",
		},
	}, nil)
	patches.ApplyFuncReturn(getUserInfoFromIAM, map[string]string{})
	patches.ApplyFuncReturn(getAllServerFromRMS, []EcsInstancesInfo{}, nil)
	var obsInfo OBSInfo
	CloudConf.Global.MetricInfoExpirationDays = 2
	CloudConf.Global.MetricInfoCleanThreshold = 50000
	labels, metrics := obsInfo.GetResourceInfo()
	assert.NotEmpty(t, metrics)
	assert.Empty(t, labels)
}

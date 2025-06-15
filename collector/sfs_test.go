package collector

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core"
	cesmodel "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ces/v1/model"
)

func TestSFSInfo_GetResourceInfo(t *testing.T) {
	metricConf = map[string]MetricConf{
		"SYS.SFS": {
			DimMetricName: map[string][]string{
				"share_id":    {"read_bandwidth"},
				"bucket_name": {"read_bytes_intranet"},
			},
		},
	}

	CloudConf.Global.RmsRetryTimes = 1
	conf.AccessKey = "test_ak"
	conf.SecretKey = "test_sk"
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	outputs := []gomonkey.OutputCell{
		{
			Values: gomonkey.Params{&ListAllShareResponse{
				Shares: []Share{
					{
						ID:   "id1",
						Name: "name1",
					},
				},
			}, nil},
		},
		{
			Values: gomonkey.Params{&ListAllShareResponse{
				Shares: []Share{},
			}, nil},
		},
	}

	hcClient := &core.HcHttpClient{}
	patches.ApplyMethodSeq(hcClient, "Sync", outputs)
	patches.ApplyFuncReturn(getHcClient, hcClient)
	patches.ApplyFuncReturn(listAllMetrics, []cesmodel.MetricInfoList{
		{
			Namespace:  "SYS.SFS",
			MetricName: "read_bytes_intranet",
			Dimensions: []cesmodel.MetricsDimension{
				{
					Name:  "bucket_name",
					Value: "bucket1",
				},
			},
		},
	}, nil)
	var sfsGetter SFSInfo
	labelInfos, filteredMetrics := sfsGetter.GetResourceInfo()
	assert.Equal(t, 1, len(labelInfos))
	assert.Equal(t, 2, len(filteredMetrics))
}

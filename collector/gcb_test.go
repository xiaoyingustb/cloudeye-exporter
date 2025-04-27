package collector

import (
	"testing"

	model2 "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ces/v1/model"

	"github.com/huaweicloud/cloudeye-exporter/logs"
	cc "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/cc/v3"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/cc/v3/model"
	"github.com/stretchr/testify/assert"
)

func TestGCBGetResourceInfo(t *testing.T) {
	conf.AccessKey = "test_ak"
	conf.SecretKey = "test_sk"
	patches := getPatches()
	defer patches.Reset()
	metricConf = map[string]MetricConf{
		"SYS.GCB": {
			DimMetricName: map[string][]string{
				"gcb_id": {"network_bandwidth"},
			},
		},
	}
	logs.InitLog("")
	defaultEpId := "0"
	id := "gcb-0001"
	name := "gcb1"
	connectionsPage := model.ListGlobalConnectionBandwidthsResponse{
		HttpStatusCode: 200,
		GlobalconnectionBandwidths: []model.GlobalConnectionBandwidth{
			{Id: id, Name: &name, EnterpriseProjectId: &defaultEpId},
		},
		PageInfo: &model.PageInfo{},
	}

	ccClient := cc.CcClient{}
	patches.ApplyFuncReturn(getCCClient, &ccClient)
	patches.ApplyMethodFunc(&ccClient, "ListGlobalConnectionBandwidths", func(req *model.ListGlobalConnectionBandwidthsRequest) (*model.ListGlobalConnectionBandwidthsResponse, error) {
		return &connectionsPage, nil
	})
	patches.ApplyFuncReturn(listAllMetrics, []model2.MetricInfoList{
		{
			Dimensions: []model2.MetricsDimension{
				{
					Name:  "gcb_id",
					Value: "gcb-0001",
				},
			},
			Namespace:  "SYS.GCB",
			MetricName: "network_bandwidth",
		},
	}, nil)

	var gcbgetter GCBInfo
	labels, metrics := gcbgetter.GetResourceInfo()
	assert.Equal(t, 1, len(labels))
	assert.Equal(t, 1, len(metrics))
	metricConf = map[string]MetricConf{}
}

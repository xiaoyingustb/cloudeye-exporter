package collector

import (
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/rms/v1/model"
	"github.com/stretchr/testify/assert"

	"github.com/huaweicloud/cloudeye-exporter/logs"
)

func TestERInfo_GetResourceInfo(t *testing.T) {
	metricConf = map[string]MetricConf{
		"SYS.ER": {
			Resource: "rms",
			DimMetricName: map[string][]string{
				"er_instance_id": {"instance_bits_rate_in"}, "er_instance_id,er_attachment_id": {"attachment_bits_rate_in"},
			},
		},
	}
	erInstanceName, erInstanceId, erAttachmentName, erAttachmentId, erEpId :=
		"demo-er", "0cbcebbf-7732-4d50-8186-45cdf2fec1f3", "demo-er-attachment", "440a1bf6-68cc-4a6a-b286-e89811a99b6", "test_epId"
	outputs := []gomonkey.OutputCell{
		{
			Values: gomonkey.Params{[]model.ResourceEntity{
				{
					Id:   &erInstanceId,
					Name: &erInstanceName,
					EpId: &erEpId,
				},
			}, nil},
		},
		{
			Values: gomonkey.Params{[]model.ResourceEntity{
				{
					Id:   &erAttachmentId,
					Name: &erAttachmentName,
					EpId: &erEpId,
					Properties: map[string]interface{}{
						"er_id": erInstanceId,
					},
				},
			}, nil},
		},
	}
	patches := getPatches()
	defer patches.Reset()
	logs.InitLog("")
	patches.ApplyFuncSeq(listResources, outputs)
	tmpErInfo := ERInfo{}
	resourceInfo, metrics := tmpErInfo.GetResourceInfo()
	assert.NotEmpty(t, resourceInfo)
	assert.NotEmpty(t, metrics)
}

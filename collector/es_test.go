package collector

import (
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/stretchr/testify/assert"
)

func TestEsGetResourceInfo(t *testing.T) {
	conf.AccessKey = "test_ak"
	conf.SecretKey = "test_sk"
	conf.Region = "cn-test-01"
	metricConf = map[string]MetricConf{
		"SYS.ES": {
			DimMetricName: map[string][]string{
				"cluster_id": {"disk_util"},
			},
		},
	}
	patches := gomonkey.ApplyFuncReturn(listResources, mockRmsResource(), nil)
	defer patches.Reset()

	var esgetter ESInfo
	labels, metrics := esgetter.GetResourceInfo()
	assert.Equal(t, 1, len(labels))
	assert.Equal(t, 1, len(metrics))
}

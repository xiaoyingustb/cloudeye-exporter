package collector

import (
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/stretchr/testify/assert"
)

func TestDcsGetResourceInfo(t *testing.T) {
	conf.AccessKey = "test_ak"
	conf.SecretKey = "test_sk"
	conf.Region = "cn-test-01"
	metricConf = map[string]MetricConf{
		"SYS.DCS": {
			DimMetricName: map[string][]string{
				"dcs_instance_id":           {"cpu_usage"},
				"dcs_memcached_instance_id": {"cpu_usage"},
			},
		},
	}
	patches := gomonkey.ApplyFuncReturn(listResources, mockRmsResource(), nil)
	defer patches.Reset()

	var dcsgetter DCSInfo
	labels, metrics := dcsgetter.GetResourceInfo()
	assert.Equal(t, 1, len(labels))
	assert.Equal(t, 2, len(metrics))
}

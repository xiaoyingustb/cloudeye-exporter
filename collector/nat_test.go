package collector

import (
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/stretchr/testify/assert"
)

func TestNatGetResourceInfo(t *testing.T) {
	conf.AccessKey = "test_ak"
	conf.SecretKey = "test_sk"
	conf.Region = "cn-test-01"
	metricConf = map[string]MetricConf{
		"SYS.NAT": {
			DimMetricName: map[string][]string{
				"nat_gateway_id": {"snat_connection"},
			},
		},
	}
	patches := gomonkey.ApplyFuncReturn(listResources, mockRmsResource(), nil)
	defer patches.Reset()

	var netgetter NATInfo
	labels, metrics := netgetter.GetResourceInfo()
	assert.Equal(t, 1, len(labels))
	assert.Equal(t, 1, len(metrics))
}

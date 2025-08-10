package collector

import (
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/stretchr/testify/assert"
)

func TestRdsGetResourceInfo(t *testing.T) {
	conf.AccessKey = "test_ak"
	conf.SecretKey = "test_sk"
	conf.Region = "cn-test-01"
	metricConf = map[string]MetricConf{
		"SYS.RDS": {
			DimMetricName: map[string][]string{
				"rds_cluster_id": {"rds001_cpu_util"},
			},
		},
	}
	volumes := mockRmsResource()
	volumes[0].Properties = map[string]interface{}{
		"cpu":        "2",
		"mem":        "16",
		"engineName": "mysql",
	}
	patches := gomonkey.ApplyFuncReturn(listResources, volumes, nil)
	defer patches.Reset()

	var rdsInfo RDSInfo
	labels, metrics := rdsInfo.GetResourceInfo()
	assert.Equal(t, 1, len(labels))
	assert.Equal(t, 1, len(metrics))
}

package collector

import (
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/cbr/v1/model"
	"github.com/stretchr/testify/assert"
)

func TestCbrGetResourceInfo(t *testing.T) {
	conf.AccessKey = "test_ak"
	conf.SecretKey = "test_sk"
	conf.Region = "cn-test-01"
	defaultEpId := "0"
	respPage1 := model.ListVaultResponse{
		HttpStatusCode: 200,
		Vaults: &[]model.Vault{
			{Id: "app-0001", Name: "app1", EnterpriseProjectId: &defaultEpId},
		},
	}
	respPage2 := model.ListVaultResponse{
		HttpStatusCode: 200,
		Vaults:         &[]model.Vault{},
	}
	metricConf = map[string]MetricConf{
		"SYS.CBR": {
			DimMetricName: map[string][]string{
				"instance_id": {"vault_util"},
			},
		},
	}

	cbrClient := getCBRClient()
	patches := gomonkey.ApplyFuncReturn(getResourceFromRMS, false)
	patches.ApplyMethodFunc(cbrClient, "ListVault", func(req *model.ListVaultRequest) (*model.ListVaultResponse, error) {
		if *req.Offset == 0 {
			return &respPage1, nil
		}
		return &respPage2, nil
	})
	defer patches.Reset()

	var cbrgetter CBRInfo
	labels, metrics := cbrgetter.GetResourceInfo()
	assert.Equal(t, 1, len(labels))
	assert.Equal(t, 1, len(metrics))
}

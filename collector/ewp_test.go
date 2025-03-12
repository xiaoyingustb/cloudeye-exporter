package collector

import (
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ces/v1/model"
	"github.com/stretchr/testify/assert"
)

func TestEWPInfo_GetResourceInfo(t *testing.T) {
	CloudConf.Global.RmsRetryTimes = 1
	conf.AccessKey = "test_ak"
	conf.SecretKey = "test_sk"
	var ewpGetter EWPInfo
	metricConf = map[string]MetricConf{
		"SYS.EWP": {
			DimMetricName: map[string][]string{"user_id": {"total_request_success_rate"}, "site_id": {"website_design_clicks"}},
		},
	}
	patches := gomonkey.NewPatches()
	defer patches.Reset()
	patches.ApplyFuncReturn(listResources, mockRmsResource(), nil)
	patches.ApplyFuncReturn(getUserInfoFromIAM, map[string]string{"test_domain_id": "test_domain_name"})
	patches.ApplyFuncReturn(listAllMetrics, []model.MetricInfoList{
		{
			Dimensions: []model.MetricsDimension{
				{
					Name:  "user_id",
					Value: "111111",
				},
			},
			Namespace:  "SYS.EWP",
			MetricName: "total_request_success_rate",
		},
		{
			Dimensions: []model.MetricsDimension{
				{
					Name:  "site_id",
					Value: "222222",
				},
			},
			Namespace:  "SYS.EWP",
			MetricName: "website_design_clicks",
		},
	}, nil)
	hcClient := &core.HcHttpClient{}
	patches.ApplyMethodReturn(hcClient, "Sync", &EwpSiteInfoResponse{}, nil)
	patches.ApplyFuncReturn(getHcClient, hcClient)
	info, lists := ewpGetter.GetResourceInfo()
	assert.NotNil(t, info)
	assert.NotNil(t, lists)
}

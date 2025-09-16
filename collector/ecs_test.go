package collector

import (
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ces/v1/model"
	ecsmodel "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ecs/v2/model"
	model2 "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/rms/v1/model"
	"github.com/stretchr/testify/assert"

	"github.com/huaweicloud/cloudeye-exporter/logs"
)

func TestECSInfo_GetResourceInfo(t *testing.T) {
	CloudConf.Global.RmsRetryTimes = 1
	conf.AccessKey = "test_ak"
	conf.SecretKey = "test_sk"
	patches := getPatches()
	logs.InitLog("")
	patches.ApplyFuncReturn(getAllServer, []EcsInstancesInfo{
		{
			ResourceBaseInfo: ResourceBaseInfo{
				ID:   "c53dbd97-c4ae-4ef7-8d12-fae2d3d38a80",
				Name: "nodelete-DRS-autotest-datamode1-test2",
				Tags: map[string]string{"wukong": "000"},
				EpId: "0",
			},
			IP: "192.168.20.53,100.93.4.203",
		},
	}, nil)
	patches.ApplyFuncReturn(getAllServerFromRMS, []EcsInstancesInfo{
		{
			ResourceBaseInfo: ResourceBaseInfo{
				ID:   "c53dbd97-c4ae-4ef7-8d12-fae2d3d38a80",
				Name: "nodelete-DRS-autotest-datamode1-test2",
				Tags: map[string]string{"wukong": "000"},
				EpId: "0",
			},
			IP: "192.168.20.53,100.93.4.203",
		},
	}, nil)
	patches.ApplyFunc(loadAgentDimensions, func(instanceID string) error { return nil })
	defer patches.Reset()
	ecsInfo1 := ECSInfo{}
	_, filterMetrics := ecsInfo1.GetResourceInfo()
	assert.NotNil(t, filterMetrics)
}

func TestAGTECSInfo_GetResourceInfo(t *testing.T) {
	CloudConf.Global.RmsRetryTimes = 1
	conf.AccessKey = "test_ak"
	conf.SecretKey = "test_sk"
	unit := "%"
	namespace := "AGT.ECS"
	metricInfosList := model.MetricInfoList{
		Unit:       unit,
		Namespace:  namespace,
		MetricName: "cpu_usage",
		Dimensions: []model.MetricsDimension{
			{
				Name:  "instance_id",
				Value: "9234ad9f-87a5-49a9-b6ec-11ecde1d943b",
			},
		},
	}
	patches := getPatches()
	logs.InitLog("")
	patches.ApplyFuncReturn(listAllMetrics, []model.MetricInfoList{metricInfosList}, nil)
	defer patches.Reset()
	agtEcsInfo1 := AGTECSInfo{}
	resourceInfo, _ := agtEcsInfo1.GetResourceInfo()
	assert.NotNil(t, resourceInfo)
}

func Test_getAllServerFromRMS(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()
	testId, testName, testEpId := "test_id", "test_name", "test_epId"
	patches.ApplyFuncReturn(listResources, []model2.ResourceEntity{
		{
			Id:   &testId,
			Name: &testName,
			EpId: &testEpId,
			Tags: map[string]string{},
			Properties: map[string]interface{}{
				"addresses": []map[string]string{
					{
						"addr":         "127.0.0.1",
						"osExtIpsType": "floating",
					},
					{
						"addr":         "127.0.0.1",
						"osExtIpsType": "fixed",
					},
				},
			},
		},
	}, nil)
	allServer, err := getAllServerFromRMS("ecs", "cloudservers")
	assert.Nil(t, err)
	assert.NotEmpty(t, allServer)
}

func Test_getIPFromEcsInfo(t *testing.T) {
	ipTypeEnum := ecsmodel.GetServerAddressOSEXTIPStypeEnum()
	addresses := map[string][]ecsmodel.ServerAddress{
		"address": {
			{
				Addr:         "127.0.0.1",
				OSEXTIPStype: &ipTypeEnum.FLOATING,
			},
			{
				Addr:         "127.0.0.1",
				OSEXTIPStype: &ipTypeEnum.FIXED,
			},
		},
	}
	ips, fixedIps, floatingIps := getIPFromEcsInfo(addresses)
	assert.NotEmpty(t, ips)
	assert.NotEmpty(t, fixedIps)
	assert.NotEmpty(t, floatingIps)
}

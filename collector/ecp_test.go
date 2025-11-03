package collector

import (
	"errors"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/rms/v1/model"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"

	"github.com/huaweicloud/cloudeye-exporter/logs"
)

func TestECPInfo_GetResourceInfo_configIsNil(t *testing.T) {
	patches := getEcpPatches()
	defer patches.Reset()
	logs.InitLog("")
	ecpInfoTest := ECPInfo{}
	labelInfos, filterMetrics := ecpInfoTest.GetResourceInfo()
	assert.Nil(t, labelInfos)
	assert.Nil(t, filterMetrics)
}

func TestECPInfo_GetResourceInfo_dimConfigIsNotExists(t *testing.T) {
	metricConfigMap := map[string][]string{}
	patches := getEcpPatches()
	defer patches.Reset()
	patches.ApplyFuncReturn(getMetricConfigMap, metricConfigMap)

	logs.InitLog("")
	ecpInfoTest := ECPInfo{}
	labelInfos, filterMetrics := ecpInfoTest.GetResourceInfo()
	assert.Nil(t, labelInfos)
	assert.Nil(t, filterMetrics)
}

func TestECPInfo_GetResourceInfo_dimConfigIsEmpty(t *testing.T) {
	patches := getEcpPatches()
	defer patches.Reset()

	metricConfigMap := map[string][]string{
		"instance_id": nil,
	}
	patches.ApplyFuncReturn(getMetricConfigMap, metricConfigMap)
	logs.InitLog("")
	ecpInfoTest := ECPInfo{}
	labelInfos, filterMetrics := ecpInfoTest.GetResourceInfo()
	assert.Nil(t, labelInfos)
	assert.Nil(t, filterMetrics)
}

func TestECPInfo_GetResourceInfo_getResourcesFromRMSFailed(t *testing.T) {
	patches := getEcpPatches()
	defer patches.Reset()

	metricConfigMap := map[string][]string{
		"instance_id": {"metric1", "metric2"},
	}
	patches.ApplyFuncReturn(getMetricConfigMap, metricConfigMap)
	patches.ApplyFuncReturn(listResources, nil, errors.New("test err"))

	logs.InitLog("")
	ecpInfoTest := ECPInfo{}
	labelInfos, filterMetrics := ecpInfoTest.GetResourceInfo()
	assert.Nil(t, labelInfos)
	assert.Nil(t, filterMetrics)
}

func TestECPInfo_GetResourceInfo_success(t *testing.T) {
	conf.AccessKey = "test_ak"
	conf.SecretKey = "test_sk"
	conf.Region = "cn-test-01"
	patches := getEcpPatches()
	defer patches.Reset()

	metricConf = map[string]MetricConf{
		"SYS.ECP": {
			DimMetricName: map[string][]string{
				"instance_id": {"metric1", "metric2"},
			},
		},
	}
	patches.ApplyFuncReturn(listResources, ecpResourceEntityInit(), nil)
	logs.InitLog("")
	ecpInfoTest := ECPInfo{}
	// 两个指标，两个资源
	labelInfos, filterMetrics := ecpInfoTest.GetResourceInfo()
	assert.NotNil(t, labelInfos)
	assert.Equal(t, 2, len(labelInfos))
	assert.NotNil(t, filterMetrics)
	assert.Equal(t, 4, len(filterMetrics))
}

func getEcpPatches() *gomonkey.Patches {
	confLoader := &logs.ConfLoader{}
	patches := gomonkey.ApplyMethodFunc(*confLoader, "LoadFile", func(fPath string, ecp interface{}) error {
		ecpTmp, _ := ecp.(*map[string][]logs.Config)
		ecpPointer := make(map[string][]logs.Config)
		ecpPointer["business"] = []logs.Config{
			{
				Level: zapcore.InfoLevel,
			},
		}
		*ecpTmp = ecpPointer
		return nil
	})
	return patches
}

func ecpResourceEntityInit() []model.ResourceEntity {
	id1 := "xxxx1"
	name1 := "test"
	epId1 := "1"
	epName1 := "测试企业1"
	checksum1 := "xxxb"
	create1 := "2025-10-13T02:08:52.000Z"
	update1 := "2025-10-13T02:08:52.000Z"
	provisioningState1 := "Succeeded"

	id2 := "xxxx2"
	name2 := "test2"
	epId2 := "2"
	epName2 := "测试企业2"
	checksum2 := "xxxc"
	create2 := "2025-10-13T02:08:52.000Z"
	update2 := "2025-10-13T02:08:52.000Z"
	provisioningState2 := "Succeeded"

	provider := "cph"
	typestr := "elasticcloudphones"
	regionId := "cn-north-7"
	projectId := "xxxx0"
	projectName := "cn-north-7"
	tags := map[string]string{}
	properties := map[string]interface{}{
		"availability_zone": "cn-north-7",
		"model_name":        "ecp1.std.a",
		"vpc_id":            "1775afc9-1b01-40XXXXXXXX",
		"keypair_name":      "KeyPair-XXX-2048",
		"subnet_id":         "c64400f2-7bcc-41cb-b09XXXXXXX",
		"image_id":          "510d6c14-ec80-4b53-bdd5-XXXXXXXX",
		"status":            "RUNNING",
	}

	response := []model.ResourceEntity{
		{
			Id:                &id1,
			Name:              &name1,
			Provider:          &provider,
			Type:              &typestr,
			RegionId:          &regionId,
			ProjectId:         &projectId,
			ProjectName:       &projectName,
			EpId:              &epId1,
			EpName:            &epName1,
			Checksum:          &checksum1,
			Created:           &create1,
			Updated:           &update1,
			ProvisioningState: &provisioningState1,
			Tags:              tags,
			Properties:        properties,
		},
		{
			Id:                &id2,
			Name:              &name2,
			Provider:          &provider,
			Type:              &typestr,
			RegionId:          &regionId,
			ProjectId:         &projectId,
			ProjectName:       &projectName,
			EpId:              &epId2,
			EpName:            &epName2,
			Checksum:          &checksum2,
			Created:           &create2,
			Updated:           &update2,
			ProvisioningState: &provisioningState2,
			Tags:              tags,
			Properties:        properties,
		},
	}
	return response
}

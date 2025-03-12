package collector

import (
	"errors"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/waf/v1/model"
	"github.com/stretchr/testify/assert"

	"github.com/huaweicloud/cloudeye-exporter/logs"
)

func TestWAFInfo_GetResourceInfo(t *testing.T) {
	tests := []struct {
		name     string
		wantNil  bool
		want1Nil bool
	}{
		{
			"getAllWafInstancesFromRMSErr",
			true,
			true,
		},
		{
			"normal",
			false,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			patches := gomonkey.NewPatches()
			if tt.name == "normal" {
				metricConf = map[string]MetricConf{
					"SYS.WAF": {
						Resource: "rms",
						DimMetricName: map[string][]string{
							"waf_instance_id": {"attacks"},
						},
					},
				}
				patches.ApplyFuncReturn(listResources, mockRmsResource(), nil)
				conf.AccessKey = "test_ak"
				conf.SecretKey = "test_sk"
				wafClient := getWAFClient()
				ID := "1"
				InstanceName := "AA"
				outputs := []gomonkey.OutputCell{
					{
						Values: gomonkey.Params{&model.ListInstanceResponse{
							HttpStatusCode: 200,
							Items: &[]model.ListInstance{
								{
									Id:           &ID,
									InstanceName: &InstanceName,
								},
							},
						}, nil},
					},
					{
						Values: gomonkey.Params{&model.ListInstanceResponse{
							HttpStatusCode: 200,
							Items:          &[]model.ListInstance{},
						}, nil},
					},
				}
				patches = patches.ApplyMethodSeq(wafClient, "ListInstance", outputs)
			}
			if tt.name == "getAllWafInstancesFromRMSErr" {
				patches = getPatches()

				sysConfig := map[string][]string{"waf_instance_id": {"attacks"}}
				patches.ApplyFuncReturn(getMetricConfigMap, sysConfig)
				patches.ApplyFuncReturn(listResources, mockRmsResource(), errors.New(""))

				logs.InitLog("")
			}
			getter := WAFInfo{}
			got, got1 := getter.GetResourceInfo()
			assert.Equalf(t, tt.wantNil, got == nil, "GetResourceInfo()")
			assert.Equalf(t, tt.want1Nil, got1 == nil, "GetResourceInfo()")
			getter.resetResourceInfo()
			if patches != nil {
				patches.Reset()
			}
		})
	}
}

func TestWafGetResourceInfo(t *testing.T) {
	conf.AccessKey = "test_ak"
	conf.SecretKey = "test_sk"
	metricConf = map[string]MetricConf{
		"SYS.WAF": {
			Resource: "rms",
			DimMetricName: map[string][]string{
				"waf_instance_id": {"attacks"},
			},
		},
	}
	patches := getPatches()
	defer patches.Reset()
	patches = patches.ApplyFuncReturn(listResources, mockRmsResource(), nil)
	wafClient := getWAFClient()
	ID := "1"
	InstanceName := "AA"
	outputs := []gomonkey.OutputCell{
		{
			Values: gomonkey.Params{&model.ListInstanceResponse{
				HttpStatusCode: 200,
				Items: &[]model.ListInstance{
					{
						Id:           &ID,
						InstanceName: &InstanceName,
					},
				},
			}, nil},
		},
		{
			Values: gomonkey.Params{&model.ListInstanceResponse{
				HttpStatusCode: 200,
				Items:          &[]model.ListInstance{},
			}, nil},
		},
	}
	patches = patches.ApplyMethodSeq(wafClient, "ListInstance", outputs)

	var wafgetter WAFInfo
	labels, metrics := wafgetter.GetResourceInfo()
	assert.Equal(t, 1, len(labels))
	assert.Equal(t, 1, len(metrics))
	wafgetter.resetResourceInfo()
}

func TestWafGetResourceInfo_getAllWafInstancesFromRMSErr(t *testing.T) {
	wafInfo.LabelInfo = nil
	wafInfo.FilterMetrics = nil
	conf.AccessKey = "test_ak"
	conf.SecretKey = "test_sk"
	patches := getPatches()
	defer patches.Reset()

	metricConf = map[string]MetricConf{
		"SYS.WAF": {
			Resource: "rms",
			DimMetricName: map[string][]string{
				"waf_instance_id": {"attacks"},
			},
		},
	}
	patches.ApplyFuncReturn(listResources, mockRmsResource(), errors.New(""))

	logs.InitLog("")
	var wafGetter WAFInfo
	labels, metrics := wafGetter.GetResourceInfo()
	assert.Nil(t, labels)
	assert.Nil(t, metrics)
	wafGetter.resetResourceInfo()
}

func TestWafGetResourceInfo_getAllPremiumWafInstancesNormal(t *testing.T) {
	conf.AccessKey = "test_ak"
	conf.SecretKey = "test_sk"
	patches := getPatches()
	defer patches.Reset()

	metricConf = map[string]MetricConf{
		"SYS.WAF": {
			Resource: "rms",
			DimMetricName: map[string][]string{
				"waf_instance_id": {"attacks"},
				"instance_id":     {"cpu_util"},
			},
		},
	}
	patches.ApplyFuncReturn(listResources, mockRmsResource(), nil)
	wafClient := getWAFClient()
	ID := "1"
	InstanceName := "AA"
	outputs := []gomonkey.OutputCell{
		{
			Values: gomonkey.Params{&model.ListInstanceResponse{
				HttpStatusCode: 200,
				Items: &[]model.ListInstance{
					{
						Id:           &ID,
						InstanceName: &InstanceName,
					},
				},
			}, nil},
		},
		{
			Values: gomonkey.Params{&model.ListInstanceResponse{
				HttpStatusCode: 200,
				Items:          &[]model.ListInstance{},
			}, nil},
		},
	}
	patches = patches.ApplyMethodSeq(wafClient, "ListInstance", outputs)
	logs.InitLog("")
	var wafGetter WAFInfo
	labels, metrics := wafGetter.GetResourceInfo()
	assert.Equal(t, 2, len(labels))
	assert.Equal(t, 2, len(metrics))
	wafGetter.resetResourceInfo()
}

func TestWafGetResourceInfo_getAllPremiumWafInstancesErr(t *testing.T) {
	wafInfo.LabelInfo = nil
	wafInfo.FilterMetrics = nil
	conf.AccessKey = "test_ak"
	conf.SecretKey = "test_sk"
	patches := getPatches()
	defer patches.Reset()

	metricConf = map[string]MetricConf{
		"SYS.WAF": {
			Resource: "rms",
			DimMetricName: map[string][]string{
				"waf_instance_id": {"attacks"},
				"instance_id":     {"cpu_util"},
			},
		},
	}
	patches.ApplyFuncReturn(listResources, mockRmsResource(), nil)
	wafClient := getWAFClient()
	ID := "1"
	InstanceName := "AA"
	resp := &model.ListInstanceResponse{
		HttpStatusCode: 200,
		Items: &[]model.ListInstance{
			{
				Id:           &ID,
				InstanceName: &InstanceName,
			},
		},
	}
	patches = gomonkey.ApplyMethodReturn(wafClient, "ListInstance", resp, errors.New(""))

	logs.InitLog("")
	var wafGetter WAFInfo
	labels, metrics := wafGetter.GetResourceInfo()
	assert.Nil(t, labels)
	assert.Nil(t, metrics)
	wafGetter.resetResourceInfo()
}

func Test_getWAFClient(t *testing.T) {
	conf.AccessKey = "test_ak"
	conf.SecretKey = "test_sk"
	tests := []struct {
		name    string
		wantNil bool
	}{
		{
			"normal",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.wantNil, getWAFClient() == nil, "getWAFClient()")
		})
	}
}

func Test_getAllPremiumWafInstances(t *testing.T) {
	conf.AccessKey = "test_ak"
	conf.SecretKey = "test_sk"
	ID := "1"
	InstanceName := "AA"
	tests := []struct {
		name    string
		wantNil bool
	}{
		{
			"normal",
			false,
		},
		{
			"HttpStatusCode ERR",
			true,
		},
		{
			"Request ERR",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			patches := getPatches()
			logs.InitLog("")
			if tt.name == "normal" {
				wafClient := getWAFClient()
				outputs := []gomonkey.OutputCell{
					{
						Values: gomonkey.Params{&model.ListInstanceResponse{
							HttpStatusCode: 200,
							Items: &[]model.ListInstance{
								{
									Id:           &ID,
									InstanceName: &InstanceName,
								},
							},
						}, nil},
					},
					{
						Values: gomonkey.Params{&model.ListInstanceResponse{
							HttpStatusCode: 200,
							Items:          &[]model.ListInstance{},
						}, nil},
					},
				}
				patches = patches.ApplyMethodSeq(wafClient, "ListInstance", outputs)
			}
			if tt.name == "HttpStatusCode ERR" {
				wafClient := getWAFClient()
				resp := &model.ListInstanceResponse{
					HttpStatusCode: 404,
					Items:          &[]model.ListInstance{},
				}
				patches = gomonkey.ApplyMethodReturn(wafClient, "ListInstance", resp, nil)
			}
			if tt.name == "Request ERR" {
				wafClient := getWAFClient()
				resp := &model.ListInstanceResponse{
					HttpStatusCode: 200,
					Items:          &[]model.ListInstance{},
				}
				patches = gomonkey.ApplyMethodReturn(wafClient, "ListInstance", resp, errors.New("aa"))
			}
			instances, _ := getAllPremiumWafInstances()
			assert.Equalf(t, tt.wantNil, len(instances) == 0, "getAllPremiumWafInstances()")
			if patches != nil {
				patches.Reset()
			}
		})
	}
}

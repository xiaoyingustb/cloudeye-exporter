package collector

import (
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/elb/v3/model"
	"github.com/stretchr/testify/assert"

	"github.com/huaweicloud/cloudeye-exporter/logs"
)

func TestELBInfo_GetResourceInfo(t *testing.T) {
	conf.AccessKey = "test_ak"
	conf.SecretKey = "test_sk"
	metricConf = map[string]MetricConf{
		"SYS.ELB": {
			DimMetricName: map[string][]string{"lbaas_instance_id": {"mc_l7_http_2xx"}, "lbaas_instance_id,lbaas_listener_id": {"mc_l7_http_2xx"},
				"lbaas_instance_id,lbaas_pool_id": {"m18_l7_upstream_2xx"}, "lbaas_instance_id,available_zone": {"mc_l7_http_2xx"}, "lbaas_instance_id,ip_address": {"upstream_tls_negotiation_error"}},
		},
	}
	patches := getPatches()
	logs.InitLog("")
	lbOutputs := []gomonkey.OutputCell{
		{
			Values: gomonkey.Params{&model.ListLoadBalancersResponse{
				Loadbalancers: &[]model.LoadBalancer{
					{Name: "test_lb", EnterpriseProjectId: "test_epId", VipAddress: "127.0.0.1", Ipv6VipAddress: "test_ipv6_addr", Provider: "test_provider",
						Listeners: []model.ListenerRef{{Id: "test_listener"}}, Pools: []model.PoolRef{{Id: "test_pool"}}, AvailabilityZoneList: []string{"cn-north-7"}},
				},
			}, nil},
		},
		{
			Values: gomonkey.Params{&model.ListLoadBalancersResponse{
				Loadbalancers: &[]model.LoadBalancer{},
			}, nil},
		},
	}
	patches.ApplyMethodSeq(getELBClient(), "ListLoadBalancers", lbOutputs)

	listenerOutputs := []gomonkey.OutputCell{
		{
			Values: gomonkey.Params{&model.ListListenersResponse{
				Listeners: &[]model.Listener{
					{Id: "test_listener", Name: "test_listener", Protocol: "http", ProtocolPort: 8087},
				},
			}, nil},
		},
		{
			Values: gomonkey.Params{&model.ListListenersResponse{
				Listeners: &[]model.Listener{},
			}, nil},
		},
	}
	patches.ApplyMethodSeq(getELBClient(), "ListListeners", listenerOutputs)

	poolOutputs := []gomonkey.OutputCell{
		{
			Values: gomonkey.Params{&model.ListPoolsResponse{
				Pools: &[]model.Pool{
					{Id: "test_pool", Name: "test_listener", Protocol: "http"},
				},
			}, nil},
		},
		{
			Values: gomonkey.Params{&model.ListPoolsResponse{
				Pools: &[]model.Pool{},
			}, nil},
		},
	}
	patches.ApplyMethodSeq(getELBClient(), "ListPools", poolOutputs)

	zoneOutputs := []gomonkey.OutputCell{
		{
			Values: gomonkey.Params{&model.ListAvailabilityZonesResponse{
				AvailabilityZones: &[][]model.AvailabilityZone{
					[]model.AvailabilityZone{
						{Code: "cn-north-7", State: "ACTIVE", Protocol: []string{"l4", "l7"}, PublicBorderGroup: "center", Category: 0},
					},
				},
			}, nil},
		},
		{
			Values: gomonkey.Params{&model.ListAvailabilityZonesResponse{
				AvailabilityZones: &[][]model.AvailabilityZone{},
			}, nil},
		},
	}
	patches.ApplyMethodSeq(getELBClient(), "ListAvailabilityZones", zoneOutputs)
	var elbGetter ELBInfo
	labelInfos, metricInfos := elbGetter.GetResourceInfo()
	assert.NotEmpty(t, labelInfos)
	assert.NotEmpty(t, metricInfos)
}

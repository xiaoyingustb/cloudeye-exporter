package collector

import (
	"fmt"
	"time"

	"github.com/huaweicloud/cloudeye-exporter/logs"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ces/v1/model"
)

var geipInfo serversInfo

type GeipInfo struct{}

type InternetBandwidthInfo struct {
	Size int64  `json:"size"`
	ID   string `json:"id"`
}

type GlobalConnectionBandwidthInfo struct {
	GcbID      string `json:"gcb_id"`
	Size       int64  `json:"size"`
	AdminState string `json:"admin_state"`
	GcbType    string `json:"gcb_type"`
}

type AssociationInstanceInfo struct {
	InstanceID   string `json:"instance_id"`
	ProjectID    string `json:"project_id"`
	Region       string `json:"region"`
	InstanceType string `json:"instance_type"`
}

type GlobalEipProperties struct {
	IpAddress                     string                        `json:"ip_address"`
	Status                        string                        `json:"status"`
	IpVersion                     int                           `json:"ip_version"`
	InternetBandwidthInfo         InternetBandwidthInfo         `json:"internet_bandwidth_info"`
	GlobalConnectionBandwidthInfo GlobalConnectionBandwidthInfo `json:"global_connection_bandwidth_info"`
	AssociationInstanceInfo       AssociationInstanceInfo       `json:"associate_instance_info"`
}

type GlobalEip struct {
	ResourceBaseInfo
	GlobalEipProperties
}

type InternetBandwidthProperties struct {
	ChargeMode string `json:"charge_mode"`
	Status     string `json:"status"`
}

type InternetBandwidth struct {
	ResourceBaseInfo
	InternetBandwidthProperties
}

func (getter GeipInfo) GetResourceInfo() (map[string]labelInfo, []model.MetricInfoList) {
	resourceInfos := map[string]labelInfo{}
	filterMetrics := make([]model.MetricInfoList, 0)
	geipInfo.Lock()
	defer geipInfo.Unlock()
	if geipInfo.LabelInfo == nil || time.Now().Unix() > geipInfo.TTL {
		err := buildGlobalEipsInfo(resourceInfos)
		if err != nil {
			logs.Logger.Error("Get all global eip info error:", err.Error())
			return geipInfo.LabelInfo, geipInfo.FilterMetrics
		}
		err = buildInternetBandwidthsInfo(resourceInfos)
		if err != nil {
			logs.Logger.Error("Get all internet bandwidth info error:", err.Error())
			return geipInfo.LabelInfo, geipInfo.FilterMetrics
		}
		allMetrics, err := listAllMetrics("SYS.GEIP")
		if err != nil {
			logs.Logger.Error("Get all global eip metrics error:", err.Error())
			return geipInfo.LabelInfo, geipInfo.FilterMetrics
		}
		for _, metricInfo := range allMetrics {
			resourceKey := GetResourceKeyFromMetricInfo(metricInfo)
			if _, ok := resourceInfos[resourceKey]; !ok {
				continue
			}

			if IsMetricInfoInWhiteList(metricInfo) {
				filterMetrics = append(filterMetrics, metricInfo)
			}
		}
		geipInfo.LabelInfo = resourceInfos
		geipInfo.FilterMetrics = filterMetrics
		geipInfo.TTL = time.Now().Add(GetResourceInfoExpirationTime()).Unix()
	}
	return geipInfo.LabelInfo, geipInfo.FilterMetrics
}

func buildGlobalEipsInfo(resourceInfos map[string]labelInfo) error {
	globalEips, err := getAllGlobalEIpsFromRMS()
	if err != nil {
		return err
	}
	for _, globalEip := range globalEips {
		info := labelInfo{
			Name: []string{"name", "epId", "ipAddress", "ipVersion", "associate_region", "associate_type"},
			Value: []string{globalEip.Name, globalEip.EpId, globalEip.IpAddress, fmt.Sprintf("%d", globalEip.IpVersion), globalEip.AssociationInstanceInfo.Region,
				globalEip.AssociationInstanceInfo.InstanceType},
		}
		keys, values := getTags(globalEip.Tags)
		info.Name = append(info.Name, keys...)
		info.Value = append(info.Value, values...)
		resourceInfos[globalEip.ID] = info
	}
	return nil
}

func buildInternetBandwidthsInfo(resourceInfos map[string]labelInfo) error {
	bandwidths, err := getAllInternetBandwidthsFromRMS()
	if err != nil {
		return err
	}
	for _, bandwidth := range bandwidths {
		info := labelInfo{
			Name:  []string{"name", "epId", "status", "chargeMode"},
			Value: []string{bandwidth.Name, bandwidth.EpId, bandwidth.Status, bandwidth.ChargeMode},
		}
		keys, values := getTags(bandwidth.Tags)
		info.Name = append(info.Name, keys...)
		info.Value = append(info.Value, values...)
		resourceInfos[bandwidth.ID] = info
	}
	return nil
}

func getAllInternetBandwidthsFromRMS() ([]InternetBandwidth, error) {
	resp, err := listResources("vpc", "internet-bandwidths", "global")
	if err != nil {
		logs.Logger.Errorf("Failed to list resource of vpc.internet-bandwidths, error: %s", err.Error())
		return nil, err
	}
	bandwidths := make([]InternetBandwidth, 0, len(resp))
	for _, resource := range resp {
		var bandwidthProperties InternetBandwidthProperties
		err := fmtResourceProperties(resource.Properties, &bandwidthProperties)
		if err != nil {
			logs.Logger.Errorf("fmt internet-bandwidths properties error: %s", err.Error())
			continue
		}
		bandwidths = append(bandwidths, InternetBandwidth{
			ResourceBaseInfo: ResourceBaseInfo{
				ID:   *resource.Id,
				Name: *resource.Name,
				EpId: *resource.EpId,
				Tags: resource.Tags,
			},
			InternetBandwidthProperties: bandwidthProperties,
		})
	}
	return bandwidths, nil
}

func getAllGlobalEIpsFromRMS() ([]GlobalEip, error) {
	resp, err := listResources("vpc", "global-eips", "global")
	if err != nil {
		logs.Logger.Errorf("Failed to list resource of vpc.global-eips, error: %s", err.Error())
		return nil, err
	}
	globalEips := make([]GlobalEip, 0, len(resp))
	for _, resource := range resp {
		var globalEipProperties GlobalEipProperties
		err := fmtResourceProperties(resource.Properties, &globalEipProperties)
		if err != nil {
			logs.Logger.Errorf("fmt global-eips properties error: %s", err.Error())
			continue
		}
		globalEips = append(globalEips, GlobalEip{
			ResourceBaseInfo: ResourceBaseInfo{
				ID:   *resource.Id,
				Name: *resource.Name,
				EpId: *resource.EpId,
				Tags: resource.Tags},
			GlobalEipProperties: globalEipProperties,
		})
	}
	return globalEips, nil
}

package collector

import (
	"fmt"
	"time"

	cesmodel "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ces/v1/model"

	"github.com/huaweicloud/cloudeye-exporter/logs"
)

var (
	egInfo         serversInfo
	egDimConfigMap = map[string]map[string]string{
		"channel_id": {
			"resourceType":      "channels",
			"providerId":        "eg",
			"nameInSubResource": "channel_name",
		},
		"source_name": {
			"resourceType":      "sources",
			"providerId":        "eg",
			"nameInSubResource": "source_name",
		},
		"subscription_id": {
			"resourceType":      "subscriptions",
			"providerId":        "eg",
			"nameInSubResource": "subscription_name",
		},
		"streaming_id": {
			"resourceType":      "streamings",
			"providerId":        "eg",
			"nameInSubResource": "streaming_name",
		},
	}
)

type EgInfo struct{}

func (getter EgInfo) GetResourceInfo() (map[string]labelInfo, []cesmodel.MetricInfoList) {
	egInfo.Lock()
	defer egInfo.Unlock()
	if egInfo.LabelInfo == nil || time.Now().Unix() > egInfo.TTL {
		resourceInfos := map[string]labelInfo{}
		allMetrics, err := listAllMetrics("SYS.EG")
		if err != nil {
			logs.Logger.Errorf("Get all eg metrics: %s", err.Error())
			return egInfo.LabelInfo, egInfo.FilterMetrics
		}
		egResourceMetricMap := getEgSubResourceMetrics(allMetrics)
		if err = getEgResourceInfoFromRms("source_name", resourceInfos, egResourceMetricMap); err != nil {
			logs.Logger.Errorf("Failed to query eg source info from rms: %s", err.Error())
			return egInfo.LabelInfo, egInfo.FilterMetrics
		}
		if err = getEgResourceInfoFromRms("channel_id", resourceInfos, egResourceMetricMap); err != nil {
			logs.Logger.Errorf("Failed to query eg channel info from rms: %s", err.Error())
			return egInfo.LabelInfo, egInfo.FilterMetrics
		}
		if err = getEgResourceInfoFromRms("subscription_id", resourceInfos, egResourceMetricMap); err != nil {
			logs.Logger.Errorf("Failed to query eg subscription info from rms: %s", err.Error())
			return egInfo.LabelInfo, egInfo.FilterMetrics
		}
		if err = getEgResourceInfoFromRms("streaming_id", resourceInfos, egResourceMetricMap); err != nil {
			logs.Logger.Errorf("Failed to query eg streaming info from rms: %s", err.Error())
			return egInfo.LabelInfo, egInfo.FilterMetrics
		}

		var filteredMetrics []cesmodel.MetricInfoList
		for _, metricInfo := range allMetrics {
			if !IsMetricInfoInWhiteList(metricInfo) {
				continue
			}
			if _, ok := resourceInfos[GetResourceKeyFromMetricInfo(metricInfo)]; ok {
				filteredMetrics = append(filteredMetrics, metricInfo)
			}
		}
		egInfo.LabelInfo = resourceInfos
		egInfo.FilterMetrics = filteredMetrics
		egInfo.TTL = time.Now().Add(GetResourceInfoExpirationTime()).Unix()
	}
	return egInfo.LabelInfo, egInfo.FilterMetrics
}

func getEgResourceInfoFromRms(dimName string, resourceInfos map[string]labelInfo, allMetricMap map[string][]cesmodel.MetricInfoList) error {
	dimConfig, ok := egDimConfigMap[dimName]
	if !ok {
		return fmt.Errorf("dimension is not in metric config map: %s", dimName)
	}
	rmsResourceInfos, err := getResourcesBaseInfoFromRMS(dimConfig["providerId"], dimConfig["resourceType"])
	if err != nil {
		return err
	}
	for _, rmsResourceInfo := range rmsResourceInfos {
		info := labelInfo{
			Name:  []string{"name", "epId"},
			Value: []string{rmsResourceInfo.Name, rmsResourceInfo.EpId},
		}
		keys, values := getTags(rmsResourceInfo.Tags)
		info.Name = append(info.Name, keys...)
		info.Value = append(info.Value, values...)
		resourceInfos[rmsResourceInfo.ID] = info

		// 设置子维度指标资源关联
		subMetricList, ok := allMetricMap[rmsResourceInfo.ID]
		if !ok {
			continue
		}
		for _, subMetric := range subMetricList {
			resourceInfos[GetResourceKeyFromMetricInfo(subMetric)] = labelInfo{
				Name:  []string{dimConfig["nameInSubResource"], "epId"},
				Value: []string{rmsResourceInfo.Name, rmsResourceInfo.EpId},
			}
		}
	}
	return nil
}

func getEgSubResourceMetrics(allMetrics []cesmodel.MetricInfoList) map[string][]cesmodel.MetricInfoList {
	resultMap := make(map[string][]cesmodel.MetricInfoList)
	for _, metric := range allMetrics {
		if len(metric.Dimensions) < 2 {
			continue
		}
		firstDimValue := metric.Dimensions[0].Value
		subMetricList, ok := resultMap[firstDimValue]
		if ok {
			subMetricList = append(subMetricList, metric)
		} else {
			subMetricList = []cesmodel.MetricInfoList{metric}
		}
		resultMap[firstDimValue] = subMetricList
	}
	return resultMap
}

package collector

import (
	"fmt"
	"time"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ces/v1/model"
	cesmodel "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ces/v1/model"

	"github.com/huaweicloud/cloudeye-exporter/logs"
)

var (
	ancInfo      serversInfo
	dimConfigMap = map[string]map[string]string{
		"anc_anc_id": {
			"resourceType":      "ancs",
			"providerId":        "anc",
			"nameInSubResource": "anc_name",
		},
		"anc_service_id": {
			"resourceType":      "services",
			"providerId":        "anc",
			"nameInSubResource": "anc_service_name",
		},
		"anc_membergroup_id": {
			"resourceType":      "member-groups",
			"providerId":        "anc",
			"nameInSubResource": "member_group_name",
		},
	}
)

type ANCInfo struct{}

func (getter ANCInfo) GetResourceInfo() (map[string]labelInfo, []cesmodel.MetricInfoList) {

	ancInfo.Lock()
	defer ancInfo.Unlock()
	if ancInfo.LabelInfo == nil || time.Now().Unix() > ancInfo.TTL {
		resourceInfos := map[string]labelInfo{}
		currentAllMetrics, err := listAllMetrics("SYS.ANC")
		if err != nil {
			logs.Logger.Errorf("Get all anc metrics: %s", err.Error())
			return ancInfo.LabelInfo, ancInfo.FilterMetrics
		}
		var filteredMetrics []model.MetricInfoList
		for _, metricInfo := range currentAllMetrics {
			if IsMetricInfoInWhiteList(metricInfo) {
				filteredMetrics = append(filteredMetrics, metricInfo)
			}
		}
		ancResourceMetricMap := getAncSubResourceMetrics(filteredMetrics)
		if err = getAncResourceInfoFromRms("anc_anc_id", resourceInfos, ancResourceMetricMap); err != nil {
			logs.Logger.Errorf("Failed to get anc resource from rms: %s", err.Error())
			return ancInfo.LabelInfo, ancInfo.FilterMetrics
		}
		if err = getAncResourceInfoFromRms("anc_service_id", resourceInfos, ancResourceMetricMap); err != nil {
			logs.Logger.Errorf("Failed to get anc service resource from rms: %s", err.Error())
			return ancInfo.LabelInfo, ancInfo.FilterMetrics
		}
		if err = getAncResourceInfoFromRms("anc_membergroup_id", resourceInfos, ancResourceMetricMap); err != nil {
			logs.Logger.Errorf("Failed to get anc member group resource from rms: %s", err.Error())
			return ancInfo.LabelInfo, ancInfo.FilterMetrics
		}
		ancInfo.LabelInfo = resourceInfos
		ancInfo.FilterMetrics = filteredMetrics
		ancInfo.TTL = time.Now().Add(GetResourceInfoExpirationTime()).Unix()
	}
	return ancInfo.LabelInfo, ancInfo.FilterMetrics
}

func getAncResourceInfoFromRms(dimName string, resourceInfos map[string]labelInfo, allMetricMap map[string][]cesmodel.MetricInfoList) error {
	dimConfig, ok := dimConfigMap[dimName]
	if !ok {
		return fmt.Errorf("dim name not in dimension config file: %s", dimName)
	}
	rmsResourceInfos, err := getResourcesBaseInfoFromRMS(dimConfig["providerId"], dimConfig["resourceType"], "global")
	if err != nil {
		return err
	}
	metricNames, ok := getMetricConfigMap("SYS.ANC")[dimName]
	if !ok {
		return fmt.Errorf("dim name not in metric config file: %s", dimName)
	}
	for _, rmsResourceInfo := range rmsResourceInfos {
		metrics := buildSingleDimensionMetrics(metricNames, "SYS.ANC", dimName, rmsResourceInfo.ID)
		info := labelInfo{
			Name:  []string{"name", "epId"},
			Value: []string{rmsResourceInfo.Name, rmsResourceInfo.EpId},
		}
		keys, values := getTags(rmsResourceInfo.Tags)
		info.Name = append(info.Name, keys...)
		info.Value = append(info.Value, values...)
		resourceInfos[GetResourceKeyFromMetricInfo(metrics[0])] = info

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

func getAncSubResourceMetrics(allMetrics []cesmodel.MetricInfoList) map[string][]cesmodel.MetricInfoList {
	resultMap := make(map[string][]cesmodel.MetricInfoList)
	for _, metric := range allMetrics {
		// 只操作多维度的指标，将其第一维度值映射到指标对象
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

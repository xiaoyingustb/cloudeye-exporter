package collector

import (
	"github.com/huaweicloud/cloudeye-exporter/logs"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ces/v1/model"
	"time"
)

type ECPInfo struct{}

var ecpInfo serversInfo

func (getter ECPInfo) GetResourceInfo() (map[string]labelInfo, []model.MetricInfoList) {
	resourceInfos := map[string]labelInfo{}
	filterMetrics := make([]model.MetricInfoList, 0)
	ecpInfo.Lock()
	defer ecpInfo.Unlock()
	if ecpInfo.LabelInfo == nil || time.Now().Unix() > ecpInfo.TTL {
		ecpConfigMap := getMetricConfigMap("SYS.ECP")
		if ecpConfigMap == nil {
			logs.Logger.Warn("Metric config is nil.")
			return ecpInfo.LabelInfo, ecpInfo.FilterMetrics
		}

		if _, ok := ecpConfigMap["instance_id"]; !ok {
			logs.Logger.Warn("Metric config is nil of SYS.ecp of elastic_cloud_phone_id.")
			return ecpInfo.LabelInfo, ecpInfo.FilterMetrics
		}

		metricNames := ecpConfigMap["instance_id"]
		if len(metricNames) == 0 {
			logs.Logger.Warn("Metric config is empty of SYS.ecp of elastic_cloud_phone_id.")
			return ecpInfo.LabelInfo, ecpInfo.FilterMetrics
		}

		servers, err := getResourcesBaseInfoFromRMS("cph", "elasticcloudphones")
		if err != nil {
			logs.Logger.Errorf("Get resource base info from RMS Server error:", err.Error())
			return ecpInfo.LabelInfo, ecpInfo.FilterMetrics
		}

		for _, server := range servers {
			metrics := buildSingleDimensionMetrics(metricNames, "SYS.ECP", "instance_id", server.ID)
			filterMetrics = append(filterMetrics, metrics...)
			info := labelInfo{
				Name:  []string{"name", "epId"},
				Value: []string{server.Name, server.EpId},
			}
			keys, values := getTags(server.Tags)
			info.Name = append(info.Name, keys...)
			info.Value = append(info.Value, values...)
			resourceInfos[GetResourceKeyFromMetricInfo(metrics[0])] = info
		}
		ecpInfo.LabelInfo = resourceInfos
		ecpInfo.FilterMetrics = filterMetrics
		ecpInfo.TTL = time.Now().Add(GetResourceInfoExpirationTime()).Unix()
	}
	return ecpInfo.LabelInfo, ecpInfo.FilterMetrics
}

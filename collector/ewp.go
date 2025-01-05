package collector

import (
	"fmt"
	"net/http"
	"time"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/def"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ces/v1/model"

	"github.com/huaweicloud/cloudeye-exporter/logs"
)

type EwpSiteInfoRequest struct{}

type EwpSiteInfoResponse struct {
	DisplayName    string `json:"displayName"`
	HttpStatusCode int    `json:"-"`
}

type EwpDomainInfoRequest struct{}

type EwpDomainInfo struct {
	Name      string `json:"name"`
	XDomainID string `json:"xdomain_id"`
}

type EwpDomainInfoResponse struct {
	Domains        []EwpDomainInfo `json:"domains"`
	HttpStatusCode int             `json:"-"`
}

var ewpInfo serversInfo

type EWPInfo struct{}

func (getter EWPInfo) GetResourceInfo() (map[string]labelInfo, []model.MetricInfoList) {
	resourceInfos := map[string]labelInfo{}
	filterMetrics := make([]model.MetricInfoList, 0)
	ewpInfo.Lock()
	defer ewpInfo.Unlock()
	if ewpInfo.LabelInfo == nil || time.Now().Unix() > ewpInfo.TTL {
		sysConfigMap := getMetricConfigMap("SYS.EWP")
		if sysConfigMap == nil {
			logs.Logger.Warn("Metric config is nil.")
			return ewpInfo.LabelInfo, ewpInfo.FilterMetrics
		}
		allMetrics, listMetricErr := listAllMetrics("SYS.EWP")
		if listMetricErr != nil {
			logs.Logger.Errorf("Failed to get site resources, error: %s", listMetricErr.Error())
			return nil, nil
		}
		metricGroups := getEwpMetricsGroups(allMetrics)
		siteMetrics, ok := metricGroups["site_id"]
		var siteIds []string
		if ok {
			siteIdMap := make(map[string]bool, 0)
			for _, siteMetric := range siteMetrics {
				siteId := siteMetric.Dimensions[0].Value
				if _, inMap := siteIdMap[siteId]; inMap {
					continue
				}
				siteIdMap[siteId] = true
				siteIds = append(siteIds, siteMetric.Dimensions[0].Value)
			}
		}

		userInfoMap := getUserInfoFromIAM()
		for domainId, domainName := range userInfoMap {
			if metricNames, ok := sysConfigMap["user_id"]; ok {
				metrics := buildSingleDimensionMetrics(metricNames, "SYS.EWP", "user_id", domainId)
				filterMetrics = append(filterMetrics, metrics...)
				info := labelInfo{
					Name:  []string{"name"},
					Value: []string{domainName},
				}
				resourceInfos[GetResourceKeyFromMetricInfo(metrics[0])] = info
			}
		}

		siteInfos, siErr := getAllSiteInfos(siteIds)
		if siErr != nil {
			logs.Logger.Errorf("Get site infos error: %s", siErr.Error())
			return ewpInfo.LabelInfo, ewpInfo.FilterMetrics
		}
		for _, siteInfo := range siteInfos {
			if metricNames, ok := sysConfigMap["site_id"]; ok {
				metrics := buildSingleDimensionMetrics(metricNames, "SYS.EWP", "site_id", siteInfo.ID)
				filterMetrics = append(filterMetrics, metrics...)
				info := labelInfo{
					Name:  []string{"name"},
					Value: []string{siteInfo.Name},
				}
				resourceInfos[GetResourceKeyFromMetricInfo(metrics[0])] = info
			}
		}

		ewpInfo.LabelInfo = resourceInfos
		ewpInfo.FilterMetrics = filterMetrics
		ewpInfo.TTL = time.Now().Add(GetResourceInfoExpirationTime()).Unix()
	}
	return ewpInfo.LabelInfo, ewpInfo.FilterMetrics
}

func getEwpMetricsGroups(metrics []model.MetricInfoList) map[string][]model.MetricInfoList {
	resultMap := make(map[string][]model.MetricInfoList, 0)
	for _, metric := range metrics {
		if len(metric.Dimensions) < 1 {
			continue
		}
		dimName := metric.Dimensions[0].Name
		if metricInfoList, ok := resultMap[dimName]; ok {
			metricInfoList = append(metricInfoList, metric)
			resultMap[dimName] = metricInfoList
		} else {
			resultMap[dimName] = []model.MetricInfoList{metric}
		}
	}
	return resultMap
}

func getAllSiteInfos(instanceIds []string) ([]ResourceBaseInfo, error) {
	var results []ResourceBaseInfo
	for _, instanceId := range instanceIds {
		urlPath := fmt.Sprintf("/v1/cloudsite/%s/display-name", instanceId)
		requestDef := def.NewHttpRequestDefBuilder().
			WithMethod(http.MethodGet).
			WithPath(urlPath).
			WithResponse(new(EwpSiteInfoResponse)).
			WithContentType("application/json").Build()
		req := EwpSiteInfoRequest{}
		resp, err := getHcClient(getEndpoint("ewp", "v1")).Sync(req, requestDef)
		if err != nil {
			logs.Logger.Errorf("Get ewp site display name error: %s", err.Error())
			return nil, err
		}
		siteInfo, ok := resp.(*EwpSiteInfoResponse)
		if !ok {
			logs.Logger.Errorf("Get all ewp site convert to EwpSiteInfoResponse failed!")
			return results, fmt.Errorf("convert to EwpSiteInfoResponse failed")
		}
		results = append(results, ResourceBaseInfo{ID: instanceId, Name: siteInfo.DisplayName})
	}

	return results, nil
}

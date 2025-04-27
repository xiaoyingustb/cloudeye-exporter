package collector

import (
	"time"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/cc/v3/model"
	cesmodel "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ces/v1/model"

	"github.com/huaweicloud/cloudeye-exporter/logs"
)

var gcbInfo serversInfo

type GCBInfo struct{}

func (getter GCBInfo) GetResourceInfo() (map[string]labelInfo, []cesmodel.MetricInfoList) {
	gcbInfo.Lock()
	defer gcbInfo.Unlock()
	if gcbInfo.LabelInfo == nil || time.Now().Unix() > gcbInfo.TTL {
		resourceInfos := map[string]labelInfo{}
		filterMetrics := make([]cesmodel.MetricInfoList, 0)
		globalConnectionBandwidths, err := listGcbInfos()
		if err != nil {
			logs.Logger.Errorf("Get gcb infos error: %s", err.Error())
			return gcbInfo.LabelInfo, gcbInfo.FilterMetrics
		}
		for _, globalConnectionBandwidth := range globalConnectionBandwidths {
			info := labelInfo{
				Name:  []string{"name", "epId"},
				Value: []string{*globalConnectionBandwidth.Name, *globalConnectionBandwidth.EnterpriseProjectId},
			}
			keys, values := getTags(fmtTags(globalConnectionBandwidth.Tags))
			info = appendNameValuePairToLabelInfo(info, keys, values)
			resourceInfos[globalConnectionBandwidth.Id] = info
		}
		allMetrics, err := listAllMetrics("SYS.GCB")
		if err != nil {
			logs.Logger.Errorf("List all metrics for SYS.GCB error: %s", err.Error())
			return gcbInfo.LabelInfo, gcbInfo.FilterMetrics
		}
		for _, metric := range allMetrics {
			resourceKey := GetResourceKeyFromMetricInfo(metric)
			if resourceKey == "" {
				continue
			}
			if _, ok := resourceInfos[resourceKey]; !ok {
				continue
			}

			if IsMetricInfoInWhiteList(metric) {
				filterMetrics = append(filterMetrics, metric)
			}
		}
		gcbInfo.LabelInfo = resourceInfos
		gcbInfo.FilterMetrics = filterMetrics
		gcbInfo.TTL = time.Now().Add(GetResourceInfoExpirationTime()).Unix()
	}
	return gcbInfo.LabelInfo, gcbInfo.FilterMetrics
}

func listGcbInfos() ([]model.GlobalConnectionBandwidth, error) {
	client := getCCClient()
	epIds := getEpIdRequestPart()
	request := &model.ListGlobalConnectionBandwidthsRequest{Limit: &limit, EnterpriseProjectId: &epIds}
	var allBandWidths []model.GlobalConnectionBandwidth
	for {
		response, err := client.ListGlobalConnectionBandwidths(request)
		if err != nil {
			return allBandWidths, err
		}
		if len(response.GlobalconnectionBandwidths) == 0 {
			break
		}
		allBandWidths = append(allBandWidths, response.GlobalconnectionBandwidths...)
		if response.PageInfo.NextMarker == nil {
			break
		}
		request.Marker = response.PageInfo.NextMarker
	}
	return allBandWidths, nil
}

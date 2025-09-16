package collector

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/def"
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
			for _, directionalConnection := range globalConnectionBandwidth.DirectionalConnections {
				dcInfo := labelInfo{
					Name:  []string{"directional_connection_name"},
					Value: []string{directionalConnection.Name},
				}
				dcInfo.Name = append(dcInfo.Name, info.Name...)
				dcInfo.Value = append(dcInfo.Value, info.Value...)
				resourceInfos[fmt.Sprintf("%s.%s", directionalConnection.Id, globalConnectionBandwidth.Id)] = dcInfo
			}
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

type ListGlobalConnectionBandwidthsResponse struct {
	// 资源ID标识符。
	RequestId string          `json:"request_id"`
	PageInfo  *model.PageInfo `json:"page_info,omitempty"`
	// 全域互联带宽列表响应体。
	GlobalconnectionBandwidths []GlobalConnectionBandwidth `json:"globalconnection_bandwidths"`
	HttpStatusCode             int                         `json:"-"`
}

type GlobalConnectionBandwidth struct {
	Id                     string                  `json:"id"`
	Name                   *string                 `json:"name,omitempty"`
	EnterpriseProjectId    *string                 `json:"enterprise_project_id,omitempty"`
	DirectionalConnections []DirectionalConnection `json:"directional_connections"`
	Tags                   *[]Tag                  `json:"tags,omitempty"`
}

type DirectionalConnection struct {
	Name string `json:"name"`
	Id   string `json:"id"`
}

func GenReqDefForListGlobalConnectionBandwidths() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodGet).
		WithPath("/v3/{domain_id}/gcb/gcbandwidths").
		WithResponse(new(ListGlobalConnectionBandwidthsResponse)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithLocationType(def.Query).
		WithName("EnterpriseProjectId").
		WithJsonTag("enterprise_project_id"))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithLocationType(def.Query).
		WithName("Limit").
		WithJsonTag("limit"))
	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithLocationType(def.Query).
		WithName("Marker").
		WithJsonTag("marker"))
	requestDef := reqDefBuilder.Build()
	return requestDef
}

func listGcbInfos() ([]GlobalConnectionBandwidth, error) {
	client := getCCClient()
	epIds := getEpIdRequestPart()
	request := &model.ListGlobalConnectionBandwidthsRequest{Limit: &limit, EnterpriseProjectId: &epIds}
	var allBandWidths []GlobalConnectionBandwidth
	for {
		requestDef := GenReqDefForListGlobalConnectionBandwidths()
		response, err := client.HcClient.Sync(request, requestDef)
		if err != nil {
			return allBandWidths, err
		}
		gcbResponse, ok := response.(*ListGlobalConnectionBandwidthsResponse)
		if !ok {
			return nil, errors.New("resp type is not ListGlobalConnectionBandwidthsResponse")
		}
		if len(gcbResponse.GlobalconnectionBandwidths) == 0 {
			break
		}
		allBandWidths = append(allBandWidths, gcbResponse.GlobalconnectionBandwidths...)
		if gcbResponse.PageInfo.NextMarker == nil {
			break
		}
		request.Marker = gcbResponse.PageInfo.NextMarker
	}
	return allBandWidths, nil
}

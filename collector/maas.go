package collector

import (
	"fmt"
	"net/http"
	"time"

	"github.com/huaweicloud/cloudeye-exporter/logs"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/def"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ces/v1/model"
)

var maasInfo serversInfo

type MaaSInfo struct{}

type ListMaasIntanceRequest struct {
	Body ListMaasInstanceRequestBody `json:"body,omitempty"`
}

type ListMaasInstanceRequestBody struct {
	Namespace string                          `json:"namespace,omitempty"`
	Start     int32                           `json:"start,omitempty"`
	Limit     int32                           `json:"limit,omitempty"`
	Query     []ListMaasInstanceSubDimReqInfo `json:"query,omitempty"`
}

type ListMaasInstanceSubDimReqInfo struct {
	DimName string `json:"dim_name,omitempty"`
}

type ListMaasInstanceResp struct {
	Instances      []ListMaasInstanceRespItem `json:"instances,omitempty"`
	Total          int32                      `json:"total,omitempty"`
	HttpStatusCode int                        `json:"-"`
}

type ListMaasInstanceRespItem struct {
	MaasApiIdDimInfo ListMaasInstanceRespDimInfo `json:"maas_api_id,omitempty"`
	MaasKeyIdDimInfo ListMaasInstanceRespDimInfo `json:"maas_key_id,omitempty"`
}

type ListMaasInstanceRespDimInfo struct {
	Name string `json:"name,omitempty"`
	Id   string `json:"id,omitempty"`
}

func GenReqDefForListMaasInstanceData() *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().
		WithMethod(http.MethodPost).
		WithPath("/v1/{project_id}/maas/monitor/instances").
		WithResponse(new(ListMaasInstanceResp)).
		WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().
		WithName("Body").
		WithLocationType(def.Body))

	requestDef := reqDefBuilder.Build()
	return requestDef
}

func buildMaasApiResourceInfo(resourceInfos map[string]labelInfo) error {
	options := &ListMaasIntanceRequest{
		Body: ListMaasInstanceRequestBody{
			Namespace: "SYS.MaaS",
			Start:     0,
			Limit:     100,
			Query: []ListMaasInstanceSubDimReqInfo{
				{
					"maas_api_id",
				},
			},
		},
	}
	maasClient := getHcClient(getEndpoint("maas", "v1.0"))
	req := GenReqDefForListMaasInstanceData()
	for {
		resp, err := maasClient.Sync(options, req)
		if err != nil {
			logs.Logger.Errorf("List maas instance error: %s", err.Error())
			return err
		}
		maasResp, ok := resp.(*ListMaasInstanceResp)
		if !ok || maasResp == nil {
			return fmt.Errorf("convert response to ListMaasInstanceResp failed")
		}
		if len(maasResp.Instances) == 0 {
			break
		}
		for _, instance := range maasResp.Instances {
			apiIdDimInfo := instance.MaasApiIdDimInfo
			resourceInfos[apiIdDimInfo.Id] = labelInfo{
				Name:  []string{"maas_api_name"},
				Value: []string{apiIdDimInfo.Name},
			}
		}
		options.Body.Start += 100
	}
	return nil
}

func buildMaasKeyResourceInfo(resourceInfos map[string]labelInfo) error {
	options := &ListMaasIntanceRequest{
		Body: ListMaasInstanceRequestBody{
			Namespace: "SYS.MaaS",
			Start:     0,
			Limit:     100,
			Query: []ListMaasInstanceSubDimReqInfo{
				{
					"maas_api_id",
				},
				{
					"maas_key_id",
				},
			},
		},
	}
	maasClient := getHcClient(getEndpoint("maas", "v1.0"))
	req := GenReqDefForListMaasInstanceData()
	for {
		resp, err := maasClient.Sync(options, req)
		if err != nil {
			logs.Logger.Errorf("List maas instance error: %s", err.Error())
			return err
		}
		maasResp, ok := resp.(*ListMaasInstanceResp)
		if !ok || maasResp == nil {
			return fmt.Errorf("convert response to ListMaasInstanceResp failed")
		}
		if len(maasResp.Instances) == 0 {
			break
		}
		for _, instance := range maasResp.Instances {
			apiIdDimInfo := instance.MaasApiIdDimInfo
			keyIdDimInfo := instance.MaasKeyIdDimInfo
			resourceInfos[fmt.Sprintf("%s.%s", apiIdDimInfo.Id, keyIdDimInfo.Id)] = labelInfo{
				Name:  []string{"maas_api_name", "maas_key_name"},
				Value: []string{apiIdDimInfo.Name, keyIdDimInfo.Name},
			}
		}
		options.Body.Start += 100
	}
	return nil
}

func (getter MaaSInfo) GetResourceInfo() (map[string]labelInfo, []model.MetricInfoList) {
	maasInfo.Lock()
	defer maasInfo.Unlock()
	if maasInfo.LabelInfo == nil || time.Now().Unix() > maasInfo.TTL {
		resourceInfos := make(map[string]labelInfo)
		err := buildMaasApiResourceInfo(resourceInfos)
		if err != nil {
			logs.Logger.Errorf("List maas api instances error: %s", err.Error())
			return maasInfo.LabelInfo, maasInfo.FilterMetrics
		}
		err = buildMaasKeyResourceInfo(resourceInfos)
		if err != nil {
			logs.Logger.Errorf("List maas key instances error: %s", err.Error())
			return maasInfo.LabelInfo, maasInfo.FilterMetrics
		}
		allMetrics, err := listAllMetrics("SYS.MaaS")
		if err != nil {
			logs.Logger.Errorf("[%s] Get all metrics of SYS.MaaS error: %s", err.Error())
			return maasInfo.LabelInfo, maasInfo.FilterMetrics
		}
		var filteredMetrics []model.MetricInfoList
		for _, metricInfo := range allMetrics {
			if IsMetricInfoInWhiteList(metricInfo) {
				filteredMetrics = append(filteredMetrics, metricInfo)
			}
		}
		maasInfo.LabelInfo = resourceInfos
		maasInfo.FilterMetrics = filteredMetrics
		maasInfo.TTL = time.Now().Add(GetResourceInfoExpirationTime()).Unix()
	}
	return maasInfo.LabelInfo, maasInfo.FilterMetrics
}

package collector

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/def"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ces/v1/model"
	dwsmodel "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/dws/v2/model"

	"github.com/huaweicloud/cloudeye-exporter/logs"
)

var dwsInfo serversInfo

type DWSInfo struct{}

type ListClustersRequest struct {
	EnterpriseProjectId *string `json:"enterprise_project_id,omitempty"`
}

func (o ListClustersRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListClustersRequest struct{}"
	}

	return strings.Join([]string{"ListClustersRequest", string(data)}, " ")
}

func (getter DWSInfo) GetResourceInfo() (map[string]labelInfo, []model.MetricInfoList) {
	resourceInfos := map[string]labelInfo{}
	filterMetrics := make([]model.MetricInfoList, 0)
	dwsInfo.Lock()
	defer dwsInfo.Unlock()
	if dwsInfo.LabelInfo == nil || time.Now().Unix() > dwsInfo.TTL {
		clusters, err := queryDwsCluster()
		if err != nil {
			logs.Logger.Errorf("List dws clusters error: %s", err.Error())
			return dwsInfo.LabelInfo, dwsInfo.FilterMetrics
		}
		for _, cluster := range clusters {
			metrics := buildSingleDimensionMetrics(getMetricConfigMap("SYS.DWS")["datastore_id"], "SYS.DWS", "datastore_id", cluster.Id)
			filterMetrics = append(filterMetrics, metrics...)
			info := labelInfo{
				Name:  []string{"clusterName", "epId"},
				Value: []string{cluster.Name, cluster.EnterpriseProjectId},
			}
			keys, values := getTags(fmtTags(cluster.Tags))
			info.Name = append(info.Name, keys...)
			info.Value = append(info.Value, values...)
			resourceInfos[GetResourceKeyFromMetricInfo(metrics[0])] = info
			for _, node := range cluster.Nodes {
				nodeMetrics := buildSingleDimensionMetrics(getMetricConfigMap("SYS.DWS")["dws_instance_id"], "SYS.DWS", "dws_instance_id", node.Id)
				filterMetrics = append(filterMetrics, nodeMetrics...)
				resourceInfos[GetResourceKeyFromMetricInfo(nodeMetrics[0])] = info
			}
		}
		dwsInfo.LabelInfo = resourceInfos
		dwsInfo.FilterMetrics = filterMetrics
		dwsInfo.TTL = time.Now().Add(GetResourceInfoExpirationTime()).Unix()
	}
	return dwsInfo.LabelInfo, dwsInfo.FilterMetrics
}

func queryDwsCluster() ([]dwsmodel.ClusterInfo, error) {
	epIds := getEpIdRequestPart()
	dwsClient := getHcClient(getEndpoint("dws", "v1.0"))
	var dwsClusters []dwsmodel.ClusterInfo
	for _, epId := range epIds {
		dwsResponse, err := queryDwsClusterByEpId(dwsClient, epId)
		if err != nil {
			return nil, err
		}
		dwsClusters = append(dwsClusters, *dwsResponse.Clusters...)
	}
	return dwsClusters, nil
}

func queryDwsClusterByEpId(dwsClient *core.HcHttpClient, epId string) (*dwsmodel.ListClustersResponse, error) {
	urlPath := fmt.Sprintf("/v1.0/%s/clusters", conf.ProjectID)
	requestDef := genDWSReqDef(urlPath, new(dwsmodel.ListClustersResponse))
	req := &ListClustersRequest{
		EnterpriseProjectId: &epId,
	}
	resp, err := dwsClient.Sync(req, requestDef)
	if err != nil {
		logs.Logger.Errorf("List dws cluster error: %s", err.Error())
		return nil, err
	}
	dwsResp, ok := resp.(*dwsmodel.ListClustersResponse)
	if !ok {
		logs.Logger.Error("Convert response to ListClustersResponse failed")
		return nil, fmt.Errorf("convert response to ListClustersResponse")
	}
	return dwsResp, nil
}

func genDWSReqDef(path string, response interface{}) *def.HttpRequestDef {
	reqDefBuilder := def.NewHttpRequestDefBuilder().WithMethod(http.MethodGet).WithPath(path).
		WithResponse(response).WithContentType("application/json")

	reqDefBuilder.WithRequestField(def.NewFieldDef().WithName("EnterpriseProjectId").WithJsonTag("enterprise_project_id").WithLocationType(def.Query))
	return reqDefBuilder.Build()
}

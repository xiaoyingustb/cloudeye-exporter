package collector

import (
	"time"

	cesmodel "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ces/v1/model"

	"github.com/huaweicloud/cloudeye-exporter/logs"
)

var erInfo serversInfo

type ERInfo struct{}

func (getter ERInfo) GetResourceInfo() (map[string]labelInfo, []cesmodel.MetricInfoList) {
	resourceInfos := map[string]labelInfo{}
	filterMetrics := make([]cesmodel.MetricInfoList, 0)
	erInfo.Lock()
	defer erInfo.Unlock()
	if erInfo.LabelInfo == nil || time.Now().Unix() > erInfo.TTL {
		erInstanceMap := map[string]ResourceBaseInfo{}
		if err := buildErInstanceMetricInfo(resourceInfos, &filterMetrics, erInstanceMap); err != nil {
			return erInfo.LabelInfo, erInfo.FilterMetrics
		}
		if err := buildErAttachmentMetricInfo(resourceInfos, &filterMetrics, erInstanceMap); err != nil {
			return erInfo.LabelInfo, erInfo.FilterMetrics
		}

		erInfo.LabelInfo = resourceInfos
		erInfo.FilterMetrics = filterMetrics
		erInfo.TTL = time.Now().Add(GetResourceInfoExpirationTime()).Unix()
	}
	return erInfo.LabelInfo, erInfo.FilterMetrics
}

func buildErInstanceMetricInfo(resourceInfos map[string]labelInfo, filterMetrics *[]cesmodel.MetricInfoList, erInstanceMap map[string]ResourceBaseInfo) error {
	if metricNames, ok := getMetricConfigMap("SYS.ER")["er_instance_id"]; ok {
		var erInstances []ResourceBaseInfo
		var err error
		if erInstances, err = getAllErInstanceFromRMS(); err != nil {
			logs.Logger.Errorf("Get enterprise router instance from rms error: %s", err.Error())
			return err
		}
		for _, erInstance := range erInstances {
			metrics := buildSingleDimensionMetrics(metricNames, "SYS.ER", "er_instance_id", erInstance.ID)
			*filterMetrics = append(*filterMetrics, metrics...)
			info := labelInfo{
				Name:  []string{"name", "epId"},
				Value: []string{erInstance.Name, erInstance.EpId},
			}
			keys, values := getTags(erInstance.Tags)
			info.Name = append(info.Name, keys...)
			info.Value = append(info.Value, values...)
			resourceInfos[GetResourceKeyFromMetricInfo(metrics[0])] = info
			erInstanceMap[erInstance.ID] = erInstance
		}
	}
	return nil
}

func buildErAttachmentMetricInfo(resourceInfos map[string]labelInfo, filterMetrics *[]cesmodel.MetricInfoList, erInstanceMap map[string]ResourceBaseInfo) error {
	if metricNames, ok := getMetricConfigMap("SYS.ER")["er_instance_id,er_attachment_id"]; ok {
		var erAttachmentInfos []ErAttachmentInfo
		var err error
		if erAttachmentInfos, err = getAllErAttachmentFromRMS(); err != nil {
			logs.Logger.Errorf("Get enterprise router instance from rms error: %s", err.Error())
			return err
		}
		for _, erAttachmentInfo := range erAttachmentInfos {
			erInstanceInfo, ok := erInstanceMap[erAttachmentInfo.ErId]
			if !ok {
				logs.Logger.Warnf("Parent er attachment instance not found: %s", erAttachmentInfo.ErId)
				continue
			}
			metrics := buildDimensionMetrics(metricNames, "SYS.ER", []cesmodel.MetricsDimension{
				{Name: "er_instance_id", Value: erAttachmentInfo.ErId},
				{Name: "er_attachment_id", Value: erAttachmentInfo.ID},
			})

			*filterMetrics = append(*filterMetrics, metrics...)
			info := labelInfo{
				Name:  []string{"attachment_name", "instance_name", "epId"},
				Value: []string{erAttachmentInfo.Name, erInstanceInfo.Name, erInstanceInfo.EpId},
			}
			keys, values := getTags(erAttachmentInfo.Tags)
			info.Name = append(info.Name, keys...)
			info.Value = append(info.Value, values...)
			resourceInfos[GetResourceKeyFromMetricInfo(metrics[0])] = info
		}
	}
	return nil
}

func getAllErInstanceFromRMS() ([]ResourceBaseInfo, error) {
	return getResourcesBaseInfoFromRMS("er", "instances")
}

func getAllErAttachmentFromRMS() ([]ErAttachmentInfo, error) {
	resp, err := listResources("er", "attachments")
	if err != nil {
		logs.Logger.Errorf("Failed to list resource of dcs.node, error: %s", err.Error())
		return nil, err
	}
	attachments := make([]ErAttachmentInfo, 0, len(resp))
	for _, resource := range resp {
		var attachmentProperties RmsErAttachmentProperties
		err := fmtResourceProperties(resource.Properties, &attachmentProperties)
		if err != nil {
			logs.Logger.Errorf("fmt er attachment properties error: %s", err.Error())
			continue
		}
		attachments = append(attachments, ErAttachmentInfo{
			ID:                        *resource.Id,
			Name:                      *resource.Name,
			Tags:                      resource.Tags,
			RmsErAttachmentProperties: attachmentProperties,
		})
	}
	return attachments, nil
}

type ErAttachmentInfo struct {
	ID   string
	Name string
	Tags map[string]string
	RmsErAttachmentProperties
}

type RmsErAttachmentProperties struct {
	VirtualSubnetId string `json:"virsubnet_id"`
	ResourceType    string `json:"resource_type"`
	ErId            string `json:"er_id"`
}

package collector

import (
	"time"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ces/v1/model"

	"github.com/huaweicloud/cloudeye-exporter/logs"
)

type DmsInstanceInfo struct {
	ResourceBaseInfo
	Mode string
}

var dmsInfo serversInfo

type DMSInfo struct{}

func (getter DMSInfo) GetResourceInfo() (map[string]labelInfo, []model.MetricInfoList) {
	dmsInfo.Lock()
	defer dmsInfo.Unlock()
	if dmsInfo.LabelInfo == nil {
		var err error
		dmsInfo.LabelInfo, dmsInfo.FilterMetrics, err = getDMSResourceAndMetrics()
		if err != nil {
			return dmsInfo.LabelInfo, dmsInfo.FilterMetrics
		}
		dmsInfo.TTL = time.Now().Add(GetResourceInfoExpirationTime()).Unix()
	}
	if time.Now().Unix() > dmsInfo.TTL {
		go func() {
			label, metrics, err := getDMSResourceAndMetrics()
			if err != nil {
				logs.Logger.Errorf("query dms resource info error: %s", err.Error())
				return
			}
			dmsInfo.Lock()
			defer dmsInfo.Unlock()
			dmsInfo.LabelInfo = label
			dmsInfo.FilterMetrics = metrics
			dmsInfo.TTL = time.Now().Add(GetResourceInfoExpirationTime()).Unix()
		}()
	}
	return dmsInfo.LabelInfo, dmsInfo.FilterMetrics
}

func getDMSResourceAndMetrics() (map[string]labelInfo, []model.MetricInfoList, error) {
	resourceInfos := map[string]labelInfo{}
	instances, err := getDMSInstanceFromRMS()
	if err != nil {
		return nil, nil, err
	}
	for _, instance := range instances {
		info := labelInfo{
			Name:  []string{"instanceName", "epId"},
			Value: []string{instance.Name, instance.EpId},
		}
		keys, values := getTags(instance.Tags)
		info.Name = append(info.Name, keys...)
		info.Value = append(info.Value, values...)
		resourceInfos[instance.ID] = info
	}
	allMetrics, err := listAllMetrics("SYS.DMS")
	if err != nil {
		logs.Logger.Errorf("[%s] Get all metrics of SYS.DMS error: %s", err.Error())
		return nil, nil, err
	}
	var filteredMetrics []model.MetricInfoList
	for _, metricInfo := range allMetrics {
		if IsMetricInfoInWhiteList(metricInfo) {
			filteredMetrics = append(filteredMetrics, metricInfo)
		}
	}

	return resourceInfos, filteredMetrics, nil
}

func getDMSInstanceFromRMS() ([]ResourceBaseInfo, error) {
	instances := make([]ResourceBaseInfo, 0)

	kafkaResp, err := getResourcesBaseInfoFromRMS("dms", "kafkas")
	if err != nil {
		logs.Logger.Errorf("Get all dms kafkas : %s", err.Error())
		return nil, err
	} else {
		instances = append(instances, kafkaResp...)
	}

	rabbitResp, err := getResourcesBaseInfoFromRMS("dms", "rabbitmqs")
	if err != nil {
		logs.Logger.Errorf("Get all dms rabbitmqs: %s", err.Error())
		return nil, err
	} else {
		instances = append(instances, rabbitResp...)
	}

	rocketMqs, err := getResourcesBaseInfoFromRMS("dms", "reliabilitys")
	if err != nil {
		logs.Logger.Errorf("Get all dms rocketmqs : %s", err.Error())
		return nil, err
	} else {
		instances = append(instances, rocketMqs...)
	}

	return instances, nil
}

package collector

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/huaweicloud/cloudeye-exporter/logs"
)

func TestAgentDimensionRefresher_GetAgentDimension(t *testing.T) {
	patches := getPatches()
	defer patches.Reset()
	logs.InitLog("")
	refresher := GetAgentDimensionRefresher()
	refresher.GetAgentDimension()
	assert.Empty(t, agentDimensions)

	ecsInstances := []EcsInstancesInfo{
		{
			ResourceBaseInfo: ResourceBaseInfo{
				ID:   "c53dbd97-c4ae-4ef7-8d12-fae2d3d38a80",
				Name: "nodelete-DRS-autotest-datamode1-test2",
				Tags: map[string]string{"wukong": "000"},
				EpId: "0",
			},
			IP: "192.168.20.53,100.93.4.203",
		},
	}
	ecsInfo.LabelInfo = map[string]labelInfo{
		"instance001": {},
	}
	bmsInfo.LabelInfo = map[string]labelInfo{
		"instance002": {},
	}
	patches.ApplyFuncReturn(getResourceFromRMS, true)
	patches.ApplyFuncReturn(listResources, mockRmsResource(), nil)
	patches.ApplyFuncReturn(getAllServer, ecsInstances, nil)
	patches.ApplyFunc(loadAgentDimensions, func(instanceID string) error { return nil })
	refresher.GetAgentDimension()
	assert.Empty(t, agentDimensions)

	patches.ApplyFuncReturn(listResources, nil, errors.New("xxx"))
	patches.ApplyFuncReturn(getAllServer, nil, errors.New("xxx"))
	refresher.GetAgentDimension()
	assert.Empty(t, agentDimensions)
}

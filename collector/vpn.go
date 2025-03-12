package collector

import (
	"strings"
	"time"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ces/v1/model"

	"github.com/huaweicloud/cloudeye-exporter/logs"
)

var vpnInfo serversInfo

type VPNInfo struct{}

func (getter VPNInfo) GetResourceInfo() (map[string]labelInfo, []model.MetricInfoList) {
	resourceInfos := map[string]labelInfo{}
	filterMetrics := make([]model.MetricInfoList, 0)
	vpnInfo.Lock()
	defer vpnInfo.Unlock()
	if vpnInfo.LabelInfo == nil || time.Now().Unix() > vpnInfo.TTL {
		allMetrics, err := listAllMetrics("SYS.VPN")
		if err != nil {
			logs.Logger.Error("Get all evs metrics error:", err.Error())
			return vpnInfo.LabelInfo, vpnInfo.FilterMetrics
		}
		metricMap := map[string][]model.MetricInfoList{}
		evpnSaMap := map[string][]model.MetricInfoList{}
		for _, metric := range allMetrics {
			if !IsMetricInfoInWhiteList(metric) {
				continue
			}
			resourceKey := GetResourceKeyFromMetricInfo(metric)
			dimNameArr := make([]string, 0, 0)
			for _, dimension := range metric.Dimensions {
				dimNameArr = append(dimNameArr, dimension.Name)
			}
			dimName := strings.Join(dimNameArr, ",")

			metrics, ok := metricMap[resourceKey]
			if !ok {
				metrics = make([]model.MetricInfoList, 0, 0)
			}
			metrics = append(metrics, metric)
			if dimName == "evpn_connection_id,evpn_sa_id" {
				evpnSaMap[resourceKey] = metrics
			} else {
				metricMap[resourceKey] = metrics
			}
		}
		if err = buildIpsecConnectionsInfo(&filterMetrics, resourceInfos, metricMap); err != nil {
			logs.Logger.Errorf("Build ipsec connections info error: %s", err.Error())
			return vpnInfo.LabelInfo, vpnInfo.FilterMetrics
		}
		if err = buildConnectionsInfo(&filterMetrics, resourceInfos, metricMap); err != nil {
			logs.Logger.Errorf("Build connections info error: %s", err.Error())
			return vpnInfo.LabelInfo, vpnInfo.FilterMetrics
		}
		if err = buildEVPNGatewaysInfo(&filterMetrics, resourceInfos, metricMap); err != nil {
			logs.Logger.Errorf("Build evpn gateways info error: %s", err.Error())
			return vpnInfo.LabelInfo, vpnInfo.FilterMetrics
		}
		if err = buildEVPNConnectionsInfo(&filterMetrics, resourceInfos, metricMap); err != nil {
			logs.Logger.Errorf("Build evpn connections info error: %s", err.Error())
			return vpnInfo.LabelInfo, vpnInfo.FilterMetrics
		}
		if err = buildP2CVpnGatewaysInfo(&filterMetrics, resourceInfos, metricMap); err != nil {
			logs.Logger.Errorf("Build p2c vpn gateways info error: %s", err.Error())
			return vpnInfo.LabelInfo, vpnInfo.FilterMetrics
		}
		if err = buildEVPNSaInfo(&filterMetrics, resourceInfos, evpnSaMap); err != nil {
			logs.Logger.Errorf("Build evpn vpn sa info error: %s", err.Error())
			return vpnInfo.LabelInfo, vpnInfo.FilterMetrics
		}
		vpnInfo.LabelInfo = resourceInfos
		vpnInfo.FilterMetrics = filterMetrics
		vpnInfo.TTL = time.Now().Add(GetResourceInfoExpirationTime()).Unix()
	}
	return vpnInfo.LabelInfo, vpnInfo.FilterMetrics
}

func buildIpsecConnectionsInfo(filterMetrics *[]model.MetricInfoList, resourceInfos map[string]labelInfo, metricsMap map[string][]model.MetricInfoList) error {
	ipsecConnections, err := getAllIpsecConnectionsFromRMS()
	if err != nil {
		return err
	}
	for _, ipsecConnection := range ipsecConnections {
		metric, ok := metricsMap[ipsecConnection.ID]
		if ok && metric[0].Dimensions[0].Name == "vgw_ipsec_connect_id" {
			*filterMetrics = append(*filterMetrics, metric...)
			info := labelInfo{
				Name:  []string{"name", "epId", "peer_address"},
				Value: []string{ipsecConnection.Name, ipsecConnection.EpId, ipsecConnection.PeerAddress},
			}
			keys, values := getTags(ipsecConnection.Tags)
			info.Name = append(info.Name, keys...)
			info.Value = append(info.Value, values...)
			resourceInfos[GetResourceKeyFromMetricInfo(metric[0])] = info
		}
	}
	return nil
}

func buildEVPNConnectionsInfo(filterMetrics *[]model.MetricInfoList, resourceInfos map[string]labelInfo, metricsMap map[string][]model.MetricInfoList) error {
	allConnections, err := getAllConnectionsFromRMS()
	if err != nil {
		return err
	}
	for _, ipsecConnection := range allConnections {
		metric, ok := metricsMap[ipsecConnection.ID]
		if ok && metric[0].Dimensions[0].Name == "evpn_connection_id" {
			*filterMetrics = append(*filterMetrics, metric...)
			info := labelInfo{
				Name:  []string{"name", "epId"},
				Value: []string{ipsecConnection.Name, ipsecConnection.EpId},
			}
			keys, values := getTags(ipsecConnection.Tags)
			info.Name = append(info.Name, keys...)
			info.Value = append(info.Value, values...)
			resourceInfos[GetResourceKeyFromMetricInfo(metric[0])] = info
		}
	}
	return nil
}

func buildEVPNGatewaysInfo(filterMetrics *[]model.MetricInfoList, resourceInfos map[string]labelInfo, metricsMap map[string][]model.MetricInfoList) error {
	allIpsecConnections, err := getAllVPNGatewaysFromRMS()
	if err != nil {
		return err
	}
	for _, ipsecConnection := range allIpsecConnections {
		metric, ok := metricsMap[ipsecConnection.ID]
		if ok && metric[0].Dimensions[0].Name == "evpn_gateway_id" {
			*filterMetrics = append(*filterMetrics, metric...)
			info := labelInfo{
				Name:  []string{"name", "epId"},
				Value: []string{ipsecConnection.Name, ipsecConnection.EpId},
			}
			keys, values := getTags(ipsecConnection.Tags)
			info.Name = append(info.Name, keys...)
			info.Value = append(info.Value, values...)
			resourceInfos[GetResourceKeyFromMetricInfo(metric[0])] = info
		}
	}
	return nil
}

func buildP2CVpnGatewaysInfo(filterMetrics *[]model.MetricInfoList, resourceInfos map[string]labelInfo, metricsMap map[string][]model.MetricInfoList) error {
	p2cVpnGateways, err := getAllP2CVPNGatewaysFromRMS()
	if err != nil {
		return err
	}
	for _, ipsecConnection := range p2cVpnGateways {
		metric, ok := metricsMap[ipsecConnection.ID]
		if ok && metric[0].Dimensions[0].Name == "p2c_vpn_gateway_id" {
			*filterMetrics = append(*filterMetrics, metric...)
			info := labelInfo{
				Name:  []string{"name", "epId"},
				Value: []string{ipsecConnection.Name, ipsecConnection.EpId},
			}
			keys, values := getTags(ipsecConnection.Tags)
			info.Name = append(info.Name, keys...)
			info.Value = append(info.Value, values...)
			resourceInfos[GetResourceKeyFromMetricInfo(metric[0])] = info
		}
	}
	return nil
}

func buildConnectionsInfo(filterMetrics *[]model.MetricInfoList, resourceInfos map[string]labelInfo, metricsMap map[string][]model.MetricInfoList) error {
	connections, err := getAllConnectionsFromRMS()
	if err != nil {
		return err
	}
	for _, connection := range connections {
		metric, ok := metricsMap[connection.ID]
		if ok && metric[0].Dimensions[0].Name == "vpn_connection_id" {
			*filterMetrics = append(*filterMetrics, metric...)
			info := labelInfo{
				Name:  []string{"name", "epId", "peer_address"},
				Value: []string{connection.Name, connection.EpId, connection.PeerAddress},
			}
			keys, values := getTags(connection.Tags)
			info.Name = append(info.Name, keys...)
			info.Value = append(info.Value, values...)
			resourceInfos[GetResourceKeyFromMetricInfo(metric[0])] = info
		}
	}
	return nil
}

func buildEVPNSaInfo(filterMetrics *[]model.MetricInfoList, resourceInfos map[string]labelInfo, evpnSaMap map[string][]model.MetricInfoList) error {
	allConnections, err := getAllConnectionsFromRMS()
	if err != nil {
		return err
	}
	connectionMap := map[string]ConnectionInfo{}
	for _, connection := range allConnections {
		connectionMap[connection.ID] = connection
	}
	for dimKey, metric := range evpnSaMap {
		dimValueArr := strings.Split(dimKey, ".")
		connectionInfo, ok := connectionMap[dimValueArr[0]]
		if ok {
			*filterMetrics = append(*filterMetrics, metric...)
			info := labelInfo{
				Name:  []string{"name", "epId", "peer_address"},
				Value: []string{connectionInfo.Name, connectionInfo.EpId, connectionInfo.PeerAddress},
			}
			keys, values := getTags(connectionInfo.Tags)
			info.Name = append(info.Name, keys...)
			info.Value = append(info.Value, values...)
			resourceInfos[GetResourceKeyFromMetricInfo(metric[0])] = info
		}
	}
	return nil
}

type ConnectionInfo struct {
	ResourceBaseInfo
	PeerAddress string
}

type ConnectionProperties struct {
	PeerAddress string `json:"peer_address"`
}

type EVPNConnectionProperties struct {
	TunnelPeerAddress string `json:"tunnel_peer_address"`
}

func getAllVPNGatewaysFromRMS() ([]ConnectionInfo, error) {
	resp, err := listResources("vpnaas", "vpnGateways")
	if err != nil {
		logs.Logger.Errorf("Failed to list resource of vpnaas.vpnGateways, error: %s", err.Error())
		return nil, err
	}
	connections := make([]ConnectionInfo, 0, len(resp))
	for _, resource := range resp {
		connections = append(connections, ConnectionInfo{
			ResourceBaseInfo: ResourceBaseInfo{
				ID:   *resource.Id,
				Name: *resource.Name,
				EpId: *resource.EpId,
				Tags: resource.Tags,
			},
		})
	}
	return connections, nil
}

func getAllIpsecConnectionsFromRMS() ([]ConnectionInfo, error) {
	resp, err := listResources("vpnaas", "ipsec-site-connections")
	if err != nil {
		logs.Logger.Errorf("Failed to list resource of vpnaas.ipsec-site-connections, error: %s", err.Error())
		return nil, err
	}
	connections := make([]ConnectionInfo, 0, len(resp))
	for _, resource := range resp {
		var properties ConnectionProperties
		err := fmtResourceProperties(resource.Properties, &properties)
		if err != nil {
			logs.Logger.Errorf("fmt vpn properties error: %s", err.Error())
			continue
		}
		connections = append(connections, ConnectionInfo{
			ResourceBaseInfo: ResourceBaseInfo{
				ID:   *resource.Id,
				Name: *resource.Name,
				EpId: *resource.EpId,
				Tags: resource.Tags,
			},
			PeerAddress: properties.PeerAddress,
		})
	}
	return connections, nil
}

func getAllConnectionsFromRMS() ([]ConnectionInfo, error) {
	resp, err := listResources("vpnaas", "vpnConnections")
	if err != nil {
		logs.Logger.Errorf("Failed to list resource of vpnaas.vpnConnections, error: %s", err.Error())
		return nil, err
	}
	connections := make([]ConnectionInfo, 0, len(resp))
	for _, resource := range resp {
		var properties ConnectionProperties
		err := fmtResourceProperties(resource.Properties, &properties)
		if err != nil {
			logs.Logger.Errorf("fmt vpn properties error: %s", err.Error())
			continue
		}
		connections = append(connections, ConnectionInfo{
			ResourceBaseInfo: ResourceBaseInfo{
				ID:   *resource.Id,
				Name: *resource.Name,
				EpId: *resource.EpId,
				Tags: resource.Tags,
			},
			PeerAddress: properties.PeerAddress,
		})
	}
	return connections, nil
}

func getAllP2CVPNGatewaysFromRMS() ([]ConnectionInfo, error) {
	resp, err := listResources("vpnaas", "p2c-vpn-gateways")
	if err != nil {
		logs.Logger.Errorf("Failed to list resource of vpnaas.p2c-vpn-gateways, error: %s", err.Error())
		return nil, err
	}
	connections := make([]ConnectionInfo, 0, len(resp))
	for _, resource := range resp {
		var properties EVPNConnectionProperties
		err := fmtResourceProperties(resource.Properties, &properties)
		if err != nil {
			logs.Logger.Errorf("fmt evpn properties error: %s", err.Error())
			continue
		}
		connections = append(connections, ConnectionInfo{
			ResourceBaseInfo: ResourceBaseInfo{
				ID:   *resource.Id,
				Name: *resource.Name,
				EpId: *resource.EpId,
				Tags: resource.Tags,
			},
		})
	}
	return connections, nil
}

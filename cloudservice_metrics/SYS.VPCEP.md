> CES Exporter支持导出的“VPC终端节点”指标如下表所示

|维度|指标名|指标描述|指标单位|
|:--|:--|:--|:--|
|ep_instance_id<br>（终端节点实例）|vpcep_bps|终端节点总带宽|bit/s|
|ep_instance_id<br>（终端节点实例）|vpcep_pps|终端节点总数据包速率|packets/s|
|ep_instance_id<br>（终端节点实例）|vpcep_drop_num|终端节点新增丢包数|packets|
|ep_instance_id<br>（终端节点实例）|vpcep_connections|终端节点连接数|count|
|ep_instance_id<br>（终端节点实例）|vpcep_rx_pps|终端节点流入包速率|packets/s|
|ep_instance_id<br>（终端节点实例）|vpcep_from_server_rst_num|服务端重置数量|packets|
|ep_instance_id<br>（终端节点实例）|vpcep_tx_byte|终端节点流出流量|byte(IEC)|
|ep_instance_id<br>（终端节点实例）|vpcep_tx_pps|终端节点流出包速率|packets/s|
|ep_instance_id<br>（终端节点实例）|vpcep_from_client_rst_num|客户端重置数量|packets|
|ep_instance_id<br>（终端节点实例）|vpcep_rx_bps|终端节点流入带宽|bit/s(IEC)|
|ep_instance_id<br>（终端节点实例）|vpcep_to_server_rst_num|终端节点重置服务端数量|packets|
|ep_instance_id<br>（终端节点实例）|vpcep_to_client_rst_num|终端节点重置客户端数量|packets|
|ep_instance_id<br>（终端节点实例）|vpcep_rx_byte|终端节点流入流量|byte(IEC)|
|ep_instance_id<br>（终端节点实例）|vpcep_tx_bps|终端节点流出带宽|bit/s(IEC)|
|ep_instance_id<br>（终端节点实例）|vpcep_act_connections|终端节点活跃连接数|count|
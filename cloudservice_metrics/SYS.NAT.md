> CES Exporter支持导出的“NAT网关”指标如下表所示

|维度|指标名|指标描述|指标单位|
|:--|:--|:--|:--|
|nat_gateway_id<br>（公网NAT网关）|snat_connection|SNAT连接数|个|
|nat_gateway_id<br>（公网NAT网关）|inbound_bandwidth|入方向带宽|bit/s|
|nat_gateway_id<br>（公网NAT网关）|outbound_bandwidth|出方向带宽|bit/s|
|nat_gateway_id<br>（公网NAT网关）|inbound_pps|入方向PPS|个|
|nat_gateway_id<br>（公网NAT网关）|outbound_pps|出方向PPS|个|
|nat_gateway_id<br>（公网NAT网关）|inbound_traffic|入方向流量|Byte|
|nat_gateway_id<br>（公网NAT网关）|outbound_traffic|出方向流量|Byte|
|nat_gateway_id<br>（公网NAT网关）|snat_connection_ratio|SNAT连接数使用率|%|
|nat_gateway_id<br>（公网NAT网关）|inbound_bandwidth_ratio|入方向带宽使用率|%|
|nat_gateway_id<br>（公网NAT网关）|outbound_bandwidth_ratio|出方向带宽使用率|%|
|nat_gateway_id<br>（公网NAT网关）|total_inbound_tcp_bandwidth|入方向TCP总带宽|bit/s(IEC)|
|nat_gateway_id<br>（公网NAT网关）|total_inbound_udp_bandwidth|入方向UDP总带宽|bit/s(IEC)|
|nat_gateway_id<br>（公网NAT网关）|total_outbound_tcp_bandwidth|出方向TCP总带宽|bit/s(IEC)|
|nat_gateway_id<br>（公网NAT网关）|total_outbound_udp_bandwidth|出方向UDP总带宽|bit/s(IEC)|
|nat_gateway_id<br>（公网NAT网关）|packets_drop_count_snat_connection_beyond|丢包数（SNAT连接数超限）|count|
|nat_gateway_id<br>（公网NAT网关）|packets_drop_count_pps_beyond|丢包数（PPS超限）|count|
|nat_gateway_id<br>（公网NAT网关）|packets_drop_count_eip_port_alloc_beyond|丢包数（EIP端口分配超限）|count|
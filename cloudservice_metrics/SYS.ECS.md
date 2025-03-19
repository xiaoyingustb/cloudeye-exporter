> CES Exporter默认导出AGT.ECS下已上报的全量指标，具体请参考[官网文档](https://support.huaweicloud.com/usermanual-ecs/ecs_03_1003.html)

> CES Exporter支持导出的“弹性云服务器”指标（SYS.ECS）如下表所示

|维度|指标名|指标描述|指标单位|
|:--|:--|:--|:--|
|instance_id<br>（云服务器）|cpu_util|CPU使用率|%|
|instance_id<br>（云服务器）|mem_util|（Windows）内存使用率|%|
|instance_id<br>（云服务器）|disk_util_inband|（Windows）磁盘使用率|%|
|instance_id<br>（云服务器）|disk_read_bytes_rate|磁盘读带宽|byte/s|
|instance_id<br>（云服务器）|disk_write_bytes_rate|磁盘写带宽|byte/s|
|instance_id<br>（云服务器）|disk_read_requests_rate|磁盘读IOPS|请求/秒|
|instance_id<br>（云服务器）|disk_write_requests_rate|磁盘写IOPS|请求/秒|
|instance_id<br>（云服务器）|network_incoming_bytes_rate_inband|带内网络流入速率（仅支持Windows操作系统）|byte/s|
|instance_id<br>（云服务器）|network_outgoing_bytes_rate_inband|带内网络流出速率（仅支持Windows操作系统）|byte/s|
|instance_id<br>（云服务器）|network_incoming_bytes_aggregate_rate|带外网络流入速率|byte/s|
|instance_id<br>（云服务器）|cpu_credit_balance|CPU积分累积量|积分|
|instance_id<br>（云服务器）|cpu_credit_usage|CPU积分使用量|积分|
|instance_id<br>（云服务器）|cpu_surplus_credit_balance|CPU超额积分累积量|积分|
|instance_id<br>（云服务器）|cpu_surplus_credit_charged|CPU超额积分收费量|积分|
|instance_id<br>（云服务器）|network_vm_connections|网络连接数|个|
|instance_id<br>（云服务器）|network_vm_bandwidth_in|虚拟机入方向带宽|Byte/s|
|instance_id<br>（云服务器）|network_vm_bandwidth_out|虚拟机出方向带宽|Byte/s|
|instance_id<br>（云服务器）|network_vm_pps_in|虚拟机入方向PPS|包/秒|
|instance_id<br>（云服务器）|network_vm_pps_out|虚拟机出方向PPS|包/秒|
|instance_id<br>（云服务器）|network_vm_newconnections|虚拟机整机新建连接数|connect/s|
|instance_id<br>（云服务器）|network_outgoing_bytes_aggregate_rate|带外网络流出速率|byte/s|
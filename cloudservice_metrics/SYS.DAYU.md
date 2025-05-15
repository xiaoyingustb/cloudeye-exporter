> CES Exporter支持导出的“数据治理中心”指标如下表所示

|维度|指标名|指标描述|指标单位|
|:--|:--|:--|:--|
|stream_id<br>（实时数据接入）|dis11_stream_record_retention_time|记录滞留时间（毫秒）|ms|
|stream_id<br>（实时数据接入）|dis01_stream_put_bytes_rate|总输入流量|byte/s|
|stream_id<br>（实时数据接入）|dis02_stream_get_bytes_rate|总输出流量|byte/s|
|stream_id<br>（实时数据接入）|dis03_stream_put_records|总输入记录数|个/秒|
|stream_id<br>（实时数据接入）|dis04_stream_get_records|总输出记录数|个/秒|
|stream_id<br>（实时数据接入）|dis05_stream_put_requests_succeed|上传请求成功数|个/秒|
|stream_id<br>（实时数据接入）|dis06_stream_get_requests_succeed|下载请求成功数|个/秒|
|stream_id<br>（实时数据接入）|dis07_stream_put_req_average_latency|上传请求平均处理时间|毫秒|
|stream_id<br>（实时数据接入）|dis08_stream_get_req_average_latency|下载请求平均处理时间|毫秒|
|stream_id<br>（实时数据接入）|dis09_stream_traffic_control_put_records|因流控拒绝的上传请求数|个/秒|
|stream_id<br>（实时数据接入）|dis10_stream_traffic_control_get_records|因流控拒绝的下载请求数|个/秒|
|stream_id<br>（实时数据接入）|dis11_stream_record_retention_time|记录滞留时间（毫秒）|ms|
|stream_id<br>（实时数据接入）|dis12_stream_accumulated_message_count|消息堆积数|count|
|stream_id<br>（实时数据接入）|dis13_stream_transfer_delay|转储时延|s|
|stream_id<br>（实时数据接入）|dis18_stream_bytes_put_max|通道峰值流量|byte/s|
|stream_id<br>（实时数据接入）|dis19_stream_readable_partition_count|通道可读分区个数|count|
|stream_id<br>（实时数据接入）|dis20_stream_writable_partition_count|通道可写分区个数|count|
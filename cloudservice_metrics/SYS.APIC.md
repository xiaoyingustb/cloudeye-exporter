> CES Exporter支持导出的“API网关专享版”指标如下表所示

|维度|指标名|指标描述|指标单位|
|:--|:--|:--|:--|
|instance_id<br>（APIG实例）|requests|接口调用次数|次/分钟|
|instance_id<br>（APIG实例）|error_4xx|4xx 异常次数|次/分钟|
|instance_id<br>（APIG实例）|error_5xx|5xx 异常次数|次/分钟|
|instance_id<br>（APIG实例）|throttled_calls|被流控的调用次数|次/分钟|
|instance_id<br>（APIG实例）|avg_latency|平均延迟毫秒数|毫秒|
|instance_id<br>（APIG实例）|max_latency|最大延迟毫秒数|毫秒|
|instance_id,api_id<br>（APIG实例,接口）|req_count|接口调用次数|次/分钟|
|instance_id,api_id<br>（APIG实例,接口）|req_count_2xx|2xx调用次数|次/分钟|
|instance_id,api_id<br>（APIG实例,接口）|req_count_4xx|4xx异常次数|次/分钟|
|instance_id,api_id<br>（APIG实例,接口）|req_count_5xx|5xx异常次数|次/分钟|
|instance_id,api_id<br>（APIG实例,接口）|req_count_error|异常次数|次/分钟|
|instance_id,api_id<br>（APIG实例,接口）|avg_latency|平均延迟毫秒数|毫秒|
|instance_id,api_id<br>（APIG实例,接口）|max_latency|最大延迟毫秒数|毫秒|
|instance_id,api_id<br>（APIG实例,接口）|input_throughput|流入流量|byte|
|instance_id,api_id<br>（APIG实例,接口）|output_throughput|流出流量|byte|
|instance_id,node_ip<br>（APIG实例,节点ip）|node_system_load|网关节点系统负载|count|
|instance_id,node_ip<br>（APIG实例,节点ip）|node_qps|网关节点qps|次/秒|
|instance_id,node_ip<br>（APIG实例,节点ip）|node_cpu_usage|网关节点cpu使用率|%|
|instance_id,node_ip<br>（APIG实例,节点ip）|node_memory_usage|网关节点内存使用率|%|
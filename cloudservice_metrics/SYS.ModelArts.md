> CES Exporter支持导出的“ModelArts”指标如下表所示

|维度|指标名|指标描述|指标单位|
|:--|:--|:--|:--|
|service_id<br>（服务）|successfully_called_times|调用成功次数|次/分钟|
|service_id<br>（服务）|failed_called_times|调用失败次数|次/分钟|
|service_id<br>（服务）|total_called_times|调用总次数|次/分钟|
|service_id<br>（服务）|req_count_2xx|2xx响应次数|次/分钟|
|service_id<br>（服务）|req_count_4xx|4xx异常次数|次/分钟|
|service_id<br>（服务）|req_count_5xx|5xx异常次数|次/分钟|
|service_id<br>（服务）|avg_latency|平均延迟毫秒数|ms|
|service_id,model_id<br>（服务,模型实例）|cpu_usage|CPU使用率|%|
|service_id,model_id<br>（服务,模型实例）|mem_usage|内存使用率|%|
|service_id,model_id<br>（服务,模型实例）|gpu_util|GPU使用率|%|
|service_id,model_id<br>（服务,模型实例）|gpu_mem_usage|GPU显存使用率|%|
|service_id,model_id<br>（服务,模型实例）|successfully_called_times|调用成功次数|次/分钟|
|service_id,model_id<br>（服务,模型实例）|failed_called_times|调用失败次数|次/分钟|
|service_id,model_id<br>（服务,模型实例）|total_called_times|调用总次数|次/分钟|
|service_id,model_id<br>（服务,模型实例）|disk_read_rate|磁盘读取速率|bit/min|
|service_id,model_id<br>（服务,模型实例）|disk_write_rate|磁盘写入速率|bit/min|
|service_id,model_id<br>（服务,模型实例）|send_bytes_rate|上行Bps|bit/min|
|service_id,model_id<br>（服务,模型实例）|recv_bytes_rate|下行Bps|bit/min|
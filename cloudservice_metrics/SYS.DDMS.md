> CES Exporter支持导出的“分布式数据库中间件”指标如下表所示

|维度|指标名|指标描述|指标单位|
|:--|:--|:--|:--|
|instance_id,node_id<br>（DDM实例,DDM节点）|ddm_cpu_util|CPU使用率|%|
|instance_id,node_id<br>（DDM实例,DDM节点）|ddm_mem_util|内存使用率|%|
|instance_id,node_id<br>（DDM实例,DDM节点）|ddm_bytes_in|网络输入吞吐量|byte/s|
|instance_id,node_id<br>（DDM实例,DDM节点）|ddm_bytes_out|网络输出吞吐量|byte/s|
|instance_id,node_id<br>（DDM实例,DDM节点）|ddm_qps|QPS|次数|
|instance_id,node_id<br>（DDM实例,DDM节点）|ddm_rw_ratio|读写比例|%|
|instance_id,node_id<br>（DDM实例,DDM节点）|ddm_read_count|读次数|次数|
|instance_id,node_id<br>（DDM实例,DDM节点）|ddm_write_count|写次数|次数|
|instance_id,node_id<br>（DDM实例,DDM节点）|ddm_slow_log|慢SQL数|条数|
|instance_id,node_id<br>（DDM实例,DDM节点）|ddm_rt_avg|平均响应时延|毫秒|
|instance_id,node_id<br>（DDM实例,DDM节点）|ddm_connections|连接数|个数|
|instance_id,node_id<br>（DDM实例,DDM节点）|ddm_backend_connection_ratio|后端连接池水位|%|
|instance_id,node_id<br>（DDM实例,DDM节点）|ddm_active_connections|活跃连接数|次数|
|instance_id,node_id<br>（DDM实例,DDM节点）|ddm_connection_util|连接数使用率|%|
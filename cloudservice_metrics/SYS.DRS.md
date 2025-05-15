> CES Exporter支持导出的“数据复制服务”指标如下表所示

|维度|指标名|指标描述|指标单位|
|:--|:--|:--|:--|
|instance_id<br>（DRS运行实例）|cpu_util|CPU使用率|%|
|instance_id<br>（DRS运行实例）|mem_util|内存使用率|%|
|instance_id<br>（DRS运行实例）|network_incoming_bytes_rate|网络输入吞吐量|byte/s|
|instance_id<br>（DRS运行实例）|network_outgoing_bytes_rate|网络输出吞吐量|byte/s|
|instance_id<br>（DRS运行实例）|disk_read_bytes_rate|磁盘读吞吐量|byte/s|
|instance_id<br>（DRS运行实例）|disk_write_bytes_rate|磁盘写吞吐量|byte/s|
|instance_id<br>（DRS运行实例）|disk_util|磁盘利用率|%|
|instance_id<br>（DRS运行实例）|extract_bytes_rate|读源库吞吐量|byte/s|
|instance_id<br>（DRS运行实例）|extract_rows_rate|读源库频率|row/s|
|instance_id<br>（DRS运行实例）|extract_latency|源库WAL抽取延迟|ms|
|instance_id<br>（DRS运行实例）|apply_bytes_rate|写目标库吞吐量|byte/s|
|instance_id<br>（DRS运行实例）|apply_rows_rate|写目标库频率|row/s|
|instance_id<br>（DRS运行实例）|apply_transactions_rate|DML TPS|transaction/s|
|instance_id<br>（DRS运行实例）|apply_ddls_rate|DDL TPS||
|instance_id<br>（DRS运行实例）|apply_latency|数据同步延迟|ms|
|instance_id<br>（DRS运行实例）|apply_average_execute_time|事务平均执行时间|ms|
|instance_id<br>（DRS运行实例）|apply_average_commit_time|事务平均提交时间|ms|
|instance_id<br>（DRS运行实例）|apply_current_state|同步状态||
|instance_id<br>（DRS运行实例）|apply_thread_workers|同步线程数量||
|instance_id<br>（DRS运行实例）|apply_job_status|任务状态||
|instance_id<br>（DRS运行实例）|apply_ddls_number|DDL总数|count|
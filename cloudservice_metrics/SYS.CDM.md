> CES Exporter支持导出的“云数据迁移服务”指标如下表所示

|维度|指标名|指标描述|指标单位|
|:--|:--|:--|:--|
|instance_id<br>（实例）|bytes_in|网络流入速率|byte/s|
|instance_id<br>（实例）|bytes_out|网络流出速率|byte/s|
|instance_id<br>（实例）|cpu_usage|CPU使用率|%|
|instance_id<br>（实例）|mem_usage|内存使用率|%|
|instance_id<br>（实例）|disk_usage|磁盘利用率|%|
|instance_id<br>（实例）|disk_io|磁盘io|byte/s|
|instance_id<br>（实例）|tomcat_heap_usage|堆内存使用率|%|
|instance_id<br>（实例）|tomcat_connect|tomcat并发连接数|个|
|instance_id<br>（实例）|tomcat_thread_count|tomcat线程数|个|
|instance_id<br>（实例）|pg_connect|数据库连接数|个|
|instance_id<br>（实例）|pg_submission_row|历史记录表行数|行|
|instance_id<br>（实例）|pg_failed_job_rate|失败作业率|%|
|instance_id<br>（实例）|inodes_usage|Inodes利用率|%|
|instance_id<br>（实例）|pg_pending_job|排队作业数|count|
|instance_id<br>（实例）|pending_threads|排队抽取并发数|count|
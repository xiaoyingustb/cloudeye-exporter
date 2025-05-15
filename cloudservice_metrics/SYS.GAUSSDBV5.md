> CES Exporter支持导出的“云数据库 GaussDB(for openGauss)”指标如下表所示

|维度|指标名|指标描述|指标单位|
|:--|:--|:--|:--|
|gaussdbv5_instance_id<br>（云数据库 GaussDB实例）|rds005_instance_disk_used_size|实例数据磁盘已使用大小|GB|
|gaussdbv5_instance_id<br>（云数据库 GaussDB实例）|rds006_instance_disk_total_size|实例数据磁盘总大小|GB|
|gaussdbv5_instance_id<br>（云数据库 GaussDB实例）|rds007_instance_disk_usage|实例数据磁盘已使用百分比|%|
|gaussdbv5_instance_id<br>（云数据库 GaussDB实例）|rds035_buffer_hit_ratio|buffer 命中率|%|
|gaussdbv5_instance_id<br>（云数据库 GaussDB实例）|rds036_deadlocks|死锁次数|Count|
|gaussdbv5_instance_id<br>（云数据库 GaussDB实例）|rds048_P80|80% SQL的响应时间|us|
|gaussdbv5_instance_id<br>（云数据库 GaussDB实例）|rds049_P95|95% SQL的响应时间|us|
|gaussdbv5_instance_id,gaussdbv5_node_id<br>（云数据库 GaussDB实例,云数据库 GaussDB节点）|rds001_cpu_util|CPU使用率|%|
|gaussdbv5_instance_id,gaussdbv5_node_id<br>（云数据库 GaussDB实例,云数据库 GaussDB节点）|rds002_mem_util|内存使用率|%|
|gaussdbv5_instance_id,gaussdbv5_node_id<br>（云数据库 GaussDB实例,云数据库 GaussDB节点）|rds003_bytes_in|数据写入量|Byte/s|
|gaussdbv5_instance_id,gaussdbv5_node_id<br>（云数据库 GaussDB实例,云数据库 GaussDB节点）|rds004_bytes_out|数据传出量|Byte/s|
|gaussdbv5_instance_id,gaussdbv5_node_id<br>（云数据库 GaussDB实例,云数据库 GaussDB节点）|rds014_iops|数据磁盘每秒读写次数|Count/s|
|gaussdbv5_instance_id,gaussdbv5_node_id<br>（云数据库 GaussDB实例,云数据库 GaussDB节点）|rds016_disk_write_throughput|数据磁盘写吞吐量|Byte/s|
|gaussdbv5_instance_id,gaussdbv5_node_id<br>（云数据库 GaussDB实例,云数据库 GaussDB节点）|rds017_disk_read_throughput|数据磁盘读吞吐量|Byte/s|
|gaussdbv5_instance_id,gaussdbv5_node_id<br>（云数据库 GaussDB实例,云数据库 GaussDB节点）|rds020_avg_disk_ms_per_write|数据磁盘单次写入花费的时间|ms|
|gaussdbv5_instance_id,gaussdbv5_node_id<br>（云数据库 GaussDB实例,云数据库 GaussDB节点）|rds021_avg_disk_ms_per_read|数据磁盘单次读取花费的时间|ms|
|gaussdbv5_instance_id,gaussdbv5_node_id<br>（云数据库 GaussDB实例,云数据库 GaussDB节点）|io_bandwidth_usage|磁盘io带宽占用率|%|
|gaussdbv5_instance_id,gaussdbv5_node_id<br>（云数据库 GaussDB实例,云数据库 GaussDB节点）|iops_usage|IOPS使用率|%|
|gaussdbv5_instance_id,gaussdbv5_node_id<br>（云数据库 GaussDB实例,云数据库 GaussDB节点）|rds068_swap_used_ratio|交换内存使用率|%|
|gaussdbv5_instance_id,gaussdbv5_node_id<br>（云数据库 GaussDB实例,云数据库 GaussDB节点）|rds069_swap_total_size|交换内存总大小|MB|
|gaussdbv5_instance_id,gaussdbv5_node_id<br>（云数据库 GaussDB实例,云数据库 GaussDB节点）|rds084_sys_database_size|系统数据库大小|byte(IEC)|
|gaussdbv5_instance_id,gaussdbv5_node_id<br>（云数据库 GaussDB实例,云数据库 GaussDB节点）|rds085_user_database_size|用户数据库总大小|byte(IEC)|
|gaussdbv5_instance_id,gaussdbv5_node_id<br>（云数据库 GaussDB实例,云数据库 GaussDB节点）|retrans_rate|重传比例|%|
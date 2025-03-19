> CES Exporter支持导出的“云硬盘”指标如下表所示

|维度|指标名|指标描述|指标单位|
|:--|:--|:--|:--|
|disk_name<br>（磁盘）|disk_device_read_bytes_rate|磁盘读带宽|byte/s|
|disk_name<br>（磁盘）|disk_device_write_bytes_rate|磁盘写带宽|byte/s|
|disk_name<br>（磁盘）|disk_device_read_requests_rate|磁盘读IOPS|请求/秒|
|disk_name<br>（磁盘）|disk_device_write_requests_rate|磁盘写IOPS|请求/秒|
|disk_name<br>（磁盘）|disk_device_queue_length|平均队列长度|个|
|disk_name<br>（磁盘）|disk_device_io_util|磁盘读写使用率|%|
|disk_name<br>（磁盘）|disk_device_write_bytes_per_operation|平均写操作大小|KB/op|
|disk_name<br>（磁盘）|disk_device_read_bytes_per_operation|平均读操作大小|KB/op|
|disk_name<br>（磁盘）|disk_device_write_await|平均写操作耗时|ms/op|
|disk_name<br>（磁盘）|disk_device_read_await|平均读操作耗时|ms/op|
|disk_name<br>（磁盘）|disk_device_io_svctm|平均IO服务时长|ms/op|
|disk_name<br>（磁盘）|disk_device_io_iops_qos_num|IOPS达到上限(次数)|count|
|disk_name<br>（磁盘）|disk_device_io_iobw_qos_num|带宽达到上限(次数)|count|
# cloudeye-exporter

Prometheus cloudeye exporter for [Huaweicloud](https://www.huaweicloud.com/).

Note: The plug-in is applicable only to the Huaweicloud regions.

[中文](./README_cn.md)

## Download
```
$ git clone https://github.com/huaweicloud/cloudeye-exporter
```

## (Option) Building The Discovery with Exact steps on clean Ubuntu 16.04 
```
$ wget https://dl.google.com/go/go1.17.6.linux-amd64.tar.gz
$ sudo tar -C /usr/local -xzf go1.17.6.linux-amd64.tar.gz
$ export PATH=$PATH:/usr/local/go/bin # You should put in your .profile or .bashrc
$ go version # to verify it runs and version #

$ go get github.com/huaweicloud/cloudeye-exporter
$ cd ~/go/src/github.com/huaweicloud/cloudeye-exporter
$ go build
```

## Usage
```
 ./cloudeye-exporter  -config=clouds.yml
```

The default port is 8087, default config file location is ./clouds.yml.

Visit metrics in http://localhost:8087/metrics?services=SYS.VPC,SYS.ELB


To avoid ak&sk leaks, you can input your ak&sk through the command line with '-s' param like this
```shell
./cloudeye-exporter -s true
```
To startup the cloudeye-exporter using script, here is an example
```shell
#!/bin/bash
# Don't write plain text ak&sk in script.
# Encrypt the ak&sk instead and use your own decrypt function to assign the return value.
huaweiCloud_AK=your_decrypt_function("your encrypted ak")
huaweiCloud_SK=your_decrypt_function("your encrypted sk")
$(./cloudeye-exporter -s true<<EOF
$huaweiCloud_AK $huaweiCloud_SK
EOF)
```

## Help
```
Usage of ./cloudeye-exporter:
  -config string
        Path to the cloud configuration file (default "./clouds.yml") 
```

## Example of config file(clouds.yml)
The "URL" value can be get from [Identity and Access Management (IAM) endpoint list (Internal)](https://developer.huaweicloud.com/intl/en-us/endpoint?IAM) and [Identity and Access Management (IAM) endpoint list (China)](https://developer.huaweicloud.com/endpoint?IAM).
```
global:
  prefix: "huaweicloud"
  port: "{private IP}:8087" # For security purposes, you are advised not to expose the Expoter service port to the public network. You are advised to set this parameter to 127.0.0.1:{port} or {private IP address}:{port}, for example, 192.168.1.100:8087. If the port needs to be exposed to the public network, ensure that the security group, firewall, and iptables access control policies are properly configured to meet the minimum access permission principle.
  metric_path: "/metrics"
  scrape_batch_size: 300
  resource_sync_interval_minutes: 20 # Update frequency of resource information: resource information is updated every 180 minutes by default; If this parameter is set to a value less than 10 minutes, the information is updated every 10 minutes.
  ep_ids: "xxx1,xxx2" # This is optional. Filter resources by enterpries project, cloudeye-exporter will get all resources when this is empty, if you need multiple enterprise project, use comma split them.
  logs_conf_path: "/root/logs.yml" # This is optional. We recommend that you use an absolute path for the log configuration file path. If this line is absent, the program will use configuration file in the directory where the startup command is executed by default.
  metrics_conf_path: "/root/metrics.yml" # This is optional. We recommend that you use an absolute path for the metrics configuration file path. If this line is absent, the program will use configuration file in the directory where the startup command is executed by default.
  endpoints_conf_path: "/root/endpoints.yml" # This is optional. We recommend that you use an absolute path for the service endpoints configuration file path. If this line is absent, the program will use configuration file in the directory where the startup command is executed by default.
  ignore_ssl_verify: false # This is optional. The SSL certificate is verified by default when the exporter queries resources or indicators. If the exporter is abnormal due to SSL certificate verification, you can set this configuration to true to skip SSL certificate verification.
  client_cn: "192.168.XX.XXX" # This is optional. Specify a list of trusted DNS names to verify the DNS names in the HTTPS certificate chain. Separate multiple DNS names in the list with commas.
  unit_standardization_enabled: false # This is optional.The unit label optimization feature for enabling metric data is now available. By default, this feature is disabled. After enabling it, the original unit label "unit" will be consistent with that on CES Console and Help Center documents. A new unit label "unit_v2" has been introduced to establish a new standard for metric units across various Huawei Cloud services, making them more readable and understandable. This makes it easier for users to select units when configuring Grafana dashboards.
  i18n_config_file_path: "/root/i18n.json" # This is optional.The ability to optimize metric data labels depends on the CES I18N internationalization configuration. By default, the exporter retrieves the internationalization configuration from the CES service I18N interface. If the interface is unavailable, it is recommended to specify the local file path for this internationalization configuration.
  unit_standardization_file_path: "/root/unit_standard_transform.json" # The ability to optimize unit labels for indicator data depends on this file. When enabling the optimization feature, it is recommended to specify the path to this configuration file.
  metric_timestamp_export_enabled: false # This is optional. used to export the original timestamp of metrics. Default is false. When enabled, the metric timestamp field will be set to the actual collection time; otherwise, it will be set to the time when Prometheus requests metrics from the exporter.
  metric_query_duration: 10 # This is optional.The time interval for querying metric data. The default value is 10 minutes. You can configure this value according to the metric reporting period, with a configurable range of 10 to 60.
  
  # This is optional. You can replace the proxy schema, host and port in the example below according to actual situation.
  proxy_schema: "http"                      
  proxy_host: "proxy.huaweicloud.com"        
  proxy_port: 8080    
  # This is optional. If the proxy requires authentication, please configure the username and password
  proxy_username: "username"                 
  proxy_password: "password" 
auth:
  auth_url: "https://iam.{region_id}.myhuaweicloud.com/v3"
  project_name: "{project_name}"
  access_key: "{access_key}" # It is strongly remommended that you use a script to decrypt the AK/SK by following the instructions provided in section 4.1 to prevent information leakage caused by plaintext AK/SK configuration in the configuration file.
  secret_key: "{secret_key}"
  region: "{region}"
```

## Example of endpoint file(endpoints.yml)
European site users should configure the domain names of the rms and eps services as follows.
```
"rms":
  "https://rms.myhuaweicloud.eu"
"eps":
  "https://eps.eu-west-101.myhuaweicloud.eu"
```
Other site users should not configure anything in this file.

## Prometheus Configuration
The huaweicloud exporter needs to be passed the address as a parameter, this can be done with relabelling.

Example config:

```
global:
  scrape_interval: 1m # Set the scrape interval to every 1 minute seconds. Default is every 1 minute.
  scrape_timeout: 1m
scrape_configs:
  - job_name: 'huaweicloud'
    static_configs:
    - targets: ['10.0.0.10:8087']
    params:
      services: ['SYS.VPC,SYS.ELB']
```

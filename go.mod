module github.com/armory-io/terraform-provider-spinnaker

go 1.12

require (
	github.com/apache/thrift v0.12.0 // indirect
	github.com/ghodss/yaml v1.0.0
	github.com/hashicorp/terraform v0.12.0
	github.com/hashicorp/yamux v0.0.0-20181012175058-2f1d1f20f75d // indirect
	github.com/mitchellh/mapstructure v1.1.2
	github.com/openzipkin/zipkin-go v0.1.6 // indirect
	github.com/prometheus/client_golang v0.9.3-0.20190127221311-3c4408c8b829 // indirect
	github.com/spf13/cobra v0.0.5
	github.com/spf13/pflag v1.0.3
	github.com/spinnaker/spin v0.4.1-0.20200818005041-2b364bb5214e
)

replace git.apache.org/thrift.git => github.com/apache/thrift v0.0.0-20180902110319-2566ecd5d999

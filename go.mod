module github.com/tidal-engineering/terraform-provider-spinnaker

go 1.14

require (
	github.com/ghodss/yaml v1.0.0
	github.com/hashicorp/terraform-plugin-sdk v1.7.0
	github.com/mitchellh/mapstructure v1.1.2
	github.com/spf13/pflag v1.0.3
	github.com/spinnaker/spin v0.0.0-20190530150642-535d2dc1b985
	golang.org/x/crypto v0.0.0-20191011191535-87dc89f01550 // indirect
)

replace git.apache.org/thrift.git => github.com/apache/thrift v0.0.0-20180902110319-2566ecd5d999

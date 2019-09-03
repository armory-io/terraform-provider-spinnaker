module github.com/armory-io/terraform-provider-spinnaker

go 1.12

require (
	cloud.google.com/go v0.37.4 // indirect
	github.com/cenkalti/backoff v2.2.1+incompatible
	github.com/ghodss/yaml v1.0.0
	github.com/hashicorp/terraform v0.12.0
	github.com/hashicorp/yamux v0.0.0-20181012175058-2f1d1f20f75d // indirect
	github.com/mitchellh/mapstructure v1.1.2
	github.com/spf13/pflag v1.0.3
	github.com/spinnaker/spin v0.0.0-20190530150642-535d2dc1b985
)

replace git.apache.org/thrift.git => github.com/apache/thrift v0.0.0-20180902110319-2566ecd5d999

module github.com/dmyerscough/terraform-provider-spinnaker

go 1.12

require (
	github.com/armory-io/terraform-provider-spinnaker v0.0.0-20190611222336-b623116edce1
	github.com/ghodss/yaml v1.0.0
	github.com/hashicorp/terraform v0.12.9
	github.com/kr/pty v1.1.3 // indirect
	github.com/mitchellh/mapstructure v1.1.2
	github.com/spf13/pflag v1.0.3
	github.com/spinnaker/spin v0.4.1-0.20190923211306-6cd29799b508
)

replace git.apache.org/thrift.git => github.com/apache/thrift v0.0.0-20180902110319-2566ecd5d999

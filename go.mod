module github.com/armory-io/terraform-provider-spinnaker

go 1.12

require (
	cloud.google.com/go v0.37.4 // indirect
	github.com/ghodss/yaml v1.0.0
	github.com/hashicorp/terraform v0.12.0
	github.com/hashicorp/yamux v0.0.0-20181012175058-2f1d1f20f75d // indirect
	github.com/mitchellh/mapstructure v1.1.2
	github.com/spf13/pflag v1.0.3
	github.com/spinnaker/spin v0.4.1-0.20190514235037-ce3904e37aaa
	k8s.io/client-go v11.0.0+incompatible // indirect
)

replace git.apache.org/thrift.git => github.com/apache/thrift v0.0.0-20180902110319-2566ecd5d999

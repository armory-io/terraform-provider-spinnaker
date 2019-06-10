module github.com/armory-io/terraform-provider-spinnaker

go 1.12

require (
	cloud.google.com/go v0.39.0
	github.com/agext/levenshtein v1.2.2
	github.com/apparentlymart/go-cidr v1.0.0
	github.com/apparentlymart/go-textseg v1.0.0
	github.com/armon/go-radix v1.0.0
	github.com/aws/aws-sdk-go v1.19.42
	github.com/bgentry/go-netrc v0.0.0-20140422174119-9fd32a8b3d3d
	github.com/bgentry/speakeasy v0.1.0
	github.com/blang/semver v0.0.0-20190414102917-ba2c2ddd8906
	github.com/estebangarcia/spin v0.0.0
	github.com/fatih/color v1.7.0
	github.com/ghodss/yaml v1.0.0
	github.com/golang/protobuf v1.3.1
	github.com/google/go-cmp v0.3.0
	github.com/googleapis/gax-go v2.0.0+incompatible
	github.com/hashicorp/errwrap v1.0.0
	github.com/hashicorp/go-cleanhttp v0.5.1
	github.com/hashicorp/go-getter v1.3.0
	github.com/hashicorp/go-hclog v0.9.2
	github.com/hashicorp/go-multierror v1.0.0
	github.com/hashicorp/go-plugin v1.0.0
	github.com/hashicorp/go-safetemp v1.0.0
	github.com/hashicorp/go-uuid v1.0.1
	github.com/hashicorp/go-version v1.2.0
	github.com/hashicorp/golang-lru v0.5.1
	github.com/hashicorp/hcl v1.0.0
	github.com/hashicorp/hcl2 v0.0.0-20190515223218-4b22149b7cef
	github.com/hashicorp/hil v0.0.0-20190212132231-97b3a9cdfa93
	github.com/hashicorp/terraform v0.11.14
	github.com/hashicorp/yamux v0.0.0-20181012175058-2f1d1f20f75d
	github.com/jmespath/go-jmespath v0.0.0-20180206201540-c2b33e8439af
	github.com/kr/pty v1.1.3 // indirect
	github.com/mattn/go-colorable v0.0.9
	github.com/mattn/go-isatty v0.0.8
	github.com/mitchellh/cli v1.0.0
	github.com/mitchellh/colorstring v0.0.0-20190213212951-d06e56a500db
	github.com/mitchellh/copystructure v1.0.0
	github.com/mitchellh/go-homedir v1.1.0
	github.com/mitchellh/go-testing-interface v1.0.0
	github.com/mitchellh/go-wordwrap v1.0.0
	github.com/mitchellh/hashstructure v1.0.0
	github.com/mitchellh/mapstructure v1.1.2
	github.com/mitchellh/reflectwalk v1.0.1
	github.com/oklog/run v1.0.0
	github.com/posener/complete v1.2.1
	github.com/spf13/pflag v1.0.3
	github.com/spf13/viper v1.4.0 // indirect
	github.com/spinnaker/spin v0.4.0 // indirect
	github.com/ulikunitz/xz v0.5.6
	github.com/zclconf/go-cty v0.0.0-20190516203816-4fecf87372ec
	go.opencensus.io v0.22.0
	golang.org/x/crypto v0.0.0-20190530122614-20be4c3c3ed5
	golang.org/x/exp v0.0.0-20190221220918-438050ddec5e // indirect
	golang.org/x/net v0.0.0-20190603091049-60506f45cf65
	golang.org/x/oauth2 v0.0.0-20190604053449-0f29369cfe45
	golang.org/x/sys v0.0.0-20190602015325-4c4f7f33c9ed
	golang.org/x/text v0.3.2
	google.golang.org/api v0.5.0
	google.golang.org/appengine v1.6.0
	google.golang.org/genproto v0.0.0-20190530194941-fb225487d101
	google.golang.org/grpc v1.21.0
	gopkg.in/yaml.v2 v2.2.2
	k8s.io/client-go v11.0.0+incompatible
)

replace github.com/estebangarcia/spin v0.0.0 => ./vendor/github.com/estebangarcia/spin

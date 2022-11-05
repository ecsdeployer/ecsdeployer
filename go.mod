module ecsdeployer.com/ecsdeployer

go 1.19

require (
	github.com/Masterminds/semver/v3 v3.1.1
	github.com/aws/aws-sdk-go-v2 v1.17.1
	github.com/aws/aws-sdk-go-v2/config v1.17.10
	github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs v1.15.20
	github.com/aws/aws-sdk-go-v2/service/ec2 v1.59.0
	github.com/aws/aws-sdk-go-v2/service/ecs v1.18.26
	github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2 v1.18.19
	github.com/aws/aws-sdk-go-v2/service/eventbridge v1.16.15
	github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi v1.13.19
	github.com/aws/aws-sdk-go-v2/service/ssm v1.28.1
	github.com/aws/aws-sdk-go-v2/service/sts v1.17.1
	github.com/caarlos0/log v0.1.6
	github.com/iancoleman/orderedmap v0.0.0-20190318233801-ac98e3ecb4b0
	github.com/invopop/jsonschema v0.6.0
	github.com/muesli/mango-cobra v1.2.0
	github.com/muesli/roff v0.1.0
	github.com/spf13/cobra v1.5.0
	github.com/stretchr/testify v1.8.0
	github.com/webdestroya/awsmocker v0.1.3
	github.com/withfig/autocomplete-tools/integrations/cobra v1.2.1
	github.com/xeipuuv/gojsonschema v1.2.0
	golang.org/x/exp v0.0.0-20220916125017-b168a2c6b86b
	golang.org/x/sync v0.0.0-20220923202941-7f9b1623fab7
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/aws/aws-sdk-go-v2/credentials v1.12.23 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.12.19 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.1.25 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.4.19 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.3.26 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.0.14 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.9.19 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.11.25 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.13.8 // indirect
	github.com/aws/smithy-go v1.13.4 // indirect
	github.com/aymanbagabas/go-osc52 v1.0.3 // indirect
	github.com/charmbracelet/lipgloss v0.6.1-0.20220911181249-6304a734e792 // indirect
	github.com/clbanning/mxj v1.8.4 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.2 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/elazarl/goproxy v0.0.0-20221015165544-a0805db90819 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/lucasb-eyer/go-colorful v1.2.0 // indirect
	github.com/mattn/go-isatty v0.0.16 // indirect
	github.com/mattn/go-runewidth v0.0.13 // indirect
	github.com/muesli/mango v0.1.0 // indirect
	github.com/muesli/mango-pflag v0.1.0 // indirect
	github.com/muesli/reflow v0.3.0 // indirect
	github.com/muesli/termenv v0.12.1-0.20220901123159-d729275e0977 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rivo/uniseg v0.4.2 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/xeipuuv/gojsonpointer v0.0.0-20180127040702-4e3ac2762d5f // indirect
	github.com/xeipuuv/gojsonreference v0.0.0-20180127040603-bd5ef7bd5415 // indirect
	golang.org/x/sys v0.0.0-20220909162455-aba9fc2a8ff2 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

replace github.com/invopop/jsonschema v0.6.0 => github.com/webdestroya/jsonschema v0.0.0-20221009071543-bd9a93154641

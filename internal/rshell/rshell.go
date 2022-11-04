package rshell

import "ecsdeployer.com/ecsdeployer/internal/util"

// TODO: This should be imported from https://github.com/webdestroya/remote-shell-client

const LabelName = "cloud87.rshell"

type DockerLabel struct {
	Cluster          string   `json:"cluster"`
	SubnetIds        []string `json:"subnets"`
	SecurityGroupIds []string `json:"security_groups"`
	AssignPublicIp   bool     `json:"public"`
	Port             int32    `json:"port"`
	Path             string   `json:"path,omitempty"`
}

func (obj *DockerLabel) ToJSON() string {
	return util.Must(util.Jsonify(obj))
}

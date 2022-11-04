//go:build generate
// +build generate

package main

import (
	"encoding/json"
	"fmt"
	"os"

	"ecsdeployer.com/ecsdeployer/pkg/config"
)

type cond struct {
	test     string
	variable string
	value    string
	values   []string
}

type IamDoc struct {
	Version   string      `json:"Version"`
	Statement []Statement `json:"Statement"`
}

type Statement struct {
	Sid       string                            `json:"Sid,omitempty"`
	Effect    string                            `json:"Effect,omitempty"`
	Action    interface{}                       `json:"Action,omitempty"`
	Resource  interface{}                       `json:"Resource,omitempty"`
	Condition map[string]map[string]interface{} `json:"Condition,omitempty"`
}

type stmt struct {
	id      string
	action  string
	actions []string
	res     string
	reslist []string
	cond    *cond
	conds   []cond
}

func (s *stmt) toPretty() *Statement {
	st := &Statement{
		Sid:    s.id,
		Effect: "Allow",
	}

	if s.res != "" {
		st.Resource = s.res
	} else if len(s.reslist) > 0 {
		st.Resource = s.reslist
	} else {
		st.Resource = "*"
	}

	if s.action != "" {
		st.Action = s.action
	} else if len(s.actions) > 0 {
		st.Action = s.actions
	} else {
		panic("missing action")
	}

	condList := make([]cond, 0, 10)
	if s.cond != nil {
		condList = append(condList, *s.cond)
	}

	if len(s.conds) > 0 {
		condList = append(condList, s.conds...)
	}

	if len(condList) > 0 {
		finalCondlist := make(map[string]map[string]interface{})
		for _, c := range condList {
			_, exists := finalCondlist[c.test]
			if !exists {
				finalCondlist[c.test] = make(map[string]interface{})
			}

			if c.value != "" {
				finalCondlist[c.test][c.variable] = c.value
			} else {
				finalCondlist[c.test][c.variable] = c.values
			}
		}
		st.Condition = finalCondlist
	}

	return st
}

const (
	clusterArn = "arn:aws:ecs:REGION:ACCOUNTID:cluster/CLUSTER_NAME"
	serviceArn = "arn:aws:ecs:REGION:ACCOUNTID:service/CLUSTER_NAME/PROJECT_NAME-*"
)

var (
	clusterCondition = cond{
		test:     "ArnEquals",
		variable: "ecs:cluster",
		value:    clusterArn,
	}

	commonStatements = []stmt{
		{
			id: "EnvInfoGathering",
			actions: []string{
				"ec2:DescribeVpcs",
				"ec2:DescribeSubnets",
				"ec2:DescribeSecurityGroups",
				"elasticloadbalancing:DescribeTargetGroups",
				"logs:DescribeLogGroups",
				"tag:GetResources",
			},
		},
		{
			id: "TaskDefinitions",
			actions: []string{
				"ecs:RegisterTaskDefinition",
				"ecs:DeregisterTaskDefinition",

				"ecs:ListTaskDefinitionFamilies",
				"ecs:ListTaskDefinitions",
				"ecs:DescribeTaskDefinition",
			},
		},
		// {
		// 	id:     "ImageVerification",
		// 	action: "ecr:DescribeImages",
		// },
	}
)

func generateIAMRestricted() []stmt {

	nameTemplates := &config.NameTemplates{}
	nameTemplates.ApplyDefaults()

	markerTagKey := *nameTemplates.MarkerTagKey

	return append(commonStatements, []stmt{
		{
			actions: []string{
				"ecs:DescribeTasks",
				"ecs:DescribeServices",
			},
			cond: &clusterCondition,
		},
		{
			id: "ServiceDeployment",
			actions: []string{
				"ecs:UpdateService",
				"ecs:DeleteService",
			},
			res:  serviceArn,
			cond: &clusterCondition,
		},
		{
			id: "ServiceCreation",
			actions: []string{
				"ecs:CreateService",
			},
			res: serviceArn,
			conds: []cond{
				clusterCondition,
				{
					test:     "StringEquals",
					variable: "aws:RequestTag/" + markerTagKey,
					values: []string{
						"PROJECT_NAME",
					},
				},
			},
		},
		{
			id: "ResourceTaggingECS",
			actions: []string{
				"ecs:TagResource",
			},
			reslist: []string{
				serviceArn,
				"arn:aws:ecs:REGION:ACCOUNTID:task-definition/PROJECT_NAME-*",
			},
		},
		{
			id: "PreDeployTasks",
			actions: []string{
				"ecs:RunTask",
			},
			cond: &clusterCondition,
		},
		{
			id: "ImportSSMParams",
			actions: []string{
				"ssm:GetParametersByPath",
			},
			res: "arn:aws:ssm:*:*:parameter/ecsdeployer/secrets/PROJECT_NAME/*",
		},
		{
			id: "Logging",
			actions: []string{
				"logs:CreateLogGroup",
				"logs:PutRetentionPolicy",
				"logs:TagLogGroup",
			},
			reslist: []string{
				"arn:aws:logs:REGION:ACCOUNTID:log-group:/ecsdeployer/app/PROJECT_NAME*",
				"arn:aws:logs:REGION:ACCOUNTID:log-group:/ecsdeployer/app/PROJECT_NAME*:*",
			},
		},
		{
			id: "CronSetup",
			actions: []string{
				"events:PutRule",
				"events:DeleteRule",
				"events:RemoveTargets",
				"events:ListTargetsByRule",
				"events:TagResource",
				"events:UntagResource",
			},
			reslist: []string{
				"arn:aws:events:REGION:ACCOUNTID:rule/PROJECT_NAME-*",
			},
		},
		{
			id:     "CronSetupTargets",
			action: "events:PutTargets",
			res:    "arn:aws:events:REGION:ACCOUNTID:rule/PROJECT_NAME-*",
			cond: &cond{
				test:     "ArnEquals",
				variable: "events:TargetArn",
				values:   []string{clusterArn},
			},
		},
		{
			id:     "EcsPassRole",
			action: "iam:PassRole",
			reslist: []string{
				"arn:aws:iam::ACCOUNTID:role/APP_ROLE",
				"arn:aws:iam::ACCOUNTID:role/ECS_EXECUTION_ROLE",
			},
			cond: &cond{
				test:     "StringLike",
				variable: "iam:PassedToService",
				values:   []string{"ecs-tasks.amazonaws.com"},
			},
		},
		{
			id:     "CronPassRole",
			action: "iam:PassRole",
			reslist: []string{
				"arn:aws:iam::ACCOUNTID:role/CRON_LAUNCHER_ROLE",
			},
			cond: &cond{
				test:     "StringLike",
				variable: "iam:PassedToService",
				values:   []string{"events.amazonaws.com"},
			},
		},
	}...)
}

func generateIAMNormal() []stmt {
	return append(commonStatements, []stmt{
		{
			actions: []string{
				"ecs:DescribeTasks",
				"ecs:DescribeServices",
				"ecs:UpdateService",
				"ecs:DeleteService",
				"ecs:CreateService",
				"ecs:RunTask",
			},
		},
		{
			id: "CronSetup",
			actions: []string{
				"events:PutRule",
				"events:DeleteRule",
				"events:RemoveTargets",
				"events:ListTargetsByRule",
				"events:PutTargets",
			},
		},
		{
			id: "ImportSSMParams",
			actions: []string{
				"ssm:GetParametersByPath",
			},
		},
		{
			id: "Logging",
			actions: []string{
				"logs:CreateLogGroup",
				"logs:PutRetentionPolicy",
			},
		},
		{
			id: "TagManagement",
			actions: []string{
				"logs:TagLogGroup",
				"ecs:TagResource",
				"events:TagResource",
				"events:UntagResource",
			},
		},
		{
			id:     "RolePassing",
			action: "iam:PassRole",
			cond: &cond{
				test:     "StringLike",
				variable: "iam:PassedToService",
				values: []string{
					"ecs-tasks.amazonaws.com",
					"events.amazonaws.com",
				},
			},
		},
	}...)
}

func makePrettyStatements(stmts []stmt) []Statement {
	pretty := make([]Statement, 0, len(stmts))

	dupeIds := make(map[string]struct{})

	for _, stmt := range stmts {
		pr := stmt.toPretty()
		if pr == nil {
			continue
		}

		if pr.Sid != "" {
			_, exists := dupeIds[pr.Sid]

			if exists {
				panic(fmt.Errorf("Duplicate SID: %s", pr.Sid))
			}

			dupeIds[pr.Sid] = struct{}{}
		}

		pretty = append(pretty, *pr)
	}

	return pretty
}

func exportPolicy(stmts []Statement, filename string) error {

	policy := IamDoc{
		Version:   "2012-10-17",
		Statement: stmts,
	}

	bts, err := json.MarshalIndent(policy, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to create policy: %w", err)
	}

	if err := os.WriteFile(filename, bts, 0o644); err != nil {
		return fmt.Errorf("failed to write policy file: %w", err)
	}

	fmt.Printf("Wrote policy file to: %s\n", filename)

	return nil
}

func main() {
	exportPolicy(makePrettyStatements(generateIAMNormal()), "./data/iam/lax.json")
	exportPolicy(makePrettyStatements(generateIAMRestricted()), "./data/iam/restricted.json")
}

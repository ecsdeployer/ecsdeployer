package taskmock

import (
	"ecsdeployer.com/ecsdeployer/internal/util"
)

func jsonify(obj any) string {
	return util.Must(util.Jsonify(obj))
}

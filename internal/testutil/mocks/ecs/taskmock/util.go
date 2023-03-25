package taskmock

import (
	"ecsdeployer.com/ecsdeployer/internal/util"
)

func jsonify(obj interface{}) string {
	return util.Must(util.Jsonify(obj))
}

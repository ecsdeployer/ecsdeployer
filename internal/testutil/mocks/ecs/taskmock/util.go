package taskmock

import (
	"crypto/rand"
	"encoding/hex"

	"ecsdeployer.com/ecsdeployer/internal/util"
)

func randomHex(n int) string {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		panic(err)
	}
	return hex.EncodeToString(bytes)
}

func jsonify(obj interface{}) string {
	return util.Must(util.Jsonify(obj))
}

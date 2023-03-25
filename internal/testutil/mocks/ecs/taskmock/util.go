package taskmock

import (
	"crypto/rand"
	"encoding/hex"

	"ecsdeployer.com/ecsdeployer/internal/util"
)

// number must always be even
func randomHex(n int) string {
	bytes := make([]byte, (n+1)/2)
	if _, err := rand.Read(bytes); err != nil {
		panic(err)
	}
	return hex.EncodeToString(bytes)
}

func jsonify(obj interface{}) string {
	return util.Must(util.Jsonify(obj))
}

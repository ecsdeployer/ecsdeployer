package version

import (
	"fmt"

	hcVersion "github.com/hashicorp/go-version"
)

// This is the default version used locally
const DevVersionID = "9999.0.0"

// the semantic version of ECS Deployer
var Version string = DevVersionID

// the commit sha used for building
var BuildSHA string = "master"

// the commit sha used for building
var ShortSHA string = "master"

var Prerelease string = "dev"

var SemVer *hcVersion.Version

func init() {
	SemVer = hcVersion.Must(hcVersion.NewVersion(Version))
}

// String returns the complete version string, including prerelease
func String() string {
	if IsPrelease() {
		return fmt.Sprintf("%s-%s", Version, Prerelease)
	}
	return Version
}

func IsPrelease() bool {
	return Prerelease != ""
}

package steps

import (
	"fmt"
)

type OldStep interface {
	fmt.Stringer
}

// type oldStep_PreApplyer interface {
// 	PreApply(*config.Context) error
// }

// type oldStep_Creator interface {
// 	Create(*config.Context) error
// }

// type oldStep_Updater interface {
// 	Update(*config.Context) error
// }

// type oldStep_Depender interface {
// 	Dependencies(*config.Context) []OldStep
// }

// type oldStep_Paralleler interface {
// 	ParallelDeps() bool
// }

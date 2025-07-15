package dns

import (
	"fmt"

	"ziniki.org/deployer/driver/pkg/errorsink"
)

type cnameModel struct {
	loc  *errorsink.Location
	name string

	pointsTo fmt.Stringer
}

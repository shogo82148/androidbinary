// +build go1.13

package apk

import (
	"errors"
	"fmt"
)

var newError = errors.New
var errorf = fmt.Errorf

// +build !go1.13

package apk

// Error wrapping comes from Go 1.13
// use a compatibility package
import "golang.org/x/xerrors"

var newError = xerrors.New
var errorf = xerrors.Errorf

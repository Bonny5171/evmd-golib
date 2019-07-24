package errors

import (
	"fmt"
	"runtime"

	"github.com/pkg/errors"
)

func Trace(err error) error {
	pc := make([]uintptr, 10) // at least 1 entry needed
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	file, line := f.FileLine(pc[0])
	return errors.Wrap(err, fmt.Sprintf("%s:%d %s\n", file, line-1, f.Name()))
}

package interp2d

import (
	"fmt"
	"unsafe"
)

/*
#cgo pkg-config: gsl
#include <gsl/gsl_errno.h>

void errHandler(char *reason, char *file, int line, int gsl_errno);
*/
import "C"

var (
	lastError error // hello?
)

func init() {
	C.gsl_set_error_handler((*C.gsl_error_handler_t)(unsafe.Pointer(C.errHandler)))
}

//export errHandler
func errHandler(reason *C.char, file *C.char, line C.int, gsl_errno C.int) {
	lastError = fmt.Errorf("gsl: %s: %s", C.GoString(C.gsl_strerror(gsl_errno)), C.GoString(reason))
}

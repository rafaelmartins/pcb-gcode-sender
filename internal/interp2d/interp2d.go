package interp2d

import (
	"errors"
	"math"

	"github.com/rafaelmartins/pcb-gcode-sender/internal/point"
)

// FIXME: port to pure go, maybe? ¯\_(ツ)_/¯

// #cgo pkg-config: gsl
// #include <gsl/gsl_errno.h>
// #include <gsl/gsl_math.h>
// #include <gsl/gsl_interp2d.h>
// #include <gsl/gsl_spline2d.h>
import "C"

type Spline struct {
	ptr *C.gsl_spline2d

	xacc *C.gsl_interp_accel
	yacc *C.gsl_interp_accel

	closed bool
}

func NewSpline(points [][]*point.Point) (*Spline, error) {
	ny := len(points)
	if ny == 0 {
		return nil, errors.New("interp2d: no points")
	}

	nx := -1
	ya := []float64{}
	for _, lp := range points {
		if nx == -1 {
			nx = len(lp)
		} else {
			if nx != len(lp) {
				return nil, errors.New("interp2d: invalid matrix of points")
			}

			if nx < 1 {
				return nil, errors.New("interp2d: no points")
			}
		}

		ya = append(ya, lp[0].Y)
	}

	xa := []float64{}
	for _, p := range points[0] {
		xa = append(xa, p.X)
	}

	// min_size could be grabbed from the type
	typ := C.gsl_interp2d_bicubic
	if nx < 4 || ny < 4 {
		typ = C.gsl_interp2d_bilinear
	} else if nx < 2 || ny < 2 {
		return nil, errors.New("interp2d: matrix of points must be at least 2x2")
	}

	spline := C.gsl_spline2d_alloc(typ, C.size_t(nx), C.size_t(ny))
	if spline == nil {
		return nil, lastError
	}

	za := make([]float64, nx*ny)
	for j, lp := range points {
		for i, p := range lp {
			C.gsl_spline2d_set(spline, (*C.double)(&za[0]), C.size_t(i), C.size_t(j), C.double(p.Z))
		}
	}

	if C.GSL_SUCCESS != C.gsl_spline2d_init(spline, (*C.double)(&xa[0]), (*C.double)(&ya[0]), (*C.double)(&za[0]), C.size_t(nx), C.size_t(ny)) {
		return nil, lastError
	}

	xacc := C.gsl_interp_accel_alloc()
	if xacc == nil {
		return nil, lastError
	}

	yacc := C.gsl_interp_accel_alloc()
	if yacc == nil {
		return nil, lastError
	}

	return &Spline{
		ptr:  spline,
		xacc: xacc,
		yacc: yacc,
	}, nil
}

func (s *Spline) Close() {
	if s == nil || s.closed {
		return
	}

	if s.ptr != nil {
		C.gsl_spline2d_free(s.ptr)
	}

	if s.xacc != nil {
		C.gsl_interp_accel_free(s.xacc)
	}

	if s.yacc != nil {
		C.gsl_interp_accel_free(s.yacc)
	}

	s.closed = true
}

func (s *Spline) At(x float64, y float64) (float64, error) {
	if s == nil || s.closed {
		return 0, errors.New("interp2d: invalid spline")
	}

	rv := float64(C.gsl_spline2d_eval(s.ptr, C.double(x), C.double(y), s.xacc, s.yacc))
	if math.IsNaN(rv) {
		return 0, lastError
	}
	return rv, nil
}

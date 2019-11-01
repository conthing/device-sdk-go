package transformer

import (
	"fmt"
)

type OverflowError struct {
	origin      interface{}
	transformed float64
}

func (e OverflowError) Error() string {
	return fmt.Sprintf("overflow failed, transformed value '%v' is not within the '%T' value type range", e.transformed, e.origin)
}

func (e OverflowError) String() string {
	return fmt.Sprintf("overflow failed, transformed value '%v' is not within the '%T' value type range", e.transformed, e.origin)
}

func NewOverflowError(origin interface{}, transformed float64) OverflowError {
	return OverflowError{origin: origin, transformed: transformed}
}

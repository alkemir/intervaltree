package intervaltree

import "fmt"

// OverlapError is returned whenever an Insert() call tries to insert a value
// previously inserted.
type OverlapError uint64

func (e OverlapError) Error() string {
	return fmt.Sprintf("Tried to insert value already inserted: %d", uint64(e))
}

// InvalidIntervalError is returned whenever an Insert() call tries to insert a
// interval [x, y] where x > y.
type InvalidIntervalError struct {
	x uint64
	y uint64
}

func (e InvalidIntervalError) Error() string {
	return fmt.Sprintf("Invalid interval: [%d, %d]", e.x, e.y)
}

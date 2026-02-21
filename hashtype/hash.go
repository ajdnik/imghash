package hashtype

// Hash is the common interface for all hash representations.
// It is implemented by Binary, UInt8, and Float64.
type Hash interface {
	String() string
	Len() int
	ValueAt(idx int) float64
}

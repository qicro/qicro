package code

//go:generate codegen -type=int
//go:generate codegen -type=int -doc -output ./error_code_generated.md

const (
	// ErrUserNotFound - 404: User not found.
	ErrUserNotFound int = iota + 100401
)

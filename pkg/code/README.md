
```
go install tools/codegen/codegen.go
```

eg. user.go

```
//go:generate codegen -type=int -doc -output ./error_code_generated.md

const (
	// ErrUserNotFound - 404: User not found
	ErrUserNotFound int = iota + 100401
)
```
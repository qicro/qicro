// Package errors provides simple error handling primitives.
//
// The traditional error handling idiom in Go is roughly akin to
//
//	if err != nil {
//	        return err
//	}
//
// which when applied recursively up the call stack results in error reports
// without context or debugging information. The errors package allows
// programmers to add context to the failure path in their code in a way
// that does not destroy the original value of the error.
//
// # Adding context to an error
//
// The errors.Wrap function returns a new error that adds context to the
// original error by recording a stack trace at the point Wrap is called,
// together with the supplied message. For example
//
//	_, err := ioutil.ReadAll(r)
//	if err != nil {
//	        return errors.Wrap(err, "read failed")
//	}
//
// If additional control is required, the errors.WithStack and
// errors.WithMessage functions destructure errors.Wrap into its component
// operations: annotating an error with a stack trace and with a message,
// respectively.
//
// # Retrieving the cause of an error
//
// Using errors.Wrap constructs a stack of errors, adding context to the
// preceding error. Depending on the nature of the error it may be necessary
// to reverse the operation of errors.Wrap to retrieve the original error
// for inspection. Any error value which implements this interface
//
//	type causer interface {
//	        Cause() error
//	}
//
// can be inspected by errors.Cause. errors.Cause will recursively retrieve
// the topmost error that does not implement causer, which is assumed to be
// the original cause. For example:
//
//	switch err := errors.Cause(err).(type) {
//	case *MyError:
//	        // handle specifically
//	default:
//	        // unknown error
//	}
//
// Although the causer interface is not exported by this package, it is
// considered a part of its stable public interface.
//
// # Formatted printing of errors
//
// All error values returned from this package implement fmt.Formatter and can
// be formatted by the fmt package. The following verbs are supported:
//
//	%s    print the error. If the error has a Cause it will be
//	      printed recursively.
//	%v    see %s
//	%+v   extended format. Each Frame of the error's StackTrace will
//	      be printed in detail.
//
// # Retrieving the stack trace of an error or wrapper
//
// New, Errorf, Wrap, and Wrapf record a stack trace at the point they are
// invoked. This information can be retrieved with the following interface:
//
//	type stackTracer interface {
//	        StackTrace() errors.StackTrace
//	}
//
// The returned errors.StackTrace type is defined as
//
//	type StackTrace []Frame
//
// The Frame type represents a call site in the stack trace. Frame supports
// the fmt.Formatter interface that can be used for printing information about
// the stack trace of this error. For example:
//
//	if err, ok := err.(stackTracer); ok {
//	        for _, f := range err.StackTrace() {
//	                fmt.Printf("%+s:%d\n", f, f)
//	        }
//	}
//
// Although the stackTracer interface is not exported by this package, it is
// considered a part of its stable public interface.
//
// See the documentation for Frame.Format for more details.
package errors

import (
	"fmt"
	"io"

	gCodes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// New 返回一条带有指定信息的错误信息
// New 同时记录调用它时的堆栈跟踪信息
func New(message string) error {
	return &fundamental{
		msg:   message,
		stack: callers(),
	}
}

// Errorf 根据格式说明符进行格式化并返回一个字符串，该字符串作为满足 error 的值。
// 在调用 Errorf 时也会记录堆栈跟踪。
func Errorf(format string, args ...interface{}) error {
	return &fundamental{
		msg:   fmt.Sprintf(format, args...),
		stack: callers(),
	}
}

// fundamental includes a message and stack, but not caller's error type
type fundamental struct {
	msg string
	*stack
}

func (f *fundamental) Error() string { return f.msg }

func (f *fundamental) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			io.WriteString(s, f.msg)
			f.stack.Format(s, verb)
			return
		}
		fallthrough
	case 's':
		io.WriteString(s, f.msg)
	case 'q':
		fmt.Fprintf(s, "%q", f.msg)
	}
}

// FromGrpcError grpcError 转换 Error
func FromGrpcError(err error) error {
	if err == nil {
		return err
	}
	st, ok := status.FromError(err)
	if !ok {
		return WithCode(100002, "unknown error")
	}

	return &withCode{
		err:  st.Err(),
		code: int(st.Code()),
	}
}

// ToGrpcError convert Error to grpcError
func ToGrpcError(err error) error {
	if err == nil {
		return err
	}
	// TODO: 只实现了 withCode 转 grpc 方法
	var perr *withCode
	// 使用 go 基本库 errors.As 方法 如果是同一类型的 error
	if As(err, &perr) {
		err := status.Error(gCodes.Code(perr.code), perr.err.Error())
		return err
	}
	return status.Error(gCodes.Unknown, err.Error())
}

// WithStack 函数会在调用 WithStack 的时候，为 err 添加一个堆栈跟踪。
// 如果 err 为 nil，则 WithStack 返回 nil。
func WithStack(err error) error {
	if err == nil {
		return nil
	}

	if e, ok := err.(*withCode); ok {
		return &withCode{
			err:   e.err,
			code:  e.code,
			cause: err,
			stack: callers(),
		}
	}

	return &withStack{
		err,
		callers(),
	}
}

type withStack struct {
	error
	*stack
}

func (w *withStack) Cause() error { return w.error }

// Unwrap 提供了Go 1.13错误链的兼容性。
func (w *withStack) Unwrap() error {
	if e, ok := w.error.(interface{ Unwrap() error }); ok {
		return e.Unwrap()
	}

	return w.error
}

func (w *withStack) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			fmt.Fprintf(s, "%+v", w.Cause())
			w.stack.Format(s, verb)
			return
		}
		fallthrough
	case 's':
		io.WriteString(s, w.Error())
	case 'q':
		fmt.Fprintf(s, "%q", w.Error())
	}
}

// Wrap 返回一个在调用Wrap时，使用堆栈跟踪注释err和提供的消息的错误。
// 如果err为nil，则Wrap返回nil。
func Wrap(err error, message string) error {
	if err == nil {
		return nil
	}
	if e, ok := err.(*withCode); ok {
		return &withCode{
			err:   fmt.Errorf(message),
			code:  e.code,
			cause: err,
			stack: callers(),
		}
	}

	err = &withMessage{
		cause: err,
		msg:   message,
	}
	return &withStack{
		err,
		callers(),
	}
}

// Wrapf 返回一个在调用Wrapf时，使用堆栈跟踪注释err和格式说明符的错误。
// 如果err为nil，则Wrapf返回nil。
func Wrapf(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}

	// 如果是 withCode 结构体
	if e, ok := err.(*withCode); ok {
		return &withCode{
			err:   fmt.Errorf(format, args...),
			code:  e.code,
			cause: err,
			stack: callers(),
		}
	}

	err = &withMessage{
		cause: err,
		msg:   fmt.Sprintf(format, args...),
	}
	return &withStack{
		err,
		callers(),
	}
}

// WithMessage 用新消息注释err。
// 如果err为nil，则WithMessage返回nil。
func WithMessage(err error, message string) error {
	if err == nil {
		return nil
	}
	return &withMessage{
		cause: err,
		msg:   message,
	}
}

// WithMessagef 使用格式说明符注释err
// 如果err为nil，则WithMessagef返回nil
func WithMessagef(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	return &withMessage{
		cause: err,
		msg:   fmt.Sprintf(format, args...),
	}
}

type withMessage struct {
	cause error
	msg   string
}

func (w *withMessage) Error() string { return w.msg }
func (w *withMessage) Cause() error  { return w.cause }

// Unwrap 提供了Go 1.13错误链的兼容性
func (w *withMessage) Unwrap() error { return w.cause }

func (w *withMessage) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			fmt.Fprintf(s, "%+v\n", w.Cause())
			io.WriteString(s, w.msg)
			return
		}
		fallthrough
	case 's', 'q':
		io.WriteString(s, w.Error())
	}
}

type withCode struct {
	err    error
	code   int
	cause  error
	*stack // 堆栈信息存储
}

// WithCode 生成 withCode 结构体
func WithCode(code int, format string, args ...interface{}) error {
	return &withCode{
		err:   fmt.Errorf(format, args...),
		code:  code,
		stack: callers(),
	}
}

// WrapC 生成 携带 子错误 的 withCode 结构体
func WrapC(err error, code int, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}

	return &withCode{
		err:   fmt.Errorf(format, args...),
		code:  code,
		cause: err,
		stack: callers(),
	}
}

// Error 返回外部安全的错误消息
func (w *withCode) Error() string { return fmt.Sprintf("%v", w) }

// Cause 返回带有代码的错误的原因
func (w *withCode) Cause() error { return w.cause }

// Unwrap 提供了Go 1.13错误链的兼容性
func (w *withCode) Unwrap() error { return w.cause }

// Cause 返回错误的根本原因，如果可能的话。
// 如果一个错误值实现了以下接口，则它有一个原因：
// interface:
//
//	type causer interface {
//	       Cause() error
//	}
//
// 如果错误没有实现Cause，则将返回原始错误。如果错误为nil，则不进行进一步的调查，直接返回nil。
func Cause(err error) error {
	type causer interface {
		Cause() error
	}

	for err != nil {
		cause, ok := err.(causer) // 判断 是否 err 结构体中 有 causer 方法
		if !ok {
			break
		}

		if cause.Cause() == nil { // 如果有方法 判断是否为空
			break
		}

		err = cause.Cause()
	}
	return err
}

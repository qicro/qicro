package errors

import (
	"fmt"
	"net/http"
	"sync"
)

var (
	unknownCoder defaultCoder = defaultCoder{1, http.StatusInternalServerError, "An internal server error occurred", "http://czc/pkg/errors/README.md"}
)

// Coder 定义了一个错误码详细信息的接口
type Coder interface {
	// HTTPStatus HTTP状态码，应该与相关的错误码一起使用
	HTTPStatus() int

	// String 面向外部用户的错误文本
	String() string

	// Reference 返回用户的详细文档
	Reference() string

	// Code 返回该错误码的代码
	Code() int
}

type defaultCoder struct {
	// C 代表ErrCode的整数代码
	C int

	// HTTP 状态，应用于相关的错误代码
	HTTP int

	// External（用户）可见的错误文本
	Ext string

	// Ref 指定参考文档。
	Ref string
}

// Code returns the integer code of the coder.
func (coder defaultCoder) Code() int {
	return coder.C

}

// String implements stringer. String returns the external error message,
// if any.
func (coder defaultCoder) String() string {
	return coder.Ext
}

// HTTPStatus returns the associated HTTP status code, if any. Otherwise,
// returns 200.
func (coder defaultCoder) HTTPStatus() int {
	if coder.HTTP == 0 {
		return 500
	}

	return coder.HTTP
}

// Reference returns the reference document.
func (coder defaultCoder) Reference() string {
	return coder.Ref
}

// codes 包含一个错误代码到元数据的映射表
var codes = map[int]Coder{}
var codeMux = &sync.Mutex{}

// Register 用于注册自定义错误码。
// 如果注册的错误码已经存在，则会覆盖之前的错误码
func Register(coder Coder) {
	// 如果 code 未赋值 直接 panic
	if coder.Code() == 0 {
		panic("code `0` is reserved by `czc/pkg/errors` as unknownCode error code")
	}

	codeMux.Lock()
	defer codeMux.Unlock()

	codes[coder.Code()] = coder
}

// MustRegister 注册一个用户定义的错误代码。
// 如果相同的 Code 已经存在，它将会 panic
func MustRegister(coder Coder) {
	// 如果 code 未赋值 直接 panic
	if coder.Code() == 0 {
		panic("code '0' is reserved by 'qicro/qicro/pkg/errors' as ErrUnknown error code")
	}
	// 上锁
	codeMux.Lock()
	defer codeMux.Unlock()

	if _, ok := codes[coder.Code()]; ok {
		panic(fmt.Sprintf("code: %d already exists", coder.Code()))
	}

	codes[coder.Code()] = coder
}

// ParseCoder 将任何错误解析为*withCode
// 空错误将直接返回nil
// 没有堆栈信息的错误将被解析为 ErrUnknown
func ParseCoder(err error) Coder {
	if err == nil {
		return nil
	}
	// 解析为统一类型 withCode
	if v, ok := err.(*withCode); ok {
		if coder, ok := codes[v.code]; ok {
			return coder
		}
	}

	return unknownCoder
}

// IsCode 函数会检查err链中是否包含给定的错误代码，如果有则返回true，否则返回false
func IsCode(err error, code int) bool {
	if v, ok := err.(*withCode); ok {
		if v.code == code {
			return true
		}

		if v.cause != nil {
			return IsCode(v.cause, code)
		}

		return false
	}

	return false
}

func init() {
	codes[unknownCoder.Code()] = unknownCoder
}

package jerrors

import (
	"fmt"

	"google.golang.org/grpc/codes"
	gs "google.golang.org/grpc/status"
)

type Errno struct {
	Code    int64
	Message string
}

func (e *Errno) Error() string {
	return e.Message
}

func (e *Errno) Add(msg string) *Errno {
	e.Message += " " + msg
	return e
}

func (e *Errno) Equal(o *Errno) bool {
	return e.Code == o.Code
}

//ToGRPC 转换为grpc错误
func (e *Errno) ToGRPC() error {
	s, _ := gs.New(codes.Unknown, "").WithDetails(&Error{
		Code:    e.Code,
		Message: e.Message,
	})
	return s.Err()
}

//ToErrno grpc转errno
//	返回值
//	Errno:具体错误
//	bool: 是否成功转换,如果server端返回的错误不是errno.ToGRPC可能不能转换
func FromGRPC(err error) (*Errno, bool) {
	var er Errno
	if e, ok := err.(interface {
		GRPCStatus() *gs.Status
	}); ok {
		s := e.GRPCStatus()
		for _, d := range s.Details() {
			v, ok := d.(*Error)
			if !ok {
				continue
			}
			er.Code = v.Code
			er.Message = v.Message
			return &er, true
		}
	}
	return &er, false
}

var codeCache = make(map[int64]struct{})

func NewError(code int64, msg string) *Errno {
	if _, ok := codeCache[code]; ok {
		panic(fmt.Sprintf("错误码:Code%d   Message:%s 已存在,请重新定义新code码", code, msg))
	}
	codeCache[code] = struct{}{}
	return &Errno{
		Code:    code,
		Message: msg,
	}
}

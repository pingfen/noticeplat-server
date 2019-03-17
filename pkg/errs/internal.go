package errs

const (
	INTERNAL_ERROR InternalErr = 1
	STORAGE_ERROR  InternalErr = 2
)

type InternalErr uint8

func (ie InternalErr) Error() string {
	return "internal error"
}

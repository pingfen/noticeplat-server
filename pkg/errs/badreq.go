package errs

import "fmt"

const (
	FIELD_EMPTY  BadRequest = 1
	FIELD_MISS   BadRequest = 2
	OBJECT_EXIST BadRequest = 3
)

type BadRequest uint8

func (br BadRequest) Error() string {
	return fmt.Sprintf("bad request, code %d", br)
}

package resp

type StringArray []string

type BoolResponse struct {
	*BaseResponse
	Data bool `json:"data"`
}

package code

type Code int

const (
	Success Code = 1000 + iota
	InvalidParams
	InvalidRequest
	BadRequest
	Unauthorized
)

var Messages = map[Code]string{
	Success:       "success",
	InvalidParams: "invalid params",
	BadRequest:    "bad request",
	Unauthorized:  "unauthorized",
}

func (c Code) GetMessage() string {
	msg, ok := Messages[c]
	if ok {
		return msg
	}
	return "unknown status"
}

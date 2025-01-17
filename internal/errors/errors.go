package errors

const (
	ErrUserAlreadyExists = "user already exists"
	ErrUserNotFound      = "user not found"
)

type AppError struct {
	Status     int
	BErrorText string
	UErrorText string
}
type ReqError struct {
	Status int `json:"status,omitempty" bson:"status"`

	Text string `json:"text" bson:"text"`
}

func (e AppError) Error() string {
	return e.BErrorText
}
func (e ReqError) Error() string {
	return e.Text
}

package usererr

var (
	ErrDeveloperError = new(DeveloperError)
	ErrNotFinished    = new(NotFinishedError)
)

type DeveloperError struct {
	Message string
}

func (e DeveloperError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return `Whoops! Some error occurred that the developer of this program should have handled. Sorry about that :(`
}

type NotFinishedError struct {
	DeveloperError
}

func (NotFinishedError) Error() string {
	return `NOT FINISHED!`
}

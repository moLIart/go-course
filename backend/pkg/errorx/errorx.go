package errorx

import "fmt"

func MustNoError(err error, message string) {
	if err != nil {
		panic(Wrap(err, message))
	}
}

func Wrap(err error, message string) error {
	if err == nil {
		return nil
	}

	return fmt.Errorf("%s: %w", message, err)
}

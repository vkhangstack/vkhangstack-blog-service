package utils

import "errors"

func ErrorAs(err error, target any) bool {
	return errors.As(err, target)

}

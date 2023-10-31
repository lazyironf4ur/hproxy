package common

import "runtime"

func Stack(skip int) (string, int, bool) {
	_, file, line, ok := runtime.Caller(skip)

	return file, line, ok
}

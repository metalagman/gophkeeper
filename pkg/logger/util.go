package logger

import (
	"encoding/json"
)

// DumpJSON string of value
func DumpJSON(v interface{}) string {
	b, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(b)
}

// CheckErr logs fatal error and exits if error exists
func CheckErr(err error) {
	Global().CheckErr(err)
}

package utils

import (
	"encoding/json"
	"io"
	"strings"
)

func JSONStringReaderFor(o interface{}) io.Reader {
	bytes, _ := json.Marshal(o)
	return strings.NewReader(string(bytes))
}

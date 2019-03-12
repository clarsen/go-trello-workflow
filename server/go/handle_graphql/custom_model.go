package handle_graphql

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"time"

	"github.com/99designs/gqlgen/graphql"
)

// if the type referenced in .gqlgen.yml is a function that returns a marshaller we can use it to encode and decode
// onto any existing go type.
func MarshalTimestamp(t time.Time) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		io.WriteString(w, strconv.FormatInt(t.Unix(), 10))
	})
}

// Unmarshal{Typename} is only required if the scalar appears as an input. The raw values have already been decoded
// from json into int/float64/bool/nil/map[string]interface/[]interface
func UnmarshalTimestamp(v interface{}) (time.Time, error) {
	i64, err := v.(json.Number).Int64()
	if err != nil {
		return time.Time{}, errors.New(fmt.Sprintf("time should be a unix timestamp, not %+v", reflect.TypeOf(v)))
	}
	return time.Unix(i64, 0), nil
}

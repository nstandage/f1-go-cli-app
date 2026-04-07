package customtype

import (
	"encoding/json"
	"fmt"
)

type FlexibleString string

func (fs *FlexibleString) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		*fs = ""
		return nil
	}

	var raw any
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}

	switch v := raw.(type) {
	case string:
		*fs = FlexibleString(v)
	case float64:
		*fs = FlexibleString(fmt.Sprintf("%.2f", v))
	default:
		*fs = ""
	}
	return nil
}

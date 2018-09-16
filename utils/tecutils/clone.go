package tecutils

import (
	"encoding/json"
	"fmt"
	"reflect"
)

func Clone(s, t interface{}) (interface{}, error) {
	if s == nil || t == nil {
		return nil, fmt.Errorf("cannot clone if source or target are nil")
	}

	if reflect.TypeOf(s).Elem().Kind() != reflect.TypeOf(t).Elem().Kind() {
		return nil, fmt.Errorf("source and target must be the same type")
	}

	data, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, t)
	if err != nil {
		return nil, err
	}
	return t, nil
}

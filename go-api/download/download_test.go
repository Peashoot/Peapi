package download_test

import (
	"encoding/json"
	"log"
	"testing"

	"github.com/mitchellh/mapstructure"
)

type Abc struct {
	Name []string `json:"name"`
}

func TestJsonUnmarshal(t *testing.T) {
	jsonStr := `{"code": 1, "message": "abc", "ioio": [{"name":["abc", "def"]}]}`
	respMap := make(map[string]interface{})
	if err := json.Unmarshal([]byte(jsonStr), &respMap); err != nil {
		t.Fail()
	}
	result, ok := respMap["ioio"].([]interface{})
	var abc Abc
	if err := mapstructure.Decode(result[0], &abc); err != nil {
		log.Println(err)
	}
	_, ok = result[0].(Abc)
	log.Println(ok)
	log.Println(result)
	log.Println(respMap)
}

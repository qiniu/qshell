package data

import (
	"encoding/json"
	"fmt"
	"testing"
)

type Data struct {
	Can  Bool   `json:"can"`
	Age  Int    `json:"age"`
	Name String `json:"name"`
}

func TestData(t *testing.T) {

	d := Data{
		Can:  NewBool(true),
		Age:  NewInt(10),
		Name: NewString("tom"),
	}

	dataString := ""
	data, err := json.Marshal(&d)
	if err != nil {
		t.Fail()
	} else {
		dataString = string(data)
	}

	fmt.Println("marshal json:" + dataString)

	d = Data{}
	err = json.Unmarshal(data, &d)
	if err != nil {
		t.Fail()
	}
}

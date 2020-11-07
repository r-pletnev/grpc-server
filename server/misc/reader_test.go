package misc

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"testing"
)


func makeCSVBuffer() *bytes.Buffer {
	buf := bytes.NewBuffer(nil)
	w  := csv.NewWriter(buf)
	w.Comma = ';'

	for i := 1; i <=3; i++ {
		val1 := fmt.Sprintf("product_%d", i)
		val2 := fmt.Sprint(i * 100)
		w.Write([]string{val1, val2})
	}
	w.Flush()
	return buf
}



func TestReadCSV(t *testing.T) {
	buf := makeCSVBuffer()
	products, err := ReadCSV(buf)
	if err != nil {
		t.Error(err)
	}

	if len(products) != 3 {
		t.Error("rows doesn't equals 3")
	}

}
package misc

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"github.com/globalsign/mgo/bson"
	"io"
	"strconv"
	"time"
)

type Product struct {
	ID                bson.ObjectId `bson:"_id,omitempty"`
	Name              string        `bson:"name,omitempty"`
	Price             int64         `bson:"price,omitempty"`
	CreatedAt         time.Time     `bson:"created_at,omitempty"`
	UpdatedAt         time.Time     `bson:"updated_at,omitempty"`
	PriceUpdatedCount int           `bson:"price_updated_count,omitempty"`
}


func ReadCSV(buffer *bytes.Buffer) (products []Product, err error) {
	r := csv.NewReader(buffer)
	r.Comma = ';'

	for {
		row, err := r.Read()
		if err != nil && err == io.EOF {
			break
		}
		if err != nil {
			return products, fmt.Errorf("read csv error: %v", err)
		}
		if len(row) != 2 {
			return products, fmt.Errorf("wrong csv format: %v", row)
		}
		price, err := strconv.ParseInt(row[1], 10, 0)
		if err != nil {
			return products, fmt.Errorf("bad price column: %s", row[1])
		}
		product := Product{
			Name:  row[0],
			Price: price,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		products = append(products, product)
	}

	return products, err
}
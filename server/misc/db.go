package misc

import (
	"fmt"
	"github.com/globalsign/mgo"
)

func NewDBSession(url string) *mgo.Session {
	session, err := mgo.Dial(url)
	if err != nil {
		panic(err)
	}
	return session
}

// Create product if product is new, otherwise only update price
func CreateOrUpdate(collection *mgo.Collection, products []Product) error {
	for _, elm := range products {
		oldProduct := Product{}
		if err := collection.Find(Product{Name: elm.Name}).One(&oldProduct); err != nil && err.Error() != "not found" {
			return err
		}
		if oldProduct.ID.Valid() && oldProduct.Price != elm.Price {
			oldProduct.PriceUpdatedCount += 1
			oldProduct.Price = elm.Price
			oldProduct.UpdatedAt = elm.UpdatedAt
			if err := collection.Update(Product{ID: oldProduct.ID}, oldProduct); err != nil {
				return err
			}
			continue
		}

		if !oldProduct.ID.Valid() {
			elm.PriceUpdatedCount = 1
			if err := collection.Insert(elm); err != nil {
				return err
			}
		}

	}
	return nil
}

type FieldType int

func (f FieldType) String() string {
	switch f {
	case 0:
		return "_id"
	case 1:
		return "name"
	case 2:
		return "price"
	case 3:
		return "price_updated_count"
	case 4:
		return "updated_at"
	default:
		return "_id"
	}
}

type OrderType int

func (o OrderType) String() string {
	if o > 0 {
		return "-"
	}
	return ""
}


type FilterParameters struct {
	PageNumber int32
	PerPage    int32
	Field FieldType
	Order OrderType
}

func (fp FilterParameters) getSortQuery() string {
	return fmt.Sprintf("%s%s",fp.Order,fp.Field)
}

func (fp FilterParameters) getPageNumber() int32 {
	if fp.PageNumber == 0 {
		return 1
	}
	return fp.PageNumber
}

func GetList(collection *mgo.Collection, params FilterParameters) (products []Product, err error) {
	q := collection.Find(nil).
		Sort(params.getSortQuery()).
		Limit(int(params.PerPage))

	q = q.Skip(int((params.getPageNumber() - 1) * params.PerPage))

	if err = q.All(&products); err != nil {
		return
	}
	return
}

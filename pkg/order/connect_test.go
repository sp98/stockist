package order

import (
	"log"
	"testing"
)

func GetData() *Order {
	ord := &Order{
		KC: getConnection(),
	}
	return ord

}
func TestGetSecondLegOrderID(t *testing.T) {
	parentID := "190701000324103"
	ord := GetData()
	slOrderID, err := ord.GetSecondLegOrderID(parentID)
	if err != nil {
		log.Println(err)
	}

	t.Error(slOrderID)

}

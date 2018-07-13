package relay

import (
	"github.com/dukfaar/goUtils/com/dukfaar/service"
)

func GetHasPreviousAndNextPage(resultLength int, firstItemID string, lastItemID string, service service.DBService) (chan bool, chan bool) {
	hasPreviousPage := make(chan bool)
	hasNextPage := make(chan bool)

	go func() {
		if resultLength > 0 {
			result, _ := service.HasElementBeforeID(firstItemID)
			hasPreviousPage <- result
		} else {
			hasPreviousPage <- false
		}
	}()

	go func() {
		if resultLength > 0 {
			result, _ := service.HasElementAfterID(lastItemID)
			hasNextPage <- result
		} else {
			hasNextPage <- false
		}
	}()

	return hasPreviousPage, hasNextPage
}

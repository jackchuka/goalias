package worker

import (
	"fmt"

	pricingsvc "example.com/shop/pricing"
)

func Invoice(items []int) string {
	return fmt.Sprintf("total: %d", pricingsvc.Total(items...))
}

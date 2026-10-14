package api

import (
	pr "example.com/shop/pricing"
)

func Checkout(items []int) int {
	return pr.Total(items...)
}

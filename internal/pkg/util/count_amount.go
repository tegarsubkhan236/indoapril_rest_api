package util

import "errors"

func CountAmount(tax, discount int8, prices, quantity []int) (int, error) {
	var totalAmount int

	if len(prices) != len(quantity) {
		return 0, errors.New("number of prices does not match number of quantities")
	}

	for i, price := range prices {
		subAmount := price * quantity[i]
		totalAmount += subAmount
	}

	countTax := int(tax) * totalAmount / 100
	countDiscount := int(discount) * totalAmount / 100

	totalAmount = totalAmount + countTax - countDiscount
	return totalAmount, nil
}

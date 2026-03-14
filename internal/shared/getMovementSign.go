package shared

import "errors"

func GetSignedMovement(accountCategory string, movementType string, value float64) (float64, error) {
	if accountCategory == "assets" || accountCategory == "expenses" {
		if movementType == "debit" {
			return value, nil
		} else {
			return -value, nil
		}
	}

	if accountCategory == "liabilities" || accountCategory == "equity" || accountCategory == "revenue" {
		if movementType == "credit" {
			return value, nil
		} else {
			return -value, nil
		}
	}

	return 0, errors.New("Invalid account category")
}

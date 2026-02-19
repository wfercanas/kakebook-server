package model

import "errors"

var ErrNoRecord = errors.New("model: no matching record found")

// Accounts
var ErrDeleteUsedAccount = errors.New("model: cannot delete account used in one or more movements")

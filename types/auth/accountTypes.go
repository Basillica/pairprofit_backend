package auth

import (
	"errors"
)

var AccountTypes = newAccountRegistry()

type AccountType struct {
	AccountTypeName       string
	AccountTypeCodeNumber string
}

func (c *AccountType) String() string {
	return c.AccountTypeName
}

func newAccountRegistry() *AccountTypeRegistry {
	regularAccount := &AccountType{"regularAccount", "001"}
	serviceProvider := &AccountType{"serviceProvider", "002"}
	blogger := &AccountType{"blogger", "003"}

	return &AccountTypeRegistry{
		ServiceProvider: serviceProvider,
		RegularAccount:  regularAccount,
		BloggerAccount:  blogger,
		accountTypes:    []*AccountType{serviceProvider, regularAccount, blogger},
	}
}

type AccountTypeRegistry struct {
	ServiceProvider *AccountType
	RegularAccount  *AccountType
	BloggerAccount  *AccountType
	accountTypes    []*AccountType
}

func (c *AccountTypeRegistry) List() []*AccountType {
	return c.accountTypes
}

func (c *AccountTypeRegistry) Parse(s string) (*AccountType, error) {
	for _, accountType := range c.List() {
		if accountType.String() == s {
			return accountType, nil
		}
	}
	return nil, errors.New("couldn't find it")
}

package model

import (
	"context"
	"errors"
	"path"

	"github.com/hashicorp/vault/sdk/logical"
)

type Account struct {
	Address string `json:"address"`
	URL     string `json:"url"`
}

func ReadAccount(ctx context.Context, req *logical.Request) (*Account, error) {

	accountPath := path.Dir(req.Path)

	entry, err := req.Storage.Get(ctx, accountPath)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	var account *Account
	err = entry.DecodeJSON(&account)
	if err != nil {
		return nil, errors.New("Fail to decode account to JSON format")
	}

	return account, nil
}

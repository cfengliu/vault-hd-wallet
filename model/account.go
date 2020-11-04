package model

import (
	"context"
	"errors"
	"path"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/hashicorp/vault/sdk/logical"
)

type Account accounts.Account

func ReadAccount(ctx context.Context, req *logical.Request) (*Account, error) {

	storagePath, _ := path.Split(req.Path)

	entry, err := req.Storage.Get(ctx, storagePath)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	var account *Account
	err = entry.DecodeJSON(&account)
	if err != nil {
		return nil, errors.New("Fail to decode wallet to JSON format")
	}

	return account, nil
}

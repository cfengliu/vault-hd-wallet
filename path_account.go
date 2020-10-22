package main

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
)

// AccountPaths aa
func AccountPaths(b *PluginBackend) []*framework.Path {
	return []*framework.Path{
		{
			Pattern:         "accounts/init",
			HelpSynopsis:    "Generate nmemonic",
			HelpDescription: `Generate nmemonic`,
			ExistenceCheck:  pathExistenceCheck,
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.CreateOperation: &framework.PathOperation{
					Callback: b.pathAccountsCreate,
					Summary:  "init the account",
				},
			},
		},
	}
}

func (b *PluginBackend) pathAccountsCreate(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {

	mnemonic, _ := hdwallet.NewMnemonic(128)

	wallet, err := hdwallet.NewFromMnemonic(mnemonic)
	if err != nil {
		log.Fatal(err)
	}

	path := hdwallet.MustParseDerivationPath("m/44'/60'/0'/0/0")
	account, err := wallet.Derive(path, false)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(account.Address.Hex())

	return &logical.Response{
		Data: map[string]interface{}{
			"address": account.Address.Hex(),
		},
	}, nil
}

func pathExistenceCheck(ctx context.Context, req *logical.Request, data *framework.FieldData) (bool, error) {
	out, err := req.Storage.Get(ctx, req.Path)
	if err != nil {
		return false, fmt.Errorf("existence check failed: %v", err)
	}

	return out != nil, nil
}

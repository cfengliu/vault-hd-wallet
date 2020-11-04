package path

import (
	"context"
	"errors"
	"vault-hd-wallet/model"
	"vault-hd-wallet/utils"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
)

// AccountPaths aa
func AccountPaths(b *PluginBackend) []*framework.Path {
	return []*framework.Path{
		{
			Pattern:         "accounts/" + framework.GenericNameRegex("name"),
			HelpSynopsis:    "create account with bip-44 path",
			HelpDescription: `create account with bip-44 path`,
			ExistenceCheck:  utils.PathExistenceCheck,
			Fields: map[string]*framework.FieldSchema{
				"name": {
					Type: framework.TypeString,
				},
				"derivationPath": {
					Type: framework.TypeString,
				},
			},
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.CreateOperation: &framework.PathOperation{
					Callback: b.createAccount,
					Summary:  "create a account",
				},
			},
		},
		{
			Pattern:         "accounts/" + framework.GenericNameRegex("name") + "/address",
			HelpSynopsis:    "get account address",
			HelpDescription: `get account address`,
			ExistenceCheck:  utils.PathExistenceCheck,
			Fields: map[string]*framework.FieldSchema{
				"name": {
					Type: framework.TypeString,
				},
			},
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback: b.readAddress,
					Summary:  "get address from an account",
				},
			},
		},
		{
			Pattern:         "accounts/" + framework.GenericNameRegex("name") + "/path",
			HelpSynopsis:    "get account derivation path",
			HelpDescription: `get account derivation path`,
			ExistenceCheck:  utils.PathExistenceCheck,
			Fields: map[string]*framework.FieldSchema{
				"name": {
					Type: framework.TypeString,
				},
			},
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback: b.readDerivationPath,
					Summary:  "read derivation path from an account",
				},
			},
		},
	}
}

func (b *PluginBackend) createAccount(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {

	derivationPathField, ok := data.Get("derivationPath").(string)
	if !ok {
		return nil, errors.New("derivationPath is not a string")
	}

	wallet, err := model.ReadWallet(ctx, req)

	derivationPath, err := hdwallet.ParseDerivationPath(derivationPathField)
	if err != nil {
		return nil, err
	}

	account, err := wallet.Derive(derivationPath)
	if err != nil {
		return nil, err
	}

	// save account
	entry, err := logical.StorageEntryJSON(req.Path, account)
	if err != nil {
		return nil, err
	}

	err = req.Storage.Put(ctx, entry)
	if err != nil {
		return nil, err
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"address": account.Address.Hex(),
		},
	}, nil
}

func (b *PluginBackend) readAddress(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {

	account, err := model.ReadAccount(ctx, req)
	if err != nil {
		return nil, err
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"address": account.Address.Hex(),
		},
	}, nil
}

func (b *PluginBackend) readDerivationPath(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {

	account, err := model.ReadAccount(ctx, req)
	if err != nil {
		return nil, err
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"derivation_path": account.URL.Path,
		},
	}, nil
}

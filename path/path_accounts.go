package path

import (
	"context"
	"encoding/hex"
	"fmt"
	"vault-hd-wallet/model"
	"vault-hd-wallet/utils"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
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
		{
			Pattern:         "accounts/" + framework.GenericNameRegex("name") + "/sign",
			HelpSynopsis:    "sign data",
			HelpDescription: `sign data`,
			ExistenceCheck:  utils.PathExistenceCheck,
			Fields: map[string]*framework.FieldSchema{
				"name": {
					Type: framework.TypeString,
				},
				"data": {
					Type:        framework.TypeString,
					Description: "The data to sign.",
				},
			},
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.CreateOperation: &framework.PathOperation{
					Callback: b.signData,
					Summary:  "read derivation path from an account",
				},
			},
		},
		{
			Pattern:         "accounts/" + framework.GenericNameRegex("name") + "/sign-tx",
			HelpSynopsis:    "sign a transaction",
			HelpDescription: `sign a transaction`,
			ExistenceCheck:  utils.PathExistenceCheck,
			Fields: map[string]*framework.FieldSchema{
				"name": {
					Type: framework.TypeString,
				},
				"address_to": {
					Type:        framework.TypeString,
					Description: "The address of the account to send tx to.",
				},
				"data": {
					Type:        framework.TypeString,
					Description: "The data to sign.",
				},
				"amount": {
					Type:        framework.TypeString,
					Description: "Amount of ETH (in wei).",
				},
				"nonce": {
					Type:        framework.TypeString,
					Description: "The transaction nonce.",
				},
				"gas_limit": {
					Type:        framework.TypeString,
					Description: "The gas limit for the transaction - defaults to 21000.",
					Default:     "21000",
				},
				"gas_price": {
					Type:        framework.TypeString,
					Description: "The gas price for the transaction in wei.",
					Default:     "0",
				},
				"chainID": {
					Type:        framework.TypeString,
					Description: "The chain ID of the blockchain network.",
				},
			},
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.CreateOperation: &framework.PathOperation{
					Callback: b.signTransaction,
					Summary:  "sign a transaction",
				},
			},
		},
	}
}

func (b *PluginBackend) createAccount(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	dataWrapper := utils.NewFieldDataWrapper(data)

	derivationPathField, err := dataWrapper.MustGetString("derivationPath")
	if err != nil {
		return nil, err
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
			"address": account.Address,
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
			"address": account.Address,
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
			"derivation_path": account.URL,
		},
	}, nil
}

func (b *PluginBackend) signTransaction(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	dataWrapper := utils.NewFieldDataWrapper(data)

	inputData := dataWrapper.GetString("data", "")

	addressToStr := dataWrapper.GetString("address_to", "")

	amount, err := dataWrapper.MustGetBigInt("amount")
	if err != nil {
		return nil, err
	}

	nonce := dataWrapper.GetUint64("nonce", 0)

	gasLimit, err := dataWrapper.MustGetUint64("gas_limit")
	if err != nil {
		return nil, err
	}

	gasPrice, err := dataWrapper.MustGetBigInt("gas_price")
	if err != nil {
		return nil, err
	}

	chainID, err := dataWrapper.MustGetBigInt("chainID")
	if err != nil {
		return nil, err
	}

	account, err := model.ReadAccount(ctx, req)

	privateKey, err := crypto.HexToECDSA(account.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("error reconstructing private key")
	}
	defer utils.ZeroKey(privateKey)

	var tx *types.Transaction
	var txDataToSign []byte

	if inputData != "" {
		txDataToSign, err = hexutil.Decode(inputData)
		if err != nil {
			return nil, err
		}
	}

	if addressToStr == "" {
		tx = types.NewContractCreation(nonce, amount, gasLimit, gasPrice, txDataToSign)
	} else {
		addressTo := common.HexToAddress(addressToStr)
		tx = types.NewTransaction(nonce, addressTo, amount, gasLimit, gasPrice, txDataToSign)
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		return nil, err
	}

	ts := types.Transactions{signedTx}
	rawTxBytes := ts.GetRlp(0)
	rawTxHex := hex.EncodeToString(rawTxBytes)

	return &logical.Response{
		Data: map[string]interface{}{
			"transaction_hash":   signedTx.Hash().Hex(),
			"address_from":       account.Address,
			"address_to":         addressToStr,
			"signed_transaction": rawTxHex,
		},
	}, nil
}

func (b *PluginBackend) signData(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	dataWrapper := utils.NewFieldDataWrapper(data)

	inputData, err := dataWrapper.MustGetString("data")
	if err != nil {
		return nil, err
	}

	account, err := model.ReadAccount(ctx, req)
	if err != nil {
		return nil, err
	}

	privateKey, err := crypto.HexToECDSA(account.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("error reconstructing private key")
	}
	defer utils.ZeroKey(privateKey)

	dataToSign, err := hexutil.Decode(inputData)
	if err != nil {
		return nil, err
	}

	dataHash := crypto.Keccak256Hash(dataToSign)

	signature, err := crypto.Sign(dataHash.Bytes(), privateKey)
	if err != nil {
		return nil, err
	}

	hexSig := hexutil.Encode(signature)

	return &logical.Response{
		Data: map[string]interface{}{
			"signature": hexSig,
		},
	}, nil
}

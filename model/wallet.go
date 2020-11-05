package model

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/tyler-smith/go-bip39"
)

type Wallet struct {
	MasterKey string `json:"masterKey"`
	Seed      []byte `json:"seed"`
}

// NewWalletFromMnemonic Generate wallet from mnemonic
func NewWalletFromMnemonic(mnemonic string, passphrase string) (*Wallet, error) {

	if mnemonic == "" {
		return nil, errors.New("mnemonic is required")
	}

	if !bip39.IsMnemonicValid(mnemonic) {
		return nil, errors.New("mnemonic is invalid")
	}

	seed, err := NewSeedFromMnemonic(mnemonic, passphrase)
	if err != nil {
		return nil, err
	}

	wallet, err := newWallet(seed)
	if err != nil {
		return nil, err
	}

	return wallet, nil
}

func newWallet(seed []byte) (*Wallet, error) {
	masterKey, err := hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
	if err != nil {
		return nil, err
	}

	masterKeyStr := masterKey.String()

	return &Wallet{
		MasterKey: masterKeyStr,
		Seed:      seed,
	}, nil
}

// NewSeedFromMnemonic returns a BIP-39 seed based on a BIP-39 mnemonic.
func NewSeedFromMnemonic(mnemonic string, passphrase string) ([]byte, error) {
	if mnemonic == "" {
		return nil, errors.New("mnemonic is required")
	}

	return bip39.NewSeedWithErrorChecking(mnemonic, passphrase)
}

func ReadWallet(ctx context.Context, req *logical.Request) (*Wallet, error) {

	walletPath := "wallet"

	entry, err := req.Storage.Get(ctx, walletPath)
	if err != nil {
		return nil, err
	}

	if entry == nil {
		return nil, fmt.Errorf("entry not existed at %v", walletPath)
	}

	var wallet *Wallet
	err = entry.DecodeJSON(&wallet)
	if err != nil {
		return nil, errors.New("Fail to decode wallet to JSON format")

	}

	return wallet, nil
}

// Derive acctount from derivation path
func (w *Wallet) Derive(path accounts.DerivationPath) (*Account, error) {

	address, err := w.deriveAddress(path)
	addressStr := address.String()

	// If an error occurred or no pinning was requested, return
	if err != nil {
		return &Account{}, err
	}

	URL := accounts.URL{
		Scheme: "",
		Path:   path.String(),
	}
	URLStr := URL.String()

	account := &Account{
		Address: addressStr,
		URL:     URLStr,
	}

	return account, nil
}

// DerivePrivateKey derives the private key of the derivation path.
func (w *Wallet) derivePrivateKey(path accounts.DerivationPath) (*ecdsa.PrivateKey, error) {
	var err error
	key, err := hdkeychain.NewKeyFromString(w.MasterKey)
	if err != nil {
		return nil, err
	}

	for _, n := range path {
		key, err = key.Child(n)
		if err != nil {
			return nil, err
		}
	}

	privateKey, err := key.ECPrivKey()
	privateKeyECDSA := privateKey.ToECDSA()
	if err != nil {
		return nil, err
	}

	return privateKeyECDSA, nil
}

// DerivePublicKey derives the public key of the derivation path.
func (w *Wallet) derivePublicKey(path accounts.DerivationPath) (*ecdsa.PublicKey, error) {
	privateKeyECDSA, err := w.derivePrivateKey(path)
	if err != nil {
		return nil, err
	}

	publicKey := privateKeyECDSA.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("failed to get public key")
	}

	return publicKeyECDSA, nil
}

// DeriveAddress derives the account address of the derivation path.
func (w *Wallet) deriveAddress(path accounts.DerivationPath) (common.Address, error) {
	publicKeyECDSA, err := w.derivePublicKey(path)
	if err != nil {
		return common.Address{}, err
	}

	address := crypto.PubkeyToAddress(*publicKeyECDSA)
	return address, nil
}

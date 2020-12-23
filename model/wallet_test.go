package model

import (
	"encoding/hex"
	"reflect"
	"testing"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/ethereum/go-ethereum/accounts"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
)

var mnemonic = "much museum inside brief load judge puppy loop quick okay gadget prevent cheese seminar beauty curve obvious sweet milk bread proud tonight wall wet"
var seedStr = "667cf4f811ff5a39be2b806e1d1e3f819ffadcfb5d16d2c9a0f1f6199fe85cc86a4b4ba7dac73594471da4ff75d5a9da732958c8d41815ee48859777d489567a"
var seedStrPassphrase = "0076324898a24676cbcc1822ad7e02d8759da7cb3bb797be2402b5d38e298791b264b30d22ac14bae56cc70cbae5d98a2f4c8d782c6380e71c9d6258d294ab09"

var seedDecoded, _ = hex.DecodeString(seedStr)
var seedDecodedPassphrase, _ = hex.DecodeString(seedStrPassphrase)

var masterKey, _ = hdkeychain.NewMaster(seedDecoded, &chaincfg.MainNetParams)
var masterKeyStr = masterKey.String()
var masterKeyPassphrase, _ = hdkeychain.NewMaster(seedDecodedPassphrase, &chaincfg.MainNetParams)
var masterKeyPassphraseStr = masterKeyPassphrase.String()

var wallet = &Wallet{
	MasterKey: masterKeyStr,
	Seed:      seedStr,
}
var walletPassphrase = &Wallet{
	MasterKey: masterKeyPassphraseStr,
	Seed:      seedStrPassphrase,
}

var account = &Account{
	Address:    "0xe5692Ff6e92c3DD89A0c35a14B17eD270Aa39881",
	URL:        "m/44'/60'/0'/0/0",
	PrivateKey: "cc30326050d6f2001b40da3a37898ad2f5b7b3062a566355b56ea0842ccbfeb8",
	PublicKey:  "ad11598ad03f09dee948401df67ae9e694e8f2b20fdf3b958007c446b95026975558078bca9b8e8298d266f6be2acf50efb70c21d756c0c11d68a5f20bf5a012",
}

var derivationPath, _ = hdwallet.ParseDerivationPath(account.URL)

func TestNewSeedFromMnemonic(t *testing.T) {
	type args struct {
		mnemonic   string
		passphrase string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "mnemonic",
			args: args{
				mnemonic:   mnemonic,
				passphrase: "",
			},
			want: seedDecoded,
		},
		{
			name: "mnemonic with passphrase",
			args: args{
				mnemonic:   seedStrPassphrase,
				passphrase: "123456",
			},
			want: seedDecodedPassphrase,
		},
		{
			name: "no mnemonic",
			args: args{
				mnemonic:   "",
				passphrase: "",
			},
			wantErr: true,
		},
		{
			name: "invalid mnemonic",
			args: args{
				mnemonic:   "apple banana",
				passphrase: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := NewSeedFromMnemonic(tt.args.mnemonic, tt.args.passphrase)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewSeedFromMnemonic() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSeedFromMnemonic() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newWallet(t *testing.T) {
	type args struct {
		seed []byte
	}
	tests := []struct {
		name    string
		args    args
		want    *Wallet
		wantErr bool
	}{
		{
			name: "seed",
			args: args{
				seed: seedDecoded,
			},
			want: wallet,
		},
		{
			name: "invalid seed",
			args: args{
				seed: seedDecoded[:6],
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := newWallet(tt.args.seed)
			if (err != nil) != tt.wantErr {
				t.Errorf("newWallet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newWallet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewWalletFromMnemonic(t *testing.T) {
	type args struct {
		mnemonic   string
		passphrase string
	}
	tests := []struct {
		name    string
		args    args
		want    *Wallet
		wantErr bool
	}{
		{
			name: "mnemonic",
			args: args{
				mnemonic:   seedStr,
				passphrase: "",
			},
			want: wallet,
		},
		{
			name: "mnemonic with passphrase",
			args: args{
				mnemonic:   seedStrPassphrase,
				passphrase: "123456",
			},
			want: walletPassphrase,
		},
		{
			name: "no mnemonic",
			args: args{
				mnemonic:   "",
				passphrase: "",
			},
			wantErr: true,
		},
		{
			name: "invalid mnemonic",
			args: args{
				mnemonic:   "apple banana",
				passphrase: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := NewWalletFromMnemonic(tt.args.mnemonic, tt.args.passphrase)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewWalletFromMnemonic() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewWalletFromMnemonic() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWallet_Derive(t *testing.T) {
	type args struct {
		path accounts.DerivationPath
	}
	tests := []struct {
		name    string
		w       *Wallet
		args    args
		want    *Account
		wantErr bool
	}{
		{
			name: "normal account",
			w:    wallet,
			args: args{
				path: derivationPath,
			},
			want: account,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.w.Derive(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("Wallet.Derive() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Wallet.Derive() = %v, want %v", got, tt.want)
			}
		})
	}
}

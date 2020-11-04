package utils

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/tyler-smith/go-bip39"
)

func ValidateMnemonic(mnemonic string) bool {
	mnemonicSlice := strings.Fields(mnemonic)

	if !(bip39.IsMnemonicValid(mnemonic) && len(mnemonicSlice) == 24) {
		return false
	}

	return true
}

func PathExistenceCheck(ctx context.Context, req *logical.Request, data *framework.FieldData) (bool, error) {
	out, err := req.Storage.Get(ctx, req.Path)
	if err != nil {
		return false, fmt.Errorf("existence check failed: %v", err)
	}

	return out != nil, nil
}

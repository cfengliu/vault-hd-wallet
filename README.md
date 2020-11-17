# vault-hd-wallet
 
vault-hd-wallet is a vault plugin that implements the Hierarchical Deterministic wallet (HD wallet), which support the ethereum mnemonic storage and signing ethereum transaction. The wallet seed and derived key will not be exposed while signing the transaction or creating the account. The mnemonic can be imported to the plugin to restore the account addresses.

## Getting Started

The vault server should be installed firstly. If not installed yet, read [installation guide](https://learn.hashicorp.com/tutorials/vault/getting-started-install). 

To learn what is custom plugin and how it works, please read [this guide](https://learn.hashicorp.com/tutorials/vault/plugin-backends).

To compile this plugin, run the following command. The compiled binary should be placed in vault server's plugin directory.

``` bash
go build -a -v -i -o hdwallet *.go
```

Get the SHA256 hash of binary file before registering the plugin:

``` bash
sha256=$(sha256sum ./hdwallet | cut -d " " -f1) >/dev/null
```

Register plugin:

``` bash
vault write sys/plugins/catalog/secret/hdwallet \
        sha_256=$sha256 \
        command="hdwallet"
```

Enable plugin:

``` bash
vault secrets enable -plugin-name='hdwallet' plugin
```

## Policy

The plugin policy is depended on your user management. This repo provides two examples: wallet and accounts. Wallet policy is for admin, which enables user to initialize wallet and all accounts. Accounts policy allows user to get account address and sign a transaction.

## Usage

1. Create a HD wallet
    If no mnemonic is provided, the HD wallet will randomly generate one.

    ``` bash
    curl --location --request POST "http://${ip}:${port}/v1/hdwallet/wallet" \
        --header "Authorization: Bearer ${token}" \
        --data-raw '{
            "mnemonic": "move mask pilot rather lion prevent reform mixture valve appear drop soap section pass jelly capital limb produce enough smooth nature cricket elevator jeans",
            "passphrase": ""
        }'
    ```

2. Read wallet
    Get wallet seed and master key. This function should be for testing only.

    ```bash
    curl --location --request GET "http://${ip}:${port}/v1/hdwallet/wallet" \
        --header "Authorization: Bearer ${token}"
    ```

3. Create an account
    The account address is derived from derivation path.

    ```bash
    curl --location --request POST "http://${ip}:${port}/v1/hdwallet/accounts/${name}" \
        --header "Authorization: Bearer ${token}" \
        --header "Content-Type: application/json; charset=utf-8" \
        --data-raw '{
            "derivationPath": "m/44'\''/60'\''/0'\''/0/0"
        }'
    ```

4. Get account address

    ```bash
    curl --location --request GET "http://${ip}:${port}/v1/hdwallet/accounts/${name}/address" \
        --header "Authorization: Bearer ${token}"
    ```

5. Sign a transaction
    To learn signing transaction and its parameters, read [this document](https://web3js.readthedocs.io/en/v1.2.0/web3-eth.html#signtransaction)

    ```bash
    curl --location --request POST "http://${ip}:${port}/v1/hdwallet/accounts/${name}/sign-tx" \
            --header "Authorization: Bearer ${token}" \
            --data-raw "{
                \"address_to\": \"\",
                \"amount\": \"100000\",
                \"nonce\": \"2\",
                \"gas_limit\": \"3000000\",
                \"gas_price\": \"1000000000\",
                \"chainID\": \"4\",
                \"data\": \"\"
            }"
    ```
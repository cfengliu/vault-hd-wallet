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

The plugin policy is depended on your [auth management](https://learn.hashicorp.com/tutorials/vault/identity?in=vault/auth-methods). This repo provides two examples: wallet and accounts. Wallet policy is for admin, which enables user to initialize wallet and all accounts. Accounts policy allows user to get account address and sign a transaction.

## Usage

### Create a HD wallet

If no mnemonic is provided, the HD wallet will randomly generate one.

``` bash
POST /hdwallet/wallet
```

Parameters
| Name       | Type   | In   | Description                                           |
| ---------- | ------ | ---- | ----------------------------------------------------- |
| mnemonic   | string | body | The mnemonic could be imported to restore the wallet. |
| passphrase | string | body | The mnemonic password to protect the wallet.          |

Code samples

``` bash
curl --request POST "http://${ip}:${port}/v1/hdwallet/wallet" \
    --header "Authorization: Bearer ${token}" \
    --data-raw '{
        "mnemonic": "move mask pilot rather lion prevent reform mixture valve appear drop soap section pass jelly capital limb produce enough smooth nature cricket elevator jeans",
        "passphrase": ""
    }'
```

### Read wallet

Get wallet seed and master key. This function should be for testing ONLY.

Code samples

```bash
curl --request GET "http://${ip}:${port}/v1/hdwallet/wallet" \
    --header "Authorization: Bearer ${token}"
```

### Create an account

The account address is derived from derivation path.

Parameters
| Name           | Type   | In   | Description                                                                   |
| -------------- | ------ | ---- | ----------------------------------------------------------------------------- |
| name           | string | url  | **Rquired.** The path of secrets engines where plugin store the account info. |
| derivationPath | string | body | **Rquired.** The BIP-44 path for generating the account address.              |

Code samples

```bash
curl --request POST "http://${ip}:${port}/v1/hdwallet/accounts/${name}" \
    --header "Authorization: Bearer ${token}" \
    --data-raw '{
        "derivationPath": "m/44'\''/60'\''/0'\''/0/0"
    }'
```

### Get account address

Parameters
| Name | Type   | In  | Description                                                                   |
| ---- | ------ | --- | ----------------------------------------------------------------------------- |
| name | string | url | **Rquired.** The path of secrets engines where plugin store the account info. |

Code samples

```bash
curl --request GET "http://${ip}:${port}/v1/hdwallet/accounts/${name}/address" \
    --header "Authorization: Bearer ${token}"
```

### Sign a transaction

To learn signing transaction and its parameters, read [this document](https://web3js.readthedocs.io/en/v1.2.0/web3-eth.html#signtransaction)

Parameters
| Name       | Type   | In   | Description                                                                                 |
| ---------- | ------ | ---- | ------------------------------------------------------------------------------------------- |
| name       | string | url  | **Rquired.** The path of secrets engines where plugin store the account info.               |
| address_to | string | body | The destination address for transaction. Leave empty if it is contract creation transaction |
| amount     | string | body | **Rquired.** The ether send to the destination address (in wei)                             |
| nonce      | string | body | **Rquired.** The transaction count of this account                                          |
| gas_limit  | string | body | **Rquired.** The estimated gas that transaction may consume                                 |
| gas_price  | string | body | **Rquired.** The price of gas (in wei)                                                      |
| chainID    | string | body | **Rquired.** The ID of etheruem network                                                     |
| data       | string | body | The bytecode of contract creation or function call. '0x' prefix is required.                |

Code samples

```bash
curl --request POST "http://${ip}:${port}/v1/hdwallet/accounts/${name}/sign-tx" \
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

### Sign data

Generate the signature for the input data

Parameters
| Name | Type   | In   | Description                                                                                |
| ---- | ------ | ---- | ------------------------------------------------------------------------------------------ |
| name | string | url  | **Rquired.** The path of secrets engines where plugin store the account info.              |
| data | string | body | **Rquired.** The data to be signed. (without `\x19Ethereum Signed Message:\n` prefix ) |


Code samples

```bash
curl --request POST "http://${ip}:${port}/v1/hdwallet/accounts/${name}/sign" \
    --header "Authorization: Bearer ${token}" \
    --data-raw "{
        \"data\": \"hello world\"
    }"
```

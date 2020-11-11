path "hdwallet/wallet" {
  capabilities = ["create", "read"]
}

path "hdwallet/accounts/*"{
    capabilities = ["create", "read"]
}
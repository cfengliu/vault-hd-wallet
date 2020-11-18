path "hdwallet/accounts/{{identity.entity.name}}/address" {
    capabilities = ["read"]
}

path "hdwallet/accounts/{{identity.entity.name}}/sign-tx"{
    capabilities = ["create"]
}

path "hdwallet/accounts/{{identity.entity.name}}/sign"{
    capabilities = ["create"]
}
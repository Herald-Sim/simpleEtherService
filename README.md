# simpleEtherService
API: https://neocool.iptime.org

## API guide

1. Sign-in
```
POST https://neocool.iptime.org/user

POST JSON Form
{"ID":"Test", "Passwd":"test"}
```

2. Login
```
POST https://neocool.iptime.org/user/[userID]

POST JSON Form
{test}
```

3. Add user wallet
```
GET https://neocool.iptime.org/addWallet/[userID]
```

4. Read Contract
```
GET https://neocool.iptime.org/contract
```

5. Write Contract
```
POST https://neocool.iptime.org/contract

POST JSON Form
{"WalletAddres":"address", "Value":"value"}
```

6. Transfer Token
```
POST https://neocool.iptime.org/token

POST JSON Form
{"ToWallet":"address", "FromWallet":"address", "Quantity":"value"}
```

All apis need auth checking using **JWT**
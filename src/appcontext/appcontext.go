package appcontext

// 각 context별, sql 쿼리 구문 메소드로 구현하여 코드 최적화
// interface 활용하여, debug-code 최적화

// UserObject : user info
type UserObject struct {
	ID     string
	Passwd string
}

// ReturnJSON : Return success or fail to client
type ReturnJSON struct {
	Status string
}

// LoginInfo : Login data which received from client
type LoginInfo struct {
	ID     string
	Passwd string
}

type RecentTransaction struct {
	TimeStamp string
	From      string
	To        string
	Value     string
	TokenName string
}

type History struct {
	List []RecentTransaction
}

type WalletMetaData struct {
	HRTbalance    string
	ETHbalance    string
	WalletHistroy History
}

type Wallet struct {
	WalletAddress string
}

type Wallets struct {
	List []Wallet
}

type Transfer struct {
	WalletAddress string
	Token         string
	Quantity      string
}

type ContractReq struct {
	WalletAddress string
	Value         string
}

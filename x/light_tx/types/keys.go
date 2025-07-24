package types

const (
	ModuleName   = "lighttx"
	StoreKey     = ModuleName
	RouterKey    = ModuleName // 이걸 추가해야 sdk.NewRoute()에서 사용할 수 있습니다
	QuerierRoute = ModuleName
)

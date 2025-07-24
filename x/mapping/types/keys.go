package types

const (
	// 모듈 이름 (AppModuleBasic.Name() 등에서 사용됨)
	ModuleName = "mapping"

	// 스토리지 프리픽스 (KVStore key-value 저장 시 사용)
	StoreKey = ModuleName

	// Router key (메시지 라우팅)
	RouterKey = ModuleName

	// Querier route (query path에 사용될 키)
	QuerierRoute = ModuleName
)

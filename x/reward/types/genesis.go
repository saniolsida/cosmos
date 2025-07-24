package types

// NewGenesisState creates a new GenesisState object
func NewGenesisState() *GenesisState {
	return &GenesisState{}
}

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{}
}

// Validate performs basic validation of the genesis state
func (gs *GenesisState) Validate() error {
	// 현재는 필드가 없으므로 항상 유효함
	return nil
}

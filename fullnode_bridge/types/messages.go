package types

type SolarData struct {
	DeviceID    string   `json:"device_id"`
	Timestamp   string   `json:"timestamp"`
	TotalEnergy float64  `json:"total_energy"`
	Location    Location `json:"location"`
}

type Location struct {
	Latitude   float64 `json:"latitude"`
	Longitutde float64 `json:"longitutde"`
}

type LightTxMessage struct {
	Original  *SolarData `json:"original,omitempty"`
	REC       *RECMeta   `json:"rec,omitempty"`
	Hash      string     `json:"hash"`
	Signature string     `json:"signature"`
	Pubkey    string     `json:"pubkey"`
}

type RECMeta struct {
	FacilityID       string `json:"facility_id"`
	FacilityName     string `json:"facility_name"`
	Location         string `json:"location"`
	TechnologyType   string `json:"technology_type"`   // 발전원
	CapacityMW       string `json:"capacity_mw"`       // 설비용량
	RegistrationDate string `json:"registration_date"` // i-REC 등록 승인일

	CertifiedId         string `json:"certified_id"`
	IssueData           string `json:"issue_data"`
	GenerationStartDate string `json:"generation_start_date"`
	GenerationEndDate   string `json:"generation_end_date"`
	MeasuredVolumeMWh   string `json:"measured_volume_MWh"`
	RetiredDate         string `json:"retired_date"`
	RetirementPurpose   string `json:"retirement_purpose"`
	Status              string `json:"status"`
	Timestamp           string `json:"timestamp"`
}

type DeviceToAddressMessage struct {
	DeviceID string `json:"device_id"`
	Address  string `json:"address"`
}

type AuthMessage struct {
	Address string `json:"user_address"`
}

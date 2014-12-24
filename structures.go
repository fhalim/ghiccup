package ghiccup

type HiccupInfo struct {
	Timestamp  string `json:"timestamp"`
	Resolution int64  `json:"resolution"`
	Threshold  int64  `json:"threshold"`
	Duration   int64  `json:"duration"`
}

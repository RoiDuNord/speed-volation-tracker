package models

type Passage struct {
	Track      []TPoint `json:"track"`
	LicenseNum string   `json:"licenseNumber"`
}

type TPoint struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	T int     `json:"t"`
}

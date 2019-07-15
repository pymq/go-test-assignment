package model

type Tv struct {
	ID           int64  `form:"id" query:"id"`
	Brand        NullString
	Manufacturer string `form:"manufacturer" query:"manufacturer"`
	Model        string `form:"model" query:"model"`
	Year         int    `form:"year" query:"year"`
}

func (t Tv) IsValid() bool {
	if len(t.Manufacturer) < 3 || len(t.Model) < 2 || t.Year < 2010 {
		return false
	}
	return true
}

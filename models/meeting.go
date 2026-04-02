package models

type Meeting struct {
	CircuitKey          uint   `json:"circuit_key"`
	CurcuitInfoURL      string `json:"circuit_info_url"`
	CircuitImage        string `json:"circuit_image"`
	CircuitShortName    string `json:"circuit_short_name"`
	CurcuitType         string `json:"circuit_type"`
	CountryCode         string `json:"country_code"`
	CountryFlag         string `json:"country_flag"`
	CountryKey          uint   `json:"country_key"`
	CountryName         string `json:"country_name"`
	DateEnd             string `json:"date_end"`
	DateStart           string `json:"date_start"`
	GMTOffset           string `json:"gmt_offset"`
	Location            string `json:"location"`
	MeetingKey          uint   `json:"meeting_key"`
	MeetingName         string `json:"meeting_name"`
	MeetingOfficialName string `json:"meeting_official_name"`
	Year                uint   `json:"year"`
}

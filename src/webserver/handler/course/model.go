package course

import (
	"database/sql"
)

type readParams struct {
	Page  string
	Total string
}

type readArgs struct {
	Page  uint16
	Total uint16
}

type readResponse struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Class     string `json:"class"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	Day       string `json:"day"`
	Status    string `json:"status"`
}

type createParams struct {
	Name        string
	Description string
	UCU         string
	Semester    string
	StartTime   string
	EndTime     string
	Class       string
	Day         string
	PlaceID     string
}

type createArgs struct {
	Name        string
	Description sql.NullString
	UCU         int8
	Semester    int8
	StartTime   int16
	EndTime     int16
	Class       string
	Day         int8
	PlaceID     string
}

type summaryResponse struct {
	Status string           `json:"status"`
	Course []courseResponse `json:"courses"`
}

type courseResponse struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	UCU      int8   `json:"ucu"`
	Semester int8   `json:"semester"`
}

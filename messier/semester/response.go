package semester

import (
	"neptune/backend/messier"
)

type GetSemestersResponse struct {
	SemesterID  string              `json:"SemesterID"`
	Description string              `json:"Description"`
	Start       messier.MessierTime `json:"Start"`
	End         messier.MessierTime `json:"End"`
}

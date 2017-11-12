package information

import (
	"net/http"
	"time"

	"github.com/melodiez14/meiko/src/util/helper"
	"github.com/melodiez14/meiko/src/webserver/template"

	"github.com/julienschmidt/httprouter"
	"github.com/melodiez14/meiko/src/module/course"
	inf "github.com/melodiez14/meiko/src/module/information"
	"github.com/melodiez14/meiko/src/util/alias"
	"github.com/melodiez14/meiko/src/util/auth"
)

// GetSummaryHandler func ...
func GetSummaryHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	sess := r.Context().Value("User").(*auth.User)

	// get enrolled course
	schedulesID, err := course.SelectIDByUserID(sess.ID)
	if err != nil {
		template.RenderJSONResponse(w, new(template.Response).
			SetCode(http.StatusInternalServerError).
			AddError("Internal server error"))
		return
	}

	// get information list
	informations, err := inf.SelectByScheduleID(schedulesID)
	if err != nil {
		template.RenderJSONResponse(w, new(template.Response).
			SetCode(http.StatusInternalServerError).
			AddError("Internal server error"))
		return
	}

	// convert informations to response
	var informationResponses []informationResponse
	t2 := time.Now()
	for _, val := range informations {
		informationResponses = append(informationResponses, informationResponse{
			Title:       val.Title,
			Date:        helper.DateToString(val.CreatedAt, t2),
			Description: val.Description.String,
		})
	}

	// if informations has only 5, so last and recent will be the same
	// else it has 5 last information and other is recent
	var res getSummaryResponse
	if len(informationResponses) <= alias.InformationMinimumLast {
		res = getSummaryResponse{
			Last:   informationResponses,
			Recent: informationResponses,
		}
	} else {
		res = getSummaryResponse{
			Last:   informationResponses[:alias.InformationMinimumLast],
			Recent: informationResponses[alias.InformationMinimumLast:],
		}
	}

	template.RenderJSONResponse(w, new(template.Response).
		SetCode(http.StatusOK).
		SetData(res))
	return
}

// GetDetailHandler func ..
func GetDetailHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	sess := r.Context().Value("User").(*auth.User)
	params := detailInfromationParams{
		ID: ps.ByName("id"),
	}
	args, err := params.validate()
	if err != nil {
		template.RenderJSONResponse(w, new(template.Response).
			SetCode(http.StatusBadRequest).AddError(err.Error()))
		return
	}
	scheduleID := inf.GetScheduleIDByID(args.ID)
	if scheduleID != 0 {
		if !course.IsEnrolled(sess.ID, scheduleID) {
			template.RenderJSONResponse(w, new(template.Response).
				SetCode(http.StatusBadRequest).
				AddError("you do not have permission to this informations"))
			return
		}
	}
	res, err := inf.GetByID(args.ID)
	if err != nil {
		template.RenderJSONResponse(w, new(template.Response).
			SetCode(http.StatusBadRequest).
			AddError("Information does not exist"))
		return
	}
	template.RenderJSONResponse(w, new(template.Response).
		SetCode(http.StatusOK).
		SetData(res))
	return
}

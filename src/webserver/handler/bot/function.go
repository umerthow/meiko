package bot

import (
	"database/sql"
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"time"

	"github.com/melodiez14/meiko/src/module/bot"
	cs "github.com/melodiez14/meiko/src/module/course"
	fl "github.com/melodiez14/meiko/src/module/file"
	"github.com/melodiez14/meiko/src/util/helper"
)

func generateMessage(confidence float64) string {
	messages := map[string][]string{
		"confident": []string{
			"Ini bro jawabannya",
		},
		"doubt": []string{
			"Kayaknya sih",
			"Mereun ya ini ge",
		},
		"notsure": []string{
			"sorry nih bro, gua ga tau",
			"kasih tau ga ya?",
			"Kamu ngomong apa sih?",
			"Gatau, coba tanya ke temen kamu!",
		},
	}

	levelConfidence := "notsure"
	if confidence >= 60 {
		levelConfidence = "confident"
	} else if confidence >= 30 {
		levelConfidence = "doubt"
	}

	msg := messages[levelConfidence]
	rand.Seed(time.Now().UTC().UnixNano())
	index := rand.Intn(len(msg))

	response := msg[index]
	return response
}

func handleAssistant(text string, userID int64) ([]map[string]interface{}, string, error) {
	args := []map[string]interface{}{}
	var filterCoursesRgx sql.NullString

	params := sEntity{
		text:   text,
		userID: userID,
	}

	// get days entity
	filterTime := params.getTime()
	filterDays := params.getDay()
	filterDays = append(filterDays, helper.TimeToDayInt(filterTime...)...)

	// get course entity
	filterCourses := params.getCourse()
	if len(filterCourses) > 0 {
		filterCoursesRgx = sql.NullString{
			Valid:  true,
			String: strings.Join(filterCourses, "|"),
		}
	}

	assistants, err := bot.SelectAssistantWithCourse(userID, filterCoursesRgx, filterDays)
	if err != nil {
		return args, "", err
	}

	if len(assistants) < 1 {
		return args, "", nil
	}

	mapAssistant := map[int64]map[string]interface{}{}
	for _, val := range assistants {

		image := fl.UsrNoPhotoURL
		if val.FileID.Valid {
			image = "/api/v1/file/profile/" + val.FileID.String + ".jpg"
		}

		if _, ok := mapAssistant[val.IdentityCode]; !ok {
			mapAssistant[val.IdentityCode] = map[string]interface{}{
				"name":    val.Name,
				"phone":   val.Phone,
				"line_id": val.LineID,
				"courses": []string{val.CourseName},
				"image":   image,
			}
			continue
		}

		courses := mapAssistant[val.IdentityCode]["courses"].([]string)
		if helper.IsStringInSlice(val.CourseName, courses) {
			continue
		}

		courses = append(courses, val.CourseName)
		mapAssistant[val.IdentityCode] = map[string]interface{}{
			"name":    val.Name,
			"phone":   val.Phone,
			"line_id": val.LineID,
			"courses": courses,
			"image":   image,
		}
	}

	for _, val := range mapAssistant {
		args = append(args, val)
	}
	if len(args) <= 0 {

	}
	rand.Seed(time.Now().UTC().UnixNano())
	index := rand.Intn(len(msgAssistant))

	nilMsg := msgAssistant[index]
	return args, nilMsg, nil
}

func handleInformation(text string, userID int64) ([]map[string]interface{}, error) {

	args := []map[string]interface{}{}

	params := sEntity{
		text:   text,
		userID: userID,
	}

	// get time entity
	filterTime := params.getTime()
	// get course entity
	filterCourses := params.getCourse()
	filterCoursesLen := len(filterCourses)
	filterCoursesRgx := regexp.MustCompile(strings.Join(filterCourses, "|"))

	scheduleID, err := cs.SelectScheduleIDByUserID(userID, cs.PStatusStudent)
	if err != nil {
		return args, nil
	}

	// select courses details by scheduleID
	courses, err := cs.SelectByScheduleID(scheduleID, cs.StatusScheduleActive)
	if err != nil {
		return args, nil
	}

	var activeScheduleID []int64
	for _, val := range courses {
		// check if course name not match with regex
		if filterCoursesLen > 0 {
			if !filterCoursesRgx.MatchString(strings.ToLower(val.Course.Name)) {
				continue
			}
		}
		activeScheduleID = append(activeScheduleID, val.Schedule.ID)
	}

	info, err := bot.SelectInfoWithFile(activeScheduleID, filterTime)
	if err != nil {
		return args, err
	}

	for _, val := range info {
		args = append(args, map[string]interface{}{
			"id":          val.ID,
			"title":       val.Title,
			"description": val.Description.String,
			"posted_at":   val.CreatedAt.Unix(),
			"image":       "/api/v1/file/default/information.png",
		})
	}

	return args, nil
}

func handleSchedule(text string, userID int64) ([]map[string]interface{}, string, error) {

	args := []map[string]interface{}{}
	params := sEntity{
		text:   text,
		userID: userID,
	}

	// get days entity
	filterTime := params.getTime()
	filterDays := params.getDay()
	filterDays = append(filterDays, helper.TimeToDayInt(filterTime...)...)

	// get course rgx
	var filterCoursesRgx sql.NullString
	filterCourses := params.getCourse()
	if len(filterCourses) > 0 {
		filterCoursesRgx = sql.NullString{
			Valid:  true,
			String: strings.Join(filterCourses, "|"),
		}
	}

	schedules, err := bot.SelectScheduleWithCourse(userID, filterCoursesRgx, filterDays)
	if err != nil {
		return args, "", err
	}

	for _, val := range schedules {
		t1 := helper.MinutesToTimeString(val.StartTime)
		t2 := helper.MinutesToTimeString(val.EndTime)
		t := t1 + " - " + t2
		day := helper.IntDayToString(val.Day)
		args = append(args, map[string]interface{}{
			"course_name": val.CourseName,
			"day":         day,
			"place":       val.Place,
			"time":        t,
		})
	}
	rand.Seed(time.Now().UTC().UnixNano())
	index := rand.Intn(len(msgAssistant))

	nilMsg := msgSchedule[index]
	return args, nilMsg, nil
}

func handleAssignment(text string, userID int64) ([]map[string]interface{}, string, error) {

	args := []map[string]interface{}{}
	params := sEntity{
		text:   text,
		userID: userID,
	}

	filterTime := params.getTime()

	// get course rgx
	var filterCoursesRgx sql.NullString
	filterCourses := params.getCourse()
	if len(filterCourses) > 0 {
		filterCoursesRgx = sql.NullString{
			Valid:  true,
			String: strings.Join(filterCourses, "|"),
		}
	}

	assignments, err := bot.SelectAssignmentWithCourse(userID, filterCoursesRgx, filterTime)
	if err != nil {
		return args, "", err
	}

	for _, val := range assignments {
		args = append(args, map[string]interface{}{
			"url":         "/api/v1/assignment/" + val.ID,
			"name":        val.Name,
			"description": val.Description,
			"due_date":    val.DueDate.Unix(),
			"course_name": val.CourseName,
		})
	}

	rand.Seed(time.Now().UTC().UnixNano())
	index := rand.Intn(len(msgAssignment))

	nilMsg := msgAssignment[index]
	return args, nilMsg, nil
}

func handleGrade(text string, userID int64) ([]map[string]interface{}, string, error) {

	args := []map[string]interface{}{}
	params := sEntity{
		text:   text,
		userID: userID,
	}

	filterTime := params.getTime()

	// get course rgx
	var filterCoursesRgx sql.NullString
	filterCourses := params.getCourse()
	if len(filterCourses) > 0 {
		filterCoursesRgx = sql.NullString{
			Valid:  true,
			String: strings.Join(filterCourses, "|"),
		}
	}

	grades, err := bot.SelectGradeWithCourse(userID, filterCoursesRgx, filterTime)
	if err != nil {
		return args, "", err
	}

	for _, val := range grades {
		args = append(args, map[string]interface{}{
			"url":         "/api/v1/assignment/" + val.AssignmentID,
			"name":        val.Name,
			"score":       fmt.Sprintf("%.3g", val.Score),
			"scored_time": val.UpdatedAt.Unix(),
			"course_name": val.CourseName,
		})
	}

	rand.Seed(time.Now().UTC().UnixNano())
	index := rand.Intn(len(msgGrade))

	nilMsg := msgGrade[index]
	return args, nilMsg, nil
}

func handleGreeting(username string) string {
	rand.Seed(time.Now().UTC().UnixNano())
	index := rand.Intn(len(msgGreet))
	msg := msgGreet[index]

	response := msg.text
	if msg.isExistName {
		response = fmt.Sprintf(msg.text, username)
	}

	return response
}

func handleAboutBot() string {
	rand.Seed(time.Now().UTC().UnixNano())
	index := rand.Intn(len(msgAboutBot))
	response := msgAboutBot[index]
	return response
}

func handleAboutStudent(username string) string {
	rand.Seed(time.Now().UTC().UnixNano())
	index := rand.Intn(len(msgAboutStudent))
	msg := msgAboutStudent[index]

	response := msg.text
	if msg.isExistName {
		response = fmt.Sprintf(msg.text, username)
	}

	return response
}

func handleAboutCreator() string {
	rand.Seed(time.Now().UTC().UnixNano())
	index := rand.Intn(len(msgAboutCreator))
	response := msgAboutCreator[index]
	return response
}

func handleKidding() string {
	rand.Seed(time.Now().UTC().UnixNano())
	index := rand.Intn(len(msgKidding))
	response := msgKidding[index]
	return response
}

package models

import "time"

type TimetableRequest struct {
	CourseCode  string `json:"courseCode"`
	CourseName  string `json:"courseName"`
	FacultyCode string `json:"facultyCode"`
	Venue       string `json:"venue"`
	Subgroup    string `json:"subgroup"`
	Time        string `json:"time"`
	Department  string `json:"department"`
	StartDate   string `json:"startDate"`
	EndDate     string `json:"endDate"`
	Day         string `json:"day"`
	Type        string `json:"type"`
}

type TimetableEntry struct {
	ID          string    `bson:"_id"`
	CourseCode  string    `bson:"courseCode"`
	CourseName  string    `bson:"courseName"`
	FacultyCode string    `bson:"facultyCode"`
	Venue       string    `bson:"venue"`
	Subgroup    string    `bson:"subgroup"`
	Department  string    `bson:"department"`
	Time        string    `bson:"time"`
	Date        string    `bson:"date"`
	Duration    int       `bson:"duration"`
	Type        string    `bson:"type"`
	CreatedAt   time.Time `bson:"createdAt"`
}

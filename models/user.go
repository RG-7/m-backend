package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID             primitive.ObjectID `json:"id" bson:"_id"`
	Email          string             `json:"email" bson:"email"`
	Password       string             `json:"password" bson:"password"`
	Role           string             `json:"role" bson:"role"`
	MobileNumber   string             `json:"mobileno" bson:"mobileno"`
	EmployeeID     string             `json:"employeeId" bson:"employeeId"`
	FacultyCode    string             `json:"facultyCode" bson:"facultyCode"`
	Department     string             `json:"department" bson:"department"`
	DepartmentCode string             `json:"departmentCode" bson:"departmentCode"`
	Designation    string             `json:"designation" bson:"designation"`
	Availability   string             `json:"availability" bson:"availability"`
}

package controllers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/RG-7/m-backend/database"
	"github.com/RG-7/m-backend/helpers"
	"github.com/RG-7/m-backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func generateCommonTimetableEntries(ctt models.CommonTimetableRequest) ([]interface{}, error) {
	startDate, err := helpers.ParseDate(ctt.StartDate)
	if err != nil {
		return nil, err
	}
	endDate, err := helpers.ParseDate(ctt.EndDate)
	if err != nil {
		return nil, err
	}

	dayMap := map[string]time.Weekday{
		"Sunday":    time.Sunday,
		"Monday":    time.Monday,
		"Tuesday":   time.Tuesday,
		"Wednesday": time.Wednesday,
		"Thursday":  time.Thursday,
		"Friday":    time.Friday,
		"Saturday":  time.Saturday,
	}

	targetDay, exists := dayMap[ctt.Day]
	if !exists {
		return nil, err
	}

	var entries []interface{}
	currentDate := startDate

	for !currentDate.After(endDate) {
		if currentDate.Weekday() == targetDay {
			newEntry := models.CommonTimetableEntry{
				ID:          primitive.NewObjectID().Hex(),
				CourseCode:  ctt.CourseCode,
				CourseName:  ctt.CourseName,
				FacultyCode: ctt.FacultyCode,
				Venue:       ctt.Venue,
				Subgroup:    ctt.Subgroup,
				Department:  ctt.Department,
				Date:        currentDate.Format("2006-01-02"),
				Type:        ctt.Type,
				Time:        ctt.Time,
				CreatedAt:   time.Now(),
			}
			entries = append(entries, newEntry)
		}
		currentDate = currentDate.AddDate(0, 0, 1)
	}
	return entries, nil
}

func CreateCommonTimetableEntry(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var ctt models.CommonTimetableRequest
	err := json.NewDecoder(r.Body).Decode(&ctt)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	entries, err := generateCommonTimetableEntries(ctt)
	if err != nil {
		http.Error(w, "Failed to generate timetable entries: "+err.Error(), http.StatusInternalServerError)
		return
	}

	facultyCollection := database.Client.Database("ttms").Collection("facultyTT")
	_, err = facultyCollection.InsertMany(context.TODO(), entries)
	if err != nil {
		http.Error(w, "Failed to create timetable entries", http.StatusInternalServerError)
		return
	}

	roomCollection := database.Client.Database("ttms").Collection("roomsTT")
	_, err = roomCollection.InsertMany(context.TODO(), entries)
	if err != nil {
		http.Error(w, "Failed to create timetable entries", http.StatusInternalServerError)
		return
	}

	commonSubGroupCollection := database.Client.Database("ttms").Collection("commonSubGroupTT")
	_, err = commonSubGroupCollection.InsertMany(context.TODO(), entries)
	if err != nil {
		http.Error(w, "Failed to create timetable entries", http.StatusInternalServerError)
		return
	}

	// Return all generated entries
	json.NewEncoder(w).Encode(map[string]interface{}{
		"facultyTT":  entries,
		"roomTT":     entries,
		"subgroupTT": entries,
	})
}

// del
func DeleteCommonTimetableEntry(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var ctt models.CommonTimetableRequest
	err := json.NewDecoder(r.Body).Decode(&ctt)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Parse the start and end dates
	startDate, err := helpers.ParseDate(ctt.StartDate)
	if err != nil {
		http.Error(w, "Invalid start date", http.StatusBadRequest)
		return
	}
	endDate, err := helpers.ParseDate(ctt.EndDate)
	if err != nil {
		http.Error(w, "Invalid end date", http.StatusBadRequest)
		return
	}

	// Map of weekdays to compare
	dayMap := map[string]time.Weekday{
		"Sunday":    time.Sunday,
		"Monday":    time.Monday,
		"Tuesday":   time.Tuesday,
		"Wednesday": time.Wednesday,
		"Thursday":  time.Thursday,
		"Friday":    time.Friday,
		"Saturday":  time.Saturday,
	}

	// Get the target day for comparison
	targetDay, exists := dayMap[ctt.Day]
	if !exists {
		http.Error(w, "Invalid day", http.StatusBadRequest)
		return
	}

	// Generate all matching dates within the given range that match the target day
	var dates []string
	currentDate := startDate
	for !currentDate.After(endDate) {
		if currentDate.Weekday() == targetDay {
			dates = append(dates, currentDate.Format("2006-01-02"))
		}
		currentDate = currentDate.AddDate(0, 0, 1)
	}

	// Delete all entries that match the course code and dates
	for _, date := range dates {
		filter := bson.M{
			"courseCode":  ctt.CourseCode,
			"courseName":  ctt.CourseName,
			"facultyCode": ctt.FacultyCode,
			"venue":       ctt.Venue,
			"department":  ctt.Department,
			"subgroup":    ctt.Subgroup,
			"time":        ctt.Time,
			"type":        ctt.Type,
			"date":        date, // Match the date for each iteration
		}

		// Delete from faculty timetable
		facultyCollection := database.Client.Database("ttms").Collection("facultyTT")
		result, err := facultyCollection.DeleteMany(context.TODO(), filter)
		if err != nil {
			log.Println("Failed to delete faculty timetable entries:", err)
			http.Error(w, "Failed to delete faculty timetable entries", http.StatusInternalServerError)
			return
		}
		log.Printf("Deleted %d entries from facultyTT for date %s\n", result.DeletedCount, date)

		// Delete from room timetable
		roomCollection := database.Client.Database("ttms").Collection("roomsTT")
		result, err = roomCollection.DeleteMany(context.TODO(), filter)
		if err != nil {
			log.Println("Failed to delete room timetable entries:", err)
			http.Error(w, "Failed to delete room timetable entries", http.StatusInternalServerError)
			return
		}
		log.Printf("Deleted %d entries from roomsTT for date %s\n", result.DeletedCount, date)

		commonSubGroupCollection := database.Client.Database("ttms").Collection("commonSubGroupTT")
		result, err = commonSubGroupCollection.DeleteMany(context.TODO(), filter)
		if err != nil {
			log.Println("Failed to delete subgroup timetable entries:", err)
			http.Error(w, "Failed to delete subgroup timetable entries", http.StatusInternalServerError)
			return
		}
		log.Printf("Deleted %d entries from subgroupTT on date %s\n", result.DeletedCount, date)
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Timetable entries deleted successfully",
	})

}

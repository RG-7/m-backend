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

// Function to generate timetable entries for each occurrence of the day between start and end date
func generateTimetableEntries(tt models.TimetableRequest) ([]interface{}, error) {
	startDate, err := helpers.ParseDate(tt.StartDate)
	if err != nil {
		return nil, err
	}
	endDate, err := helpers.ParseDate(tt.EndDate)
	if err != nil {
		return nil, err
	}

	// Convert weekday string to time.Weekday
	dayMap := map[string]time.Weekday{
		"Sunday":    time.Sunday,
		"Monday":    time.Monday,
		"Tuesday":   time.Tuesday,
		"Wednesday": time.Wednesday,
		"Thursday":  time.Thursday,
		"Friday":    time.Friday,
		"Saturday":  time.Saturday,
	}

	targetDay, exists := dayMap[tt.Day]
	if !exists {
		return nil, err
	}

	var entries []interface{}
	currentDate := startDate

	// Loop through the date range and find matching weekdays
	for !currentDate.After(endDate) {
		if currentDate.Weekday() == targetDay {
			newEntry := models.TimetableEntry{
				ID:          primitive.NewObjectID().Hex(), // Generate MongoDB ObjectID and convert to string
				CourseCode:  tt.CourseCode,
				CourseName:  tt.CourseName,
				FacultyCode: tt.FacultyCode,
				Venue:       tt.Venue,
				Subgroup:    tt.Subgroup,
				Department:  tt.Department,
				Date:        currentDate.Format("2006-01-02"), // Store the date in "YYYY-MM-DD" format
				Type:        tt.Type,
				Time:        tt.Time,
				Duration:    helpers.GetDuration(tt.Type), // Get the duration based on session type
				CreatedAt:   time.Now(),                   // Store the creation time
			}
			entries = append(entries, newEntry)
		}
		currentDate = currentDate.AddDate(0, 0, 1) // Move to the next day
	}

	return entries, nil
}

func CreateTimetableEntry(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var tt models.TimetableRequest
	err := json.NewDecoder(r.Body).Decode(&tt)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	entries, err := generateTimetableEntries(tt)
	if err != nil {
		http.Error(w, "Error processing timetable entries", http.StatusInternalServerError)
		return
	}

	// Save the timetable entries to MongoDB
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

	// Generate subgroup-specific entries
	var subgroupEntries []interface{}
	for _, entry := range entries {
		ttEntry := entry.(models.TimetableEntry) // Type assertion

		// Generate individual subgroup entries
		subgroups := helpers.GenerateSubgroups(ttEntry.Subgroup)
		for _, subgroup := range subgroups {
			newEntry := models.TimetableEntry{
				ID:          primitive.NewObjectID().Hex(),
				CourseCode:  ttEntry.CourseCode,
				CourseName:  ttEntry.CourseName,
				FacultyCode: ttEntry.FacultyCode,
				Venue:       ttEntry.Venue,
				Subgroup:    subgroup, // Assign individual subgroup
				Department:  ttEntry.Department,
				Date:        ttEntry.Date,
				Type:        ttEntry.Type,
				Time:        tt.Time,
				Duration:    helpers.GetDuration(tt.Type),
				CreatedAt:   time.Now(),
			}
			subgroupEntries = append(subgroupEntries, newEntry)
		}
	}

	// Save subgroup-specific timetable entries
	subgroupCollection := database.Client.Database("ttms").Collection("subgroupTT")
	_, err = subgroupCollection.InsertMany(context.TODO(), subgroupEntries)
	if err != nil {
		http.Error(w, "Failed to create subgroup timetable entries", http.StatusInternalServerError)
		return
	}

	// Return all generated entries
	json.NewEncoder(w).Encode(map[string]interface{}{
		"facultyTT":  entries,
		"roomTT":     entries,
		"subgroupTT": subgroupEntries,
	})
}

// del
func DeleteTimetableEntry(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var tt models.TimetableRequest
	err := json.NewDecoder(r.Body).Decode(&tt)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Parse the start and end dates
	startDate, err := helpers.ParseDate(tt.StartDate)
	if err != nil {
		http.Error(w, "Invalid start date", http.StatusBadRequest)
		return
	}
	endDate, err := helpers.ParseDate(tt.EndDate)
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
	targetDay, exists := dayMap[tt.Day]
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

	// Loop through each date and delete the entries
	for _, date := range dates {
		filter := bson.M{
			"courseCode":  tt.CourseCode,
			"courseName":  tt.CourseName,
			"facultyCode": tt.FacultyCode,
			"venue":       tt.Venue,
			"department":  tt.Department,
			"subgroup":    tt.Subgroup,
			"time":        tt.Time,
			"type":        tt.Type,
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

		// Generate subgroups
		subgroups := helpers.GenerateSubgroups(tt.Subgroup)

		// Delete entries from subgroup timetable for each subgroup
		for _, subgroup := range subgroups {
			subgroupFilter := bson.M{
				"courseCode":  tt.CourseCode,
				"facultyCode": tt.FacultyCode,
				"venue":       tt.Venue,
				"department":  tt.Department,
				"subgroup":    subgroup, // Replace with the specific subgroup
				"time":        tt.Time,
				"type":        tt.Type,
				"date":        date, // Match the date for each iteration
			}

			subgroupCollection := database.Client.Database("ttms").Collection("subgroupTT")
			result, err := subgroupCollection.DeleteMany(context.TODO(), subgroupFilter)
			if err != nil {
				log.Println("Failed to delete subgroup timetable entries:", err)
				http.Error(w, "Failed to delete subgroup timetable entries", http.StatusInternalServerError)
				return
			}
			log.Printf("Deleted %d entries from subgroupTT for subgroup %s on date %s\n", result.DeletedCount, subgroup, date)
		}
	}

	// Return success message
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Timetable entries deleted successfully",
	})
}

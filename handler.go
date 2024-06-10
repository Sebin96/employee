package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func CreateEmployeeHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var emp *Employee
		err := json.NewDecoder(r.Body).Decode(&emp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Call a function to insert the employee data into the database
		emp, err = CreateEmployeeAPI(db, emp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Set the appropriate response status and content type
		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")

		// Encode the response JSON
		json.NewEncoder(w).Encode(emp)
	}
}

func ReadEmployeeHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the employee ID from the request parameters
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, "Invalid employee ID", http.StatusBadRequest)
			return
		}

		var emp *Employee

		// Call the API function to retrieve the employee by ID
		emp, apiErr := ReadEmployeeAPI(db, id, emp)
		if apiErr != nil {
			if apiErr == sql.ErrNoRows {
				http.Error(w, "Employee not found", http.StatusNotFound)
			} else {
				http.Error(w, apiErr.Error(), http.StatusInternalServerError)
			}
			return
		}

		// Set the appropriate content type header
		w.Header().Set("Content-Type", "application/json")

		// Encode the retrieved employee as JSON and send it in the response
		err = json.NewEncoder(w).Encode(emp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func ReadEmployeeListHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse query parameters for pagination
		page, err := strconv.Atoi(r.URL.Query().Get("page"))
		if err != nil || page < 1 {
			page = 1 // Default to page 1 if page parameter is invalid
		}
		limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
		if err != nil || limit < 1 {
			limit = 10 // Default to 10 records per page if limit parameter is invalid
		}
		offset := (page - 1) * limit

		// Call the API function to retrieve paginated employees
		employees, apiErr := ReadEmployeeListAPI(db, limit, offset)
		if apiErr != nil {
			http.Error(w, apiErr.Error(), http.StatusInternalServerError)
			return
		}

		// Set the appropriate content type header
		w.Header().Set("Content-Type", "application/json")

		// Encode the retrieved employees as JSON and send them in the response
		err = json.NewEncoder(w).Encode(employees)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func UpdateEmployeeHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the employee ID from the request parameters
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, "Invalid employee ID", http.StatusBadRequest)
			return
		}

		var empReq *Employee
		err = json.NewDecoder(r.Body).Decode(&empReq)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Call the API function to update the employee by ID
		emp, apiErr := UpdateEmployeeAPI(db, id, empReq)
		if apiErr != nil {
			if apiErr == sql.ErrNoRows {
				http.Error(w, "Employee not found", http.StatusNotFound)
			} else {
				http.Error(w, apiErr.Error(), http.StatusInternalServerError)
			}
			return
		}

		// Set the appropriate content type header
		w.Header().Set("Content-Type", "application/json")

		// Encode the updated employee as JSON and send it in the response
		err = json.NewEncoder(w).Encode(emp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func DeleteEmployeeHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the employee ID from the request parameters
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, "Invalid employee ID", http.StatusBadRequest)
			return
		}

		// Call the API function to delete the employee by ID
		apiErr := DeleteEmployeeAPI(db, id)
		if apiErr != nil {
			if apiErr == sql.ErrNoRows {
				http.Error(w, "Employee not found", http.StatusNotFound)
			} else {
				http.Error(w, apiErr.Error(), http.StatusInternalServerError)
			}
			return
		}

		// Respond with success message
		successMessage := map[string]string{"message": "Employee deleted successfully"}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(successMessage); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

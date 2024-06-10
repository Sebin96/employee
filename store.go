package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"
)

var db *sql.DB

func initDB() *sql.DB {
	const (
		host     = "localhost"
		port     = 5432
		user     = "postgres"
		password = "9101996"
		dbname   = "organisation"
	)

	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error connecting to the database:", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Error testing database connection:", err)
	}

	createTableSQL := `
	CREATE TABLE IF NOT EXISTS employee (
		ID SERIAL PRIMARY KEY,
		Name VARCHAR(100) NOT NULL,
		Designation VARCHAR(100) NOT NULL,
		Salary FLOAT8 NOT NULL,
		CreatedAt TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		UpdatedAt TIMESTAMPTZ NOT NULL DEFAULT NOW()

	);
	`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("Unable to create table: %v", err)
	}

	fmt.Println("Successfully connected to the database and ensured employee table exists!")

	return db
}

func CreateEmployeeStore(db *sql.DB, emp *Employee) error {
	insertEmployeeSQL := `
        INSERT INTO employee (ID, Name, Designation, Salary, CreatedAt, UpdatedAt)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING ID;
    `

	var empID int
	err := db.QueryRow(insertEmployeeSQL, emp.ID, emp.Name, emp.Designation, emp.Salary, time.Now(), time.Now()).Scan(&empID)
	if err != nil {
		return err
	}

	// Update the Employee ID in the passed struct
	emp.ID = empID
	return nil
}

func ReadEmployeeStore(db *sql.DB, id int) (*Employee, error) {

	emp := &Employee{}

	err := db.QueryRow("SELECT ID, Name, Designation, Salary, CreatedAt,UpdatedAt FROM employee WHERE ID = $1", id).
		Scan(&emp.ID, &emp.Name, &emp.Designation, &emp.Salary, &emp.CreatedAt, &emp.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return emp, nil
}

func ReadEmployeeListStore(db *sql.DB, limit, offset int) ([]Employee, error) {
	// Execute the query to fetch paginated employees
	rows, err := db.Query("SELECT ID, Name, Designation, Salary, CreatedAt, UpdatedAt FROM employee ORDER BY ID LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate through the result set and populate the employees slice
	var employees []Employee
	for rows.Next() {
		var emp Employee
		err := rows.Scan(&emp.ID, &emp.Name, &emp.Designation, &emp.Salary, &emp.CreatedAt, &emp.UpdatedAt)
		if err != nil {
			return nil, err
		}
		employees = append(employees, emp)
	}

	// Check for any errors during rows iteration
	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(employees) == 0 {
		return nil, errors.New("no records available")
	}

	return employees, nil
}

func UpdateEmployeeStore(db *sql.DB, id int, updatedEmp *Employee) (*Employee, error) {
	// If updatedEmp is not provided, perform only read operation
	emp := &Employee{}
	err := db.QueryRow("SELECT ID, Name, Designation, Salary, CreatedAt, UpdatedAt FROM employee WHERE ID = $1", id).
		Scan(&emp.ID, &emp.Name, &emp.Designation, &emp.Salary, &emp.CreatedAt, &emp.UpdatedAt)
	if err != nil {
		return nil, err
	}

	if updatedEmp != nil {
		if updatedEmp.Name == "" {
			updatedEmp.Name = emp.Name
		}
		if updatedEmp.Salary == 0 {
			updatedEmp.Salary = emp.Salary
		}
		if updatedEmp.Designation == "" {
			updatedEmp.Designation = emp.Designation
		}
		if updatedEmp.ID == 0 {
			updatedEmp.ID = emp.ID
		}
	}

	// If updatedEmp is provided, perform update operation
	_, err = db.Exec("UPDATE employee SET Name = $1, Designation = $2, Salary = $3, UpdatedAt = $4 WHERE ID = $5",
		updatedEmp.Name, updatedEmp.Designation, updatedEmp.Salary, time.Now(), id)
	if err != nil {
		return nil, err
	}

	// Return the updated employee data
	return updatedEmp, nil
}

func DeleteEmployeeStore(db *sql.DB, id int) error {

	emp := &Employee{}
	err := db.QueryRow("SELECT ID, Name, Designation, Salary, CreatedAt, UpdatedAt FROM employee WHERE ID = $1", id).
		Scan(&emp.ID, &emp.Name, &emp.Designation, &emp.Salary, &emp.CreatedAt, &emp.UpdatedAt)
	if err != nil {
		return err
	}
	if emp != nil {
		_, err := db.Exec("DELETE FROM employee WHERE ID = $1", id)
		if err != nil {
			return err
		}
	} else {
		return errors.New("no result received before timeout")

	}
	// No need to fetch employee data after deletion
	return nil

}

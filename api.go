package main

import (
	"database/sql"
	"errors"
	"time"
)

type Employee struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Designation string    `json:"designation"`
	Salary      float64   `json:"salary"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

var ErrTimeoutCreatingEmployee = errors.New("timeout occurred while creating employee")
var ErrTimeoutReadingEmployee = errors.New("timeout occurred while reading employee")
var ErrTimeoutUpdatingEmployee = errors.New("timeout occurred while updating employee")
var ErrTimeoutDeletingEmployee = errors.New("timeout occurred while deleting employee")

func CreateEmployeeAPI(db *sql.DB, emp *Employee) (*Employee, error) {
	// Channel to receive errors from the store operation
	errChan := make(chan error, 1)

	// Asynchronously call the CreateEmployeeStore function
	go func() {
		err := CreateEmployeeStore(db, emp)
		errChan <- err
	}()

	// Wait for either a timeout or an error from the store operation
	select {
	case <-time.After(5 * time.Second): // Timeout after 5 seconds
		return nil, ErrTimeoutCreatingEmployee
	case err := <-errChan:
		if err != nil {
			return nil, err
		}
	}

	// If we reach this point, it means employee creation was successful
	return emp, nil
}

func ReadEmployeeAPI(db *sql.DB, id int, emp *Employee) (*Employee, error) {
	// Channel to receive errors from the store operation
	errChan := make(chan error, 1)
	empChan := make(chan *Employee, 1)

	// Asynchronously call the ReadEmployeeStore function
	go func() {
		emp, err := ReadEmployeeStore(db, id)
		if err != nil {
			errChan <- err
			return
		}
		empChan <- emp
	}()

	// Wait for either a timeout or an error from the store operation
	select {
	case <-time.After(5 * time.Second): // Timeout after 5 seconds
		return nil, ErrTimeoutReadingEmployee
	case err := <-errChan:
		if err != nil {
			return nil, err
		}
	case emp := <-empChan:
		return emp, nil
	}

	// If we reach this point, it means no result was received before the timeout
	return nil, errors.New("no result received before timeout")
}

func ReadEmployeeListAPI(db *sql.DB, limit, offset int) ([]Employee, error) {
	// Channel to receive errors from the store operation
	errChan := make(chan error, 1)
	empChan := make(chan []Employee, 1)

	// Asynchronously call the ReadEmployeeListStore function
	go func() {
		emp, err := ReadEmployeeListStore(db, limit, offset)
		if err != nil {
			errChan <- err
			return
		}
		empChan <- emp
	}()

	// Wait for either a timeout or an error from the store operation
	select {
	case <-time.After(5 * time.Second): // Timeout after 5 seconds
		return nil, ErrTimeoutReadingEmployee
	case err := <-errChan:
		if err != nil {
			return nil, err
		}
	case emp := <-empChan:
		return emp, nil
	}

	// If we reach this point, it means no result was received before the timeout
	return nil, errors.New("no result received before timeout")
}

func UpdateEmployeeAPI(db *sql.DB, id int, emp *Employee) (*Employee, error) {
	// Channel to receive errors from the store operation
	errChan := make(chan error, 1)
	empChan := make(chan *Employee, 1)

	// Asynchronously call the ReadEmployeeStore function
	go func() {
		emp, err := UpdateEmployeeStore(db, id, emp)
		if err != nil {
			errChan <- err
			return
		}
		empChan <- emp
	}()

	// Wait for either a timeout or an error from the store operation
	select {
	case <-time.After(5 * time.Second): // Timeout after 5 seconds
		return nil, ErrTimeoutUpdatingEmployee
	case err := <-errChan:
		if err != nil {
			return nil, err
		}
	case emp := <-empChan:
		return emp, nil
	}

	// If we reach this point, it means no result was received before the timeout
	return nil, errors.New("no result received before timeout")
}

func DeleteEmployeeAPI(db *sql.DB, id int) error {
	// Channel to receive errors from the store operation
	errChan := make(chan error, 1)

	// Asynchronously call the ReadEmployeeStore function
	go func() {
		err := DeleteEmployeeStore(db, id)
		errChan <- err
	}()

	// Wait for either a timeout or an error from the store operation
	select {
	case <-time.After(5 * time.Second): // Timeout after 5 seconds
		return ErrTimeoutDeletingEmployee
	case err := <-errChan:
		return err
	}
}

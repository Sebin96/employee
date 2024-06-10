package main

import (
	"database/sql"
	"fmt"
	"reflect"
	"testing"
	"time"

	_ "github.com/lib/pq"
)

// Helper function to initialize a test database
func initTestDB(t *testing.T) *sql.DB {
	const (
		host     = "localhost"
		port     = 5432
		user     = "postgres"
		password = "9101996"
		dbname   = "testDB" // Use a different database for testing
	)

	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		t.Fatalf("Error connecting to the test database: %v", err)
	}

	return db
}

// CreateTableEmployee creates the employee table
func CreateTableEmployee(db *sql.DB) error {
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

	// Execute the SQL statement to create the table
	_, err := db.Exec(createTableSQL)
	if err != nil {
		return fmt.Errorf("Unable to create employee table: %v", err)
	}

	return nil
}

// InsertTableEmployee inserts the employee table
func InsertTableEmployee(db *sql.DB, employees []Employee) error {
	insertSQL := `
        INSERT INTO employee (Name, Designation, Salary, CreatedAt, UpdatedAt)
        VALUES ($1, $2, $3, NOW(), NOW())
        RETURNING ID;
    `

	for i, emp := range employees {
		err := db.QueryRow(insertSQL, emp.Name, emp.Designation, emp.Salary).Scan(&employees[i].ID)
		if err != nil {
			return fmt.Errorf("Unable to insert employee: %v", err)
		}
	}
	return nil
}

// DeleteTableEmployee deletes the employee table from the provided database
func DeleteTableEmployee(db *sql.DB) error {
	// SQL statement to delete the employee table
	deleteTableSQL := `DROP TABLE IF EXISTS employee;`

	// Execute the SQL statement to delete the table
	_, err := db.Exec(deleteTableSQL)
	if err != nil {
		return fmt.Errorf("Unable to delete employee table: %v", err)
	}

	return nil
}

func TestCreateEmployeeAPI(t *testing.T) {

	// Set up a test database connection
	db := initTestDB(t)
	defer db.Close()

	err := CreateTableEmployee(db)
	if err != nil {
		t.Fatalf("Unable to create employee table: %v", err)
	}
	type args struct {
		db  *sql.DB
		emp *Employee
	}
	tests := []struct {
		name    string
		args    args
		want    *Employee
		wantErr bool
	}{
		{
			name: "Successfull Creation of employee",
			args: args{
				db: db,
				emp: &Employee{
					ID:          1,
					Name:        "John Doe",
					Designation: "Engineer",
					Salary:      50000,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				},
			},
			want: &Employee{
				ID:          1,
				Name:        "John Doe",
				Designation: "Engineer",
				Salary:      50000,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			wantErr: false,
		},
		{
			name: "Error Creating Duplicate Employee",
			args: args{
				db: db,
				emp: &Employee{
					ID:          1,
					Name:        "John Doe",
					Designation: "Engineer",
					Salary:      50000,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Error Creating Employe without name",
			args: args{
				db: db,
				emp: &Employee{
					ID:          1,
					Name:        "",
					Designation: "Engineer",
					Salary:      50000,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				},
			},
			want:    nil,
			wantErr: true,
		},

		{
			name: "Error Creating Employee without designation",
			args: args{
				db: db,
				emp: &Employee{
					ID:          1,
					Name:        "Kiran",
					Designation: "",
					Salary:      50000,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				},
			},
			want:    nil,
			wantErr: true,
		},

		{
			name: "Error Creating Employee without createdAt time",
			args: args{
				db: db,
				emp: &Employee{
					ID:          1,
					Name:        "Kiran",
					Designation: "",
					Salary:      12340,
					UpdatedAt:   time.Now(),
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	// Inside the for loop of the TestCreateEmployeeAPI function
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateEmployeeAPI(tt.args.db, tt.args.emp)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateEmployeeAPI() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.want != nil {
				tt.want.ID = got.ID
				tt.want.CreatedAt = got.CreatedAt
				tt.want.UpdatedAt = got.UpdatedAt
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateEmployeeAPI() = %v, want %v", got, tt.want)
			}
		})

		// Clean up the test database after each test case
		DeleteTableEmployee(db)
	}

}

func TestReadEmployeeAPI(t *testing.T) {
	// Set up a test database connection
	db := initTestDB(t)
	defer db.Close()

	err := CreateTableEmployee(db)
	if err != nil {
		t.Fatalf("Unable to create employee table: %v", err)
	}

	// Prepare test data
	employees := []Employee{
		{Name: "Dan", Designation: "Software Developer", Salary: 23456.00},
	}

	err = InsertTableEmployee(db, employees)
	if err != nil {
		t.Fatalf("Unable to insert employee table: %v", err)
	}

	type args struct {
		db  *sql.DB
		id  int
		emp *Employee
	}
	tests := []struct {
		name    string
		args    args
		want    *Employee
		wantErr bool
	}{
		{
			name: "Successfully read employee",
			args: args{
				db:  db,
				id:  employees[0].ID,
				emp: &Employee{},
			},
			want: &Employee{
				ID:          1,
				Name:        "Dan",
				Designation: "Software Developer",
				Salary:      23456.00,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			wantErr: false,
		},
		{
			name: "Read employee with invalid id",
			args: args{
				db:  db,
				id:  2,
				emp: &Employee{},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Read employee with no id",
			args: args{
				db:  db,
				emp: &Employee{},
			},
			want:    nil,
			wantErr: true,
		},

		{
			name: "Read employee with no response",
			args: args{
				db: db,
				id: 1,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadEmployeeAPI(tt.args.db, tt.args.id, tt.args.emp)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadEmployeeAPI() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.want != nil {
				tt.want.CreatedAt = got.CreatedAt
				tt.want.UpdatedAt = got.UpdatedAt
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadEmployeeAPI() = %v, want %v", got, tt.want)
			}
		})
		// Clean up the test database after each test case
		DeleteTableEmployee(db)
	}

}

func TestReadEmployeeListAPI(t *testing.T) {

	// Set up a test database connection
	db := initTestDB(t)
	defer db.Close()

	err := CreateTableEmployee(db)
	if err != nil {
		t.Fatalf("Unable to create employee table: %v", err)
	}

	// Prepare test data
	employees := []Employee{
		{Name: "Sen", Designation: "Account Manager", Salary: 44566.00},
		{Name: "Dan", Designation: "Software Developer", Salary: 23456.00},
	}

	err = InsertTableEmployee(db, employees)
	if err != nil {
		t.Fatalf("Unable to insert employee table: %v", err)
	}

	type args struct {
		db     *sql.DB
		limit  int
		offset int
	}
	tests := []struct {
		name    string
		args    args
		want    []Employee
		wantErr bool
	}{

		{
			name: "Successfully read employee list",
			args: args{
				db:     db,
				limit:  2,
				offset: 0,
			},
			want: []Employee{
				Employee{
					ID:          employees[0].ID,
					Name:        "Sen",
					Designation: "Account Manager",
					Salary:      44566.00,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				},
				Employee{
					ID:          employees[1].ID,
					Name:        "Dan",
					Designation: "Software Developer",
					Salary:      23456.00,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				},
			},
			wantErr: false,
		},
		{
			name: "Read employee list with invalid offset value",
			args: args{
				db:     db,
				limit:  2,
				offset: 2,
			},
			want:    nil,
			wantErr: false,
			//NB: error will be false ,because if there is no object in response, it will return empty array
		},
		{
			name: "Read employee list with invalid limit value",
			args: args{
				db:     db,
				limit:  0,
				offset: 2,
			},
			want:    nil,
			wantErr: false,
			//NB: error will be false ,because if there no object in response,it will return empty array

		},
		{
			name: "Read employee list with invalid limit and offset value",
			args: args{
				db:     db,
				limit:  0,
				offset: 0,
			},
			want:    nil,
			wantErr: false,
			//NB: error will be false ,because if there no object in response,it will return empty array

		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadEmployeeListAPI(tt.args.db, tt.args.limit, tt.args.offset)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadEmployeeListAPI() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			fmt.Println(got)
			for i := range tt.want {
				tt.want[i].ID = got[i].ID
				tt.want[i].CreatedAt = got[i].CreatedAt
				tt.want[i].UpdatedAt = got[i].UpdatedAt
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadEmployeeListAPI() = %v, want %v", got, tt.want)
			}
		})
		// Clean up the test database after each test case
	}
	DeleteTableEmployee(db)

}

func TestUpdateEmployeeAPI(t *testing.T) {

	// Set up a test database connection
	db := initTestDB(t)
	defer db.Close()

	err := CreateTableEmployee(db)
	if err != nil {
		t.Fatalf("Unable to create employee table: %v", err)
	}

	// Prepare test data
	employees := []Employee{
		{Name: "Dan", Designation: "Software Developer", Salary: 23456.00},
	}

	err = InsertTableEmployee(db, employees)
	if err != nil {
		t.Fatalf("Unable to insert employee table: %v", err)
	}

	type args struct {
		db  *sql.DB
		id  int
		emp *Employee
	}
	tests := []struct {
		name    string
		args    args
		want    *Employee
		wantErr bool
	}{
		{
			name: "Successfully update employee",
			args: args{
				db: db,
				id: employees[0].ID,
				emp: &Employee{
					ID:          employees[0].ID,
					Name:        "Dan",
					Designation: "Software Developer",
					Salary:      11111.00,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				},
			},
			want: &Employee{
				ID:          employees[0].ID,
				Name:        "Dan",
				Designation: "Software Developer",
				Salary:      11111.00,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			wantErr: false,
		},
		{
			name: "Update employee with invalid id",
			args: args{
				db: db,
				id: 2,
				emp: &Employee{
					ID:          employees[0].ID,
					Name:        "Dan",
					Designation: "Software Developer",
					Salary:      11111.00,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Update employee with invalid payload",
			args: args{
				db: db,
				id: employees[0].ID,
				emp: &Employee{
					ID:          2,
					Name:        "Ben",
					Designation: "Software Developer",
					Salary:      23333.00,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Update employee with invalid payload and id",
			args: args{
				db: db,
				id: 3,
				emp: &Employee{
					ID:          2,
					Name:        "Ben",
					Designation: "Software Developer",
					Salary:      23333.00,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UpdateEmployeeAPI(tt.args.db, tt.args.id, tt.args.emp)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateEmployeeAPI() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.want != nil {
				tt.want.CreatedAt = got.CreatedAt
				tt.want.UpdatedAt = got.UpdatedAt
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UpdateEmployeeAPI() = %v, want %v", got, tt.want)
			}
		})
		DeleteTableEmployee(db)

	}
}

func TestDeleteEmployeeAPI(t *testing.T) {
	// Set up a test database connection
	db := initTestDB(t)
	defer db.Close()

	err := CreateTableEmployee(db)
	if err != nil {
		t.Fatalf("Unable to create employee table: %v", err)
	}

	// Prepare test data
	employees := []Employee{
		{Name: "Dan", Designation: "Software Developer", Salary: 23456.00},
	}

	err = InsertTableEmployee(db, employees)
	if err != nil {
		t.Fatalf("Unable to insert employee table: %v", err)
	}
	type args struct {
		db *sql.DB
		id int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Successfully delete employee",
			args: args{
				db: db,
				id: employees[0].ID,
			},
			wantErr: false,
		},
		{
			name: "Delete employee with invalid id",
			args: args{
				db: db,
				id: 2,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DeleteEmployeeAPI(tt.args.db, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("DeleteEmployeeAPI() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

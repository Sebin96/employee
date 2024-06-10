package main

func (cr *CustomRouter) SetupRouter() {

	cr.HandleFunc("/employees", CreateEmployeeHandler(cr.DB)).Methods("POST")
	cr.HandleFunc("/employees/{id}", ReadEmployeeHandler(cr.DB)).Methods("GET")
	cr.HandleFunc("/employeeList", ReadEmployeeListHandler(cr.DB)).Methods("GET")
	cr.HandleFunc("/employees/{id}", UpdateEmployeeHandler(cr.DB)).Methods("PUT")
	cr.HandleFunc("/employees/{id}", DeleteEmployeeHandler(cr.DB)).Methods("DELETE")
}

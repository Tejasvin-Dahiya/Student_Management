package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type Student struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Grade string `json:"grade"`
}

var db *sql.DB

func connectDB() {
	var err error
	db, err = sql.Open("mysql", "root:Alphacamero@7@tcp(127.0.0.1:3306)/studentdb")
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

}

// Get All Students function
func getStudents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	rows, err := db.Query("SELECT * FROM students")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	var students []Student
	for rows.Next() {
		var s Student
		if err := rows.Scan(&s.ID, &s.Name, &s.Age, &s.Grade); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		students = append(students, s)
	}
	json.NewEncoder(w).Encode(students)
}

func getStudentsId(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	urlpara := mux.Vars(r)
	id := urlpara["id"]
	var s Student
	err := db.QueryRow("SELECT id, name, age, grade FROM students WHERE id=?", id).Scan(&s.ID, &s.Name, &s.Age, &s.Grade)
	if err != nil {
		http.Error(w, "Student not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(s)
}

func createStudent(w http.ResponseWriter, r *http.Request) {
	var s Student
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	result, err := db.Exec("INSERT INTO students (name, age, grade) VALUES (?, ?, ?)", s.Name, s.Age, s.Grade)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	id, _ := result.LastInsertId()
	s.ID = int(id)
	json.NewEncoder(w).Encode(s)
}

func updateStudent(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	var s Student
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	_, err := db.Exec("UPDATE students SET name = ?, age = ?, grade = ? WHERE id = ?", s.Name, s.Age, s.Grade, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.ID = int(id[0])
	json.NewEncoder(w).Encode(s)
}

func deleteStudent(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	_, err := db.Exec("DELETE FROM students WHERE id = ?", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"message": "Student deleted successfully"})
}

func main() {
	connectDB()
	r := mux.NewRouter()
	r.HandleFunc("/students", getStudents).Methods("GET")
	r.HandleFunc("/students/{id}", getStudentsId).Methods("GET")
	r.HandleFunc("/students", createStudent).Methods("POST")
	r.HandleFunc("/students/{id}", updateStudent).Methods("PUT")
	r.HandleFunc("/students/{id}", deleteStudent).Methods("DELETE")

	log.Println("Server started on :8080")
	http.ListenAndServe(":8080", r)
}

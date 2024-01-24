package main

import (
	"database/sql"
	"time"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func main () {
	db, err := sql.Open("mysql", "temp_user:password@/portfolio")
	if err != nil {
		panic(err)
	}
	// Default tutorial settings, consider revisiting
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	createPortfolioTable := `
		CREATE TABLE IF NOT EXISTS companies (
			id int AUTO_INCREMENT,
			name varchar(255),
			background_color varchar(7),
			text_color varchar(7),
			header_color varchar(7),
			button_color varchar(7),
			button_text_color varchar(7),
			PRIMARY KEY (id)
		);
	`

	_, err = db.Exec(createPortfolioTable)

	if err != nil {
		panic(err)
	}

	createReasonsTable := `
		CREATE TABLE IF NOT EXISTS reasons (
			id int AUTO_INCREMENT PRIMARY KEY,
			reason varchar(255),
			company_id int,
			FOREIGN KEY (company_id) REFERENCES companies(id)
		)
	`

	_, err = db.Exec(createReasonsTable)

	if err != nil {
		panic(err)
	}

	r := mux.NewRouter()

	r.HandleFunc(`/companies/search`, SearchCompanies)

	log.Fatal(http.ListenAndServe(":8080", r))

	defer db.Close()
}

func SearchCompanies(http.ResponseWriter, *http.Request){
	log.Println("HIT!!!!")
}
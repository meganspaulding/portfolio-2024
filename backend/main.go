package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("mysql", "temp_user:password@/portfolio")
	if err != nil {
		panic(err)
	}
	// Default tutorial settings, consider revisiting
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	createPortfolioTable := `
		CREATE TABLE IF NOT EXISTS companies (
			id int AUTO_INCREMENT NOT NULL,
			name varchar(255) NOT NULL,
			background_color varchar(7) NOT NULL,
			text_color varchar(7) NOT NULL,
			header_color varchar(7) NOT NULL,
			button_color varchar(7) NOT NULL,
			button_text_color varchar(7) NOT NULL,
			PRIMARY KEY (id)
		);
	`

	_, err = db.Exec(createPortfolioTable)

	if err != nil {
		panic(err)
	}

	createReasonsTable := `
		CREATE TABLE IF NOT EXISTS reasons (
			id int AUTO_INCREMENT PRIMARY KEY NOT NULL,
			reason varchar(255) NOT NULL,
			company_id int NOT NULL,
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

type CompanySearchResult struct {
	Name string `json:"name,omitempty"`
	Id   int    `json:"id"`
}

func SearchCompanies(w http.ResponseWriter, r *http.Request) {
	selectBusinessNames := `
		SELECT name FROM companies;
	`
	results, err := db.Query(selectBusinessNames)
	if err != nil {
		panic(err)
	}

	companies := make([]CompanySearchResult, 0)

	for results.Next() {
		company := CompanySearchResult{}
		err = results.Scan(&company.Name)
		if err != nil {
			panic(err)
		}
		companies = append(companies, company)
	}

	json.NewEncoder(w).Encode(companies)
}

package main

import (
	"os"
	"html/template"
	"log"
	"net/http"

	"gotps/database"

	_ "github.com/mattn/go-sqlite3"
)

var templates = template.Must(template.ParseGlob("templates/*.html"))

// Middleware for basic authentication
func basicAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		if !ok || !checkCredentials(user, pass) {
			w.Header().Set("WWW-Authenticate", `Basic realm="gotps - Admin Portal"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	}
}

func checkCredentials(user, pass string) bool {
	adminUser := os.Getenv("ADMIN_USERNAME")
	adminPass := os.Getenv("ADMIN_PASSWORD")
	return user == adminUser && pass == adminPass
}

// Handler to process and execute SQL query
func adminHandler(writer http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		sqlQuery := req.FormValue("sqlQuery")
		if sqlQuery == "" {
			http.Error(writer, "SQL query is required", http.StatusBadRequest)
			return
		}

		rows, err := database.Query(sqlQuery)
		if err != nil {
			http.Error(writer, "Error executing query: "+err.Error(),
				http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		columns, err := rows.Columns()
		if err != nil {
			http.Error(writer, "Error retrieving columns: "+err.Error(),
				http.StatusInternalServerError)
			return
		}

		var results []map[string]interface{}
		for rows.Next() {
			cols := make([]interface{}, len(columns))
			colsPointers := make([]interface{}, len(columns))
			for i := range cols {
				colsPointers[i] = &cols[i]
			}
			if err := rows.Scan(colsPointers...); err != nil {
				log.Println("Error scanning row: ", err)
				continue
			}

			row := make(map[string]interface{})
			for i, colName := range columns {
				row[colName] = cols[i]
			}
			results = append(results, row)
		}

		// Prepare data for template
		data := struct {
			Columns []string
			Rows [][]interface{}
		}{
			Columns: columns,
			Rows: formatRows(results, columns),
		}

		if err := templates.ExecuteTemplate(writer, "admin.html", data); err != nil {
			http.Error(writer, "Template error: "+err.Error(), http.StatusInternalServerError)
		}
	} else {
		templates.ExecuteTemplate(writer, "admin.html", nil)
	}
}

func logout(writer http.ResponseWriter, req *http.Request) {
	http.SetCookie(writer, &http.Cookie{
		Name:   "session_token",
		Value:  "",
		Path:   "/",
		// Immediately expire the cookie
		MaxAge: -1, 
	})
}

func formatRows(results []map[string]interface{}, columns []string) [][]interface{} {
	rows := make([][]interface{}, 0, len(results))
	for _, row := range results {
		rowData := make([]interface{}, len(columns))
		for i, col := range columns {
			rowData[i] = row[col]
		}
		rows = append(rows, rowData)
	}
	return rows
}

func main() {
	database.Constructor(os.Getenv("DATABASE_PATH"), os.Getenv("SQL_FILEPATH"))
	defer database.Deconstructor()

	http.HandleFunc("/admin", basicAuth(adminHandler))
	http.HandleFunc("/logout", basicAuth(logout))

	var ipPort string = "0.0.0.0:"+os.Getenv("ADMIN_PORT")
	log.Printf("info: admin portal is running at http://%s\n", ipPort)
	err := http.ListenAndServe(ipPort, nil)
	if err != nil {
		log.Fatal("error: failed to start admin portal: ", err)
	}
}

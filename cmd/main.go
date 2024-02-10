package main

import (
	"database/sql"
	"embed"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

	"github.com/till/s9y-hugo/hugo"
	"github.com/till/s9y-hugo/s9y"
)

var (
	// blog post entry template
	tplFS embed.FS

	// database creds/settings
	dbUser, dbPass, dbName, dbPrefix string

	// blog url prefix
	blogUrl string

	// build vars
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

const (
	tplFilename string = "post.tmpl"
)

// poor man's guard
func init() {
	// setup templates
	tplFS = hugo.Template

	// init db
	var status bool
	dbUser, status = os.LookupEnv("DB_USER")
	if !status {
		log.Fatal(fmt.Errorf("missing DB_USER"))
	}

	dbPass, status = os.LookupEnv("DB_PASS")
	if !status {
		log.Fatal(fmt.Errorf("missing DB_PASS"))
	}

	dbName, status = os.LookupEnv("DB_NAME")
	if !status {
		log.Fatal(fmt.Errorf("missing DB_NAME"))
	}

	dbPrefix, status = os.LookupEnv("DB_TABLE_PREFIX")
	if !status {
		dbPrefix = "s9y_"
	}

	blogUrl, status = os.LookupEnv("BLOG_URL")
	if !status {
		blogUrl = "/blog"
	}

	log.Printf("Starting %s (%s) (build date: %s)", version, commit, date)
}

func main() {

	// construct a local dsn
	dsn := fmt.Sprintf(
		`%s:%s@unix(%s)/%s?allowNativePasswords=true`,
		dbUser, dbPass,
		"/var/run/mysql/mysql.sock",
		dbName,
	)

	// log.Printf("dsn: %s", dsn)

	dbConn, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}

	db := sqlx.NewDb(dbConn, "mysql")
	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")

	s := s9y.New(db, dbPrefix)

	entries, err := s.Entries(blogUrl)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Found %d entries", len(entries))

	postsDir := filepath.Join(".", "content", "posts")

	err = os.MkdirAll(postsDir, 0755)
	if err != nil {
		log.Fatalf("os.MkdirAll: %v", err)
	}
	log.Printf("Created %s directory", postsDir)

	section := hugo.New(postsDir, tplFS)

	for _, e := range entries {
		var fileName string
		if e.CategoryName.Valid {
			categoryName := strings.ReplaceAll(strings.Trim(strings.ToLower(e.CategoryName.String), " "), " ", "-")
			link := fmt.Sprintf("%s/%s", blogUrl, e.CategoryLink.String)
			err = section.Create(categoryName, e.CategoryDescription.String, link)
			if err != nil {
				log.Fatal(err)
			}

			fileName = filepath.Join(postsDir, categoryName, fmt.Sprintf("%s.md", e.URL()))
		} else {
			fileName = filepath.Join(postsDir, fmt.Sprintf("%s.md", e.URL()))
		}

		file, err := os.Create(fileName)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		log.Printf("Created: %s", fileName)

		tpl := template.Must(template.New(tplFilename).ParseFS(tplFS, tplFilename))
		err = tpl.Execute(file, e)
		if err != nil {
			log.Fatal(fmt.Errorf("tpl.Execute (%s): %v", tplFilename, err))
		}

		log.Printf("Migrated entry: %s", fileName)
	}
	log.Print("Finished!")
}

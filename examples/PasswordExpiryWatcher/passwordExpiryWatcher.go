package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/TheJumpCloud/jcapi"
)

func main() {
	// Input parameters
	var apiKey string
	var csvFile string

	// Obtain the input parameters
	flag.StringVar(&csvFile, "output", "o", "-output=<filename>")
	flag.StringVar(&apiKey, "key", "k", "-key=<API-key-value>")
	flag.Parse()

	if csvFile == "" || apiKey == "" {
		fmt.Println("Usage of ./CSVImporter:")
		fmt.Println("  -output=\"\": -output=<filename>")
		fmt.Println("  -key=\"\": -key=<API-key-value>")
		return
	}

	// Attach to JumpCloud
	// jc := jcapi.NewJCAPI(apiKey, jcapi.StdUrlBase)
	jc := jcapi.NewJCAPI(apiKey, "http://dev.local:3004/api")

	// Fetch all users who's password expires between given dates in
	userList, err := jc.GetSystemUsers(false)

	if err != nil {
		fmt.Printf("Could not read system users, err='%s'\n", err)
		return
	}

	// Setup access the CSV file specified
	file, err := os.Create(csvFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	w := csv.NewWriter(file)

	if err := w.Write([]string{"FIRSTNAME", "LASTNAME", "EMAIL", "PASSWORD EXPIRY DATE"}); err != nil {
		log.Fatalln("error writing record to csv:", err)
	}

	for _, record := range userList {
		if err := w.Write([]string{record.FirstName, record.LastName, record.Email, record.PasswordExpirationDate.String()}); err != nil {
			log.Fatalln("error writing record to csv:", err)
		}
	}
	w.Flush()

	if err := w.Error(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Finished")

	return
}

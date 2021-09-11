package main

import (
	"fmt"
	"os"
)

func main() {
	// Set an Environment Variable
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_USER", "admin")
	os.Setenv("DB_PASS", "test123")

	// Get the value of an Environment Variable
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	fmt.Printf("Host: %s, Port: %s\n", user, pass)

	// Unset an Environment Variable
	os.Unsetenv("DB_HOST")
	fmt.Printf("Try to get host: %s\n", os.Getenv("DB_HOST"))

	/*
		Get the value of an environment variable and a boolean indicating whether the
		environment variable is present or not.
	*/
	database, ok := os.LookupEnv("DB_NAME")
	if !ok {
		fmt.Println("DB_NAME is not present")
	} else {
		fmt.Printf("Database Name: %s\n", database)
	}
}

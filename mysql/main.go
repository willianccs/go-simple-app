package main

import (
	"fmt"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"log"
)

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}

// insert query generation
type FirstType struct {
    FirstParam  string
    SecondParam int
}

type SecondType struct {
    FirstParam  string
    SecondParam string
    ThirdParam  int
}

func generateInsetMethod(input interface{}) (string, error) {
    if reflect.ValueOf(q).Kind() == reflect.Struct {
        query := fmt.Sprintf("insert into %s values(", reflect.TypeOf(q).Name())
        v := reflect.ValueOf(q)
        for i := 0; i < v.NumField(); i++ {
            switch v.Field(i).Kind() {
            case reflect.Int:
                if i == 0 {
                    query = fmt.Sprintf("%s%d", query, v.Field(i).Int())
                } else {
                    query = fmt.Sprintf("%s, %d", query, v.Field(i).Int())
                }
            case reflect.String:
                if i == 0 {
                    query = fmt.Sprintf("%s\"%s\"", query, v.Field(i).String())
                } else {
                    query = fmt.Sprintf("%s, \"%s\"", query, v.Field(i).String())
                }
            default:
                fmt.Println("Unsupported type")
            }
        }
        query = fmt.Sprintf("%s)", query)
        return query, nil
    }
    return ``, QueryGenerationError{}
}

func main() {
	rootPassword := getEnv("MYSQL_ROOT_PASSWORD", "password")
	dbDatabase := getEnv("MYSQL_DATABASE", "cat")
	dbHost := getEnv("MYSQL_HOST", "127.0.0.1")
	dbPort := getEnv("MYSQL_PORT", "3306")


	fmt.Println("Open MySQL database")
	db, err := sql.Open("mysql", "root:"+rootPassword+"@tcp("+dbHost+":"+dbPort+")/"+dbDatabase+"?charset=utf8mb4&parseTime=true")

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
	    panic(err.Error()) // proper error handling instead of panic in your app
	}

	log.Printf("Successfully connected to MySQL database")
}
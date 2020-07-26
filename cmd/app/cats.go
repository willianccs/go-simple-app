package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
    "os"

	// used in mysql
	"github.com/willianccs/go-simple-app/lib/utils/mysql"
	_ "github.com/go-sql-driver/mysql"
)

// Cats struct to map a response
type Cats struct {
    ID string `json:"id"`
    Name string `json:"name"`
	Description string `json:"description"`
    Origin      string `json:"origin"`
    CountryCode string `json:"country_code"`
    Temperament string `json:"temperament"`
}

// CatsImages struct to fetch URL
type CatsImages struct {
    ID string `json:"id"`
    URL string `json:"url"`
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}

func main() {
	user := getEnv("MYSQL_ROOT_USER", "root")
	password := getEnv("MYSQL_ROOT_PASSWORD", "password")
	database := getEnv("MYSQL_DATABASE", "cats")
	host := getEnv("MYSQL_HOST", "127.0.0.1")
    port := getEnv("MYSQL_PORT", "3306")

	response, err := http.Get("https://api.thecatapi.com/v1/breeds/")
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
		os.Exit(1)
    }
    
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
    }

	//Connect to database
	db, err := mysql.NewDB(user, password, host, port, database)
	if err != nil {
		log.Fatal(err)
    }
    
	var responseObject []Cats
    json.Unmarshal([]byte(responseData), &responseObject)
    if err := json.Unmarshal([]byte(responseData), &responseObject); err != nil {
        log.Fatal(err)
    }
    
	for _, val := range responseObject {
        srcimage, err := http.Get(fmt.Sprintf("https://api.thecatapi.com/v1/images/search?breed_id=%s&limit=3", val.ID))
    	if err != nil {
	    	fmt.Printf("The HTTP request failed with error %s\n", err)
		    os.Exit(1)
        } else {
            data, _ := ioutil.ReadAll(srcimage.Body)
            var responseImg []CatsImages
            if err := json.Unmarshal([]byte(data), &responseImg); err != nil {
                log.Fatal(err)
            }
            for _, img := range responseImg {
                // INSERT INTO DB
                // prepare
                // fmt.Printf("prepare")
		        stmt, err := db.Prepare("INSERT INTO breeds_images(id,img) VALUES(?,?)")
                if err != nil {
                    log.Fatal(err)
                }

                //execute
                // fmt.Printf("execute images")
                stmt.Exec(val.ID, img.URL)
                if err != nil {
                    log.Fatal(err)
                }
                defer db.Close()
            }
        }

        // INSERT INTO DB
        fmt.Println(val)
        // prepare
		stmt, err := db.Prepare("INSERT INTO breeds(name,origin,country_code,temperament,description) VALUES(?,?,?,?,?)")
        if err != nil {
            log.Fatal(err)
        }
        
        //execute
        res, err := stmt.Exec(val.Name, val.Origin, val.CountryCode, val.Temperament, val.Description)
        if err != nil {
            log.Fatal(err)
        }
        
        id, err := res.LastInsertId()
        if err != nil {
            log.Fatal(err)
        }
        
        fmt.Println("Insert id", id)
		if err != nil {
            fmt.Print(err)
		}
        defer db.Close()

    }

    // Cats images unsing Glasses and/or Hats (category_ids=1,4)
    fmt.Println("Insert data into tables sunglasses and hats")
    for n :=0; n < 3; n++ {
        glasses, err := http.Get("https://api.thecatapi.com/v1/images/search?category_ids=4")
        if err != nil {
            fmt.Printf("The HTTP request failed with error %s\n", err)
            os.Exit(1)
            } else {
                data, _ := ioutil.ReadAll(glasses.Body)
                var responseImg []CatsImages
                if err := json.Unmarshal([]byte(data), &responseImg); err != nil {
                    log.Fatal(err)
                }
                for _, img := range responseImg {
                    // INSERT INTO DB
                    // prepare
                    // fmt.Printf("prepare")
                    stmt, err := db.Prepare("INSERT INTO cats_glasses(id,img) VALUES(?,?)")
                    if err != nil {
                        log.Fatal(err)
                    }
                    //execute
                    // fmt.Printf("execute cats glasses")
                    stmt.Exec(img.ID, img.URL)
                    if err != nil {
                        log.Fatal(err)
                    }
                    defer db.Close()

                }
            }
    }

    for n :=0; n < 3; n++ {
        hats, err := http.Get("https://api.thecatapi.com/v1/images/search?category_ids=1")
        if err != nil {
            fmt.Printf("The HTTP request failed with error %s\n", err)
            os.Exit(1)
            } else {
                data, _ := ioutil.ReadAll(hats.Body)
                var responseImg []CatsImages
                if err := json.Unmarshal([]byte(data), &responseImg); err != nil {
                    log.Fatal(err)
                }
                for _, img := range responseImg {
                    // INSERT INTO DB
                    // prepare
                    // fmt.Printf("prepare")
                    stmt, err := db.Prepare("INSERT INTO cats_hat(id,img) VALUES(?,?)")
                    if err != nil {
                        log.Fatal(err)
                    }
                    //execute
                    // fmt.Printf("execute cats hats")
                    stmt.Exec(img.ID, img.URL)
                    if err != nil {
                        log.Fatal(err)
                    }
                    defer db.Close()

                }
            }
    }
    
    db.Close()

}

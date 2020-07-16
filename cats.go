package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "os"
)

// A Response struct to map the Entire Response
type Response struct {
    Nome string `json:"name"`
    Origem string `json:"origin"`
    Temperamento string `json:"temperament"`
}

// A Cat Struct to map every cat to.
// type Cat struct {
	// Nome string `json:"name"`
    // Origem string `json:"origin"`
    // Temperamento string `json:"temperament"`
// }

// A struct to map our Cat's names
// type CatsNames struct {
    // Name string `json:"name"`
// }

func main() {
    response, err := http.Get("https://api.thecatapi.com/v1/breeds/")
    if err != nil {
        fmt.Printf("The HTTP request failed with error %s\n", err)
        os.Exit(1)
    } else {
        data, _ := ioutil.ReadAll(response.Body)
        fmt.Println(string(data))
    }

    responseData, err := ioutil.ReadAll(response.Body)
    if err != nil {
        log.Fatal(err)
    }

	var responseObject []Response
	if err := json.Unmarshal([]byte(responseData), &responseObject); err != nil {
        panic(err)
    }
    for _, val := range responseObject {
        fmt.Println("Nome: "+val.Nome)
        fmt.Println("Origem: "+val.Origem)
        fmt.Println("Temperamento: "+val.Temperamento)
	}
	

}
package main

import (
	"strconv"
    "fmt"
	"log"
	"os"
	"net/http"
	"encoding/json"
	"time"

	"github.com/gorilla/mux" 
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/olivere/elastic/v7"
	"github.com/sirupsen/logrus"
	"gopkg.in/sohlich/elogrus.v7"
)

var db *gorm.DB
var err error

func homePage(histogram *prometheus.HistogramVec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		//monitoring how long it takes to respond
		start := time.Now()
		defer r.Body.Close()
		code := 500
		handler := r.RequestURI
		method := r.Method
		
		defer func() {
			httpDuration := time.Since(start)
			scode := strconv.Itoa(code)
			histogram.WithLabelValues(handler, method, scode).Observe(httpDuration.Seconds())
		}()

		code = http.StatusBadRequest // if req is not GET
		if r.Method == "GET" {
			code = http.StatusOK

			greet := fmt.Sprintf("Welcome to the HomePage!")
			w.Write([]byte(greet))
		} else {
			w.WriteHeader(code)
		}
    	fmt.Println("Endpoint Hit: HomePage")
	}
}

// func myLogs() {
// 	// Logging
// 	os.OpenFile(logPath, os.O_RDONLY|os.O_CREATE, 0666)
// 	c := zap.NewProductionConfig()
// 	c.OutputPaths = []string{"stdout", logPath}
// 	l, err := c.Build()
// 	if err != nil {
// 		panic(err)
// 	}
// }

const (
	logPath = "./logs/go.log"
)

func handleRequests() {
	// Prometheus: Histogram to collect required metrics
	histogram := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "http",
		Name:      "request_duration_seconds",
		Help:      "The latency of the HTTP requests.",
		Buckets: []float64{1, 2, 5, 6, 10}, //defining small buckets as this app should not take more than 1 sec to respond
	}, []string{"handler", "method", "code"})

	log.Println("Starting development server at http://127.0.0.1:10000/")
    log.Println("Quit the server with CONTROL-C.")
    // creates a new instance of a mux router
    myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage(histogram))
	myRouter.HandleFunc("/all-breeds", returnAllBreeds(histogram))
	myRouter.HandleFunc("/breeds/{id}", returnInfoBreed(histogram))
	myRouter.HandleFunc("/breeds/{temperament}", returnBreedsFromTemperament(histogram))
	myRouter.Handle("/metrics", promhttp.Handler())
	myRouter.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	//Registering the defined metric with Prometheus
	prometheus.Register(histogram)

    log.Fatal(http.ListenAndServe(":10000", myRouter))
}

type Breed struct {
	ID string `json:"id"`
	Description string `json:"description"`
	Origin      string `json:"origin"`
    Temperament string `json:"temperament"`
}

func returnAllBreeds(histogram *prometheus.HistogramVec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request){
		//monitoring how long it takes to respond
		start := time.Now()
		defer r.Body.Close()
		code := 500
		handler := r.RequestURI
		method := r.Method
		
		defer func() {
			httpDuration := time.Since(start)
			scode := strconv.Itoa(code)
			histogram.WithLabelValues(handler, method, scode).Observe(httpDuration.Seconds())
		}()
			
		if r.Method == "GET" {
			code = http.StatusOK
			breeds := []Breed{}
			db.Find(&breeds)

			s := json.NewEncoder(w).Encode(breeds)
			greet := fmt.Sprintf("%s", s)
			w.Write([]byte(greet))
		} else {
			w.WriteHeader(code)
		}
		fmt.Println("Endpoint Hit: returnAllBreeds")
	}

}

func returnInfoBreed(histogram *prometheus.HistogramVec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request){
		//monitoring how long it takes to respond
		start := time.Now()
		defer r.Body.Close()
		code := 500
		handler := r.RequestURI
		method := r.Method
		
		defer func() {
			httpDuration := time.Since(start)
			scode := strconv.Itoa(code)
			histogram.WithLabelValues(handler, method, scode).Observe(httpDuration.Seconds())
		}()
			
		if r.Method == "GET" {
			code = http.StatusOK
			vars := mux.Vars(r)
			key := vars["id"]
			breeds := []Breed{}
			db.Find(&breeds)
					
			for _, breed := range breeds {
				if err == nil {
					if breed.ID == key {
						fmt.Println("Endpoint Hit: Breed No:", key)
						s := json.NewEncoder(w).Encode(breed)
						greet := fmt.Sprintf("%s", s)
						w.Write([]byte(greet))
					}
				}
			} 
		} else {
			code = http.StatusNotFound
			w.WriteHeader(code)
		}

	}
}

func returnBreedsFromTemperament(histogram *prometheus.HistogramVec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request){
		start := time.Now()
		defer r.Body.Close()
		code := 500
		handler := r.RequestURI
		method := r.Method
		
		defer func() {
			httpDuration := time.Since(start)
			scode := strconv.Itoa(code)
			histogram.WithLabelValues(handler, method, scode).Observe(httpDuration.Seconds())
		}()
			
		if r.Method == "GET" {
			code = http.StatusOK
			vars := mux.Vars(r)
			key := vars["temperament"]
			breeds := []Breed{}
			db.Where("temperament LIKE (?)", "%s", key).Select("name").Find(&breeds)
			
			for _, breed := range breeds {
				if err == nil {
					if breed.Temperament == key {
						// fmt.Println(breed)
						fmt.Println("Endpoint Hit: Temperament:", key)
						s := json.NewEncoder(w).Encode(breed)
						greet := fmt.Sprintf("%s", s)
						w.Write([]byte(greet))
					}
				}
			}
		} else {
			code = http.StatusNotFound
			w.WriteHeader(code)
		}

	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}

func main() {
	user := getEnv("MYSQL_ROOT_USER", "topcat")
	password := getEnv("MYSQL_ROOT_PASSWORD", "Zaq!@wsx34")
	database := getEnv("MYSQL_DATABASE", "cats")
	host := getEnv("MYSQL_HOST", "127.0.0.1")
	port := getEnv("MYSQL_PORT", "3306")

	db, err = gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True", user, password, host, port, database))
	if err!=nil{
		log.Println("Connection Failed to Open")
	} else { 
		log.Println("Connection Established")
	}
	defer db.Close()

	log := logrus.New()
	client, err := elastic.NewClient(elastic.SetURL("http://elasticsearch:9200"))
	if err != nil {
		log.Panic(err)
	}
	hook, err := elogrus.NewAsyncElasticHook(client, "localhost", logrus.DebugLevel, "cat-api-log")
	if err != nil {
		log.Panic(err)
	}
	log.Hooks.Add(hook)

	db.AutoMigrate(&Breed{})
    handleRequests()
	
}
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
)

var db *gorm.DB
var err error

var (
    WarningLogger *log.Logger
    InfoLogger    *log.Logger
    ErrorLogger   *log.Logger
)

func setupLog() {
    file, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
    if err != nil {
        log.Fatal(err)
    }

    InfoLogger = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
    WarningLogger = log.New(file, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
    ErrorLogger = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

const (
	logPath = "/logs/go.log"
)

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
    	InfoLogger.Println("Endpoint Hit: HomePage")
	}
}

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
	myRouter.HandleFunc("/breeds/temperament/{temperament}", returnBreedsFromTemperament(histogram))
	myRouter.HandleFunc("/breeds/origin/{country_code}", returnBreedsFromOrigin(histogram))
	myRouter.Handle("/metrics", promhttp.Handler())
	//Registering the defined metric with Prometheus
	prometheus.Register(histogram)

	er := http.ListenAndServe(":10000", myRouter)
	if er != nil {
		ErrorLogger.Fatal(er)
	}
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
		InfoLogger.Println("Endpoint Hit: returnAllBreeds")
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
						InfoLogger.Println("Endpoint Hit: Breed No:", key)
						s := json.NewEncoder(w).Encode(breed)
						greet := fmt.Sprintf("%s", s)
						w.Write([]byte(greet))
						return
					}
				}

			}
			WarningLogger.Println("Endpoint Not Found: Breed No:", key)
			w.WriteHeader(400)
		} else {
			code = 400
			WarningLogger.Println("Endpoint or Method Not Found")
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

		type breed struct {
			Name string `json:"name"`
		}

		if r.Method == "GET" {
			code = http.StatusOK
			vars := mux.Vars(r)
			key := vars["temperament"]
			breeds := []breed{}
			stmt := db.Where("temperament LIKE ?", fmt.Sprintf("%%%s%%", key)).Select("name").Find(&breeds)

			if stmt.Error != nil {
				ErrorLogger.Println(err)
			} else {
				json.NewEncoder(w).Encode(breeds)
				InfoLogger.Println("Endpoint Hit: Temperament:", key)
			}

		} else {
			code = 400
			WarningLogger.Println("Method Not Found. Use GET.")
			w.WriteHeader(code)
		}

	}
}

func returnBreedsFromOrigin(histogram *prometheus.HistogramVec) http.HandlerFunc {
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

		type breed struct {
			Name string `json:"name"`
		}

		if r.Method == "GET" {
			code = http.StatusOK
			vars := mux.Vars(r)
			key := vars["country_code"]
			breeds := []breed{}
			stmt := db.Where("country_code = ?", key).Select("name").Find(&breeds)
			if stmt.Error != nil {
				ErrorLogger.Println(err)
			} else {
				json.NewEncoder(w).Encode(breeds)
				InfoLogger.Println("Endpoint Hit: Origin:", key)
			}
		} else {
			code = 400
			WarningLogger.Println("Method Not Found. Use GET.")
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
	setupLog()
	user := getEnv("MYSQL_ROOT_USER", "topcat")
	password := getEnv("MYSQL_ROOT_PASSWORD", "Zaq!@wsx34")
	database := getEnv("MYSQL_DATABASE", "cats")
	host := getEnv("MYSQL_HOST", "127.0.0.1")
	port := getEnv("MYSQL_PORT", "3306")

	db, err = gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True", user, password, host, port, database))
	if err!=nil{
		ErrorLogger.Fatalln("Connection Failed to Open", err)
	} else {
		InfoLogger.Println("Connection Established")
	}
	defer db.Close()

	db.AutoMigrate(&Breed{})
    handleRequests()

}
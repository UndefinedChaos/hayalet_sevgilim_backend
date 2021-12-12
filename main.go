package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type myInfo struct {
	Arr []ResultsStruct `json:"results"`
}
type ResultsStruct struct {
	Distance int `json:"distance"`
}
type Response struct {
	LocationFound bool
	Message       string
}

type CounterFile struct {
	Counter int `json:"counter"`
}

var router = mux.NewRouter()

func YourHandler(w http.ResponseWriter, r *http.Request) {
	//Counter check
	counterFile, err := os.Open("counter.json")
	if err != nil {
		log.Fatal("NewRequest: ", err)
		return
	}
	defer counterFile.Close()
	var counter CounterFile
	json.NewDecoder(counterFile).Decode(&counter)
	//
	if counter.Counter < 35000 {

		lat := r.FormValue("lat")
		long := r.FormValue("long")
		lat_long := lat + "," + long
		url := "https://api.foursquare.com/v3/places/search?query=starbucks&ll=" + lat_long + "&radius=50&sort=DISTANCE&limit=1"

		req, _ := http.NewRequest("GET", url, nil)

		req.Header.Add("Accept", "application/json")
		req.Header.Add("Authorization", "fsq3IpOjjwcDiE17Itn4wFCJmSRC3TdlI7JsaiJzdh+kLV8=")

		res, err := http.DefaultClient.Do(req)
		if res == nil {
			log.Fatal("NewRequest: ", err)
			return
		}
		defer res.Body.Close()
		//Counter increasing
		counter.Counter = counter.Counter + 1
		file, _ := json.Marshal(counter)
		_ = ioutil.WriteFile("counter.json", file, 0644)
		//
		var record myInfo
		json.NewDecoder(res.Body).Decode(&record)
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		var response = Response{}
		if len(record.Arr) > 0 {
			response.LocationFound = true
			response.Message = "Starbucks Bulundu"
			js, _err := json.Marshal(response)
			if _err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Write(js)
		} else {
			response.LocationFound = false
			response.Message = "Starbucks Bulunamadı"
			js, _err := json.Marshal(response)
			if _err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Write(js)
		}
	} else {
		w.Header().Set("Content-Type", "application/json")
		var response = Response{}
		response.LocationFound = false
		response.Message = "Foursquare request limit achieved"
		js, _err := json.Marshal(response)
		if _err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(js)

	}

}

func main() {
	router.Path("/foundStarbucks").HandlerFunc(YourHandler)
	port := os.Getenv("PORT")

	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatal(err)
	}

}

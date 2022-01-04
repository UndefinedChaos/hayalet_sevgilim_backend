package main

import (
	"encoding/json"
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

var router = mux.NewRouter()

func YourHandler(w http.ResponseWriter, r *http.Request) {

	lat := r.FormValue("lat")
	long := r.FormValue("long")
	lat_long := lat + "," + long
	url := "https://api.foursquare.com/v3/places/search?query=starbucks&ll=" + lat_long + "&radius=50&sort=DISTANCE&limit=1"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "fsq3IpOjjwcDiE17Itn4wFCJmSRC3TdlI7JsaiJzdh+kLV8=")
	var response = Response{}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		response.Message = "Not Found"
	}
	defer res.Body.Close()
	if res.StatusCode >= 200 && res.StatusCode <= 299 {
		var record myInfo
		json.NewDecoder(res.Body).Decode(&record)
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		if len(record.Arr) > 0 {
			response.LocationFound = true
			response.Message = "Starbucks Bulundu"
			js, _err := json.Marshal(response)
			if _err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				response.Message = "Server Error"
			}
			w.Write(js)
		} else {
			response.LocationFound = false
			response.Message = "Starbucks Bulunamadı"
			js, _err := json.Marshal(response)
			if _err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				response.Message = "Server Error"
			}
			w.Write(js)
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
		response.LocationFound = false
		response.Message = "API limit achieved"
		js, _err := json.Marshal(response)
		if _err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			response.Message = "Server Error"
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

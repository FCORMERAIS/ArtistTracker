package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Page struct {
	Name         string
	Image 		 string
	Members      []string
	CreationDate string
	FirstAlbum   string
	Locations    string
	ConcertDates string
	Relations    string
}

type ArtistAPI struct {
	Id int `json:"id"`
	Image string `json:"image"`
	Name string `json:"name"`
	Members []string `json:"members"`
	CreationDate int `json:"creationDate"`
	FirstAlbum string `json:"firstAlbum"`
	AddressLocation string `json:"locations"`
	ConcertDatesaddress string `json:"concertDates"`
	RelationsAdress string `json:"relations"`
	Location []string
	ConcertDates []string
	Relations interface{}
}

type Location struct {
	Id int `json:"id"`
	Location []string `json:"locations"`
	Dates string `json:"dates"`
}

type Dates struct {
	Id int `json:"id"`
	Dates []string `json:"dates"`
}

type Relation struct {
	Id int `json:"id"`
	DatesLocations interface{} `json:"datesLocations"`
}

func HomePage(adress string, nbPage int) interface{} {
	fmt.Println("1. Performing Http Get...")
	var idArtist = (nbPage-1)*12 +1
	var url = ""
	var artists []ArtistAPI
	var oneartist ArtistAPI
	fmt.Println("1. Performing Http Get...")
	fmt.Println("2. Le serveur est lancé sur le port 3000")
	for idArtist != nbPage*12 + 1 {
		url = "/"+strconv.Itoa(idArtist)
		resp, err := http.Get(adress+url)
		if err != nil {
			log.Fatalln(err)
		}
		defer resp.Body.Close()
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal(bodyBytes, &oneartist)
		idArtist = oneartist.Id 
		if idArtist == 0 {break}
		oneartist.Location = location(oneartist.AddressLocation)
		oneartist.ConcertDates = concertdate(oneartist.ConcertDatesaddress)
		oneartist.Relations = relation(oneartist.RelationsAdress)
		fmt.Println(oneartist.Relations)
		artists = append(artists,oneartist)
		idArtist++
	}
	return artists
}

func relation(adress string) interface{} {
	var relation Relation
	resp, err := http.Get(adress)
	if err != nil {
		log.Fatalln(err)
	}
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(bodyBytes, &relation)
	return relation.DatesLocations
}

func concertdate(adress string) []string{
	var dates Dates
	resp, err := http.Get(adress)
	if err != nil {
		log.Fatalln(err)
	}
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(bodyBytes, &dates)
	for i := 0 ; i< len(dates.Dates);i++ {
		dates.Dates[i] = gooddate(dates.Dates[i])
	}
	return dates.Dates
}
func gooddate(mois string) string {
	var date string
	tempo := strings.Split(mois, "-")
	if tempo[1] == "01" {date = strings.Replace(tempo[1],"01","Janvier", -1)
	}else if tempo[1] == "02" {date = strings.Replace(tempo[1],"02","Fevrier", -1)
	}else if tempo[1] == "03" {date = strings.Replace(tempo[1],"03","Mars", -1)
	}else if tempo[1] == "04" {date = strings.Replace(tempo[1],"04","Avril", -1)
	}else if tempo[1] == "05" {date = strings.Replace(tempo[1],"05","Mai", -1)
	}else if tempo[1] == "06" {date = strings.Replace(tempo[1],"06","Juin", -1)
	}else if tempo[1] == "07" {date = strings.Replace(tempo[1],"07","Juillet", -1)
	}else if tempo[1] == "08" {date = strings.Replace(tempo[1],"08","Aout", -1)
	}else if tempo[1] == "09" {date = strings.Replace(tempo[1],"09","Septembre", -1)
	}else if tempo[1] == "10" {date = strings.Replace(tempo[1],"10","Octobre", -1)
	}else if tempo[1] == "11" {date = strings.Replace(tempo[1],"11","Novembre", -1)
	}else if tempo[1] == "12" {date = strings.Replace(tempo[1],"12","Decembre", -1)}
	tempo[1] = date
	date = strings.Join(tempo, " ")
	date = strings.Replace(date,"*","",-1)
	return date
}

func location(adress string) []string {
	var locations Location
	var location string
	resp, err := http.Get(adress)
	if err != nil {
		log.Fatalln(err)
	}
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(bodyBytes, &locations)
	for i := 0 ; i< len(locations.Location);i++ {
		location = locations.Location[i]
		location = strings.Replace(location, "_", " ", -1)
		location = strings.Replace(location, "-", " ", -1)
		locations.Location[i] = location
	}
	return locations.Location
}

func main() {
	lien := "https://groupietrackers.herokuapp.com/api"
	fileServer := http.FileServer(http.Dir("assets")) //Envoie des fichiers aux serveurs (CSS, sons, images)
	http.Handle("/assets/", http.StripPrefix("/assets/", fileServer))
	// affiche l'html
	tmpl, err := template.ParseFiles("./assets/navPage.gohtml")
	if err != nil {
	}
	page := 1

	http.HandleFunc("/Groupie-tracker", func(w http.ResponseWriter, r *http.Request) {
		data := HomePage(lien+"/artists", page)
		tmpl.ExecuteTemplate(w, "index", data)
	})

	http.HandleFunc("/Groupie-tracker/PageSuivante", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			tmpl, err = template.ParseFiles("./assets/navPage.gohtml")
		}
		page+=1
		data := HomePage(lien+"/artists",page)
		tmpl.ExecuteTemplate(w, "index", data)
	})

	http.HandleFunc("/Groupie-tracker/PagePrecedente", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			tmpl, err = template.ParseFiles("./assets/navPage.gohtml")
		}
		page-=1
		data := HomePage(lien+"/artists",page)
		tmpl.ExecuteTemplate(w, "index", data)
	})

	fmt.Println("le serveur est en cours d'éxécution a l'adresse http://localhost:3000/Groupie-tracker")
	http.ListenAndServe("localhost:3000", nil) //lancement du serveur
}

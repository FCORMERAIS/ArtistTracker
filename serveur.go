package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

//on Importe toute les bibliothèques que l'on a besoin

func main() {
	fileServer := http.FileServer(http.Dir("static/")) //Envoie des fichiers aux serveurs (CSS, sons, images)
	http.Handle("/static/", http.StripPrefix("/static/", fileServer))
	// affiche l'html
	http.HandleFunc("/", error404)
	http.HandleFunc("/Groupie-tracker", groupieTracker)
	http.HandleFunc("/Groupie-tracker/listartist", listartist)
	http.HandleFunc("/Groupie-tracker/nbArtist", nbArtist)
	http.HandleFunc("/Groupie-tracker/artist", artist)
	http.HandleFunc("/Groupie-tracker/Recherche", rechercher)
	http.HandleFunc("/Groupie-tracker/cart", carte)
	fmt.Println("le serveur est en cours d'éxécution a l'adresse http://localhost:3000/Groupie-tracker")
	http.ListenAndServe("localhost:3000", nil) //lancement du serveur
}

//#####################################################################################################################################//

func error404(w http.ResponseWriter, r *http.Request) {
	var data interface{}
	tmpl, err := template.ParseFiles("./templates/Error404.html") // utilisation du fichier navPage.gohtml pour le template
	if err != nil {
		print(err)
	}
	tmpl.ExecuteTemplate(w, "error404", data)
}

//########################################################################################################################################//

func rechercher(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		recherche := r.FormValue("recherche")
		data, errr := rechercheFind("https://groupietrackers.herokuapp.com/api/artists", strings.Split(strings.ToUpper(recherche), "")) // récupération des artiste pour les artists recherché
		if errr != nil {
			print(errr)
		}
		tmpl, err := template.ParseFiles("./templates/navbar.html", "./templates/footer.html", "./templates/pagelistartists.html", "./templates/listartist.html") // utilisation du fichier navPage.gohtml pour le template
		if err != nil {
			print(err)
		}
		tmpl.ExecuteTemplate(w, "listartists", data)
	}
}

//##########################################################################################################################################//

func groupieTracker(w http.ResponseWriter, r *http.Request) { //récupération des donnée a envoyer sur la page html
	tmpl, err := template.ParseFiles("./templates/home.html", "./templates/navbar.html", "./templates/footer.html", "./templates/pageaccueil.html", "./templates/pagelistartists.html") // utilisation du fichier navPage.gohtml pour le template
	if err != nil {
		fmt.Println(err)
		tmpl, err = template.ParseFiles("./templates/Error500.gohtml") //utilisation du fichier Error500.gohtml pour le template
		print(err)
	}
	tmpl.ExecuteTemplate(w, "home", "") //exécution du template
}

//#######################################################################################################################################//

func artist(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		id_artiste := r.FormValue("id")
		data, err := clicked(id_artiste) //récupération des donnée a envoyer sur la page html
		if err != nil {
			fmt.Println(err, "UWU")
		}
		tmpl, err := template.ParseFiles("./templates/artist.html", "./templates/navbar.html", "./templates/footer.html", "./templates/pageartist.html") // utilisation du fichier navPage.gohtml pour le template
		if err != nil {
			fmt.Println(err, "UWU")
			tmpl, err = template.ParseFiles("./templates/Error500.gohtml") //utilisation du fichier Error500.gohtml pour le template
			print(err)
		}
		tmpl.ExecuteTemplate(w, "artist", data) //exécution du template
	}
}

//#########################################################################################################################################//

func nbArtist(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		page := 1
		lien := "https://groupietrackers.herokuapp.com/api"
		nbArtist, errr := strconv.Atoi(r.FormValue("Artists"))
		tmpl, err := template.ParseFiles("./templates/navbar.html", "./templates/footer.html", "./templates/pagelistartists.html", "./templates/listartist.html") // utilisation du fichier navPage.gohtml pour le template
		if errr != nil || err != nil {
			fmt.Println(err)
		}
		function := r.FormValue("function")
		data, err := ArtistPage(lien+"/artists", page, nbArtist, function) //récupération des donnée a envoyer sur la page html
		if err != nil {
			fmt.Println(err)
		}
		tmpl.ExecuteTemplate(w, "listartists", data) //exécution du template
	}
}

//####################################################################################################################################

func carte(w http.ResponseWriter, r *http.Request) {
	var Page Carte
	lien := "https:www.google.com/maps/embed/v1/place?key=AIzaSyAXXPpGp3CYZDcUSiE2YRlNID4ybzoZa7o&q="
	tmpl, err := template.ParseFiles("./templates/navbar.html", "./templates/footer.html", "./templates/pagelistartists.html", "./templates/listartist.html", "./templates/pagecart.html")
	if err != nil {
	}
	value := r.FormValue("carte")
	if value == "" {
		lien = "https:www.google.com/maps/embed/v1/place?key=AIzaSyAXXPpGp3CYZDcUSiE2YRlNID4ybzoZa7o&q=Paris"
	}
	Page.Location = lieux("https://groupietrackers.herokuapp.com/api/locations/")
	Page.Valeur = lien + value
	tmpl.ExecuteTemplate(w, "pagecart", Page)
}

//##############################################################################################################################//

func listartist(w http.ResponseWriter, r *http.Request) {
	lien := "https://groupietrackers.herokuapp.com/api"
	page := 1
	var data interface{}
	nbArtist := 12
	tmpl, err := template.ParseFiles("./templates/navbar.html", "./templates/footer.html", "./templates/pagelistartists.html", "./templates/listartist.html") // utilisation du fichier navPage.gohtml pour le template
	if err != nil {
		fmt.Println(err, "/")
		tmpl, err = template.ParseFiles("./assets/Error500.gohtml") //utilisation du fichier Error500.gohtml pour le template
		if err != nil {
			fmt.Println(err)
		}
	}
	if r.FormValue("FilterAlpha") == "TRUE" { // récupération de la variable qui nous indique quelle filtre on utilise
		nbArtist, err = strconv.Atoi(r.FormValue("Artists")) // récupération du nombre d'artists à afficher
		if err != nil {
			fmt.Println(err)
		}
		function := "FilterAlpha"
		data, err = FilterAlpha(lien+"/artists", page, nbArtist, function) //  récupération des artists en fonction du filtre
		if err != nil {
			fmt.Println(err)
		}
	} else if r.FormValue("FilterDate") == "TRUE" { // récupération de la variable qui nous indique quelle filtre on utilise
		nbArtist, err = strconv.Atoi(r.FormValue("Artists")) // récupération du nombre d'artists à afficher
		if err != nil {
			fmt.Println(err)
		}
		function := "FilterDate"
		data, err = FilterDate(lien+"/artists", page, nbArtist, function) //  récupération des artists en fonction du filtre
		if err != nil {
			fmt.Println(err)
		}
	} else if r.FormValue("FilterSolo") == "TRUE" { // récupération de la variable qui nous indique quelle filtre on utilise
		nbArtist, err = strconv.Atoi(r.FormValue("Artists")) // récupération du nombre d'artists à afficher
		fmt.Println(nbArtist)
		if err != nil {
			fmt.Println(err)
		}
		function := "FilterSolo"
		data, err = FilterSolo(lien+"/artists", page, nbArtist, function) //  récupération des artists en fonction du filtre
		if err != nil {
			fmt.Println(err)
		}
	} else if r.FormValue("FilterGroups") == "TRUE" { // récupération de la variable qui nous indique quelle filtre on utilise
		nbArtist, err = strconv.Atoi(r.FormValue("Artists")) // récupération du nombre d'artists à afficher
		fmt.Println(nbArtist)
		function := "FilterGroups"
		if err != nil {
			fmt.Println(err)
		}
		data, err = FilterGroups(lien+"/artists", page, nbArtist, function) //  récupération des artists en fonction du filtre
		if err != nil {
			fmt.Println(err)
		}
	} else if r.FormValue("pageSuivante") == "TRUE" { // avancé a la page suivante en gardant les filtre si utilisé grâce a la variable function
		nbArtist, err = strconv.Atoi(r.FormValue("Artists")) // récupération du nombre d'artists à afficher
		if err != nil {
			fmt.Println(err)
		}
		page, err = strconv.Atoi(r.FormValue("page"))
		if err != nil {
			fmt.Println(err)
		}
		page += 1
		if r.FormValue("function") == "FilterSolo" { // récupération de la variable qui nous indique quelle filtre on utilise
			function := r.FormValue("function")
			data, err = FilterSolo(lien+"/artists", page, nbArtist, function) //  récupération des artists en fonction du filtre
			if err != nil {
				fmt.Println(err)
			}
		} else if r.FormValue("function") == "FilterGroups" { // récupération de la variable qui nous indique quelle filtre on utilise
			function := r.FormValue("function")
			data, err = FilterGroups(lien+"/artists", page, nbArtist, function) //  récupération des artists en fonction du filtre
			if err != nil {
				fmt.Println(err)
			}
		} else if r.FormValue("function") == "FilterAlpha" { // récupération de la variable qui nous indique quelle filtre on utilise
			function := r.FormValue("function")
			data, err = FilterAlpha(lien+"/artists", page, nbArtist, function) //  récupération des artists en fonction du filtre
			if err != nil {
				fmt.Println(err)
			}
		} else if r.FormValue("function") == "FilterDate" { // récupération de la variable qui nous indique quelle filtre on utilise
			function := r.FormValue("function")
			data, err = FilterDate(lien+"/artists", page, nbArtist, function) //  récupération des artists en fonction du filtre
			if err != nil {
				fmt.Println(err)
			}
		} else {
			function := "normal"
			data, err = ArtistPage(lien+"/artists", page, nbArtist, function) //  récupération des artists en fonction du filtre
		}
		if err != nil {
			fmt.Println(err)
		}
	} else if r.FormValue("pagePrecedente") == "TRUE" { // avancé a la page précédente en gardant les filtre si utilisé grâce a la variable function
		nbArtist, err = strconv.Atoi(r.FormValue("Artists")) // récupération du nombre d'artists à afficher
		if err != nil {
			fmt.Println(err)
		}
		page, err = strconv.Atoi(r.FormValue("page"))
		if err != nil {
			fmt.Println(err)
		}
		page -= 1
		if r.FormValue("function") == "FilterSolo" { // récupération de la variable qui nous indique quelle filtre on utilise
			function := r.FormValue("function")
			data, err = FilterSolo(lien+"/artists", page, nbArtist, function) //  récupération des artists en fonction du filtre
			if err != nil {
				fmt.Println(err)
			}
		} else if r.FormValue("function") == "FilterGroups" { // récupération de la variable qui nous indique quelle filtre on utilise
			function := r.FormValue("function")
			data, err = FilterGroups(lien+"/artists", page, nbArtist, function) //  récupération des artists en fonction du filtre
			if err != nil {
				fmt.Println(err)
			}
		} else if r.FormValue("function") == "FilterAlpha" { // récupération de la variable qui nous indique quelle filtre on utilise
			function := r.FormValue("function")
			data, err = FilterAlpha(lien+"/artists", page, nbArtist, function) //  récupération des artists en fonction du filtre
			if err != nil {
				fmt.Println(err)
			}
		} else if r.FormValue("function") == "FilterDate" { // récupération de la variable qui nous indique quelle filtre on utilise
			function := r.FormValue("function")
			data, err = FilterDate(lien+"/artists", page, nbArtist, function) //  récupération des artists en fonction du filtre
			if err != nil {
				fmt.Println(err)
			}
		} else {
			function := "normal"
			data, err = ArtistPage(lien+"/artists", page, nbArtist, function) //  récupération des artists en fonction du filtre
		}
		if err != nil {
			fmt.Println(err)
		}
	} else {
		function := "normal"
		data, err = ArtistPage(lien+"/artists", page, nbArtist, function) //  récupération des artists en fonction du filtre
	}
	if err != nil {
		fmt.Println(err, "/")
		tmpl, err = template.ParseFiles("./assets/Error500.gohtml") //utilisation du fichier Error500.gohtml pour le template
		fmt.Println(err)
	} //récupération des donnée a envoyer sur la page html
	tmpl.ExecuteTemplate(w, "listartists", data) //exécution du template
}

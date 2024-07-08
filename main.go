package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/melbahja/goph"
	"github.com/webguerilla/ftps"
	"github.com/xuri/excelize/v2"
)

var (
	id          int
	name        string
	description string
	price       int
)

func main() {
	initialisation()

	var choice string

	dbInfos := "root:@tcp(127.0.0.1:3306)/tp"

	db, err := sql.Open("mysql", dbInfos)

	if err != nil {
		panic(err)
	}

	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	intro(db)

	for {
		fmt.Println("")
		fmt.Println("==========================================================================")
		fmt.Println("==========================================================================")
		fmt.Println("")
		fmt.Println("Quelle tâche voulez-vous effectuer? \nEntrez le numéro correspondant à votre option:")
		fmt.Println("")
		fmt.Println("[1] Ajouter un produit")
		fmt.Println("[2] Afficher la liste des produits")
		fmt.Println("[3] Modifier un produit")
		fmt.Println("[4] Supprimer un produit")
		fmt.Println("[5] Exporter les informations produits dans un fichier Excel (en .xlsx)")
		fmt.Println("[6] Lancer un serveur http avec une page web")
		fmt.Println("[7] Se connecter à une VM en ssh")
		fmt.Println("[8] Se connecter à un serveur FTP")
		fmt.Println("[9] Quitter")
		fmt.Println("")
		time.Sleep(time.Millisecond * 500)
		fmt.Println("==========================================================================")
		fmt.Println("==========================================================================")
		fmt.Println("")

		fmt.Print(">  ")
		fmt.Scan(&choice)

		switch choice {
		case "1":
			insertProduct(db)
		case "2":
			selectProducts(db)
		case "3":
			updateProduct(db)
		case "4":
			deleteProduct(db)
		case "5":
			export2Excel(db)
		case "6":
			runServer()
		case "7":
			ssh2VM()
		case "8":
			connect2FTPServer()
		case "9":
			fmt.Println("A bientôt!")
			return
		default:
			fmt.Print("Erreur! Veuillez entrer une option valide: ")

		}
	}
}

func initialisation() {
	fmt.Println("")
	fmt.Println("")
	fmt.Println("_____ ____                ")
	fmt.Println("|_   _|  _ \\    __ _  ___  ")
	fmt.Println("  | | | |_) |  / _` |/ _ \\ ")
	fmt.Println("  | | |  __/  | (_| | (_) |")
	fmt.Println("  |_| |_|      \\__, |\\___/ ")
	fmt.Println("               |___/       ")

	fmt.Println("")
	time.Sleep(time.Millisecond * 500)
	fmt.Println("===================================")
	time.Sleep(time.Millisecond * 500)
	fmt.Println("===================================")
	time.Sleep(time.Millisecond * 500)
	fmt.Println("===================================")
}

func intro(db *sql.DB) {
	/* Function allowing prior deletion of table 'products'
	then its creation. */
	query := "USE tp;"

	_, err := db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}

	// query = "DROP TABLE IF EXISTS products;"

	// _, err = db.Exec(query)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	query = `CREATE TABLE IF NOT EXISTS products (    id INT AUTO_INCREMENT,    name TEXT NOT NULL,    description TEXT NOT NULL,    price INT,    PRIMARY KEY (id));`
	_, err = db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
}

func insertProduct(db *sql.DB) {

	var temp string

	fmt.Print("Veuillez entrer un nom: ")
	fmt.Scan(&name)
	fmt.Print("Veuillez entrer une description: ")
	fmt.Scan(&description)
	fmt.Print("Veuillez entrer un prix: ")
	fmt.Scan(&temp)

	price, err := strconv.Atoi(temp)
	if err != nil {
		fmt.Println("Erreur, veuillez renseigner un entier pour le prix !")
	}

	query := "INSERT INTO products (name, description, price) VALUES (?, ?, ?)"

	_, err = db.Exec(query, name, description, price)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Enregistrement réussi !")
}

func selectProducts(db *sql.DB) {
	query := "SELECT * FROM products"

	rows, err := db.Query(query)

	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	fmt.Println("")
	fmt.Println("LISTE DE PRODUITS: ")
	fmt.Println("============================================================")
	fmt.Println("============================================================")

	for rows.Next() {

		err = rows.Scan(&id, &name, &description, &price)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("ID : %d, Name : %s, Description : %s, Price : %d \n", id, name, description, price)

	}
	fmt.Print("")
	fmt.Println("============================================================")
	fmt.Println("============================================================")
}

func updateProduct(db *sql.DB) {
	var temp string

	fmt.Println("Renseignez l'Id du produit à modifier :")
	fmt.Scan(&temp)

	id, err := strconv.Atoi(temp)
	if err != nil {
		fmt.Println("Erreur : l'id doit être un entier ! ")
	}

	fmt.Print("Veuillez renseigner le nouveau nom: ")
	fmt.Scan(&name)
	fmt.Print("Veuillez renseigner la nouvelle description: ")
	fmt.Scan(&description)
	fmt.Print("Veuillez renseigner le nouveau prix: ")
	fmt.Scan(&temp)

	fmt.Println("")

	price, err = strconv.Atoi(temp)
	if err != nil {
		fmt.Println("Erreur: le prix doit être un entier.")
	}

	query := "UPDATE products SET name = ?, description = ?, price = ? WHERE id = ?"

	_, err = db.Exec(query, name, description, price, id)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Modification réussie!")

}

func deleteProduct(db *sql.DB) {
	var temp string

	fmt.Print("Entrez l'id du produit à supprimer: ")
	fmt.Scan(&temp)

	fmt.Println("")

	id, err := strconv.Atoi(temp)
	if err != nil {
		fmt.Println("Erreur: l'id doit être un entier.")
	}

	query := "DELETE FROM products WHERE id = ?"

	_, err = db.Exec(query, id)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Suppression réussie!")
}

func export2Excel(db *sql.DB) {
	query := `SELECT *
	FROM products`

	results, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}

	defer results.Close()

	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	index, err := f.NewSheet("Sheet1")
	if err != nil {
		fmt.Println(err)
		return
	}

	f.SetCellValue("Sheet1", "A1", "ID")
	f.SetCellValue("Sheet1", "B1", "Name")
	f.SetCellValue("Sheet1", "C1", "Description")
	f.SetCellValue("Sheet1", "D1", "Price")

	f.SetActiveSheet(index)

	i := 2
	for results.Next() {
		err = results.Scan(&id, &name, &description, &price)
		if err != nil {
			log.Fatal(err)
		}
		f.SetCellValue("Sheet1", "A"+strconv.Itoa(i), id)
		f.SetCellValue("Sheet1", "B"+strconv.Itoa(i), name)
		f.SetCellValue("Sheet1", "C"+strconv.Itoa(i), description)
		f.SetCellValue("Sheet1", "D"+strconv.Itoa(i), price)
		fmt.Printf("ID : %d, Name : %s, Description : %s, Price : %d \n", id, name, description, price)
		i += 1
	}

	if err := f.SaveAs("products.xlsx"); err != nil {
		fmt.Println(err)
	}

	fmt.Println("Données exportées dans un fichier nommé 'products.xlsx'")
}

func runServer() {
	http.HandleFunc("/", writeIntoServer)
	http.HandleFunc("/submit", methodTraitements)
	fmt.Println("Démarrage du serveur sur le port 4242")

	err := http.ListenAndServe(":4242", nil)
	if err != nil {
		fmt.Println("Erreur serveur: ", err)
	}
}

func methodTraitements(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:

		name := r.FormValue("name")
		email := r.FormValue("email")

		response := fmt.Sprintf("Données du formulaire: \nNom %s, Email: %s", name, email)

		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(response))

	case http.MethodGet:

		search := r.FormValue("search")

		response := fmt.Sprintf("Données du formulaire GET: \n%s", search)

		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(response))

	default:
		http.Error(w, "Méthode invalide", http.StatusInternalServerError)
	}

}

func writeIntoServer(w http.ResponseWriter, r *http.Request) {
	htmlFile := "index.html"

	data, err := os.ReadFile(htmlFile)

	if err != nil {
		http.Error(w, "Fichier "+htmlFile+" introuvable!", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.Write(data)
}

func ssh2VM() {
	var username, address, password string
	fmt.Println("")
	fmt.Println("Bienvenue dans le module ConnectViaSSH2VM!")
	fmt.Println("")

	fmt.Print("Veuillez entrer le nom d'utilisateur: ")
	fmt.Scan(&username)
	fmt.Print("Veuillez entrer l'adresse IP de destination: ")
	fmt.Scan(&address)
	fmt.Print("Veuillez entrer le mot de passe: ")
	fmt.Scan(&password)
	time.Sleep(time.Millisecond * 1000)
	fmt.Println("")
	fmt.Println("Connexion...")

	client, err := goph.New(username, address, goph.Password(password))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connecté!")

	time.Sleep(time.Millisecond * 1000)
	fmt.Println("Lancement de la commande 'bash -c whoami'...")
	fmt.Println("")
	out, err := client.Run("bash -c 'whoami'")
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(time.Millisecond * 1000)
	fmt.Println("Reponse: ", string(out))

	time.Sleep(time.Millisecond * 1000)
	fmt.Println("")
	fmt.Println("Déconnexion...")
	fmt.Println("")
	// Close client net connection
	defer client.Close()
}

func connect2FTPServer() {
	var option, address, port, username, password string

	ftps := new(ftps.FTPS)
	ftps.TLSConfig.InsecureSkipVerify = true

	fmt.Println("")
	fmt.Println("Bienvenue dans le module Connect2FTPServer!")
	fmt.Println("")
	fmt.Println("Veuillez choisir entre les deux options:")
	fmt.Println("[1] Connexion manuelle à un serveur FTP")
	fmt.Println("[2] Connexion automatique à un serveur FTP par défaut")

	fmt.Println("")
	fmt.Print("> ")
	fmt.Scan(&option)

	if option == "1" {
		fmt.Println("")
		fmt.Print("Veuillez entrer l'adresse de destination: ")
		fmt.Scan(&address)
		fmt.Println("")
		fmt.Print("Veuillez entrer le port: ")
		fmt.Scan(&port)
		fmt.Println("")
		fmt.Print("Veuillez entrer votre identifiant: ")
		fmt.Scan(&username)
		fmt.Println("")
		fmt.Print("Veuillez entrer votre mot de passe: ")
		fmt.Scan(&password)
		fmt.Println("")

		fmt.Println("Connexion à l'adresse '" + username + ":" + password + "@" + address + port + "'...")
		fmt.Println("")

		nbPort, err := strconv.Atoi(port)
		if err != nil {
			fmt.Println("Erreur, veuillez renseigner un entier pour le port!")
		}

		err = ftps.Connect(address, nbPort)

		if err != nil {
			panic(err)
		}

		err = ftps.Login(username, password)
		if err != nil {
			panic(err)
		}

	} else {
		fmt.Println("Connexion au serveur FTP 'test.rebex.net' au port 21...")
		time.Sleep(time.Millisecond * 500)
		fmt.Println("")
		err := ftps.Connect("test.rebex.net", 21)
		if err != nil {
			panic(err)
		}

		err = ftps.Login("demo", "password")
		if err != nil {
			panic(err)
		}

		time.Sleep(time.Millisecond * 500)
		fmt.Println("Connecté au serveur FTP!")
		time.Sleep(time.Millisecond * 500)
		fmt.Println("")

		directory, err := ftps.PrintWorkingDirectory()
		if err != nil {
			panic(err)
		}

		time.Sleep(time.Millisecond * 500)
		fmt.Printf("Current working directory %s", directory)
	}
	fmt.Println("")
	fmt.Println("")
	time.Sleep(time.Millisecond * 500)
	fmt.Println("Déconnexion du serveur...")

	err := ftps.Quit()
	if err != nil {
		panic(err)
	}

	time.Sleep(time.Millisecond * 500)
	fmt.Println("")
	fmt.Println("Déconnecté!")
}

package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strings"
	"time"
)

func main() {

	listener, err := net.Listen("tcp", ":80")

	if err != nil {
		fmt.Println(err)
		return
	}

	for {

		// ATTENTE DE CONNEXION TCP
		connexion, err := listener.Accept()

		if err == nil {

			FileToSend := os.Args[1] // Argument 1 = Dossier du site

			// LECTURE DE LA REQUEST
			scanner := bufio.NewScanner(connexion)
			for scanner.Scan() {

				line := scanner.Text()

				if strings.Contains(line, "HTTP/1.1") {
					lineSplit := strings.Split(line, " ")

					// Récupération du fichier a charger ou index.html
					if lineSplit[1] != "/" {
						FileToSend += lineSplit[1]
					} else {
						FileToSend += "/index.html"
					}
				}

				// Si fin du HEADER on quitte la boucle de lecture
				if line == "" {
					break
				}
			}

			// GESTION DE LA REPONSE
			header := ""
			file := strings.Split(FileToSend, "?")[0]

			//LOG
			fmt.Println(connexion.RemoteAddr().String(), "-->", file)

			//RECHERCHE DU FICHIER A ENVOYER
			data, errFile := ioutil.ReadFile(file)
			if errFile != nil {
				data = []byte("")

				// HEADER SI FICHIER NON TROUVE
				header = "HTTP/1.1 404 Not Found\n"

			} else {

				// HEADER SI FICHIER TROUVE
				header = "HTTP/1.1 200 OK\n"

			}

			//HEADER DE LA REPONSE HTTP
			header += "Date: " + (time.Now()).Format(time.RFC1123) + "\n"
			header += "Server: SimpleGoHTTP\n"
			header += "Connection: Closed\n\n"

			//AJOUT DES DONNEES
			header += string(data)

			//ON ENVOIT LES DONNEES
			wr := bufio.NewWriter(connexion)
			wr.WriteString(header)
			err := wr.Flush()

			//ON FERME LA CONNECTION
			connexion.Close()

			if err != nil {
				fmt.Println("/!\\ Erreur serveur --> Message non envoyé")
			}
		}
	}
}

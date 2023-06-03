package main

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"os"
	"time"
)

type Server struct {
	ServerName string
	ServerURL  string
	LoadTime   float64
	Status     int
	FailDate   string
}

func createServerList(serverlist *os.File) []Server {
	csvReader := csv.NewReader(serverlist)
	data, err := csvReader.ReadAll()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	var servers []Server
	for i, line := range data {
		if i > 0 {
			server := Server{
				ServerName: line[0],
				ServerURL:  line[1],
			}
			servers = append(servers, server)
		}
	}
	return servers
}

func checkServer(servers []Server) []Server {
	var downServers []Server
	now := time.Now()

	for _, server := range servers {
		initialTime := time.Now()
		get, err := http.Get(server.ServerURL)
		if err != nil {
			fmt.Printf("Server [%s] is down [%s]\n", server.ServerName, err.Error())
			server.Status = 0
			server.FailDate = now.Format("02/01/2006 15:04:05")
			downServers = append(downServers, server)
			continue
		}
		server.Status = get.StatusCode

		if server.Status != 200 {
			server.FailDate = now.Format("02/01/2006 15:04:05")
			downServers = append(downServers, server)
		}

		server.LoadTime = time.Since(initialTime).Seconds()

		fmt.Printf("Name: [%s] Url: [%s] Status Code: [%d] Load time: [%f seconds]\n", server.ServerName, server.ServerURL, server.Status, server.LoadTime)
	}
	return downServers
}

func openFiles(serverListFile, downtimeFile string) (*os.File, *os.File) {
	serverList, err := os.OpenFile(serverListFile, os.O_RDONLY, 0666)
	if err != nil {
		fmt.Println("Ocorreu um erro ao abrir o arquivo: ", err)
		os.Exit(1)
	}
	downtimeList, err := os.OpenFile(downtimeFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println("Ocorreu um erro ao abrir o arquivo: ", err)
		os.Exit(1)
	}
	return serverList, downtimeList
}

func generateDowntimeFile(downtimeList *os.File, downServers []Server) {
	csvWriter := csv.NewWriter(downtimeList)
	for _, servidor := range downServers {
		line := []string{servidor.ServerName, servidor.ServerURL, string(rune(servidor.Status)), servidor.FailDate, string(rune(servidor.LoadTime))}
		csvWriter.Write(line)
	}
	csvWriter.Flush()
}

func main() {
	serverList, downtimeList := openFiles(os.Args[1], os.Args[2])
	defer serverList.Close()
	defer downtimeList.Close()
	servers := createServerList(serverList)
	downServers := checkServer(servers)
	generateDowntimeFile(downtimeList, downServers)
}

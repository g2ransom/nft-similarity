package main

import (
	"bytes"
	// "encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"log"
	"os"
	"time"
	"github.com/joho/godotenv"
	"github.com/g2ransom/nft-similarity/subgraph/query"
)

// add retries
// use goroutines to append the csv results to the main goroutine, which will write all results to csv file

// var backoffSchedule = []time.Duration{
// 	1 * time.Second,
// 	3 * time.Second,
// 	10 * time.Second,
// }

// type Client struct {
// 	*http.Client
// }

func NewSlice(start, end, step int) []int {
	if step <= 0 || end < start {
		return []int{}
	}
	s := make([]int, 0, 1+(end-start)/step)
	for start <= end {
		s = append(s, start)
		start += step
	}
	return s
}

// func (c *Client) DoWithRetries(request *http.Request) (response *http.Response, err error) {
// 	for _, backoff := range backoffSchedule {
// 		response, err = c.Do(request)		
// 		if err == nil {
// 			break
// 		}

// 		log.Printf("Request error: %+v\n", err)
// 		log.Printf("Retrying in %v\n", backoff)
// 		time.Sleep(backoff)
// 	}	
// 	if err != nil {
// 		return nil, err
// 	}

// 	return response, nil
// }

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file: ", err)
	}
	graphApiKey := os.Getenv("GRAPH_API_KEY")
	graphUrl := fmt.Sprintf("https://gateway.thegraph.com/api/%s/subgraphs/id/AVZ1dGwmRGKsbDAbwvxNmXzeEkD48voB3LfGqj5w7FUS", graphApiKey)
	client := http.Client{Timeout: time.Second * 60}
	records := [][]string{
		{"contract_id", "name", "symbol", "identifier", "owner", "uri"},
	}

	// step := 3600*24
	// start := 1633284771
	// end := 1589587200
	// end := start + (step*1)
	// timeRange := NewSlice(start, end, step)
	// fmt.Printf("len timeRange: %d", len(timeRange))
	// resultChan := make(chan [][]string)

	// for i := 0; i < len(timeRange)-1; i++ {
	// log.Printf("query [start]: %d; [end]: %d\n", timeRange[0], timeRange[1])
	queryStatement := query.EIP721Query(10, 1633284771, 1667671971, 10)
	jsonQuery, err := json.Marshal(queryStatement)
	if err != nil {
		log.Printf("Marshalling queryStatement failed with error: %s\n", err)
	}
	
	request, err := http.NewRequest("POST", graphUrl, bytes.NewBuffer(jsonQuery))
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	if err != nil {
		log.Printf("Creating a new request failed with error: %s\n", err)
	}
	
	response, err := client.Do(request)
	if err != nil {
		log.Printf("The HTTP request failed with error: %s\n", err)
	}
	fmt.Printf("Returned request")

	defer response.Body.Close()

	data, _ := ioutil.ReadAll(response.Body)

	var queryResponse query.EIP721QueryResponse
	err = json.Unmarshal(data, &queryResponse)
	if err != nil {
		panic(err)
	}
	fmt.Println(queryResponse)
	fmt.Printf("starting data transformation")
	for i, contract := range queryResponse.Data.Contracts {
		for _, token := range queryResponse.Data.Contracts[i].Tokens {
			record := query.NewERC721Token(
				contract.ID,
				contract.Name,
				contract.Symbol,
				token.Identifier,
				token.Owner.ID,
				token.Uri,
			)
			fmt.Printf("result %d: %+v", i, record)
			records = append(records, []string{record.ContractID, record.Name, record.Symbol, record.Identifier, record.Owner, record.Uri})

		}
	}

	// f, err := os.Create("tokens.csv")
 //    defer f.Close()

 //    if err != nil {

 //        log.Fatalln("failed to open file", err)
 //    }

 //    w := csv.NewWriter(f)
 //    err = w.WriteAll(records)

	// if err != nil {
	//     log.Fatal(err)
	// }
}







	// for i := 0; i < len(timeRange)-1; i++ {
	// 	queryStatement := query.EIP721Query(1000, timeRange[i], timeRange[i+1], 5)		
	// 	go func () {
	// 		results := [][]string{{}}
	// 		log.Printf("started a goroutine")			
	// 		jsonQuery, err := json.Marshal(queryStatement)
	// 		if err != nil {
	// 			log.Printf("Marshalling queryStatement failed with error: %s\n", err)
	// 		}
			
	// 		request, err := http.NewRequest("POST", graphUrl, bytes.NewBuffer(jsonQuery))
	// 		request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	// 		if err != nil {
	// 			log.Printf("Creating a new request failed with error: %s\n", err)
	// 		}
			
	// 		response, err := client.DoWithRetries(request)
	// 		if err != nil {
	// 			log.Printf("The HTTP request failed with error: %s\n", err)
	// 		}

	// 		defer response.Body.Close()

	// 		data, _ := ioutil.ReadAll(response.Body)

	// 		var queryResponse query.EIP721QueryResponse
	// 		err = json.Unmarshal(data, &queryResponse)
	// 		if err != nil {
	// 			panic(err)
	// 		}
	// 		log.Printf("starting data transformation")
	// 		for i, contract := range queryResponse.Data.Contracts {
	// 			for _, token := range queryResponse.Data.Contracts[i].Tokens {
	// 				record := query.NewERC721Token(
	// 					contract.ID,
	// 					contract.Name,
	// 					contract.Symbol,
	// 					token.Identifier,
	// 					token.Owner.ID,
	// 					token.Uri,
	// 				)
	// 				log.Printf("result %d: %+v", i, record)
	// 				results = append(results, []string{record.ContractID, record.Name, record.Symbol, record.Identifier, record.Owner, record.Uri})

	// 			}
	// 		}
			
	// 		resultChan <- results
	// 	}()
	// }

	// f, err := os.Create("tokens.csv")
 //    defer f.Close()

 //    if err != nil {

 //        log.Fatalln("failed to open file", err)
 //    }

 //    w := csv.NewWriter(f)

	// for i := 0; i < len(timeRange)-1; i++ {
	// 	log.Printf("logging another result")
	// 	result := <-resultChan
	// 	fmt.Println("result %d: %+v", i, result)
	// 	records = append(records, result...)
	// 	err = w.WriteAll(result)

	// 	if err != nil {
	// 	    log.Fatal(err)
	// 	}

	// }
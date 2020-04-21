package quotescsvparser

import (
	"encoding/csv"
	"log"
	"math/rand"
	"os"
	"time"
)

func ReadQuotesCsvFile(filePath string) [][]string {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Unable to read file: "+filePath, err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Parsing failed for file: "+filePath, err)
	}

	return records
}

func GetRandomQuote(quotes [][]string) (string, string) {
	rand.Seed(time.Now().UnixNano())
	min := 1
	max := len(quotes) - 1
	quoteNumber := rand.Intn(max-min+1) + min
	quoteAuthor := quotes[quoteNumber][0]
	quoteText := quotes[quoteNumber][1]
	if quoteAuthor == "" {
		quoteAuthor = "Anonymous"
	}
	return quoteAuthor, quoteText
}

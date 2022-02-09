package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

const apiBaseUrl = "https://www.googleapis.com/books/v1/volumes?q=isbn:"

type Book struct {
	Items      []Item
	Kind       string
	TotalItems int
}

type Item struct {
	VolumeInfo VolumeInfos
}

type VolumeInfos struct {
	Title       string
	Authors     []string
	ImageLinks  ImageLink
	PreviewLink string
}

type ImageLink struct {
	Thumbnail string
}

type Output struct {
	Title  string `json:"title"`
	Author string `json:"author"`
	Img    string `json:"image"`
	Link   string `json:"url"`
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	fileName := flag.String("i", "./temp/isbns.txt", "isbns input text file")
	outputFileName := flag.String("o", "./temp/output.json", "output json file")

	if len(os.Args) < 2 {
		log.Fatal("please provide a api key")
	}

	apiKey := os.Args[1]

	flag.Parse()

	file, err := os.Open(*fileName)
	check(err)

	defer file.Close()

	scanner := bufio.NewScanner(file)

	var output []Output

	for scanner.Scan() {
		isbnNo := scanner.Text()

		apiUrl := apiBaseUrl + isbnNo + "&" + apiKey

		res, err := http.Get(apiUrl)
		check(err)

		body, err := ioutil.ReadAll(res.Body)
		check(err)

		var book Book
		json.Unmarshal(body, &book)

		if len(book.Items) > 0 {
			title := book.Items[0].VolumeInfo.Title
			author := book.Items[0].VolumeInfo.Authors[0]
			img := book.Items[0].VolumeInfo.ImageLinks.Thumbnail
			url := book.Items[0].VolumeInfo.PreviewLink
			var outputItem = Output{
				title,
				author,
				img,
				url,
			}
			output = append(output, outputItem)
		}
	}

	j, err := json.MarshalIndent(output, "", " ")

	content := strings.Replace(string(j), "\\u0026", "&", -1)
	check(err)

	f, err := os.Create(*outputFileName)
	defer f.Close()

	w := bufio.NewWriter(f)
	_, err = w.WriteString(content)
	check(err)
	w.Flush()
}

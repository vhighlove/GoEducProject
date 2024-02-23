// Задачка
// Прочитать файл CSV, XML в которых будет лежать информация про товары.
// В файле CSV нужно будет считать name и seasons, а в файле XML нужно будет считать остальные данные,
// в итоге нужно сформеровать список продуктов. Общим полем для этих двух файлов будет поле "Артикул" (sku),
// что бы мапить данные. Для начала будем парсить Категорию "Юбки". Результат кода загрузи на github пожалуйста,
// что бы я там смог посмотреть.
package main

import (
	"bytes"
	"encoding/csv"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"os"
)

type Catalog struct {
	Shop struct {
		Name   string `xml:"name"`
		Offers struct {
			Offer []Offer `xml:"offer"`
		} `xml:"offers"`
	} `xml:"shop"`
}

type Offer struct {
	Available   bool     `xml:"available,attr"`
	GroupID     int      `xml:"group_id,attr"`
	ID          int      `xml:"id,attr"`
	URL         string   `xml:"url"`
	Price       int      `xml:"price"`
	OldPrice    int      `xml:"old_price"`
	Currency    string   `xml:"currencyId"`
	Pictures    []string `xml:"picture"`
	Name        string   `xml:"name"`
	Description string   `xml:"description"`
	Vendor      string   `xml:"vendor"`
	Sku         string   `xml:"vendorCode"`
	CategoryID  int      `xml:"categoryId"`
	Params      []Param  `xml:"param"`
}

type Param struct {
	Name  string `xml:"name,attr"`
	Value string `xml:",chardata"`
}

type Skirt struct {
	Sku    string
	Name   string
	Season string
	Offer  Offer
}

func main() {
	csvFileData := Parsecsv("Files/2024.02.13.csv")
	xmlFileData := Parsexml("Files/export_rozetka.xml")
	skirts := MapCategory(csvFileData, xmlFileData)
	for _, value := range skirts {
		fmt.Println(value)
	}
}

func Parsecsv(filestr string) [][]string {
	// Open
	file, err := os.Open(filestr)

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	// Read + Parse
	reader := csv.NewReader(file)
	reader.Comma = ';'
	reader.LazyQuotes = true

	data, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}
	return data[1:]
}

func Parsexml(filestr string) *Catalog {
	// Open
	file, err := os.Open(filestr)

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	// Read + Parse
	b, err := io.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	var catalog Catalog
	err = xml.NewDecoder(bytes.NewReader(b)).Decode(&catalog)
	if err != nil {
		log.Fatal(err)
	}
	return &catalog
}

func MapCategory(csvData [][]string, catalog *Catalog) map[string]Skirt {
	skirts := make(map[string]Skirt)

	for _, csvRow := range csvData {
		name := csvRow[2]
		sku := csvRow[4]
		season := csvRow[9]
		for _, offer := range catalog.Shop.Offers.Offer {
			for _, param := range offer.Params {
				if offer.Sku == sku && param.Name == "Вид" && param.Value == "Юбки" {
					skirt := Skirt{
						Sku:    sku,
						Name:   name,
						Season: season,
						Offer:  offer,
					}
					skirts[sku] = skirt
				}
			}
		}
	}
	return skirts
}

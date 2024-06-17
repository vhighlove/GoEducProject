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

type csvDataForProduct struct {
	Name   string
	Season string
}

type Product struct {
	Sku         string
	Name        string
	Season      string
	Available   bool
	GroupID     int
	ID          int
	URL         string
	Price       int
	OldPrice    int
	Currency    string
	Pictures    []string
	OldXmlName  string
	Description string
	Vendor      string
	CategoryID  int
	Params      []Param
}

const skirtsCategoryId int = 28

func main() {
	// if len(os.Args) < 3 {
	// 	fmt.Println("Usage: program_name <csv_file_path> <xml_file_path>")
	// 	os.Exit(1)
	// }

	// csvFilePath := os.Args[1]
	// xmlFilePath := os.Args[2]
	csvFilePath := "Files/2024.02.13.csv"
	xmlFilePath := "Files/export_rozetka.xml"

	csvBytes, err := OpenFile(csvFilePath)
	if err != nil {
		log.Fatal(err)
	}

	xmlBytes, err := OpenFile(xmlFilePath)
	if err != nil {
		log.Fatal(err)
	}

	csvFileData, err := Parsecsv(csvBytes)
	if err != nil {
		log.Fatal(err)
	}

	xmlFileData, err := Parsexml(xmlBytes)
	if err != nil {
		log.Fatal(err)
	}

	skirts := MapCategory(csvFileData, xmlFileData)
	for _, skirt := range skirts {
		fmt.Println(skirt)
	}
}

func OpenFile(filePath string) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func Parsecsv(data []byte) (map[string]csvDataForProduct, error) {
	reader := csv.NewReader(bytes.NewReader(data))
	reader.Comma = ';'
	reader.LazyQuotes = true

	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	mapData := make(map[string]csvDataForProduct)
	for _, row := range records {
		name := row[2]
		sku := row[4]
		season := row[9]
		mapData[sku] = csvDataForProduct{
			Name:   name,
			Season: season,
		}
	}
	return mapData, nil
}

func Parsexml(data []byte) (*Catalog, error) {
	var catalog Catalog
	err := xml.NewDecoder(bytes.NewReader(data)).Decode(&catalog)
	if err != nil {
		return nil, err
	}
	return &catalog, nil
}

func MapCategory(csvData map[string]csvDataForProduct, catalog *Catalog) map[string]Product {
	skirts := make(map[string]Product)
	for _, offer := range catalog.Shop.Offers.Offer {
		if offer.CategoryID == skirtsCategoryId {
			skirts[offer.Sku] = Product{
				Sku:         offer.Sku,
				Name:        csvData[offer.Sku].Name,
				Season:      csvData[offer.Sku].Season,
				Available:   offer.Available,
				GroupID:     offer.GroupID,
				ID:          offer.ID,
				URL:         offer.URL,
				Price:       offer.Price,
				OldPrice:    offer.OldPrice,
				Currency:    offer.Currency,
				Pictures:    offer.Pictures,
				OldXmlName:  offer.Name,
				Description: offer.Description,
				Vendor:      offer.Vendor,
				CategoryID:  offer.CategoryID,
				Params:      offer.Params,
			}
		}
	}
	return skirts
}

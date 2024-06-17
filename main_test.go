package main

import (
	"reflect"
	"testing"
)

func TestParsecsv(t *testing.T) {
	filePath := "Tests/test1.csv"
	data, err := OpenFile(filePath)
	if err != nil {
		t.Fatalf("failed to open file: %v", err)
	}

	want := map[string]csvDataForProduct{
		"Classik_синий": {
			Name:   "Капри средней посадки синего цвета",
			Season: "Лето",
		},
		"Energy20_черный": {
			Name:   "Колготки Energy 20 den черного цвета",
			Season: "Демисезон",
		},
		"BodySlim40_черный": {
			Name:   "Колготки Body Slim 40 den черного цвета",
			Season: "Демисезон",
		},
		"Diamant20_черный": {
			Name:   "Колготки Diamant 20 den черного цвета",
			Season: "Демисезон",
		},
		"11094_черный": {
			Name:   "Черная юбка-карандаш с запахом",
			Season: "Демисезон",
		},
		"11094_зеленый": {
			Name:   "Зеленая юбка-карандаш с запахом",
			Season: "Демисезон",
		},
	}

	tests := []struct {
		name    string
		args    []byte
		want    map[string]csvDataForProduct
		wantErr bool
	}{
		{
			name:    "sample csv parse test1",
			args:    data,
			want:    want,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parsecsv(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parsecsv() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parsecsv() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParsexml(t *testing.T) {
	filePath := "Tests/test1.xml"
	data, err := OpenFile(filePath)
	if err != nil {
		t.Fatalf("failed to open file: %v", err)
	}

	want := &Catalog{
		Shop: struct {
			Name   string `xml:"name"`
			Offers struct {
				Offer []Offer `xml:"offer"`
			} `xml:"offers"`
		}{
			Name: "Market",
			Offers: struct {
				Offer []Offer `xml:"offer"`
			}{
				Offer: []Offer{
					{
						Available: true,
						GroupID:   68713,
						ID:        1375930,
						URL:       "https://issaplus.com/yubka-11094-11094_chernyy",
						Price:     548,
						OldPrice:  783,
						Currency:  "UAH",
						Pictures: []string{
							"https://issaplus.com/wa-data/public/shop/products/13/87/68713/images/157499/157499.602x0.jpg",
							"https://issaplus.com/wa-data/public/shop/products/13/87/68713/images/157500/157500.602x0.jpg",
							"https://issaplus.com/wa-data/public/shop/products/13/87/68713/images/157501/157501.602x0.jpg",
						},
						Name:        "Юбки ISSA PLUS 11094  S черный",
						Description: "<![CDATA[<p>Элегантная черная юбка на молнии с высокой посадкой и втачным поясом оснащена потайной резинкой на спинке. Изделие кроя карандаш с миди длиной и декоративным запахом образующим клин.</p>]]>",
						Vendor:      "ISSA PLUS",
						Sku:         "11094_черный",
						CategoryID:  28,
						Params: []Param{
							{Name: "Стиль", Value: "Деловой"},
							{Name: "Вид", Value: "Юбки"},
							{Name: "Размер", Value: "S"},
							{Name: "Цвет", Value: "черный"},
							{Name: "Коллекция", Value: "Ultrafashionable"},
							{Name: "Состав", Value: "50% хлопок, 50% полиэстер"},
							{Name: "Материал", Value: "Стрейч коттон"},
							{Name: "Замеры", Value: "Талия: S (60), M (64), L (68), XL (72),  Бедра: S (92), M (96), L (100), XL (104),   Длина - 72. Ткань: Низкой эластичности"},
						},
					},
					// TODO: add more offers
				},
			},
		},
	}

	tests := []struct {
		name    string
		args    []byte
		want    *Catalog
		wantErr bool
	}{
		{
			name:    "sample xml parse test1",
			args:    data,
			want:    want,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parsexml(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parsexml() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parsexml() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMapCategory(t *testing.T) {
	csvFilePath := "Tests/test1.csv"
	csvBytes, err := OpenFile(csvFilePath)
	if err != nil {
		t.Fatalf("failed to open CSV file: %v", err)
	}

	xmlFilePath := "Tests/test1.xml"
	xmlBytes, err := OpenFile(xmlFilePath)
	if err != nil {
		t.Fatalf("failed to open XML file: %v", err)
	}

	csvFileData, err := Parsecsv(csvBytes)
	if err != nil {
		t.Fatal(err)
	}

	xmlFileData, err := Parsexml(xmlBytes)
	if err != nil {
		t.Fatal(err)
	}

	type args struct {
		csvData map[string]csvDataForProduct
		catalog *Catalog
	}
	tests := []struct {
		name string
		args args
		want map[string]Product
	}{
		{
			name: "sample map skirt category test1",
			args: args{
				csvData: csvFileData,
				catalog: xmlFileData,
			},
			want: map[string]Product{
				"11094_черный": {
					Sku:       "11094_черный",
					Name:      "Черная юбка-карандаш с запахом",
					Season:    "Демисезон",
					Available: true,
					GroupID:   68713,
					ID:        1375930,
					URL:       "https://issaplus.com/yubka-11094-11094_chernyy",
					Price:     548,
					OldPrice:  783,
					Currency:  "UAH",
					Pictures: []string{
						"https://issaplus.com/wa-data/public/shop/products/13/87/68713/images/157499/157499.602x0.jpg",
						"https://issaplus.com/wa-data/public/shop/products/13/87/68713/images/157500/157500.602x0.jpg",
						"https://issaplus.com/wa-data/public/shop/products/13/87/68713/images/157501/157501.602x0.jpg",
					},
					OldXmlName:  "Юбки ISSA PLUS 11094  S черный",
					Description: "<![CDATA[<p>Элегантная черная юбка на молнии с высокой посадкой и втачным поясом оснащена потайной резинкой на спинке. Изделие кроя карандаш с миди длиной и декоративным запахом образующим клин.</p>]]>",
					Vendor:      "ISSA PLUS",
					CategoryID:  28,
					Params: []Param{
						{Name: "Стиль", Value: "Деловой"},
						{Name: "Вид", Value: "Юбки"},
						{Name: "Размер", Value: "S"},
						{Name: "Цвет", Value: "черный"},
						{Name: "Коллекция", Value: "Ultrafashionable"},
						{Name: "Состав", Value: "50% хлопок, 50% полиэстер"},
						{Name: "Материал", Value: "Стрейч коттон"},
						{Name: "Замеры", Value: "Талия: S (60), M (64), L (68), XL (72),  Бедра: S (92), M (96), L (100), XL (104),   Длина - 72. Ткань: Низкой эластичности"},
					},
				},
				// TODO: add more products
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MapCategory(tt.args.csvData, tt.args.catalog)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MapCategory() = %v, want %v", got, tt.want)
			}
		})
	}
}

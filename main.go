package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"github.com/maxmind/mmdbwriter"
	"github.com/maxmind/mmdbwriter/mmdbtype"
	log "github.com/sirupsen/logrus"
)

var (
	srcFile      string
	srcDir       string
	dstFile      string
	databaseType string
	cnRecord     = mmdbtype.Map{
		"country": mmdbtype.Map{
			"geoname_id":           mmdbtype.Uint32(1814991),
			"is_in_european_union": mmdbtype.Bool(false),
			"iso_code":             mmdbtype.String("CN"),
			"names": mmdbtype.Map{
				"de":    mmdbtype.String("China"),
				"en":    mmdbtype.String("China"),
				"es":    mmdbtype.String("China"),
				"fr":    mmdbtype.String("Chine"),
				"ja":    mmdbtype.String("中国"),
				"pt-BR": mmdbtype.String("China"),
				"ru":    mmdbtype.String("Китай"),
				"zh-CN": mmdbtype.String("中国"),
			},
		},
	}
	extraCountries        = []string{"usa", "japan", "korea", "hongkong", "taiwan"}
	extraCountriesRecords = map[string]mmdbtype.Map{
		"usa": {
			"country": mmdbtype.Map{
				"geoname_id":           mmdbtype.Uint32(6252001),
				"is_in_european_union": mmdbtype.Bool(false),
				"iso_code":             mmdbtype.String("US"),
				"names": mmdbtype.Map{
					"de":    mmdbtype.String("Vereinigte Staaten von Amerika"),
					"en":    mmdbtype.String("United States of America"),
					"es":    mmdbtype.String("Estados Unidos de América"),
					"fr":    mmdbtype.String("États-Unis d'Amérique"),
					"ja":    mmdbtype.String("アメリカ合衆国"),
					"pt-BR": mmdbtype.String("Estados Unidos da América"),
					"ru":    mmdbtype.String("Соединенные Штаты Америки"),
					"zh-CN": mmdbtype.String("美国"),
				},
			},
		},
		"japan": {
			"country": mmdbtype.Map{
				"geoname_id":           mmdbtype.Uint32(1861060),
				"is_in_european_union": mmdbtype.Bool(false),
				"iso_code":             mmdbtype.String("JP"),
				"names": mmdbtype.Map{
					"de":    mmdbtype.String("Japan"),
					"en":    mmdbtype.String("Japan"),
					"es":    mmdbtype.String("Japón"),
					"fr":    mmdbtype.String("Japon"),
					"ja":    mmdbtype.String("日本"),
					"pt-BR": mmdbtype.String("Japão"),
					"ru":    mmdbtype.String("Япония"),
					"zh-CN": mmdbtype.String("日本"),
				},
			},
		},
		"korea": {
			"country": mmdbtype.Map{
				"geoname_id":           mmdbtype.Uint32(1835841),
				"is_in_european_union": mmdbtype.Bool(false),
				"iso_code":             mmdbtype.String("KR"),
				"names": mmdbtype.Map{
					"de":    mmdbtype.String("Korea"),
					"en":    mmdbtype.String("Korea"),
					"es":    mmdbtype.String("Corea"),
					"fr":    mmdbtype.String("Corée"),
					"ja":    mmdbtype.String("韓国"),
					"pt-BR": mmdbtype.String("Coreia"),
					"ru":    mmdbtype.String("Корея"),
					"zh-CN": mmdbtype.String("韩国"),
				},
			},
		},
		"hongkong": {
			"country": mmdbtype.Map{
				"geoname_id":           mmdbtype.Uint32(1819729),
				"is_in_european_union": mmdbtype.Bool(false),
				"iso_code":             mmdbtype.String("HK"),
				"names": mmdbtype.Map{
					"de":    mmdbtype.String("Hong Kong"),
					"en":    mmdbtype.String("Hong Kong"),
					"es":    mmdbtype.String("Hong Kong"),
					"fr":    mmdbtype.String("Hong Kong"),
					"ja":    mmdbtype.String("香港"),
					"pt-BR": mmdbtype.String("Hong Kong"),
					"ru":    mmdbtype.String("Гонконг"),
					"zh-CN": mmdbtype.String("香港"),
				},
			},
		},
		"taiwan": {
			"country": mmdbtype.Map{
				"geoname_id":           mmdbtype.Uint32(1668284),
				"is_in_european_union": mmdbtype.Bool(false),
				"iso_code":             mmdbtype.String("TW"),
				"names": mmdbtype.Map{
					"de":    mmdbtype.String("Taiwan"),
					"en":    mmdbtype.String("Taiwan"),
					"es":    mmdbtype.String("Taiwan"),
					"fr":    mmdbtype.String("Taïwan"),
					"ja":    mmdbtype.String("台湾"),
					"pt-BR": mmdbtype.String("Taiwan"),
					"ru":    mmdbtype.String("Тайвань"),
					"zh-CN": mmdbtype.String("台湾"),
				},
			},
		},
	}
)

func init() {
	flag.StringVar(&srcFile, "s", "ipip_cn.txt", "specify source ip list file")
	flag.StringVar(&srcDir, "sd", "./dist", "specify extra country directory")
	flag.StringVar(&dstFile, "d", "Country.mmdb", "specify destination mmdb file")
	flag.StringVar(&databaseType, "t", "GeoIP2-Country", "specify MaxMind database type")
	flag.Parse()
}

func scan(srcFile string) []string {
	var ipTxtList []string
	fh, err := os.Open(srcFile)
	if err != nil {
		log.Fatalf("fail to open %s\n", err)
		os.Exit(-1)
	}
	scanner := bufio.NewScanner(fh)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		ipTxtList = append(ipTxtList, scanner.Text())
	}
	return ipTxtList
}

func main() {
	writer, err := mmdbwriter.New(
		mmdbwriter.Options{
			DatabaseType: databaseType,
			RecordSize:   24,
		},
	)
	if err != nil {
		log.Fatalf("fail to new writer %v\n", err)
	}

	ipTxtList := scan(srcFile)
	ipList := parseCIDRs(ipTxtList)
	for _, ip := range ipList {
		err = writer.Insert(ip, cnRecord)
		if err != nil {
			log.Fatalf("fail to insert to writer %v\n", err)
		}
	}
	for _, country := range extraCountries {
		ipTxtList := scan(fmt.Sprintf("%s/%s/ip.txt", srcDir, country))
		ipList := parseCIDRs(ipTxtList)
		for _, ip := range ipList {
			err = writer.Insert(ip, extraCountriesRecords[country])
			if err != nil {
				log.Fatalf("fail to insert to writer %v\n", err)
			}
		}
	}

	outFh, err := os.Create(dstFile)
	if err != nil {
		log.Fatalf("fail to create output file %v\n", err)
	}

	_, err = writer.WriteTo(outFh)
	if err != nil {
		log.Fatalf("fail to write to file %v\n", err)
	}

}

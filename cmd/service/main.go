package main

import (
	"flag"
	"fmt"
	"github.com/nkmtn/zipfetcher"
	"log"
	"zipservice/cmd/pg"
)

func main() {
	var lastUpdate string
	flag.StringVar(&lastUpdate, "date", "2020-04-01", "last known update date in format yyyy-mm-dd")
	flag.Parse()

	zf := zipfetcher.Create()
	zips, err := zf.GetAllZipsIfModifiedSince(lastUpdate)
	if err != nil {
		log.Fatal(err)
	}

	writer, err := CreateUpserter()
	if err != nil {
		log.Fatal(err)
	}
	err = writer.Write(zips)
	if err != nil {
		log.Fatal(err)
	}

	reader, err := CreateGetter()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(reader.Read())
}

type Upserter interface {
	Write([]zipfetcher.ZipCode) error
}

func CreateUpserter() (Upserter, error) {
	return pg.CreatePostgres()
}

type Getter interface {
	Read() []pg.Zip // не хорошо, что тут всплыл постгре
}

func CreateGetter() (Getter, error) {
	return pg.CreatePostgres()
}

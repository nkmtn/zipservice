package pg

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/nkmtn/zipfetcher"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"os"
	"time"
)

type Postgres struct {
	con *gorm.DB
}

func CreatePostgres() (*Postgres, error) {
	db, err := Connect()
	if err != nil {
		return nil, err
	}
	return &Postgres{con: db}, nil
}

func (p *Postgres) Write(zipCodes []zipfetcher.ZipCode) error {
	err, _ := p.addZips(p.convert(zipCodes))
	if err != nil {
		return err
	}
	return nil
}

func (p *Postgres) Read() []Zip {
	return p.GetAllZips()
}

func Connect() (*gorm.DB, error) {
	db, err := connectToDB()
	if err != nil {
		return nil, err
	}
	db.AutoMigrate(&Zip{})
	return db, nil
}

func connectToDB() (*gorm.DB, error) {
	counts := 0

	err := godotenv.Load(".env")
	if err != nil {
		return nil, err
	}
	dsn := os.Getenv("POSTGRES_CONNECTION")

	for {
		connection, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
		if err == nil {
			//log.Print("connected to database!")
			return connection, nil
		}

		if counts < 10 {
			time.Sleep(1 * time.Second)
			counts++
			continue
		}

		return nil, fmt.Errorf("can't connect to database")
	}
}

func (p *Postgres) convert(input []zipfetcher.ZipCode) []Zip {
	output := []Zip{}
	for _, i := range input {
		zip := Zip{
			ZipCode: i.Code,
			Source:  "usps",
			Location: Location{
				State: i.State,
				City:  i.City,
			},
		}
		output = append(output, zip)
	}

	return output
}

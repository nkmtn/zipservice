package pg

import (
	"fmt"
	"golang.org/x/exp/slices"
	"gorm.io/gorm"
)

type Zip struct {
	gorm.Model
	ZipCode  string   `gorm:"uniqueIndex;size:16"`
	Location Location `gorm:"embedded"`
	Source   string
	Type     string `gorm:"size:16"`
}

type Location struct {
	Country       string      `gorm:"size:16"`
	ActualCountry string      `gorm:"size:16"`
	State         string      `gorm:"size:16"`
	County        string      `gorm:"size:16"`
	City          string      `gorm:"size:32"`
	Coordinates   Coordinates `gorm:"embedded"`
}

type Coordinates struct {
	Latitude  string `gorm:"size:32"`
	Longitude string `gorm:"size:32"`
}

// addZips Add amount of zips in Database skipping already known
func (p *Postgres) addZips(zips []Zip) (error, int) {
	existedZips := p.GetOnlyZips()
	fmt.Println("Existed zips: ", len(existedZips))

	newZips := []Zip{}
	for _, z := range zips {
		if !slices.Contains(existedZips, z.ZipCode) {
			newZips = append(newZips, z)
		}
	}

	if len(newZips) == 0 {
		return nil, 0
	}

	cursor, count := 0, 0
	step := 1000
	err := error(nil)
	for cursor < len(newZips)-1 {
		chunk := newZips
		if len(newZips) > cursor+step-1 {
			chunk = newZips[cursor : cursor+step-1]
		}

		result := p.con.Create(&chunk)
		if result.Error != nil {
			err = result.Error
		}
		count += result.CreateBatchSize
		cursor += step
	}

	return err, count
}

func (p *Postgres) GetOnlyZips() []string {
	zips := []string{}
	p.con.Table("zips").Select([]string{"zip_code"}).Pluck("zip_code", &zips)
	return zips
}

func (p *Postgres) GetAllZips() []Zip {
	zips := []Zip{}
	p.con.Table("zips").Find(&zips)
	return zips
}

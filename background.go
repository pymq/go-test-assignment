package main

import (
	"encoding/xml"
	"github.com/jmoiron/sqlx"
	"github.com/pymq/go-test-assignment/model"
	"log"
	"os"
	"sync"
)
type Tvs struct {
	XMLName xml.Name `xml:"tvs"`
	Tvs   []Tv   `xml:"tv"`
}

type Tv struct {
	XMLName xml.Name `xml:"tv"`
	Id    int   `xml:"id"`
}

func processReturnedTvs(db *sqlx.DB, filename string) {
	xmlFile, err := os.Open(filename)
	if err != nil {
		log.Println(err)
		return
	}

	defer xmlFile.Close()

	var tvs Tvs
	decoder := xml.NewDecoder(xmlFile)
	err = decoder.Decode(&tvs)
	if err != nil {
		log.Println(err)
		return
	}
	idsChan := make(chan int, 100)
	var wg sync.WaitGroup
	workersNum := 10
	wg.Add(workersNum)
	for w := 1; w <= workersNum; w++ {
		go ProcessTv(idsChan, &wg, db)
	}

	for i := 0; i < len(tvs.Tvs); i++ {
		idsChan <- tvs.Tvs[i].Id
	}
	close(idsChan)
	wg.Wait()
	log.Println("file processed")

}

func ProcessTv(ids <-chan int, wg *sync.WaitGroup, db *sqlx.DB) {
	defer wg.Done()
	for id := range ids {
		tv := model.SoldTv{}
		err := db.Get(&tv, `SELECT * 
 				FROM "sold_tv" WHERE id=$1`, id)
		if err != nil {
			log.Println("soldtv id", id, "not found")
			continue
		}
		if tv.Returned {
			log.Println("soldtv id", id, "already returned")
			continue
		}
		tv.Returned = true
		query := `UPDATE "sold_tv" SET 
				"returned"=:returned WHERE id=:id`
		_, err = db.NamedExec(query, &tv)
		if err != nil {
			log.Println("error while updating soldtv id", id, err)
		}
		log.Println("tv successfully returned", id, tv)

	}
}

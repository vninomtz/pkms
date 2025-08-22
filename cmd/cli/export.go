package main

import (
	"log"
	"os"

	"github.com/vninomtz/pkms/internal"
)

const PKMS_EXPORT_DB = "PKMS_EXPORT_DB"

func ExportToDB(nodes []internal.Node) {
	db_path := os.Getenv(PKMS_EXPORT_DB)

	repo := internal.NewSqliteNodeRepo(db_path)

	err := repo.Restore()
	if err != nil {
		log.Fatal("Error restoring db: %s\n", err)
	}

	saved := 0
	for _, n := range nodes {
		err = repo.Save(n)
		if err != nil {
			log.Printf("Error saving to db note %s\n", n.Title)
			continue
		}
		saved++
	}
	log.Printf("%d Notes exported successfully\n", saved)
}

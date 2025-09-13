package main

import (
	"flag"
	"log"

	"github.com/vninomtz/pkms/internal"
)

func DocsCommand(cmd *flag.FlagSet, args []string) {
	db := cmd.String("db", "", "Database directory")
	ls := cmd.Bool("ls", false, "List all documents")
	load := cmd.String("load", "", "Load all documents from a path")
	cmd.Parse(args[2:])

	if *load != "" {
		LoadDocuments(*load, *db)
	}
	if *ls {
		List(*db)
	}
}

func List(db string) {
	repo := internal.NewRepository(db)
	docs, err := repo.All()
	if err != nil {
		log.Fatalf("Erro reading docs: %v", err)
	}
	for i, d := range docs {
		log.Printf("%d: %s %d %s", i, d.Name, d.Size, d.UpdatedAt.Format("02/01/2006"))
	}
}

func LoadDocuments(dir, db string) {
	repo := internal.NewRepository(db)
	noteRepo := internal.NewNoteRepository(db)

	err := repo.Init()
	if err != nil {
		log.Fatalf("Error to Init database: %v\n", err)
		return
	}

	docs, err := internal.Load(internal.Config{Extensions: []string{".md", ".png"}, Dir: dir})
	if err != nil {
		log.Fatalf("Error to load documents of dir: %v\n", err)
		return
	}
	size := len(docs)

	log.Printf("Documents loaded %d\n", size)

	for _, d := range docs {
		err := repo.Save(d)
		if err != nil {
			log.Fatalf("Error to save %s file: %v", d.Name, err)
			return
		}
		note, err := internal.ParseDocument(d)
		if err != nil {
			log.Fatalf("Error to parse document %s: %v", d.Name, err)
			return
		}
		err = noteRepo.Save(note)
		if err != nil {
			log.Fatalf("Error to save note %s: %v", d.Name, err)
			return
		}

	}

}

package internal

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type SQLiteRepository struct {
	db *sql.DB
}

func NewRepository(path string) *SQLiteRepository {
	db, err := sql.Open("sqlite3", path)

	if err != nil {
		log.Fatalf("Error opening DB: %s", err)
	}
	return &SQLiteRepository{
		db: db,
	}
}

func (r *SQLiteRepository) Save(node Node) error {
	q := "INSERT INTO Nodes (Title, Description, NodeType) VALUES (?,?,?)"

	_, err := r.db.Exec(q, node.Title, node.Description, node.Type)

	if err != nil {
		log.Printf("Error saving node: %s", err)
		return err
	}
	return nil
}

func (r *SQLiteRepository) GetNodes() ([]Node, error) {
	var nodes []Node
	rows, err := r.db.Query("SELECT * FROM Nodes")
	if err != nil {
		log.Printf("Error retrieving nodes: %s", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var n Node
		if err := rows.Scan(&n.Id, &n.Title, &n.Description, &n.Type); err != nil {
			log.Printf("Error scanning row: %s", err)
			return nil, err
		}
		nodes = append(nodes, n)
	}
	if err := rows.Err(); err != nil {
		log.Printf("Error scanning row: %s", err)
		return nil, err
	}
	return nodes, nil
}

func (r *SQLiteRepository) Clean() error {
	_, err := r.db.Exec("DELETE FROM Nodes")
	if err != nil {
		log.Printf("Error cleaning DB: %s", err)
		return err
	}
	return nil
}

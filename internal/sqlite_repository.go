package internal

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type sqliteNodeRepo struct {
	db *sql.DB
}

func NewSqliteNodeRepo(path string) NodeRepository {
	db, err := sql.Open("sqlite3", path)

	if err != nil {
		log.Fatalf("Error opening DB: %s", err)
	}
	return &sqliteNodeRepo{
		db: db,
	}
}

func (r *sqliteNodeRepo) Save(node Node) error {
	q := "INSERT INTO Nodes (Title, Content, Html) VALUES (?,?,?)"

	_, err := r.db.Exec(q, node.Title, node.Content, node.Html)

	if err != nil {
		log.Printf("Error saving node: %s", err)
		return err
	}
	return nil
}

func (r *sqliteNodeRepo) GetNodes() ([]Node, error) {
	var nodes []Node
	rows, err := r.db.Query("SELECT * FROM Nodes")
	if err != nil {
		log.Printf("Error retrieving nodes: %s", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var n Node
		if err := rows.Scan(&n.Id, &n.Title, &n.Content, &n.Type); err != nil {
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

func (r *sqliteNodeRepo) Restore() error {
	sql := `
DROP TABLE IF EXISTS Nodes;

CREATE TABLE Nodes(
  Title TEXT NOT NULL PRIMARY KEY,
  Content TEXT,
  Html TEXT
)`
	_, err := r.db.Exec(sql)
	return err
}

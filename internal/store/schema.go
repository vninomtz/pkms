package store

func (s *Store) Setup() error {
	schema := `
/*
PRAGMA journal_mode = WAL;
PRAGMA synchronous = NORMAL;
PRAGMA journal_size_limit = 67108864; -- 64 megabytes
PRAGMA mmap_size = 134217728; -- 128 megabytes
PRAGMA cache_size = 2000;
PRAGMA busy_timeout = 5000;
*/

PRAGMA foreign_keys = ON;

DROP TABLE IF EXISTS documents;

CREATE TABLE documents (
  name TEXT NOT NULL,
  content TEXT,
  size INTEGER,
  bytes BLOB,
  path TEXT,
  ext TEXT,
  updated_at TEXT
);

DROP TABLE IF EXISTS notes;

CREATE TABLE parse_documents (
  document_id INTEGER NOT NULL,
  title TEXT NOT NULL,
  content TEXT,
  is_public BIT,
  tags TEXT,
  type TEXT,
  FOREIGN KEY (document_id) REFERENCES documents(rowid) ON DELETE CASCADE
);

DROP TABLE IF EXISTS bookmarks;

CREATE TABLE bookmarks (
  url TEXT,
  note_title TEXT
)
	`
	_, err := s.db.Exec(schema)
	return err
}

/*
PRAGMA journal_mode = WAL;
PRAGMA synchronous = NORMAL;
PRAGMA journal_size_limit = 67108864; -- 64 megabytes
PRAGMA mmap_size = 134217728; -- 128 megabytes
PRAGMA cache_size = 2000;
PRAGMA busy_timeout = 5000;
*/


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

CREATE TABLE notes (
  title TEXT NOT NULL,
  content TEXT,
  is_public BIT,
  tags TEXT,
  type TEXT
);

DROP TABLE IF EXISTS bookmarks;

CREATE TABLE bookmarks (
  url TEXT,
  note_title TEXT
)


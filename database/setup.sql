DROP TABLE IF EXISTS Nodes;

CREATE TABLE Nodes(
  Id INTEGER PRIMARY KEY,
  Title TEXT NOT NULL,
  Description TEXT,
  NodeType TEXT
)

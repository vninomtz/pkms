package internal

type Node struct {
	Id      int32
	Title   string
	Content string
	Type    string
}

type NodeRepository interface {
	Save(Node) error
	GetNodes() ([]Node, error)
}

type NoteService interface {
	New(title, content string) error
	ListAll() ([]Node, error)
}

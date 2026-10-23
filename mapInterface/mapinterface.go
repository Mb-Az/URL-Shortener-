package mapInterface

type Map interface {
	Get(string) (string, error)
	Add(string, string) error
	Save() error
	Load()
}

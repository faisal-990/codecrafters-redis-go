package db

// the storage is supposed to be a generic type , that can be swapped if needed
// currently only supporting in memory db in the form of map[string]string
type Storage interface {
	Set(key, value string) error
	Get(key string) (string, error)
}

type DB struct {
	data map[string]string
}

var Instance *DB

func Init() {
	Instance = &DB{
		data: make(map[string]string),
	}
}

func (d *DB) Set(key, value string) error {
	d.data[key] = value
	return nil
}

func (d *DB) Get(key string) (string, error) {
	var res string
	res = d.data[key]

	return res, nil
}

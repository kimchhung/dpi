package database

type DBConn struct {
	name string
}

func New(name string) *DBConn {
	return &DBConn{
		name: name,
	}
}

func (d *DBConn) Name() string {
	return d.name
}

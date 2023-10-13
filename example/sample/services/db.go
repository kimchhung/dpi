package services

type DBConn struct {
}

func (d *DBConn) Name() string {
	return "this is DB"
}

package models

const (
	SrvSSH = iota
	SrvFTP
	SrvHTTP
	SrvSSBD
)

type Server struct {
	ServerID int64
	Name     string
	Proto    int
	Address  []byte
	Port     int
}

// GetServers returns all Servers.
func (c *Client) GetServers() ([]Server, error) {
	s := `select * from servers order by serverid asc`

	rows, err := c.DB.Query(s)
	if err != nil {
		return []Server{}, err
	}
	defer rows.Close()

	srvs := make([]Server, 0, 16)
	for rows.Next() {
		srv := Server{}
		err = rows.Scan(
			&srv.ServerID,
			&srv.Name,
			&srv.Proto,
			&srv.Address,
			&srv.Port)
		if err != nil {
			return []Server{}, err
		}
		srvs = append(srvs, srv)
	}

	return srvs, nil
}

// InsertServer inserts Server s.
func (c *Client) InsertServer(s Server) (int64, error) {
	st := `insert into servers values (?, ?, ?, ?, ?)`

	res, err := c.DB.Exec(st, nil,
		s.Name,
		s.Proto,
		s.Address,
		s.Port)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

// UpdateServer updates a Server s identified by s.ServerID.
func (c *Client) UpdateServer(s Server) error {
	st := `update servers
	set name=?,proto=?,address=?,port=?
	where serverid=?`

	_, err := c.DB.Exec(st,
		s.Name,
		s.Proto,
		s.Address,
		s.Port,
		s.ServerID)

	return err
}

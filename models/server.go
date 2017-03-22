package models

type Server struct {
	ServerID int64
	Name     string
	Address  []byte
	Port     int
}

// GetServers returns all servers.
func (c *Client) GetServers() ([]Server, error) {
	s := `select * from servers order by serverid asc`

	rows, err := c.DB.Query(s)
	if err != nil {
		return []Server{}, ErrQueryFailed
	}
	defer rows.Close()

	srvs := make([]Server, 0, 16)
	for rows.Next() {
		srv := Server{}
		err = rows.Scan(
			&srv.ServerID,
			&srv.Name,
			&srv.Address,
			&srv.Port)
		if err != nil {
			return []Server{}, ErrScan
		}
		srvs = append(srvs, srv)
	}

	return srvs, nil
}

// InsertServer inserts Server s.
func (c *Client) InsertServer(s Server) error {
	st := `insert into servers values (?, ?, ?, ?)`

	_, err := c.DB.Exec(st, nil,
		s.Name,
		s.Address,
		s.Port)

	return err
}

// UpdateServer updates a Server s identified by s.ServerID.
func (c *Client) UpdateServer(s Server) error {
	st := `update servers
	set name=?,address=?,port=?
	where serverid=?`

	_, err := c.DB.Exec(st,
		s.Name,
		s.Address,
		s.Port,
		s.ServerID)

	return err
}

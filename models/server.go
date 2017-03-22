package models

type Server struct {
	ServerID int64
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
	return nil
}

// UpdateServer updates a Server s identified by s.ServerID.
func (c *Client) UpdateServer(s Server) error {
	return nil
}

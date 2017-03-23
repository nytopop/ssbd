package models

const (
	FileDir = iota
	SSH
	AWS
	RBD
)

type Volume struct {
	VolumeID int64
	Name     string
	Backend  int
	AuthUser string
	AuthPW   string
	Capacity int64
	Free     int64
	Used     int64
}

// GetVolumes returns all volumes.
func (c *Client) GetVolumes() ([]Volume, error) {
	s := `select * from volumes order by volumeid asc`

	rows, err := c.DB.Query(s)
	if err != nil {
		return []Volume{}, ErrQueryFailed
	}
	defer rows.Close()

	vols := make([]Volume, 0, 16)
	for rows.Next() {
		v := Volume{}
		err = rows.Scan(
			&v.VolumeID,
			&v.Name,
			&v.Backend,
			&v.AuthUser,
			&v.AuthPW,
			&v.Capacity,
			&v.Free,
			&v.Used)
		if err != nil {
			return []Volume{}, ErrScan
		}
		vols = append(vols, v)
	}

	return vols, nil
}

// InsertVolume inserts a new volume v.
func (c *Client) InsertVolume(v Volume) error {
	s := `insert into volumes values (?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := c.DB.Exec(s,
		nil,
		v.Name,
		v.Backend,
		v.AuthUser,
		v.AuthPW,
		v.Capacity,
		v.Free,
		v.Used)

	return err
}

// UpdateVolume updates a volume v identified by v.VolumeID.
func (c *Client) UpdateVolume(v Volume) error {
	s := `update volumes
	set name=?,backend=?,authuser=?,authpw=?,capacity=?,free=?,used=?
	where volumeid=?`

	_, err := c.DB.Exec(s,
		v.Name,
		v.Backend,
		v.AuthUser,
		v.AuthPW,
		v.Capacity,
		v.Free,
		v.Used,
		v.VolumeID)

	return err
}

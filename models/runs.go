package models

type Run struct {
	RunID  int64
	JobID  int64
	Status int64
}

// GetRuns returns all Runs.
func (c *Client) GetRuns() ([]Run, error) {
	s := `select * from runs order by runid desc`

	rows, err := c.DB.Query(s)
	if err != nil {
		return []Run{}, ErrQueryFailed
	}
	defer rows.Close()

	runs := make([]Run, 0, 16)
	for rows.Next() {
		run := Run{}
		err = rows.Scan(
			&run.RunID,
			&run.JobID,
			&run.Status)
		if err != nil {
			return []Run{}, ErrScan
		}
		runs = append(runs, run)
	}

	return runs, nil
}

// InsertRun inserts Run r.
func (c *Client) InsertRun(r Run) error {
	s := `insert into runs values (?, ?, ?)`

	_, err := c.DB.Exec(s, nil,
		r.JobID,
		r.Status)

	return err
}

// UpdateRun updates a Run r identified by r.RunID.
func (c *Client) UpdateRun(r Run) error {
	s := `update runs
	set jobid=?,status=?
	where runid=?`

	_, err := c.DB.Exec(s,
		r.JobID,
		r.Status,
		r.RunID)

	return err
}

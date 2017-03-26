package models

const (
	StatusGood = iota
	StatusWait
	StatusFail
)

type Run struct {
	RunID  int64
	JobID  int64
	Status int
}

// GetRuns returns all Runs.
func (c *Client) GetRuns() ([]Run, error) {
	s := `select * from runs order by runid desc`

	rows, err := c.DB.Query(s)
	if err != nil {
		return []Run{}, err
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
			return []Run{}, err
		}
		runs = append(runs, run)
	}

	return runs, nil
}

func (c *Client) GetLastFullRunID(sid int64, dir string) (int64, error) {
	s := `
	SELECT
	  *
	FROM runs

	INNER JOIN jobs ON
	  runs.jobid=jobs.jobid,
	  runs.status=?,
	  jobs.serverid=?,
	  jobs.directory=?,
	  jobs.style=?
	
	ORDER BY
	  runs.runid 
	DESC

	LIMIT 1
	`

	rows, err := c.DB.Query(s,
		StatusGood,
		sid,
		dir,
		Full)
	if err != nil {
		return 0, err
	}

	r := Run{}
	err = rows.Scan(&r.RunID, &r.JobID, &r.Status)
	if err != nil {
		return 0, err
	}

	return r.RunID, nil
}

// InsertRun inserts Run r.
func (c *Client) InsertRun(r Run) (int64, error) {
	s := `insert into runs values (?, ?, ?)`

	res, err := c.DB.Exec(s, nil,
		r.JobID,
		r.Status)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
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

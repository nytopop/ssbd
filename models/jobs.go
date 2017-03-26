package models

const (
	Full = iota
	Diff
)

type Job struct {
	JobID     int64
	ServerID  int64
	VolumeID  int64
	Style     int
	Cron      string
	Directory string
	Squash    bool
	Encrypt   bool
	Key       string
	// job type (full / differential)
	// job scheduled time
	// times executed
}

// GetJobs returns all jobs.
func (c *Client) GetJobs() ([]Job, error) {
	s := `select * from jobs order by jobid asc`

	rows, err := c.DB.Query(s)
	if err != nil {
		return []Job{}, err
	}
	defer rows.Close()

	jobs := make([]Job, 0, 16)
	for rows.Next() {
		job := Job{}
		err = rows.Scan(
			&job.JobID,
			&job.ServerID,
			&job.VolumeID,
			&job.Style,
			&job.Cron,
			&job.Directory,
			&job.Squash,
			&job.Encrypt,
			&job.Key)
		if err != nil {
			return []Job{}, err
		}
		jobs = append(jobs, job)
	}

	return jobs, nil
}

// InsertJob inserts Job j.
func (c *Client) InsertJob(j Job) (int64, error) {
	s := `insert into jobs values (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	res, err := c.DB.Exec(s, nil,
		j.ServerID,
		j.VolumeID,
		j.Style,
		j.Cron,
		j.Directory,
		j.Squash,
		j.Encrypt,
		j.Key)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

// UpdateJob updates a Job j identified by j.JobID.
func (c *Client) UpdateJob(j Job) error {
	s := `update jobs
	set serverid=?,volumeid=?,style=?,cron=?directory=?,squash=?,encrypt=?,key=?
	where jobid=?`

	_, err := c.DB.Exec(s,
		j.ServerID,
		j.VolumeID,
		j.Style,
		j.Cron,
		j.Directory,
		j.Squash,
		j.Encrypt,
		j.Key,
		j.JobID)

	return err
}

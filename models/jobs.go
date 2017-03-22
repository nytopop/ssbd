package models

type Job struct {
	JobID     int64
	ServerID  int64
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
		return []Job{}, ErrQueryFailed
	}
	defer rows.Close()

	jobs := make([]Job, 0, 16)
	for rows.Next() {
		job := Job{}
		err = rows.Scan(
			&job.JobID,
			&job.ServerID,
			&job.Directory,
			&job.Squash,
			&job.Encrypt,
			&job.Key)
		if err != nil {
			return []Job{}, ErrScan
		}
		jobs = append(jobs, job)
	}

	return jobs, nil
}

// InsertJob inserts Job j.
func (c *Client) InsertJob(j Job) error {
	s := `insert into jobs values (?, ?, ?, ?, ?, ?)`

	_, err := c.DB.Exec(s, nil,
		j.ServerID,
		j.Directory,
		j.Squash,
		j.Encrypt,
		j.Key)

	return err
}

// UpdateJob updates a Job j identified by j.JobID.
func (c *Client) UpdateJob(j Job) error {
	s := `update jobs
	set serverid=?,directory=?,squash=?,encrypt=?,key=?
	where jobid=?`

	_, err := c.DB.Exec(s,
		j.ServerID,
		j.Directory,
		j.Squash,
		j.Encrypt,
		j.Key,
		j.JobID)

	return err
}

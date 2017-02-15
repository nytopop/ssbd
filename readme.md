# ssbd - simple ssh backup daemon
[![Build Status](https://travis-ci.org/nytopop/ssbd.svg?branch=master)](https://travis-ci.org/nytopop/ssbd)

ssbd is a simple daemon that backs up linux computers using ssh and tar. No client software is required excepting an ssh server (tested with OpenSSH) and GNU tar. A web interface is exposed for management of backups. It's super cool.

End to end encryption is used during transfer of files, and data at rest can be optionally encrypted using AES256. Data integrity is guaranteed with sha256 checksums at multiple stages, and old backups are periodically verified. Your data is safe with ssbd.

## Status

In progress / early dev.

## How it works
ssbd signs into a listening server over ssh, runs tar locally, and pipes the resulting datastream back through the encrypted ssh channel for safe archival. Diffs are performed on the remote server, not in ssbd - making ssbd very fast even with lots of servers.

Backups are verified for integrity before encryption, after encryption, and in periodic health checks.

ssbd boasts an advanced job scheduler, maximizing performance and reliability of backup jobs.

## Workflow:
- Login to ssbd web interface.
- Register a server to backup.
- Configure a scheduled backup job.
- Rest easy.

## TODO
- [ ] Job queuer
- [ ] Web interface
- [ ] Cloud integration
- [ ] Automatic TLS certs; self-signed and let's encrypt

## Web interface endpoints
### /servers
- [ ] GET /servers
- [ ] GET /servers/add
- [ ] POST /servers/add
- [ ] GET /servers/del/:serverid

### /jobs
- [ ] GET /jobs
- [ ] GET /jobs/queue
- [ ] GET /jobs/add
- [ ] POST /jobs/add
- [ ] GET /jobs/del/:jobid

### /history
- [ ] GET /history/:page
- [ ] GET /history/:jobid

### /browse
- [ ] GET /browse/:jobid

### /admin
- [ ] GET /admin
- [ ] GET /admin/users
- [ ] POST /admin/users/add
- [ ] GET /admin/users/del/:userid

##
- .tar         no compress, no encrypt
- .tar.aes     no compress, encrypt
- .tar.gz      compress, no encrypt
- .tar.gz.aes  compress, encrypt

- for full backup, index tree and download/tar it

- Mount tar as directory to examine files
- mount several and use rsync library 

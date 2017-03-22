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
- [ ] volumes logic, forms/validation
- [ ] servers logic, forms/validation

- [ ] Job queuer
- [ ] Cloud integration
- [ ] Automatic TLS certs; self-signed and let's encrypt

## Feature Creep List
- [ ] Storage volume support for obj storage [ceph, swift, etc], S3, dbox
- [ ] API endpoint that builds an auto-configuring binary for a persistent
      client connection

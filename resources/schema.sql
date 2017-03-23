--
-- SSBD Database Schema
--
-- Text encoding used: UTF-8
--
PRAGMA foreign_keys = off;
BEGIN TRANSACTION;

-- Table: servers
CREATE TABLE IF NOT EXISTS servers (
    serverid INTEGER PRIMARY KEY ASC AUTOINCREMENT
                     NOT NULL
                     UNIQUE,
    name     TEXT    NOT NULL
                     UNIQUE,
    address  TEXT    NOT NULL,
    port     INTEGER NOT NULL
);

-- Table: jobs
CREATE TABLE IF NOT EXISTS jobs (
    jobid     INTEGER PRIMARY KEY ASC AUTOINCREMENT
                      NOT NULL
                      UNIQUE,
    serverid  INTEGER REFERENCES servers (serverid)
                      ON DELETE CASCADE
                      NOT NULL,
    volumeid  INTEGER REFERENCES volumes (volumeid)
                      NOT NUll,
    style     INTEGER NOT NULL,
    cron      TEXT    NOT NULL,
    directory TEXT    NOT NULL,
    squash    BOOLEAN NOT NULL,
    encrypt   BOOLEAN NOT NULL,
    key       TEXT    NOT NULL
);

-- Table: runs
CREATE TABLE IF NOT EXISTS runs (
    runid  INTEGER PRIMARY KEY ASC AUTOINCREMENT
                   NOT NULL
                   UNIQUE,
    jobid  INTEGER REFERENCES jobs (jobid)
                   ON DELETE SET NULL,
    status INTEGER NOT NULL
);

-- Table: actionhistory
CREATE TABLE IF NOT EXISTS actionhistory (
    actionid INTEGER PRIMARY KEY ASC AUTOINCREMENT
                     NOT NULL
                     UNIQUE,
    userid   INTEGER REFERENCES users (userid)
                     ON DELETE SET NULL
);

-- Table: users
CREATE TABLE IF NOT EXISTS users (
    userid   INTEGER PRIMARY KEY ASC AUTOINCREMENT
                     NOT NULL
                     UNIQUE,
    username TEXT    NOT NULL
                     UNIQUE,
    salt     TEXT    NOT NULL,
    token    TEXT    NOT NULL
);

--Table: volumes
CREATE TABLE IF NOT EXISTS volumes (
    volumeid INTEGER PRIMARY KEY ASC AUTOINCREMENT
                     NOT NULL
                     UNIQUE,
    name     TEXT    NOT NULL
                     UNIQUE,
    backend  INTEGER NOT NULL,
    authuser TEXT    NOT NULL,
    authpw   TEXT    NOT NULL,
    capacity INTEGER NOT NULL,
    free     INTEGER NOT NULL,
    used     INTEGER NOT NULL
);

COMMIT TRANSACTION;
PRAGMA foreign_keys = on;

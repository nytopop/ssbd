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
    address  TEXT    NOT NULL,
	username TEXT    NOT NULL
);

-- Table: jobs
CREATE TABLE IF NOT EXISTS jobs (
    jobid     INTEGER PRIMARY KEY ASC AUTOINCREMENT
                      NOT NULL
                      UNIQUE,
    serverid  INTEGER REFERENCES servers (serverid)
                      ON DELETE CASCADE
                      NOT NULL,
    directory TEXT    NOT NULL,
    encrypt   BOOLEAN NOT NULL,
    key       TEXT    NOT NULL
);

-- Table: runhistory
CREATE TABLE IF NOT EXISTS runhistory (
    runid INTEGER PRIMARY KEY ASC AUTOINCREMENT
                  NOT NULL
                  UNIQUE,
    jobid INTEGER REFERENCES jobs (jobid)
                  ON DELETE SET NULL
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

COMMIT TRANSACTION;
PRAGMA foreign_keys = on;

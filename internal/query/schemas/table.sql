# This table is for general query
# location will be USA by default 
CREATE TABLE jobs (
  id SERIAL PRIMARY KEY, 
  title TEXT NOT NULL, 
  role TEXT NOT NULL, 
  location TEXT, 
  time_period TEXT
);

CREATE TABLE companies (
  id SERIAL PRIMARY KEY,
  name TEXT UNIQUE NOT NULL
);

CREATE TABLE skills (
  id SERIAL PRIMARY KEY,
  name TEXT UNIQUE NOT NULL
);

CREATE TABLE jobs_skills (
  job_id SERIAL REFERENCES jobs(id) ON DELETE CASCADE,
  skill_id SERIAL REFERENCES skills(id)
  PRIMARY KEY (job_id, skill_id)
);

CREATE TABLE role_skill_stats (
  role TEXT,
  skill_id SERIAL,
  frequency INT,
  last_updated DATE
);


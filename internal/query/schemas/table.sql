# This table is for general query
# location will be USA by default 
CREATE TABLE jobs (
  country TEXT,
  month TEXT,
  year TEXT, 
  id TEXT PRIMARY KEY, 
  title TEXT NOT NULL, 
  company TEXT,
  job_function TEXT,
  seniority TEXT[],
  role TEXT NOT NULL, 
  locations TEXT[], 
  skills TEXT[]
);

# Required skills
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


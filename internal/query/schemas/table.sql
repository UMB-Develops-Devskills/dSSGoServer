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

SELECT
    skill,
    COUNT(*) AS total_number
FROM jobs,
UNNEST(skills) AS skill
WHERE country = 'usa'
  AND month = 'june'
  AND year = '2026'
GROUP BY skill
ORDER BY total_number DESC

SELECT
    location,
    COUNT(*) AS total_number
FROM jobs,
UNNEST(locations) AS location
WHERE country = 'usa'
  AND month = 'june'
  AND year = '2026'
GROUP BY location
ORDER BY total_number DESC


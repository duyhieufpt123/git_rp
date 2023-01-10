-- 1. CREATE DATABASE 
CREATE DATABASE git_report;
CREATE USER db_user WITH PASSWORD 'Abc!123';
GRANT CONNECT ON DATABASE git_report TO db_user;
GRANT ALL PRIVILEGES ON DATABASE git_report to db_user;

--2. SUPPORT UUID (DB_USER)
\c git_report
create extension if not exists "uuid-ossp";

--3. Query
-- Sum code added and line_removed
Select commits.user_name, sum(commits.line_added) as Sum from commits
GROUP BY commits.user_name
order by sum DESC;


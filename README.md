# Wikipedia Path Finder

The purpose of this project is to find you a path from one wikipedia page to another other wikipedia page.
It works by querying the wikipedia API and extracting all of the links out of it.

A breadth first search is performed to find a path. Articles that have already been queried are saved in a persistent postgresql database that supports 5000 readers per second. And paths that have already been found are cached in a redis
server.

Both redis and postgresql are run from a docker container.

## Redis

Run `docker run -d --name redis-stack -p 6379:6379 -p 8001:8001 redis/redis-stack:latest`

## Postgres

Run `docker run --name some-postgres -e POSTGRES_PASSWORD=mysecretpassword -e POSTGRES_DB=wiki-db -d -p 5432:5432 postgres`

Connect to the database (I use pgAdmin) and run the following sql to set up the table that will be needed.

```sql
CREATE TABLE IF NOT EXISTS public.article
(
    title TEXT PRIMARY KEY,
    links TEXT[]
);
```

## Build Instructions

To build and run the project navigate to the client directory and run `go build` or `go run .`

Please also consider putting your email into the wikicalls.go file (in the variable named apiUserAgent). This will
allow wikipedia to contact you if for some reason the program causing them an issue. Do not increase the apiRequestSec!
Wikipedia requests that any single agent makes less than 200 calls per second.

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

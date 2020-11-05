### Postgres dumps
Simple utility to take full backups from Postgres. It also puts backups in `tree` structured folders, ex: `/base_path.../year/month/day/timestamp.sql.gz`. DB credentials will be built in final executable thus hidden from others from reading.
```
|-- 2020
|-- 2021
|-- 01
|-- 02
|-- (month numbers...)
|-- 12
    |-- 01
    |-- 02
    |-- (dates...)
    |-- 31
        | -- 2020-01-01T01-00-00.sql.gz
```

### Quick start
Implement `NewPostgresCredentials` function as in example in `pg_creds.go.example` & build
```sh
# clone and cd into repo
# create `pg_creds.go` file 
# set all necessary fields in pg_creds.go
cp pg_creds.go.example pg_creds.go

# build
go build -o pgdump
sudo chmod +x pgdump
# run
./pgdump
```

### Cron
Automate taking backups via cronjob.
```sh
crontab -e
# add the following line to take backups every hour
# https://crontab.guru/every-1-hour
0 * * * * /path/to/pgdump
```

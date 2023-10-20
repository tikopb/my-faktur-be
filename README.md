# my-faktur-be

This app is a vital backend component of My-Faktur app, employing Go programming, PostgreSQL for databases, and the Echo Golang framework. It ensures seamless user authentication, transaction processing, and data storage. Go guarantees speed and stability, PostgreSQL ensures secure data management, and Echo expedites API development. Together, they power My-Faktur's efficiency and reliability, offering users a responsive financial management experience.

## Environment Variables

To run this project, you will need to add the following environment variables to your .env file

database setup

`db_host`
`db_port`
`db_user`
`db_password`
`db_dbname`
`db_sslmode`
`key_secret set in 32 character! without $ or #`

port app use

`be_port`

## Installation

Install my-projec

prepare the database first use .env for setup

```bash
  run go install
```

first before use

```
   go run main.go
```

for start the app

## Authors

- [@tiko_pb](https://github.com/tikopb)

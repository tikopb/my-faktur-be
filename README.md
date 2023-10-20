# my-faktur-be

This app is a vital backend component of My-Faktur app, employing Go programming, PostgreSQL for databases, and the Echo Golang framework. It ensures seamless user authentication, transaction processing, and data storage. Go guarantees speed and stability, PostgreSQL ensures secure data management, and Echo expedites API development. Together, they power My-Faktur's efficiency and reliability, offering users a responsive financial management experience.

## Environment Variables

To run this project, you will need to add the following environment variables to your .env file

### database setup

`db_host`
`db_port`
`db_user`
`db_password`
`db_dbname`
`db_sslmode`

### Auth Setup

`key_secret (set in 32 character! without $ or #)`
`accessTimeout (in second)`
`refreshTimeout (in second)`

### port app use

`be_port`

## Installation

Install my-projec

prepare the database on postgresqk first use .env for setup

```bash
  run go install
```

first before use

```
   go run main.go
```

for start the app

---

### Program Tech Stack

<img src="https://cdn.worldvectorlogo.com/logos/golang-1.svg" alt="Golang Logo" width="45" height="25"> Go-Language

<img src="https://echo.labstack.com/img/logo-light.svg" alt="Golang Logo" width="65" height="25"> Echo Framework

<img src="https://e7.pngegg.com/pngimages/173/36/png-clipart-postgresql-logo-computer-software-database-open-source-s-text-head.png" alt="Golang Logo" width="45" height="45"> Postgresql

<img src="https://cdn.worldvectorlogo.com/logos/jwtio-json-web-token.svg" alt="Golang Logo" width="65" height="65"> Jwt

---

## Authors

- [@tiko_pb](https://github.com/tikopb)

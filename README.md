# my-faktur-be

This app is a vital backend component of My-Faktur app, employing Go programming, PostgreSQL for databases, and the Echo Golang framework. It ensures seamless user authentication, transaction processing, and data storage. Go guarantees speed and stability, PostgreSQL ensures secure data management, and Echo expedites API development. Together, they power My-Faktur's efficiency and reliability, offering users a responsive financial management experience.

## Environment Variables

To run this project, you will need to add the following environment variables to your .env file

### Database Setup

- `db_host`
- `db_port`
- `db_user`
- `db_password`
- `db_dbname`
- `db_sslmode`

### Authentication Setup

- `key_secret` (a 32-character secret, excluding $ or #)
- `accessTimeout` (in seconds)
- `refreshTimeout` (in seconds)

### Port Configuration

- `be_port` (the port on which the app will run)

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

<table>
  <tr>
    <th>Logo</th>
    <th>Technology/Framework</th>
  </tr>
  <tr>
    <td><img src="https://cdn.worldvectorlogo.com/logos/golang-1.svg" alt="Golang Logo" width="45" height="25"></td>
    <td>Go Language</td>
  </tr>
  <tr>
    <td><img src="https://echo.labstack.com/img/logo-light.svg" alt="Echo Logo" width="65" height="25"></td>
    <td>Echo Framework</td>
  </tr>
  <tr>
    <td><img src="https://e7.pngegg.com/pngimages/173/36/png-clipart-postgresql-logo-computer-software-database-open-source-s-text-head.png" alt="PostgreSQL Logo" width="45" height="45"></td>
    <td>PostgreSQL</td>
  </tr>
  <tr>
    <td><img src="https://cdn.worldvectorlogo.com/logos/jwtio-json-web-token.svg" alt="JWT Logo" width="65" height="65"></td>
    <td>JWT</td>
  </tr>
</table>

---

## Authors

- [@tiko_pb](https://github.com/tikopb)

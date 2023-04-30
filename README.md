# Transactions API

Golang API to that manages accounts creation and the transactions for the accounts.
The API was build using [go-chi](https://github.com/go-chi), a lightweight mux implementation wich provides a better REST support, allowing easy management of URL values, routing based on methods and simple middleware management.

The default database is sqlite3 using [go-sqlite3](https://github.com/mattn/go-sqlite3), and since it's self contained, there is no dependencies on external services.

## Running the API

If you have just cloned this repository, you can run the application using `docker-compose` by running:
```bash
$ docker-compose up
```

The api will be available at port `8088` http://localhost:8088

The `docker-compose.yml` definitions will take care of building the container and creating the database schema.
>For a more detailed steps on how to run the application for development, see the `Developing` section.

## API Endpoints

The documentation for the api endpoins, schemas and requests can be found on the file `docs/openapi.yaml`.

You can go to [https://editor.swagger.io/](https://editor.swagger.io/) and import the file there to try out the requests as well.

## Environment variables

The application expects the following envs:

|Name|Description|Default Value|
|-|-|-|
|PORT|Server Listen Port|80|
|DATABASE_FILE|Sqlite file location for the database|transactions.db|


## Developing

#### Dependencies

You will need to install sqlite3 to run the application on your local environment. For OSX you can install it with:

```bash
$ brew install sqlite3
```

#### Preparing

Create the `.env` file with the default values with:
```bash
$ cp env.default .env
```

Donwload the project dependencies with:
```bash
$ go mod download
```

With sqlitte3 installed, you need to create the database file with the defined schema located at `scripts/database_tables.sql`. You can do so running the following command:

```bash
$ sqlite3 transactions.db < scripts/database_tables.sql
````

#### Running the application

You can use the command:

```bash
$ go run cmd/api/main.go
```

#### Building the application

To build the application we need to define `CGO_ENABLED=1` as the `go-sqlite3` is a `CGO` enable package.

```bash
$ CGO_ENABLED=1 GOOS=<your_os> go build -o transactions-api cmd/api/main.go
```

#### Running the tests

To run the unit tests the command `go test ./...` should be used. You can specify the param `-cover` to show the coverage data.

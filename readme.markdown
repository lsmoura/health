# health

`health` is a command-line utility to convert Apple's
health export in xml format into a PostgreSQL database

## Usage

    health [options]
      -database string
        database name (default "health")
      -dbhost string
        database host (default "localhost")
      -dbpass string
        database password
      -dbport int
        database port (default 5432)
      -dbssl
        database ssl
      -dbuser string
        database user (default "postgres")
      -help
        show help
      -version
        show version and exit
      -input string
        input file (default "export.xml")
      -apply-schema
        apply schema (this will recreate all tables. It should be used on the first run.)

## Author

**[Sergio Moura](https://sergio.moura.ca)**

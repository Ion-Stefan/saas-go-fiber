# Kickstart your SaaS backend with Go and Fiber
Simple boilerplate for small SaaS projects using Go and Fiber to quickly setup a backend using Postgres and GORM.

It includes:
- Oauth login with Google
- JWT Authentication
- Lemonsqueezy API integration
- Database migrations with GORM
- Safe environment variables using godotenv

## Installation
Make sure you have [go](https://go.dev/doc/install) installed.

## Clone the repository
```
git clone https://github.com/Ion-Stefan/saas-go-fiber.git
```

## Setup the database

I'm using Postgres for this project but the code can easily be modified to use any other database.
```
sudo apt-get install postgresql
```
Login to the Postgres shell:
```
sudo -u postgres psql
```
Create the database:
```
CREATE DATABASE <dbname>;
```
Create the db user:
```
CREATE USER <username> WITH PASSWORD '<password>';
```
Grant permissions to the user:
```
\c <dbname>;
GRANT USAGE ON SCHEMA public TO <username>;
GRANT CREATE ON SCHEMA public TO <username>;
GRANT ALL PRIVILEGES ON DATABASE <dbname> TO <username>;
\q
```


## Run the project
Change into the project directory:
```
cd saas-go-fiber
```

Change into the cmd directory
```
cd cmd
```

Create the dotenv file
```
cp .env.example .env
```

Update the .env file with your database credentials.
Run the project with:
```
go run main.go
```

To enter the database shell:
```
sudo psql -U <username> -d <dbname>
\dt                   // to list the tables
SELECT * FROM users;  // to list the users
```

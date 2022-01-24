all:
	echo "CREATE DATABASE mydatabase;" | sudo su postgres psql 
	echo "CREATE TABLE strings(ID INT PRIMARY KEY NOT NULL, STRING TEXT NOT NULL);" | sudo su postgres psql 
	echo "CREATE USER user1 WITH PASSWORD 'password';" | sudo su postgres psql 
	go mod tidy
	go run main.go

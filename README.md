# user-crud

Work in progress 

Simple CRUD implementation utilizing the power of Go and performance of MongoDB to effectively maintain a database.

# Instructions

### Install Go Programming language latest version

[![N|Solid](https://sdtimes.com/wp-content/uploads/2018/02/golang.sh_-490x490.png)](https://golang.org/dl/)

### To get basic external modules for REST API

 ```sh
go get github.com/gorilla/mux go.mongodb.org/mongo-driver
```

* [mux](https://github.com/gorilla/mux) - Request router and dispatcher for matching incoming requests to their respective handler(stdlib in use currently)
* [mgo](https://pkg.go.dev/go.mongodb.org/mongo-driver) - MongoDB driver

### Configuration .json
Besure to recreate your own conf.json using the example provided

### To get this repository and run

 ```sh
$ git clone https://github.com/BryanSouza91/user-crud.git
$ go run app.go
```

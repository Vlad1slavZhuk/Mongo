# Mongo - Task 7

The service contains a database of advertisements for the sale of cars and information about them (`ID`, `Brand`, `Model`, `Color`, `Price`).

The folowing technologies have been introduced in this project:

`GraphQL` - For alternative request `CRUD`.

Database on `gRPC`:  
**NoSQL:**  
- [x] `Redis`;  
- [x] `Mongo`;  
- [ ] `Elasticsearch`.  

**SQL:**  
- [x] `Postgres`.

## Structure 
```
.
├── api
│   └── protoc
├── cmd
│   ├── grpc
│   └── http
├── internal
│   └── pkg
│       ├── auth
│       ├── config
│       ├── constErr
│       ├── data
│       ├── graphql
│       │   └── graph
│       │       ├── generated
│       │       └── model
│       ├── grpc
│       ├── migrations
│       ├── server
│       ├── service
│       └── storage
│           ├── grpc
│           ├── in-memory
│           ├── mongo
│           ├── postgres
│           └── redis
└── tests

26 directories
```

**Adress:**  
`http://localhost:[port]/login` (`POST`)  
`http://localhost:[port]/signup` (`POST`)  
`http://localhost:[port]/logout` (`POST`)  
`http://localhost:[port]/ads` (`GET`)  
`http://localhost:[port]/ad` (`POST`)  
`http://localhost:[port]/ad/{id:[1-9]\\d*}` (`GET`)  
`http://localhost:[port]/ad/{id:[1-9]\\d*}` (`PUT`)  
`http://localhost:[port]/ad/{id:[1-9]\\d*}` (`DELETE`)  
`http://localhost:[port]/gql` (`POST`)

## How to start:
### `Linux` / `Mac OS`:

`make` - Start command: `docker-compose up --build`  
`make docker-clean` - Start command: `docker-compose down --rmi=all`

---

## Environment
`TYPE_STORAGE` - Storage type for storing data. (default `in-memory`)  
`SERVER_PORT` - Port which `HTTP` server will use to listen connections. (default `8000`)  
`GRPC_PORT` - Port which `gRPC` server will use to listen connections. (default `8001`)  

### Redis
`REDIS_PORT` - Port which `Redis` server will use to listen conections. (default `8002`)  
`REDIS_PASSWORD` - `Redis` password. (default `vlad`)

### Postgres
`POSTGRES_PORT` - Port which `Postgres SQL` server will use to listen connections. (default `8003`)  
`POSTGRES_USER` - `Postgres` user/username. (default `postgres`)  
`POSTGRES_PASSWORD` - `Postgres` password. (default `postgres_password`)

### Mongo
`MONGO_PORT` - Port which `Mongo` server will use to listen connections. (default `8004`)  
`MONGO_USERNAME` - `Mongo` user/username. (default `mongo`)  
`MONGO_PASSWORD` - `Mongo` password. (default `pass`)

### Elasticsearch
`ELASTIC_PORT` - Port which `Elasticsearch` server will use to listen connections. (default `8005`)

---

## HTTP CRUD API

### `http://localhost:[port]/login` (`POST`)
Create new account and return `JWT` token.

#### Request:  
**Body:** 
```json5
{
    "username": "vlad",
    "password": "pass"
}
```
#### Response:
**Status:**
```
201 Created
```
**Header:**
```
Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDY0NDcwNTEsIm51bWJlciI6MTExMjIyMzMzLCJ1c2VybmFtZSI6InZsYWQifQ.99e6Xa9JDhyV-re1jyB-Yx4qPmtnfPDwsVbOu5Gmcy8
```

---

### `http://localhost:[port]/signup` (`POST`) 
Updates and returns a new token if an account exists in the database/storage.

#### Request:  
**Body:** 
```json5
{
    "username": "vlad",
    "password": "pass"
}
```
#### Response:
**Status:**
```
200 OK
```
**Header:**
```
Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDY0NDc0NDYsIm51bWJlciI6MTExMjIyMzMzLCJ1c2VybmFtZSI6InZsYWQifQ.pe10imrGex8JsE381fE9i8sg5g5jjFsTtS2GaGgKIUQ
```

---

### `http://localhost:[port]/logout` (`POST`)
Sets an empty token of the current account.

#### Request:  
**Authorization:** 
```
Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDY0NDc0NDYsIm51bWJlciI6MTExMjIyMzMzLCJ1c2VybmFtZSI6InZsYWQifQ.pe10imrGex8JsE381fE9i8sg5g5jjFsTtS2GaGgKIUQ
```
#### Response:
**Status:**
```
200 OK
```
**Header:**
```
You Logout! Bye-bye...
```

---

### `http://localhost:[port]/ads` (`GET`)
Get all ads.

#### Request:  
**nothing**
#### Response:
**Status:**
```
200 OK
```
**Header:**
```json5
[
    {
        "id": 1,
        "brand": "Nissan",
        "model": "Z500",
        "color": "Black",
        "price": 50000
    }
]
```

---

### `http://localhost:[port]/ad` (`POST`)
Create a new ad.

#### Request: 
**Body:** 
```json5
{
    "brand": "Nissan",
    "model": "Z500",
    "color": "Black",
    "price": 50000
}
```
#### Response:
**Status:**
```
201 Created
```
**Header:**
```
Create new Ad
```

---

### `http://localhost:[port]/ad/{id:[1-9]\\d*}` (`GET`) 
Get ad by `id`.

#### Request: 
**nothing**
#### Response:
**Status:**
```
200 OK
```
**Header:**
```json5
{
    "id": 1,
    "brand": "Nissan",
    "model": "Z500",
    "color": "Black",
    "price": 50000
}
```

---
### `http://localhost:[port]/ad/{id:[1-9]\\d*}` (`PUT`)
Update ad by `id`.

#### Request: 
**Body:**
```json5
{
    "brand": "Mazda",
    "model": "6",
    "color": "Yellow",
    "price": 100
}
```
#### Response:
**Status:**
```
200 OK
```
**Header:**
```
Update
```

---

### `http://localhost:[port]/ad/{id:[1-9]\\d*}` (`DELETE`) 
Delete ad by `id`.

#### Request: 
**nothing**
#### Response:
**Status:**
```
200 OK
```
**Header:**
```
Delete
```

---

## `GraphQL` Method

### `http://localhost:[port]/query` (`POST`)
Create a new ad.

#### Request: 
**Body:** 
```graphql
mutation {
  createAd(ad: {
    brand: "Mazda",
    model: "CX-5",
    color: "Red",
    price: 50000
  })
}
```
#### Response:
```graphql
{
  "data": {
    "createAd": "Create"
  }
}
```

---

### `http://localhost:[port]/query` (`POST`)
Get ad by `id`.

#### Request: 
**Body:** 
```graphql
query{
  get(id: 1){
    id
    brand
    model
    color
    price
  }
}
```
#### Response:
```graphql
{
  "data": {
    "get": {
      "id": "1",
      "brand": "Mazda",
      "model": "6",
      "color": "Yellow",
      "price": 100
    }
  }
}
```

---

### `http://localhost:[port]/query` (`POST`)
Get all ads.

#### Request: 
**Body:** 
```graphql
query{
  get(id: 1){
    id
    brand
    model
    color
    price
  }
}
```
#### Response:
```graphql
{
  "data": {
    "getall": [
      {
        "id": "1",
        "brand": "Mazda",
        "model": "6",
        "color": "Yellow",
        "price": 100
      }
    ]
  }
}
```

---

### `http://localhost:[port]/query` (`POST`)
Update ad by `id`.

#### Request: 
**Body:** 
```graphql
mutation {
  updateAd(ad: {
    brand: "BMW",
    model: "X5",
    color: "Black",
    price: 50000
  } id: 1)
}
```
#### Response:
```graphql
{
  "data": {
    "updateAd": "Update"
  }
}
```

---

### `http://localhost:[port]/query` (`POST`)
Delete ad by `id`.

#### Request: 
**Body:** 
```graphql
mutation {
  deleteAd(id:1)
}
```
#### Response:
```graphql
{
  "data": {
    "deleteAd": "Delete"
  }
}
```
# backend

middleware backend

The middleware has the following modules
- Login
- Server
- Controller
- Storage
- Database
- Contract

## login
DID

## Server 
---
The Server module provides api, and calls the corresponding method of the Controller according to the api


## Controller 
---
The Controller receives the Server request, uses the database, storage, and contract interfaces according to the request, and returns the result


## Storage
---
Storage module provides storage-related interfaces: upload, download, delete


## Database
--- 
The Database module caches some data and provides services such as storage lists


## Contract
---
The contract module is responsible for all calls to the contract

## Interface
### Login
```
Method: POST
URL http://xxx.xxx:8000/login
```

### PutObject
```
Method: POST
URL http://xxx.xxx:8000/$storage/

support storage:
mefs
ipfs
```
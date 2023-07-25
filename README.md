# F01 - Single Service (BE)

## Documentation
You can read all Redhub API documentation from here [Documentation](EXPLANATION.md)

## Brief Description
This service serves as a backend for the Admin Page and callbacks from the monolith system.

## Specification
1. CRUD Barang
2. CRUD Perusahaan
3. Login Admin

## Tech Stack
- Go
- Gorm
- Muxt
- PostgreSQL
- Docker

## How to Run
Before you run this project locally, you can copy .env.example into .env then set the environment. After that, run it with docker.
```sh
    docker-compose build
    docker-compose up -d
```
It will automatically migrate and seed the database.

## Copyright
2023 © bangkitdc. All Rights Reserved.

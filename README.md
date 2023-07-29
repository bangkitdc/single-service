# F01 - Single Service (BE)

https://github.com/bangkitdc/single-service/assets/87227379/698b392b-1180-46fc-a756-3df43ab01f0f

## Documentation
You can read all technical and API documentation from here [Documentation](EXPLANATION.md)

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
Before you run this project locally, you can copy the .env.example into .env and then set the environment. After that, run it with docker.
```sh
    docker-compose build
    docker-compose up -d
```
It will automatically migrate and seed the database.

## Copyright
2023 © bangkitdc. All Rights Reserved.

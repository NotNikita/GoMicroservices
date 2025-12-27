# Udemy Golang course on the topic of Microservices

## Setup

`docker-compose up`

## Working with Microservices in Go (Golang)
## System Architecture Diagram

┌─────────────────────────────────────────────────────────┐
│                   Frontend (Port 3000)                   │
│              (Fiber + HTML Templates)                    │
└────────────────────────────┬────────────────────────────┘
                             │ HTTP
                             ▼
┌─────────────────────────────────────────────────────────┐
│              Broker Service (Port 8080)                  │
│            (Central Request Router)                      │
└────┬──────────────────────────────────────┬──────────────┘
     │ HTTP                                 │ RabbitMQ
     ▼                                      ▼
┌─────────────────────────┐    ┌──────────────────────────┐
│  Auth Service           │    │  Listener Service        │
│  (Port 9090)            │    │  (Port 6060)             │
│  - /login               │    │  - Consumes from RabbitMQ│
│  - User authentication  │    │  - Routes to Logger      │
│  - PostgreSQL DB        │    │  - Event handlers        │
└─────────────────────────┘    └──────────────────────────┘
                                      │ HTTP
                                      ▼
                        ┌──────────────────────────┐
                        │  Logger Service          │
                        │  (Port 7070)             │
                        │  - /log (POST)           │
                        │  - /logs (GET)           │
                        │  - /log/:id (GET)        │
                        │  - /log (PUT)            │
                        │  - /logs/drop (DELETE)   │
                        │  - MongoDB storage       │
                        └──────────────────────────┘
# Dependency Dashboard Project

This project is a full-stack solution for managing software dependencies and their OpenSSF scores (fetched from deps.dev). It consists of a Go-based Backend and a React-based Frontend.

---

## Global Project Overview

The project is divided into two main components:
- **/backend**: Go REST API with SQLite storage and OpenAPI documentation.
- **/frontend**: React application for data visualization and package management.

---

## Backend Service

The Backend is a REST API that manages dependency metadata, security scores, and local caching.

### Service Modes
The behavior of the API is determined by the `SERVICE_MODE` environment variable:

1. **AUTHORITATIVE_DB**: 
   - Full CRUD capabilities.
   - Fetches missing data from deps.dev and persists it in the local database.
   - `POST`: Create/Update from upstream. `PUT`: Overwrite record. `PATCH`: Partial update enabled.

2. **UPSTREAM_CACHE**: 
   - Read-only cache behavior with proxy sync.
   - `POST` and `PATCH`: Disabled (405 Method Not Allowed).
   - `PUT`: Used exclusively to trigger a manual content sync from upstream.

### Environment Variables
| Variable | Default | Description |
|----------|---------|------------|
| `SERVICE_MODE` | `UPSTREAM_CACHE` | Modes: `AUTHORITATIVE_DB` or `UPSTREAM_CACHE` |
| `DB_PATH` | `/data/app.db` | Path to the SQLite database file |
| `HTTP_CLIENT_TIMEOUT_SEC` | `10` | Timeout for external API calls (deps.dev) |
| `CORS_ADDRESS` | `*` | Allowed Origin for CORS policy |
| `LOG_LEVEL` | `INFO` | Logging verbosity (DEBUG, INFO, WARN, ERROR) |

### Backend Quick Start
1. **Prerequisites**: Go 1.25.6, [Swag CLI](https://github.com/swaggo/swag).
2. **Generate Documentation**:
```bash
   cd backend
   swag init -g cmd/api/main.go -o ./docs -d .

```

3. **Run Locally**:
```bash
cd backend
go run cmd/api/main.go

```

## Docker Deployment

The project is optimized for Docker. You can build and run the service using the provided Dockerfile.

For running backend and frontend use docker-compose.yml file.

### Building the Image

```bash
cd backend
docker build -t dependency-dashboard-backend .

```

### Running the Container

```bash
docker run -p 8080:8080 \
  -e LOG_LEVEL=DEBUG \
  -e SERVICE_MODE=UPSTREAM_CACHE \
  -e CORS_ADDRESS=* \
  -v ./tmp:/data \
  dependency-dashboard-backend

```

*Note: Recommend mounting a volume to local `./tmp` to persist the SQLite database. You need to grant permisions on this local tmp directory*


## API Documentation

Once the server is running, you can access the interactive Swagger UI at:
`http://localhost:8080/api/v1/docs/index.html`


## Frontend
Frontend works ONLY for backend in UPSTREAM_CACHE mode. And is chatGPT generated since I do not have frontend experience.

## Overview

This is a React-based frontend application for visualizing and managing software package metadata stored in the backend service.

The application allows users to:

* View stored dependencies
* Synchronize a package with the backend (which fetches data from deps.dev)
* Delete a stored package
* Search dependencies by name
* Visualize OpenSSF score using a bar chart

The frontend communicates with a REST API exposed by the backend service at:

```
http://localhost:8080/api/v1
```

---

## Tech Stack

* React 18
* Vite
* Axios
* Recharts
* JavaScript (ES6+)

---

## Dependencies

The project depends on the following npm packages:

* react
* react-dom
* axios
* recharts
* vite
* @vitejs/plugin-react

All dependencies are defined in `package.json`.

---

## Project Structure

```
frontend/
  ├── index.html
  ├── package.json
  ├── vite.config.js
  └── src/
      ├── main.jsx
      └── App.jsx
```

* `App.jsx` – main application component
* `main.jsx` – React entry point
* `vite.config.js` – Vite configuration

---

## Prerequisites

Before running the frontend, ensure you have:

* Node.js (recommended v18 or newer)
* npm (comes with Node.js)

You can verify your installation:

```
node -v
npm -v
```

---

## Installation

After cloning the repository:

```
git clone <REPOSITORY_URL>
cd <REPOSITORY_NAME>/frontend
```

Install dependencies:

```
npm install
```

---

## Running in Development Mode

Start the development server:

```
npm run dev
```

By default, the app will be available at:

```
http://localhost:5173
```

Make sure the backend API is running at:

```
http://localhost:8080
```

---

## Features Implemented

* Fetch all dependencies from backend
* Search dependency by name
* Synchronize package using PUT request
* Delete package using DELETE request
* Display data in table
* Visualize OpenSSF score using bar chart
* Basic loading state handling
* Defensive handling of empty API responses

---

## Orchestration

For a full environment (Backend + Frontend), use the `docker-compose.yml` file located in the root directory:

```bash
docker-compose up --build

```
# PROJECT Future improvements

* TESTS!!! I would like to use AI agent to speed tests process but I do not have any private subscription at the moment.
* pagination of response with list of dependencies
* metrics
* better logging for debugging
* address other TODOs for more maintainable code

---

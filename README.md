# Sammy-PO Stadium App

A web application that shows upcoming football matches at Sammy Ofer Stadium in Haifa, Israel.

## Development

### Hot Reloading with Air

The project includes configuration for [Air](https://github.com/air-verse/air), a live reload tool for Go applications.

#### Installation

Install Air globally:

```bash
go install github.com/air-verse/air@latest
```

Make sure your Go bin directory is in your PATH.

#### Usage

Run the server with hot reloading:

- On Windows: `dev.bat` or `air -c .air.toml`
- On Mac/Linux: `./dev.sh` or `air -c .air.toml`

This will watch for file changes and automatically rebuild and restart the Go server.

### Backend API

The backend provides the following endpoints:

- `GET /api/stadium/sammyofer` - Get information about Sammy Ofer Stadium
- `GET /api/fotmob/sammyofer` - Get upcoming matches at Sammy Ofer Stadium
- `GET /api/refresh-token` - Manually refresh the Fotmob API token

### Frontend

The React frontend is in the `frontend` directory. To start it in development mode:

```bash
cd frontend
npm install
npm start
```

The frontend will be available at http://localhost:3000 and will proxy API requests to the Go backend at http://localhost:8000.

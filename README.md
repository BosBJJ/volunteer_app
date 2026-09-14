# Volunteer App

A REST API for managing volunteers and volunteer opportunities, built with Go and PostgreSQL.

## Requirements

- Go
- PostgreSQL

### Environment Variables

Create a `.env` file or otherwise configure the following environment variable:

`DATABASE_URL=postgres://user:password@localhost:5432/volunteer_app`

Replace `user` and `password` with your PostgreSQL credentials.

### Database 

Create a PostgreSQL database named `volunteer_app`.

The application creates the required tables when it starts.

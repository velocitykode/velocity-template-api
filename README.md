# Velocity API Template

API-only project template for the [Velocity](https://github.com/velocitykode/velocity) Go web framework.

## Stack

- **Backend**: Velocity Go Framework
- **API Format**: JSON REST API

## Usage

This template is used automatically by the Velocity CLI with the `--api` flag:

```bash
velocity new myapi --api
cd myapi
./vel serve
```

## Development Commands

```bash
# Start development server with hot reload
./vel serve

# Run database migrations
./vel migrate

# Generate a new controller
./vel make:controller User

# Build for production
./vel build
```

## API Endpoints

The template includes example API endpoints:

- `GET /api/health` - Health check
- `GET /api/users` - List users (requires auth)
- `POST /api/users` - Create user
- `GET /api/users/:id` - Get user (requires auth)
- `GET /api/me` - Get current user (requires auth)

## Documentation

Full documentation at **[velocity.velocitykode.com/docs](https://velocity.velocitykode.com/docs)**

## License

MIT

# Velocity Template - API

API-only starter template for the [Velocity](https://github.com/velocitykode/velocity) Go web framework. JSON REST, no frontend bundle.

This repo is a **template** consumed by the Velocity installer. To start a new project:

```bash
velocity new myapi --api
```

The installer clones this template, rewrites the module placeholders, installs dependencies, builds the project's `vel` binary, and runs the initial migrations.

## Stack

- **Backend**: Velocity Go framework
- **API format**: JSON REST
- **Auth**: JWT (`AUTH_JWT_*` env)

## Example Endpoints

The scaffolded project ships with these routes wired up:

- `GET /api/health` - health check
- `GET /api/users` - list users (auth required)
- `POST /api/users` - create user
- `GET /api/users/:id` - fetch user (auth required)
- `GET /api/me` - current user (auth required)

## Documentation

Full documentation at **[velocity.velocitykode.com/docs](https://velocity.velocitykode.com/docs)**

## Sibling Templates

- [`velocity-template-react`](https://github.com/velocitykode/velocity-template-react) - React 19 + Inertia.js
- [`velocity-template-vue`](https://github.com/velocitykode/velocity-template-vue) - Vue 3 + Inertia.js (with SSR)

## License

MIT

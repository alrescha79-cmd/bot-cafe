# Documentation Index

Welcome to the Bot Café documentation! 📚

## 🚀 Getting Started

New to the project? Start here:

1. **[Development Setup Guide](guides/development-setup.md)** - Complete setup untuk development environment
2. **[Quick Start (README)](../README.md)** - 3-step quick start guide

## 📖 Guides

Step-by-step tutorials dan how-to guides:

- **[ Development Setup](guides/development-setup.md)** - Comprehensive local development guide
  - Prerequisites
  - Initial setup
  - 3 development modes (local hot reload, Docker, no hot reload)
  - Development workflow
  - Debugging tips

- **[🚀 VPS Deployment](guides/vps-deployment.md)** - Production deployment ke VPS
  - VPS preparation
  - Docker installation
  - Systemd service setup
  - Backups & maintenance
  - Security best practices

- **[🐛 Troubleshooting](guides/troubleshooting.md)** - Solutions untuk masalah umum
  - Hot reload issues
  - Port conflicts
  - Bot not responding
  - Database issues
  - Docker problems

## 📚 Reference

Technical reference documentation:

- **[🏗️ Architecture](reference/architecture.md)** - System architecture & design
  - Microservices overview
  - Service responsibilities
  - Communication patterns
  - Database design

- **[📡 API Reference](reference/api.md)** - API documentation
  - Request/response format
  - Endpoint specifications
  - Service APIs

- **[🔧 Makefile Commands](reference/makefile-commands.md)** - Complete command reference
  - Development commands
  - Build commands
  - Docker commands
  - Utility commands

## 💡 Examples

Practical examples and workflows:

- **[👨‍💼 Admin Workflows](examples/admin-workflows.md)** - Admin feature usage examples
  - CRUD menu
  - CRUD promo
  - CRUD info café
  - Manage categories

## 🎯 Quick Links

### Development
- Start developing: `make dev-local-hot`
- Run with Docker: `make dev`
- Stop services: `make stop`

### Deployment
- Deploy to VPS: [VPS Deployment Guide](guides/vps-deployment.md)
- Setup systemd service: [Systemd Section](guides/vps-deployment.md#running-as-system-service)

### Help
- Common issues: [Troubleshooting](guides/troubleshooting.md)
- Hot reload not working: [Hot Reload Issues](guides/troubleshooting.md#hot-reload-issues)
- Bot not responding: [Bot Issues](guides/troubleshooting.md#bot-issues)

## 🤝 Contributing

Want to contribute? Check out:
1. [Development Setup](guides/development-setup.md) - Setup environment
2. [Architecture](reference/architecture.md) - Understand the system
3. [API Reference](reference/api.md) - API contracts

## 📞 Support

- **Issues**: Open GitHub issue
- **Questions**: Check [Troubleshooting Guide](guides/troubleshooting.md) first
- **Documentation**: This docs folder!

---

**Happy Learning!** 🎉

Start with the [Development Setup Guide](guides/development-setup.md) and you'll be coding in minutes!

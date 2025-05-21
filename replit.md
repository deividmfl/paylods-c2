# C2 Server Project Guide

## Overview

This project is a Command and Control (C2) server built with Flask, designed for managing and interacting with PowerShell clients on remote systems. The server provides a web-based admin dashboard for monitoring connected hosts, sending commands, and viewing logs. It's intended for educational and authorized security testing purposes.

## User Preferences

Preferred communication style: Simple, everyday language.

## System Architecture

The C2 server uses a client-server architecture where:

1. **Server**: A Flask application that provides API endpoints for client connections and an admin dashboard for management.
2. **Clients**: PowerShell scripts running on target systems that connect to the server and execute commands.
3. **Storage**: Currently uses in-memory storage for maintaining state, with the possibility to migrate to a database in the future.

The application uses basic HTTP for client-server communication, with authentication for the admin interface.

## Key Components

### Backend (Flask Application)

- **API Module** (`api.py`): Contains Flask routes for client reporting and command distribution.
- **Storage Module** (`storage.py`): Provides in-memory storage for host data, logs, commands, and configuration.
- **Authentication**: Uses basic HTTP authentication for admin access.

### Frontend

- **Admin Dashboard** (`templates/admin.html`): Web interface for monitoring and controlling connected hosts.
- **JavaScript** (`static/main.js`): Client-side code for the admin dashboard.
- **CSS** (`static/style.css`): Styling for the admin interface with a dark, terminal-like theme.

### Client

- **PowerShell Script**: Template script embedded in the server that can be customized and distributed to target systems.

## Data Flow

1. **Client Registration**:
   - PowerShell clients collect host information (hostname, username, IP, OS)
   - Clients periodically connect to the server to report status
   - Server stores this information in the in-memory storage

2. **Command Execution**:
   - Admin sends commands through the dashboard interface
   - Commands are stored in a queue for each host
   - Clients poll for new commands and execute them
   - Results are sent back to the server and displayed in the dashboard

3. **Logging**:
   - Client actions and errors are logged on the server
   - Logs are viewable through the admin interface

## External Dependencies

### Backend Dependencies
- Flask: Web framework for the server
- Flask-SQLAlchemy: ORM for potential database integration
- Gunicorn: WSGI HTTP server for production deployment
- Psycopg2-binary: PostgreSQL adapter (for potential database migration)
- Email-validator: For input validation

### Frontend Dependencies
- Bootstrap: UI framework (loaded from CDN)
- Bootstrap Icons: Icon library (loaded from CDN)

## Deployment Strategy

The application is configured for deployment on Replit with the following considerations:

1. **Server**: Uses Gunicorn as the WSGI server, binding to port 5000
2. **Environment**: Uses Python 3.11 with required packages
3. **Dependencies**: OpenSSL and PostgreSQL are included in the Nix configuration for future database integration
4. **Scaling**: Configured for autoscaling deployment

### Development Workflow

The project includes a configured workflow for starting the application:
- Run button executes "Project" workflow
- Gunicorn serves the application with auto-reload for development

### Future Enhancements

1. **Database Migration**: The current in-memory storage can be migrated to PostgreSQL
2. **Encryption**: Add TLS/SSL for secure client-server communication
3. **Authentication Improvements**: Enhance admin authentication beyond basic auth
4. **Client Features**: Expand PowerShell client capabilities

## Database Structure (Planned)

While currently using in-memory storage, the system is designed to be migrated to a PostgreSQL database with the following potential schema:

1. **Hosts**: Store connected host information
2. **Commands**: Track commands sent to each host
3. **Logs**: Store execution logs and client reports
4. **Errors**: Track error messages from clients
5. **Config**: Store server configuration

## Security Considerations

- Admin authentication is currently basic; should be enhanced for production
- Client-server communication should be encrypted in production environments
- Proper access controls should be implemented based on deployment context
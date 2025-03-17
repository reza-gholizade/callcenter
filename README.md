# AI-Powered Support Chatbot

An intelligent chatbot system for handling flight support inquiries, ticket management, and refund processing.

## Features

- Natural Language Processing for understanding user queries
- Multi-platform support (Web, Telegram, WhatsApp)
- Ticket lookup and management
- Ticket cancellation and refund processing
- Baggage policy information
- Admin dashboard for monitoring and management

## Technology Stack

- Backend: Go (Gin) + PostgreSQL
- NLP Engine: Dialogflow/Rasa/LLM
- Frontend: React
- External Integrations:
  - MyTicket API
  - Faranegar API
  - Payment Gateway API
  - SMS/Email Services
  - Telegram/WhatsApp APIs

## Prerequisites

- Go 1.16 or later
- PostgreSQL 12 or later
- Node.js 14 or later (for frontend)
- Docker (optional)

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/callcenter.git
   cd callcenter
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Copy the environment file and update the values:
   ```bash
   cp .env.example .env
   ```

4. Set up the database:
   ```bash
   createdb callcenter
   ```

5. Run database migrations:
   ```bash
   go run cmd/migrate/main.go
   ```

6. Start the server:
   ```bash
   go run cmd/api/main.go
   ```

## Configuration

The application is configured through environment variables. See `.env.example` for all available options.

Required environment variables:
- `DB_PASSWORD`: Database password
- `JWT_SECRET`: Secret key for JWT token generation

## API Endpoints

### Chat Endpoints

- `POST /api/v1/chat/message`: Send a message to the chatbot
- `GET /api/v1/chat/history/:sessionId`: Get chat history for a session

### Ticket Endpoints

- `GET /api/v1/tickets/:ticketNumber`: Get ticket details
- `POST /api/v1/tickets/:ticketNumber/cancel`: Cancel a ticket
- `GET /api/v1/tickets/:ticketNumber/refund-status`: Get refund status

## Development

### Project Structure

```
.
├── cmd/
│   ├── api/        # Main application entry point
│   └── migrate/    # Database migration tool
├── internal/
│   ├── config/     # Configuration management
│   ├── handlers/   # HTTP request handlers
│   ├── models/     # Database models
│   └── services/   # Business logic
├── pkg/            # Reusable packages
└── web/           # Frontend application
```

### Running Tests

```bash
go test ./...
```

### Code Style

The project follows the standard Go code style. Run the following command to format code:

```bash
go fmt ./...
```

## Deployment

### Docker

Build the Docker image:
```bash
docker build -t callcenter .
```

Run the container:
```bash
docker run -p 8080:8080 --env-file .env callcenter
```

### Kubernetes

See the `k8s/` directory for Kubernetes deployment manifests.

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Support

For support, please contact the development team or create an issue in the repository. 
# Neptune Backend

## Environment Setup

Create a `.env` file in the root directory with the following variables:

```env
# Server Configuration
PORT=8080

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=neptune_db

# Messier API Configuration (Binus Authentication Service)
MESSIER_API_URL=https://bluejack.binus.ac.id/lapi/api

# JWT Configuration
JWT_SECRET=your-secret-key-here

# Environment
ENV=development
```

## Important Notes

1. **MESSIER_API_URL**: This must point to the correct Binus authentication service URL
2. **Database**: Make sure PostgreSQL is running and the database exists
3. **Timeout**: The HTTP client timeout has been increased to 30 seconds to handle slow responses from the Messier API

## Running the Application

1. Install dependencies: `go mod tidy`
2. Set up your `.env` file
3. Run the application: `go run main.go`

## Troubleshooting

If you encounter timeout errors:

1. Check that the MESSIER_API_URL is correct
2. Verify network connectivity to the Binus servers
3. Check if the Binus authentication service is available

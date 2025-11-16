# RBAC & ReBAC API with Go, GORM, and PostgreSQL

A demonstration API showcasing both Role-Based Access Control (RBAC) and Relationship-Based Access Control (ReBAC) patterns using Go, Gin, GORM, and PostgreSQL.

## Features

- **RBAC (Role-Based Access Control)**: Users have roles, roles have permissions
- **ReBAC (Relationship-Based Access Control)**: Access control based on relationships between entities
- JWT-based authentication
- RESTful API with proper error handling
- Database migrations with GORM
- CORS support
- Health check endpoint

## Prerequisites

### For Docker Setup
- Docker 20.10 or higher
- Docker Compose 2.0 or higher

### For Local Setup
- Go 1.21 or higher
- PostgreSQL 15 or higher
- Make (optional, for convenience commands)

---

## 🐳 Docker Setup (Recommended)

The easiest way to run the application is using Docker Compose. This will set up both the database and the application in containers.

### Step 1: Clone and Navigate

```bash
cd go-rebac-rbac-postgres
```

### Step 2: Build and Start Services

```bash
# Build and start all services (database + app)
docker-compose up -d --build

# View logs
docker-compose logs -f app

# Or view logs for both services
docker-compose logs -f
```

The application will be available at `http://localhost:4000`

### Step 3: Seed the Database

```bash
# Seed the database with sample data
curl -X POST http://localhost:4000/seed
```

Or using Docker:

```bash
docker-compose exec app curl -X POST http://localhost:4000/seed
```

### Step 4: Verify Installation

```bash
# Check health endpoint
curl http://localhost:4000/health
```

### Docker Commands

```bash
# Start services
docker-compose up -d

# Stop services
docker-compose down

# Stop and remove volumes (clears database)
docker-compose down -v

# View logs
docker-compose logs -f app

# Rebuild after code changes
docker-compose up -d --build

# Execute commands in app container
docker-compose exec app <command>

# Access database directly
docker-compose exec db psql -U postgres -d rebac_rbac_db
```

### Docker Environment Variables

The Docker setup uses these default environment variables (defined in `docker-compose.yml`):

- `DATABASE_URL`: Automatically configured to connect to the `db` service
- `PORT`: 4000
- `GIN_MODE`: release
- `JWT_SECRET`: docker-secret-key-change-in-production

To customize, edit `docker-compose.yml` or use a `.env` file.

---

## 💻 Local Development Setup

### Step 1: Install Dependencies

```bash
# Install Go dependencies
go mod download
```

### Step 2: Start PostgreSQL Database

**Option A: Using Docker (Database Only)**

```bash
# Start only the database container
docker-compose up -d db

# Wait for database to be ready (usually 5-10 seconds)
sleep 5
```

**Option B: Using Local PostgreSQL**

Make sure PostgreSQL is running locally and create the database:

```bash
createdb rebac_rbac_db
# Or using psql:
psql -U postgres -c "CREATE DATABASE rebac_rbac_db;"
```

**Note:** The application will automatically create the database if it doesn't exist when you run it locally (requires PostgreSQL to be running and accessible).

### Step 3: Configure Environment Variables

```bash
# Copy example environment file
cp .env.example .env

# Edit .env file with your database credentials
# For Docker database:
DATABASE_URL=host=localhost user=postgres password=password dbname=rebac_rbac_db port=5432 sslmode=disable

# For local PostgreSQL, update accordingly:
# DATABASE_URL=host=localhost user=youruser password=yourpassword dbname=rebac_rbac_db port=5432 sslmode=disable
```

Or export environment variables directly:

```bash
export DATABASE_URL="host=localhost user=postgres password=password dbname=rebac_rbac_db port=5432 sslmode=disable"
export PORT=4000
export JWT_SECRET="your-secret-key-change-in-production"
```

### Step 4: Run the Application

```bash
# Run directly with go run
go run .

# Or build and run
go build -o app .
./app
```

The server will start on port 4000 (or the port specified in `PORT` environment variable).

### Step 5: Seed the Database

In a new terminal:

```bash
curl -X POST http://localhost:4000/seed
```

This creates:
- **Permissions**: `record:read`, `record:write`
- **Roles**: `doctor`, `nurse`, `admin`
- **Users**:
  - `doc@example.com` (doctor role) - password: `password`
  - `nurse@example.com` (nurse role) - password: `password`
  - `patient@example.com` (no role) - password: `password`
  - `outsider@example.com` (no role) - password: `password`
- **Patient Records**: 2 records owned by the patient
- **Relationships**: Doctor assigned to patient and record

### Step 6: Verify Installation

```bash
# Check health endpoint
curl http://localhost:4000/health

# Expected response:
# {"status":"ok","service":"rbac-rebac-api"}
```

---

## 📮 Postman Collection

A complete Postman collection is included in `postman_collection.json` with all API endpoints pre-configured.

### Import into Postman

1. Open Postman
2. Click **Import** button
3. Select `postman_collection.json` file
4. The collection will be imported with:
   - All endpoints configured
   - Environment variables set up
   - Test scripts to automatically save tokens

### Postman Environment Variables

The collection uses these variables:
- `base_url`: `http://localhost:4000` (default)
- `doctor_token`: Auto-populated after login
- `nurse_token`: Auto-populated after login
- `patient_token`: Auto-populated after login
- `outsider_token`: Auto-populated after login

### Using the Collection

1. **First, seed the database**: Run the "Seed Database" request
2. **Login as different users**: Use the login requests in the "Authentication" folder
3. **Test RBAC**: Try accessing records with different user tokens
4. **Test ReBAC**: Try accessing records with different user tokens

The collection includes test scripts that automatically save tokens to environment variables after successful login.

## API Endpoints

### Authentication

#### Register
```bash
POST /auth/register
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123",
  "name": "John Doe"
}
```

#### Login
```bash
POST /auth/login
Content-Type: application/json

{
  "email": "doc@example.com",
  "password": "password"
}
```

Response:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "email": "doc@example.com",
    "name": "Dr Alice"
  }
}
```

### Protected Endpoints

All protected endpoints require an `Authorization` header:
```
Authorization: Bearer <token>
```

#### Get Record (RBAC)
```bash
GET /records/rbac/:id
Authorization: Bearer <token>
```

Requires `record:read` permission through user's roles.

#### Get Record (ReBAC)
```bash
GET /records/rebac/:id
Authorization: Bearer <token>
```

Requires a relationship (`assigned_to`) between the user and the record or record owner.

### Utility Endpoints

#### Health Check
```bash
GET /health
```

#### Seed Database
```bash
POST /seed
```

## Access Control Examples

### RBAC Example

1. Login as doctor (has `record:read` permission):
```bash
curl -X POST http://localhost:4000/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"doc@example.com","password":"password"}'
```

2. Access record using RBAC:
```bash
curl http://localhost:4000/records/rbac/1 \
  -H "Authorization: Bearer <token>"
```

### ReBAC Example

1. Login as doctor (has relationship with patient):
```bash
curl -X POST http://localhost:4000/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"doc@example.com","password":"password"}'
```

2. Access record using ReBAC:
```bash
curl http://localhost:4000/records/rebac/1 \
  -H "Authorization: Bearer <token>"
```

The doctor can access because:
- Direct relationship: doctor → record (assigned_to)
- Indirect relationship: doctor → patient (assigned_to), and patient owns the record
- Owner access: if the user owns the record, they can always access it

## Architecture

### Models

- **Permission**: Defines what actions can be performed
- **Role**: Groups permissions together
- **User**: System users with email/password authentication
- **PatientRecord**: Example resource that needs access control
- **Relationship**: Defines relationships between entities (for ReBAC)

### Access Control Flow

**RBAC Flow:**
1. User authenticates and receives JWT token
2. Token is validated in `RequireAuth()` middleware
3. `RBAC()` middleware checks if user's roles have the required permission
4. Request proceeds if permission exists

**ReBAC Flow:**
1. User authenticates and receives JWT token
2. Token is validated in `RequireAuth()` middleware
3. `ReBACResource()` middleware checks for relationships:
   - Direct: user → resource
   - Indirect: user → owner (where owner owns resource)
   - Owner: user owns the resource
4. Request proceeds if any relationship exists

## Environment Variables

- `DATABASE_URL`: PostgreSQL connection string
- `PORT`: Server port (default: 4000)
- `GIN_MODE`: Gin mode (debug/release/test)
- `JWT_SECRET`: Secret key for JWT signing

## Development

### Project Structure

```
.
├── main.go                  # Application entry point and routes
├── db.go                    # Database connection and migration
├── models.go                # GORM models
├── handlers.go             # HTTP request handlers
├── middleware.go            # Authentication and authorization middleware
├── auth.go                  # Authentication utilities (JWT, password hashing)
├── seed.go                  # Database seeding function
├── Dockerfile               # Docker image definition
├── docker-compose.yml       # Docker Compose configuration
├── .dockerignore           # Docker ignore patterns
├── .env.example            # Environment variables template
├── postman_collection.json  # Postman API collection
├── go.mod                   # Go module dependencies
├── go.sum                   # Go module checksums
└── README.md                # This file
```

### Testing the API

#### Using Postman (Recommended)

1. Import `postman_collection.json` into Postman
2. Run "Seed Database" request first
3. Use the pre-configured requests to test all endpoints

#### Using cURL

```bash
# 1. Seed the database
curl -X POST http://localhost:4000/seed

# 2. Register a new user
curl -X POST http://localhost:4000/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123","name":"Test User"}'

# 3. Login as doctor
TOKEN=$(curl -s -X POST http://localhost:4000/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"doc@example.com","password":"password"}' | jq -r '.token')

# 4. Access protected endpoint (RBAC)
curl http://localhost:4000/records/rbac/1 \
  -H "Authorization: Bearer $TOKEN"

# 5. Access protected endpoint (ReBAC)
curl http://localhost:4000/records/rebac/1 \
  -H "Authorization: Bearer $TOKEN"
```

#### Using HTTPie

```bash
# Seed database
http POST http://localhost:4000/seed

# Login
http POST http://localhost:4000/auth/login email=doc@example.com password=password

# Access record (replace TOKEN with actual token)
http GET http://localhost:4000/records/rbac/1 Authorization:"Bearer TOKEN"
```

---

## 🔧 Troubleshooting

### Docker Issues

**Port already in use:**
```bash
# Check what's using port 4000
lsof -i :4000

# Or change port in docker-compose.yml
ports:
  - "4001:4000"  # Use 4001 instead
```

**Database connection errors:**
```bash
# Check if database container is running
docker-compose ps

# Check database logs
docker-compose logs db

# Restart services
docker-compose restart
```

**App won't start:**
```bash
# Check app logs
docker-compose logs app

# Rebuild containers
docker-compose up -d --build --force-recreate
```

### Local Development Issues

**Database connection refused:**
- Ensure PostgreSQL is running: `pg_isready` or `docker-compose ps`
- Check `DATABASE_URL` environment variable
- Verify database exists: `psql -U postgres -l`

**Port already in use:**
```bash
# Change PORT environment variable
export PORT=4001
go run .
```

**Migration errors:**
- Drop and recreate database: `dropdb rebac_rbac_db && createdb rebac_rbac_db`
- Or truncate tables manually
- The app will automatically create the database if it doesn't exist (when PostgreSQL is accessible)

**Go module errors:**
```bash
# Clean module cache
go clean -modcache

# Download dependencies again
go mod download

# Verify
go mod verify
```

---

## 🚀 Production Deployment

### Docker Production Setup

1. **Update environment variables** in `docker-compose.yml`:
   - Set strong `JWT_SECRET`
   - Use production database credentials
   - Set `GIN_MODE=release`

2. **Use Docker secrets** or environment files:
   ```bash
   docker-compose --env-file .env.production up -d
   ```

3. **Use reverse proxy** (nginx/traefik) for HTTPS

4. **Remove seed endpoint** or protect it with authentication

### Environment Variables for Production

```bash
DATABASE_URL=<production-database-url>
PORT=4000
GIN_MODE=release
JWT_SECRET=<strong-random-secret-key>
```

---

## 📚 Additional Resources

- [Gin Framework Documentation](https://gin-gonic.com/docs/)
- [GORM Documentation](https://gorm.io/docs/)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [JWT.io](https://jwt.io/) - JWT Debugger

---

## Security Notes

- ✅ Change `JWT_SECRET` in production
- ✅ Use strong passwords in production
- ⚠️ Consider rate limiting for authentication endpoints
- ⚠️ Protect `/seed` endpoint in production (remove or add auth)
- ✅ Use HTTPS in production
- ✅ Validate and sanitize all inputs
- ✅ Consider adding request logging and monitoring
- ✅ Use environment variables for sensitive data
- ✅ Keep dependencies updated: `go get -u ./...`

---

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test thoroughly
5. Submit a pull request

---

## 📝 License

MIT License - feel free to use this project for learning and development purposes.

---

## 🙏 Acknowledgments

- Built with [Gin](https://gin-gonic.com/) web framework
- Database ORM: [GORM](https://gorm.io/)
- Authentication: [JWT](https://jwt.io/)
- Database: [PostgreSQL](https://www.postgresql.org/)


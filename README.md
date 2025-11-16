# Social Media Feed

A full-featured social media feed application built with Go, featuring posts, comments, hashtags, notifications, and real-time updates. The application follows clean architecture principles and implements several design patterns for maintainability and scalability.

## Features

### Core Functionality
- **User Management**: Registration, authentication, and role-based access control (user, moderator, admin)
- **Posts**: Create, read, update, and delete posts with media support
- **Comments**: Threaded comments with reply functionality
- **Hashtags**: Automatic hashtag extraction and trending hashtag tracking
- **Notifications**: Real-time notifications for likes, comments, mentions, and replies
- **Feed**: Multiple feed sorting strategies (date, likes, engagement, trending, controversial, random)
- **Post Filters**: Decorative filters for posts (emoji overlay, glitter, frames, uppercase)
- **Likes/Dislikes**: Post engagement tracking

### Advanced Features
- **Observer Pattern**: Notification system with extensible observers
- **Strategy Pattern**: Flexible post sorting algorithms
- **Decorator Pattern**: Post content filtering and decoration
- **Adapter Pattern**: External post integration support
- **Facade Pattern**: Unified API interface
- **WebSocket Support**: Real-time updates (prepared infrastructure)

## Tech Stack

- **Language**: Go 1.25.4
- **Database**: SQLite3
- **Web Framework**: Standard `net/http` with custom routing
- **Templating**: Go `html/template`
- **Authentication**: Cookie-based session management
- **Logging**: Custom logger with file and console output
- **Password Hashing**: `golang.org/x/crypto/bcrypt`

## Project Structure

```
socialMediaFeed/
├── cmd/
│   └── app/
│       └── main.go                 # Application entry point
├── internal/
│   ├── api/                        # API layer
│   │   ├── adapter.go              # External post adapter
│   │   └── facade.go               # API facade pattern
│   ├── comment/                    # Comment domain
│   │   ├── comment.go              # Comment model
│   │   ├── handler.go              # HTTP handlers
│   │   ├── repository.go           # Repository interface
│   │   └── service.go              # Business logic
│   ├── hashtag/                    # Hashtag domain
│   │   ├── hashtag.go              # Hashtag model
│   │   ├── handler.go              # HTTP handlers
│   │   ├── repository.go           # Repository interface
│   │   └── service.go              # Business logic
│   ├── infrastructure/
│   │   ├── database/
│   │   │   ├── connection.go       # DB connection
│   │   │   ├── database.go         # DB utilities
│   │   │   └── migrations/
│   │   │       └── 0001_initial_schema.sql
│   │   └── repository/             # Repository implementations
│   │       ├── comment_repository.go
│   │       ├── hashtag_repository.go
│   │       ├── notification_repository.go
│   │       ├── post_repository.go
│   │       └── user_repository.go
│   ├── notification/               # Notification domain
│   │   ├── notification.go         # Notification model
│   │   ├── handler.go              # HTTP handlers
│   │   ├── observer.go             # Observer pattern
│   │   ├── repository.go           # Repository interface
│   │   └── service.go              # Business logic
│   ├── post/                       # Post domain
│   │   ├── post.go                 # Post model
│   │   ├── handler.go              # HTTP handlers
│   │   ├── service.go              # Business logic
│   │   ├── repository.go           # Repository interface
│   │   ├── strategy.go             # Sorting strategies
│   │   └── decorator.go            # Post decorators
│   ├── user/                       # User domain
│   │   ├── user.go                 # User model
│   │   ├── handler.go              # HTTP handlers
│   │   ├── repository.go           # Repository interface
│   │   └── service.go              # Business logic
│   └── web/                        # Web handlers
│       ├── auth.go                 # Authentication middleware
│       └── handler.go              # Web page handlers
├── pkg/                            # Shared packages
│   ├── logger/                     # Logging utilities
│   │   ├── logger.go
│   │   └── middleware.go
│   ├── responce/                   # Response utilities
│   │   └── responce.go
│   ├── types/                      # Shared types
│   │   └── feed.go                 # Feed item types
│   └── validator/                  # Validation utilities
│       └── validator.go
├── web/                            # Frontend assets
│   ├── static/
│   │   ├── css/
│   │   ├── images/
│   │   └── js/
│   └── templates/                  # HTML templates
│       ├── components/
│       ├── layout/
│       └── pages/
├── data/                           # Database files
├── logs/                           # Log files
├── go.mod
├── go.sum
└── README.md
```

## Architecture

The application follows **Clean Architecture** principles with clear separation of concerns:

1. **Domain Layer** (`internal/*/`): Core business logic and models
2. **Application Layer** (`internal/*/service.go`): Use cases and business rules
3. **Infrastructure Layer** (`internal/infrastructure/`): Database, external services
4. **Presentation Layer** (`internal/*/handler.go`, `internal/web/`): HTTP handlers and web pages
5. **API Layer** (`internal/api/`): Unified API facade

### Design Patterns

- **Repository Pattern**: Data access abstraction
- **Service Layer Pattern**: Business logic encapsulation
- **Facade Pattern**: Simplified API interface
- **Strategy Pattern**: Flexible sorting algorithms
- **Decorator Pattern**: Post content filtering
- **Observer Pattern**: Notification system
- **Adapter Pattern**: External system integration

## Getting Started

### Prerequisites

- Go 1.25.4 or higher
- SQLite3 (usually included with Go)

### Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd socialMediaFeed
```

2. Install dependencies:
```bash
go mod download
```

3. Set environment variables (optional):
```bash
export PORT=8080
export DB_PATH=data/app.db
export LOG_LEVEL=INFO
```

4. Run the application:
```bash
go run cmd/app/main.go
```

The application will:
- Create the database directory if it doesn't exist
- Run database migrations automatically
- Start the server on `http://localhost:8080` (or the configured PORT)

### Database

The application uses SQLite3 and automatically runs migrations on startup. The database file is created at `data/app.db` (configurable via `DB_PATH` environment variable).

## API Endpoints

### Health Check
- `GET /health` - Health check endpoint

### Authentication
- `POST /api/register` - Register a new user
- `POST /api/login` - Login user
- `POST /logout` - Logout user

### Users
- `GET /api/users` - Get all users
- `GET /api/users/{id}` - Get user by ID
- `PUT /api/users/{id}` - Update user
- `DELETE /api/users/{id}` - Delete user

### Posts
- `POST /api/posts` - Create a new post
- `GET /api/posts` - Get all posts
- `GET /api/posts/{id}` - Get post by ID
- `PUT /api/posts/{id}` - Update post
- `DELETE /api/posts/{id}` - Delete post
- `GET /api/feed` - Get feed (supports `sort` query parameter)
- `GET /api/trending` - Get trending posts
- `GET /api/users/{authorId}/posts` - Get posts by author
- `GET /api/hashtags/{tag}/posts` - Get posts by hashtag
- `POST /api/posts/{id}/like` - Like a post
- `POST /api/posts/{id}/dislike` - Dislike a post
- `POST /api/posts/{id}/filters` - Apply filters to a post

### Comments
- `POST /api/comments` - Create a comment
- `GET /api/comments` - Get all comments
- `GET /api/comments/{id}` - Get comment by ID
- `GET /api/posts/{postId}/comments` - Get comments for a post
- `PUT /api/comments/{id}` - Update comment
- `DELETE /api/comments/{id}` - Delete comment
- `POST /api/comments/{id}/reply` - Reply to a comment

### Hashtags
- `GET /api/hashtags` - Get all hashtags
- `GET /api/hashtags/{tag}` - Get hashtag by tag
- `GET /api/hashtags/trending` - Get trending hashtags

### Notifications
- `GET /api/notifications` - Get user notifications
- `GET /api/notifications/{id}` - Get notification by ID
- `PUT /api/notifications/{id}/read` - Mark notification as read
- `DELETE /api/notifications/{id}` - Delete notification

### Web Pages
- `GET /` - Home page (feed)
- `GET /login` - Login page
- `GET /register` - Registration page
- `GET /post/{id}` - Post detail page
- `GET /profile/{id}` - User profile page
- `GET /profile` - Current user's profile
- `GET /create-post` - Create post page

## Feed Sorting Strategies

The feed supports multiple sorting strategies via the `sort` query parameter:

- `date` or `newest` - Sort by creation date (newest first)
- `likes` or `popular` - Sort by number of likes
- `engagement` - Sort by total engagement (likes + dislikes)
- `trending` or `hot` - Sort by trending score (engagement with time decay)
- `controversial` - Sort by controversy score (balanced likes/dislikes)
- `random` - Random order

Example: `GET /api/feed?sort=trending`

## Post Filters

Posts can be decorated with various filters:

- **Emoji Overlay**: Add emoji prefix to content
- **Glitter**: Add sparkle decorations
- **Frames**: Add decorative frames (stars, hearts, brackets)
- **Uppercase**: Convert content to uppercase

Filters can be applied via the `POST /api/posts/{id}/filters` endpoint.

## Database Schema

### Users
- `id` (INTEGER PRIMARY KEY)
- `username` (TEXT UNIQUE)
- `email` (TEXT UNIQUE)
- `password_hash` (TEXT)
- `role` (TEXT DEFAULT 'user')
- `created_at` (DATETIME)
- `updated_at` (DATETIME)

### Posts
- `id` (INTEGER PRIMARY KEY)
- `author_id` (INTEGER, FOREIGN KEY)
- `content` (TEXT)
- `media_url` (TEXT)
- `likes` (INTEGER DEFAULT 0)
- `dislikes` (INTEGER DEFAULT 0)
- `created_at` (DATETIME)
- `updated_at` (DATETIME)

### Comments
- `id` (INTEGER PRIMARY KEY)
- `post_id` (INTEGER, FOREIGN KEY)
- `user_id` (INTEGER, FOREIGN KEY)
- `parent_comment_id` (INTEGER, FOREIGN KEY, nullable)
- `content` (TEXT)
- `created_at` (DATETIME)
- `updated_at` (DATETIME)

### Hashtags
- `id` (INTEGER PRIMARY KEY)
- `tag` (TEXT UNIQUE)
- `usage_count` (INTEGER DEFAULT 0)
- `created_at` (DATETIME)
- `updated_at` (DATETIME)

### Post-Hashtag Relations
- `post_id` (INTEGER, FOREIGN KEY)
- `hashtag_id` (INTEGER, FOREIGN KEY)
- PRIMARY KEY (post_id, hashtag_id)

### Notifications
- `id` (INTEGER PRIMARY KEY)
- `user_id` (INTEGER, FOREIGN KEY)
- `type` (TEXT)
- `title` (TEXT)
- `message` (TEXT)
- `is_read` (BOOLEAN DEFAULT FALSE)
- `related_entity_id` (INTEGER, nullable)
- `related_entity_type` (TEXT, nullable)
- `created_at` (DATETIME)

## Configuration

### Environment Variables

- `PORT` - Server port (default: `8080`)
- `DB_PATH` - Database file path (default: `data/app.db`)
- `LOG_LEVEL` - Logging level: DEBUG, INFO, WARNING, ERROR, FATAL (default: `INFO`)

## Logging

The application includes a custom logging system with:
- Multiple log levels (DEBUG, INFO, WARNING, ERROR, FATAL)
- File logging to `logs/app.log`
- Console output with color coding
- HTTP request logging middleware

## Development

### Running Tests
```bash
go test ./...
```

### Building
```bash
go build -o bin/app cmd/app/main.go
```

### Code Structure Guidelines

- **Models**: Domain entities with business logic methods
- **Services**: Business logic and use cases
- **Handlers**: HTTP request/response handling
- **Repositories**: Data access layer
- **Infrastructure**: External dependencies (database, etc.)

## Security Features

- Password hashing using bcrypt
- Cookie-based session management
- Role-based access control (RBAC)
- SQL injection prevention via parameterized queries
- CSRF protection (via authentication middleware)

## Future Enhancements

- WebSocket real-time updates
- Image upload and storage
- User following/followers system
- Advanced search functionality
- Rate limiting
- API rate limiting
- Docker containerization
- Unit and integration tests

## License

[Specify your license here]

## Contributing

[Contributing guidelines if applicable]


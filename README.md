Social Media Feed API: Design Patterns Demonstration

This project is a modular backend application developed in Go to demonstrate the application of core Software Design Patterns within the context of a simplified Social Media Feed architecture.

The API successfully implements the Factory, Builder, Strategy, and Adapter patterns to manage user creation, post construction, interaction handling, and third-party content integration.

🎯 Design Patterns Implemented

The following Go code implements the required patterns:

Pattern

Goal

Implementation in API

1. Factory

User Creation by Role

The NewUserFactory function creates distinct User objects (Admin, Standard) based on a provided role, centralizing object instantiation logic.

2. Builder

Complex Post Construction

The PostBuilder interface and its concrete implementation allow for step-by-step creation of complex Post objects, ensuring a valid and complete structure before being saved.

3. Strategy

Interaction Behavior

The InteractionStrategy interface allows the system to define a family of algorithms (e.g., LikeStrategy, CommentStrategy) and select one dynamically at runtime to execute interactions on a post.

4. Adapter

External Content Integration

The ExternalPostAdapter takes an incompatible external data structure (ExternalPost) and wraps it to conform to our internal FeedItemTarget interface, allowing foreign content to be displayed seamlessly in the main feed.

🏗️ Architecture and File Structure

The application is structured into a modular package, where each design pattern is contained within its own file for clarity and maintainability.

File

Primary Pattern

Description

user_factory.go

Factory

Defines the User struct and the logic for creating users based on roles.

post_builder.go

Builder

Defines the complex Post struct and the interfaces/implementation for its step-by-step construction.

strategy.go

Strategy

Defines the InteractionStrategy interface and the concrete strategies (Like, Dislike, Comment).

adapter.go

Adapter

Defines the internal FeedItem target and the adapter logic for integrating the foreign ExternalPost struct.

main.go

Orchestration

Sets up the HTTP server, defines routes, and orchestrates the calls to the various design pattern components within the request handlers.

🚀 Getting Started (Running the API)

This project requires Go 1.18+ to run.

Initialize the Module: Navigate to the project directory in your terminal and initialize the Go module.

go mod init socialmediafeed


Run the Server: Execute all the source files in the current directory.

go run .


The server will start and be accessible at: http://localhost:8080

🧪 API Usage Examples

Use tools like curl, Postman, or Thunder Client to interact with the API endpoints and observe how each pattern is executed.

1. Factory Pattern: Create a User (POST /users)

Creates a new user object based on the requested role.

# Request Body: {"name": "Test User", "role": "Standard"}
curl -X POST http://localhost:8080/users \
     -H "Content-Type: application/json" \
     -d '{"name": "Alice Developer", "role": "Verified"}'


Example Response (Status 201):

{
    "id": "u-1678888888888",
    "name": "Alice Developer",
    "role": "Verified"
}


2. Builder Pattern: Create a Post (POST /posts)

Assembles a new post using a user ID and content fields. (Use an ID returned from the /users endpoint, e.g., u-1678888888888).

# Request Body: {"authorId": "...", "content": "...", "mediaURL": "..."}
curl -X POST http://localhost:8080/posts \
     -H "Content-Type: application/json" \
     -d '{"authorId": "u-1678888888888", "content": "Demonstrating the Post Builder!", "mediaURL": ""}'


Example Response (Status 201):

{
    "id": "post-987654321",
    "author": { ... },
    "content": "Demonstrating the Post Builder!",
    "createdAt": "2025-11-12T18:00:00Z",
    "likes": 0,
    "comments": []
}


3. Strategy Pattern: Interact with a Post (POST /posts/{id}/interact)

Chooses the appropriate strategy (like, dislike, or comment) to execute based on the action field. (Use a Post ID returned from the /posts endpoint, e.g., post-987654321).

A. Like Strategy:

# Request Body: {"action": "like"}
curl -X POST http://localhost:8080/posts/post-987654321/interact \
     -H "Content-Type: application/json" \
     -d '{"action": "like"}'


B. Comment Strategy:

# Request Body: {"action": "comment", "data": "New comment text"}
curl -X POST http://localhost:8080/posts/post-987654321/interact \
     -H "Content-Type: application/json" \
     -d '{"action": "comment", "data": "Great implementation!"}'


Example Response (Comment Strategy - Status 200):

{
    "comments": [
        "Great implementation!"
    ],
    "message": "comment action successful on post post-987654321",
    "new_likes": 1
}


4. Adapter Pattern: Get the Feed (GET /feed)

Retrieves the aggregated feed, which includes internal posts and data adapted from an external source.

curl http://localhost:8080/feed


Example Snippet (Status 200):
The last item in the array demonstrates the adapter converting an external format (Handle, TweetContent) into the internal FeedItem structure.

[
    ... (Internal Posts) ...,
    {
        "type": "External",
        "author": "@GeminiAI (External Source)",
        "content": "Design patterns simplify complex systems, especially in Go!",
        "metrics": "999 Views",
        "postId": ""
    }
]

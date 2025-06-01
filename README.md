# Property Listing System

![Property Listing](PropertListing.png)

A RESTful property listing platform built with Go, Gin, MongoDB, and Redis.
This system allows users to register, list properties, search/filter listings, manage favorites, and recommend properties to others.

---

## Features

- **User Registration & Login**
- **Property CRUD** (Create, Read, Update, Delete)
- **Advanced Property Search** with filters
- **Favorites**: Mark/unmark properties as favorites
- **Recommendations**: Recommend properties to other users
- **Caching**: Redis-backed caching for performance

---

## API Endpoints

### Public Endpoints

| Method | Path             | Description                       |
|--------|------------------|-----------------------------------|
| GET    | `/`              | Health check/Hello                |
| GET    | `/properties`    | List all properties               |
| GET    | `/prop-search`   | Search/filter properties          |
| GET    | `/users`         | List all users (filter by email)  |
| POST   | `/register`      | Register a new user               |
| POST   | `/login`         | Login and receive a JWT           |

### Protected Endpoints (`/auth` prefix, JWT required)

| Method | Path             | Description                       |
|--------|------------------|-----------------------------------|
| GET    | `/auth/my-props` | List properties created by user   |
| POST   | `/auth/add-prop` | Add a new property                |
| PUT    | `/auth/update-prop` | Update a property you own      |
| GET    | `/auth/fav`      | Get user's favorite properties    |
| POST   | `/auth/fav-prop` | Mark/unmark a property as favorite|
| GET    | `/auth/my-rec`   | Get recommendations for user      |
| POST   | `/auth/rec-prop` | Recommend a property to a user    |

---


## Getting Started

1. **Clone the repository**
2. **Set up MongoDB and Redis**  
   Create a `.env` file in the project root with the following format:
   ```env
   MONGODB_URI=mongodb+srv://<username>:<password>@<cluster-url>/<params>
   DB_NAME=property-db

   REDIS_URL=redis://default:<password>@<host>:<port>
   ```
   - **MONGODB_URI**: Your MongoDB connection string.
   - **DB_NAME**: The MongoDB database name.
   - **REDIS_URL**: Your Redis connection string.
3. **Run the server:**
   ```bash
   go run cmd/main.go
   ```

# AppDrop - Mini App Config API (Backend)

This is a high-performance REST API built in **Go** and **PostgreSQL** for managing dynamic mobile app configurations. It allows developers to define screens (Pages) and UI components (Widgets) that can be instantly updated in a mobile application without a new app release.

---

## 🚀 Key Backend Features (Assignment Ready)

### 1. Robust Core API
- **Pages CRUD**: Create, Read, Update, and Delete mobile app screens.
- **Widgets CRUD**: Add, edit, remove, and reorder UI components on any page.
- **Home Page Logic**: Intelligent handling of the "Home" page. Only one home page can exist. Setting a new one automatically unsets the old one.
- **Cascade Deletion**: Deleting a page automatically cleans up all associated widgets.

### 2. Advanced Data Strategy
- **JSONB Implementation**: Uses PostgreSQL `JSONB` for widget configurations. This allows for total flexibility—a `banner` can have different fields than a `product_grid` without database schema changes.
- **Strict Validation**: Type-safe widget types (`banner`, `product_grid`, `text`, `image`, `spacer`) and unique route enforcement.

### 3. ✨ Bonus Features (Implemented)
- **✅ Pagination**: `GET /pages` supports `page` and `per_page` query parameters for optimized data fetching.
- **✅ Widget Filtering**: `GET /pages/:id/widgets` allows filtering by type (e.g., `?type=banner`).
- **✅ Reordering Logic**: Dedicated endpoint to batch reorder widgets using a transaction for data integrity.
- **✅ Request Logging**: Custom middleware that logs method, path, and latency for every request.
- **✅ Unit Testing**: Comprehensive tests for model validation and business logic.

---

## 🛠️ Tech Stack
- **Go 1.21+** (Gin Framework)
- **PostgreSQL 12+** (with UUID support)
- **Middleware**: Logging, CORS, and Panic Recovery.

---

## ⚙️ How to Setup & Run

### 1. Database Setup (PostgreSQL)
Ensure you have PostgreSQL installed and running.
1. Create a database named `appdrop`:
   ```sql
   CREATE DATABASE appdrop;
   ```
2. Run the migration script to create tables and triggers:
   *   Copy the code from `backend/migrations/001_create_tables.sql` and execute it in **pgAdmin** or your SQL tool of choice.

### 2. Environment Configuration
1. Open the file `backend/.env`.
2. Update the `DB_PASSWORD` to match your PostgreSQL password:
   ```env
   DB_PASSWORD=your_password_here
   ```

### 3. Start the Server
Navigate to the `backend` folder and run:
```bash
go mod tidy
go run main.go
```
The server will start on **`http://localhost:8080`**.

---

## 📖 API Documentation & Testing (CURL)

Use these commands to verify the endpoints as required by the assignment.

### 1. Health & Setup
```bash
# Check if API is running
curl http://localhost:8080/health
```

### 2. Page Management
```bash
# A. Create a Home Page
curl -X POST http://localhost:8080/pages \
  -H "Content-Type: application/json" \
  -d '{"name": "Home Screen", "route": "/home", "is_home": true}'

# B. Create a Sale Page
curl -X POST http://localhost:8080/pages \
  -H "Content-Type: application/json" \
  -d '{"name": "Summer Sale", "route": "/sale", "is_home": false}'

# C. List Pages (PAGINATED - Bonus)
curl "http://localhost:8080/pages?page=1&per_page=5"

# D. Update Page
# Replace :id with the UUID from step A
curl -X PUT http://localhost:8080/pages/:id \
  -H "Content-Type: application/json" \
  -d '{"name": "Main Dashboard"}'
```

### 3. Widget Management
```bash
# A. Add a Banner Widget
# Replace :page_id with the UUID from the Home Page
curl -X POST http://localhost:8080/pages/:page_id/widgets \
  -H "Content-Type: application/json" \
  -d '{
    "type": "banner",
    "position": 1,
    "config": {"title": "Huge Sale", "image_url": "https://picsum.photos/800/400"}
  }'

# B. Add a Text Widget
curl -X POST http://localhost:8080/pages/:page_id/widgets \
  -H "Content-Type: application/json" \
  -d '{
    "type": "text",
    "position": 2,
    "config": {"content": "Check our latest arrivals!", "style": "heading"}
  }'

# C. Get Widgets for Page (FILTERED - Bonus)
curl "http://localhost:8080/pages/:page_id/widgets?type=banner"

# D. Reorder Widgets (Bonus)
curl -X POST http://localhost:8080/pages/:page_id/widgets/reorder \
  -H "Content-Type: application/json" \
  -d '{"widget_ids": ["WIDGET_ID_2", "WIDGET_ID_1"]}'
```

---

## 🧪 Error Codes
| Code | Reason |
|------|--------|
| `VALIDATION_ERROR` | Missing fields or invalid widget type. |
| `CONFLICT` | The route (URL) is already taken. |
| `NOT_FOUND` | Page or Widget ID doesn't exist. |

---

## 📊 Performance Indicators
- **Database Indexing**: Optimized indices on `pages(route)`, `widgets(page_id)`, and `widgets(position)` for sub-millisecond query performance.
- **Connection Pooling**: Configured for high-concurrency mobile traffic.

# ToDo Service API Documentation

Base URL: `http://localhost:3002`  
Authorization: All endpoints require a **Bearer Token** in the headers.

---

### Global Headers
**Key:** `Authorization`  
**Value:** `Bearer <YOUR_JWT_TOKEN>`

---

## 1. Create a ToDo
Creates a new ToDo assigned to the currently authenticated user.

- **Method:** `POST`
- **URL:** `http://localhost:3002/todos`
- **Body (raw JSON):**
  ```json
  {
    "title": "Buy groceries",
    "completed": false
  }
  ```
- **Success Response:** `201 Created`
  ```json
  {
    "id": 1,
    "user_id": "your-uuid",
    "title": "Buy groceries",
    "completed": false
  }
  ```

---

## 2. Get All ToDos
Retrieves a list of all ToDo items exclusively owned by the authenticated user.

- **Method:** `GET`
- **URL:** `http://localhost:3002/todos`
- **Success Response:** `200 OK`
  ```json
  [
    {
      "id": 1,
      "user_id": "your-uuid",
      "title": "Buy groceries",
      "completed": false
    }
  ]
  ```

---

## 3. Update a ToDo
Updates the `title` and/or `completed` status of a specific ToDo item. You can only update your own items.

- **Method:** `PUT`
- **URL:** `http://localhost:3002/todos/{id}` (Replace `{id}` with the ToDo's ID)
- **Body (raw JSON):**
  ```json
  {
    "title": "Buy groceries and cook dinner",
    "completed": true
  }
  ```
- **Success Response:** `200 OK`
  ```json
  {
    "status": "updated"
  }
  ```

---

## 4. Delete a ToDo
Deletes a specific ToDo by its ID. You can only delete your own items.

- **Method:** `DELETE`
- **URL:** `http://localhost:3002/todos/{id}` (Replace `{id}` with the ToDo's ID)
- **Success Response:** `200 OK`
  ```json
  {
    "status": "deleted"
  }
  ```

# TaskSphere

TaskSphere is a full-stack task management platform featuring a high-performance **Go backend** and a reactive **React frontend**. The project is fully containerized using Docker, focusing on RESTful API design, modular architecture, and environment-driven configuration.

## Features

- **Task Management:** Full CRUD (Create, Read, Update, Delete) operations for managing workflows.
- **State Tracking:** Dynamic task status management (Pending, In-Progress, Completed).
- **Environment Management:** Secure configuration handling via `.env` files.
- **Containerized Stack:** Multi-service orchestration with Docker and Docker Compose.
- **Type Safety:** Strict data modeling using Go structs to ensure consistency between the API and the UI.

## Technical Highlights

- **RESTful API:** Developed a modular backend in Go, utilizing the standard library and idiomatic error handling.
- **JSON Serialization:** Efficient data exchange between Go and React using struct-to-JSON mapping.
- **Docker Networking:** Configured internal service communication between the backend API and the frontend client.


# เลี้ยงกาแฟผู้พัฒนา

!["Alt text"](https://warathepj.github.io/js-ai-gallery/public/image/promptpay-20.png)

# Task Manager Application

A full-stack task management application built with React (frontend) and Go (backend). The application allows users to create, manage, and prioritize tasks with detailed attributes including deadlines, priorities, dependencies, and more.

## Tech Stack

### Frontend

- Vite
- TypeScript
- React
- shadcn-ui
- Tailwind CSS

### Backend

- Go
- MongoDB
- Gorilla Mux
- CORS support

## Prerequisites

- Node.js & npm - [install with nvm](https://github.com/nvm-sh/nvm#installing-and-updating)
- Go 1.21.4 or later
- MongoDB

## Code

### Backend

```bash
git clone https://github.com/warathepj/go-task-manager-backend.git
```

### Frontend

```bash
git clone https://github.com/warathepj/go-task-manager-frontend.git
```

## Getting Started

### Backend Setup

1. Navigate to the backend directory:

```bash
cd go-task-manager-backend
```

2. Install Go dependencies:

```bash
go mod download
```

3. Start the backend server:

```bash
go run main.go
```

The server will start on `http://localhost:8000`

### Frontend Setup

1. Navigate to the frontend directory:

```bash
cd go-task-manager-frontend
```

2. Install dependencies:

```bash
npm install
```

3. Start the development server:

```bash
npm run dev
```

The application will be available at `http://localhost:5173`

## API Endpoints

- `GET /api/tasks` - Retrieve all tasks
- `POST /api/tasks` - Create a new task
- `GET /api/tasks/{id}` - Retrieve a specific task
- `PUT /api/tasks/{id}` - Update a task
- `DELETE /api/tasks/{id}` - Delete a task

## Task Structure

Tasks include the following attributes:

```typescript
{
  id: string
  description: string
  deadline: string
  timeRequired: string
  priority: string
  urgency: number
  dependencies: string[]
  resources: string[]
  subtasks: string[]
  group?: string
}
```

## Development

- Backend code is in the `backend/` directory
- Frontend code is in the `task-manager/` directory
- Frontend API calls are handled in `task-manager/src/lib/api.ts`
- Main task interface is defined in `task-manager/src/types.ts`

## Building for Production

### Backend

```bash
cd backend
go build
```

### Frontend

```bash
cd task-manager
npm run build
```

## License

This project is licensed under the MIT License - see the LICENSE file for details





# GoPloy

<img width="780" height="844" alt="Screenshot 2025-07-16 at 18 27 40" src="https://github.com/user-attachments/assets/c92b3c2d-ed1c-4010-b037-9d1eb71b7d95" />
--

**GoPloy** is a simple one-click deployment platform for React applications hosted on GitHub. It provides a seamless interface to trigger builds, monitor live logs, and access the deployed project via a custom slug.

---

## Features

- Deploy any public React GitHub repository with a single click  
- Smart slug generation for unique project URLs  
- Custom slug option for personalized project links  
- Real-time WebSocket log streaming  
- Fully serverless — powered by AWS ECS (Fargate), Redis, and Upstash  
- Modern UI built with React, Tailwind CSS, and ShadCN UI components  

---

## Tech Stack

### Frontend
- React + Vite
- Tailwind CSS + ShadCN UI
- WebSockets for live logs

### Backend
- Go (Golang) API server
- AWS ECS for task execution
- Upstash Redis Pub/Sub for log streaming

---

## Getting Started

### 1. Clone the repository

```bash
git clone https://github.com/py-xis/goploy.git
cd goploy
```

2. Start the Backend Server

```
cd backend
go run main.go
```

Make sure to:
	•	Configure your AWS credentials
	•	Set correct ECS cluster, task definition, subnets, and security groups
	•	Configure Redis using Upstash or a local instance

3. Start the Frontend

```
cd frontend
npm install
npm run dev
```

Your app should now be available at:
http://localhost:5173


## Deploying a Project
1. Enter the GitHub repository URL (e.g., https://github.com/user/my-app.git)
2. Optionally specify a custom slug
3. Click Trigger Build
4. Watch live logs stream in as the project builds
5. Access your deployed app at http://<your-slug>.localhost:8000


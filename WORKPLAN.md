# Project Tech Stack

### Frontend

* React + Vite
* React Router
* Axios
* Context API (or Zustand if allowed)
* Socket.IO client equivalent not allowed, so use native WebSocket
* CSS Modules / Tailwind (if allowed)

### Backend

* Go
* Gorilla WebSocket
* net/http
* bcrypt
* UUID

### Database

* SQLite
* golang-migrate

### DevOps

* Docker
* Docker Compose

---

# Team Structure

## Member 1: Backend Lead & Authentication

### Responsibilities

#### Authentication

* Register
* Login
* Logout
* Session Management
* Cookies

#### User Management

* Profile retrieval
* Profile editing
* Public/Private toggle

#### Image Upload

* Avatar uploads
* Post image uploads
* Validation

#### APIs

```txt
POST /register
POST /login
POST /logout

GET /profile/:id
PUT /profile

POST /upload
```

### Deliverables

* Authentication middleware
* Session package
* User APIs
* Image handling

---

## Member 2: Database Architect & Backend APIs

### Responsibilities

#### Database Design

Design ERD for:

```txt
users
sessions
followers
follow_requests
posts
comments
groups
group_members
group_invites
events
event_responses
notifications
messages
group_messages
```

#### Migrations

Create:

```txt
000001_users
000002_sessions
000003_followers
...
```

#### Repository Layer

Database functions:

```go
CreateUser()
GetUser()

CreatePost()
GetPosts()

CreateGroup()
```

### Deliverables

* Database schema
* All migrations
* Database helper functions
* Query optimization

---

## Member 3: Frontend Lead

### Responsibilities

#### Authentication Pages

* Login
* Register

#### User Pages

* Profile
* Edit Profile

#### Feed

* Home feed
* Post creation
* Comments

#### Notifications UI

* Notification dropdown
* Notification badges

#### State Management

```txt
AuthContext
NotificationContext
```

### Deliverables

* Main UI
* Responsive layout
* Routing
* API integration

---

## Member 4: Real-Time Systems & Groups

### Responsibilities

#### WebSockets

Private Chat

```txt
User A <-> User B
```

#### Group Chat

```txt
Group Room
```

#### Notifications

Real-time notifications

```txt
Follow Request
Group Invite
Event Creation
```

#### Groups

* Create group
* Invite members
* Join requests
* Events

### Deliverables

* WebSocket server
* Chat UI
* Group features
* Event system

---

# Suggested Folder Structure

```txt
social-network/

backend/
│
├── cmd/
│   └── server.go
│
├── pkg/
│   ├── auth/
│   ├── handlers/
│   ├── middleware/
│   ├── websocket/
│   ├── db/
│   │   ├── migrations/
│   │   └── sqlite/
│   └── repository/
│
├── uploads/
│
└── Dockerfile

frontend/
│
├── src/
│   ├── pages/
│   ├── components/
│   ├── services/
│   ├── contexts/
│   ├── hooks/
│   └── routes/
│
└── Dockerfile
```

---

# Database Tables

## Core

```sql
users
sessions
followers
follow_requests
```

## Posts

```sql
posts
post_visibility
comments
```

## Groups

```sql
groups
group_members
group_invites
group_join_requests
```

## Events

```sql
events
event_responses
```

## Messaging

```sql
messages
group_messages
```

## Notifications

```sql
notifications
```

---

# 14-Day Sprint Plan

## Day 1

### Entire Team

* Read project
* Design ERD
* Define APIs
* Create Git workflow

Deliverable:

```txt
ERD completed
API specification completed
```

---

## Day 2

### Backend

* Database setup
* Migrations

### Frontend

* React setup
* Routing setup

---

## Day 3

### Backend

* Register
* Login

### Frontend

* Register page
* Login page

---

## Day 4

### Backend

* Sessions
* Cookies

### Frontend

* Protected routes

---

## Day 5

### Backend

* Profiles

### Frontend

* Profile pages

---

## Day 6

### Backend

* Follow requests

### Frontend

* Follow UI

---

## Day 7

### Backend

* Posts
* Comments

### Frontend

* Feed

---

## Day 8

### Backend

* Privacy system

### Frontend

* Post visibility UI

---

## Day 9

### Backend

* Groups

### Frontend

* Group pages

---

## Day 10

### Backend

* Group invitations
* Join requests

### Frontend

* Group management

---

## Day 11

### Backend

* Events

### Frontend

* Event pages

---

## Day 12

### Backend

* WebSockets
* Messaging

### Frontend

* Chat UI

---

## Day 13

### Entire Team

Integration Day

Tasks:

* Fix broken APIs
* Fix database issues
* Fix WebSocket issues
* Dockerize

---

## Day 14

Testing & Presentation

Tasks:

```txt
Authentication testing
Groups testing
Chats testing
Notifications testing
Docker testing
Peer review preparation
```

---

# Git Workflow

Branches:

```txt
main

develop

feature/auth
feature/posts
feature/groups
feature/chat
feature/frontend
feature/database
```

Rules:

1. No direct push to main.
2. Every feature through Pull Request.
3. At least one team review before merge.
4. Merge to develop first.
5. Stable releases go to main.

---

# Critical Risk Areas

These will consume most time:

### 1. Sessions & Cookies

Do this by Day 4.
Everything depends on authentication.

### 2. Post Visibility Logic

```txt
Public
Followers Only
Selected Followers
```

This becomes messy fast if not designed carefully.

### 3. Group Permissions

```txt
Creator
Member
Invited
Requested
```

Define states early.

### 4. WebSockets

Do NOT leave chats until the final 2 days.

Start WebSocket proof-of-concept by Day 5 or Day 6.  

 
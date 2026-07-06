# Frontend Integration Checklist

This is a current snapshot of frontend-to-backend integration, using `backend/internal/routes/routes.go` as the source of truth for backend endpoints.

Main frontend integration points:

- `frontend/src/lib/api.ts`: primary API wrapper, shared request helper, response/payload types, asset URL resolution.
- `frontend/src/contexts/WebSocketContext.tsx`: WebSocket connection setup for live chat/notifications.
- `frontend/src/app/(main)/*`: pages that call API wrappers or, in one case, raw `fetch`.
- `frontend/.env.local` and `frontend/.env.example`: `NEXT_PUBLIC_API_URL=http://localhost:8080`.

Status meanings:

- Done: frontend helper exists and is used by a page/component against a matching backend route.
- Partial: some integration exists, but route coverage, UI wiring, or route matching is incomplete.
- Pending: backend route exists but frontend integration is missing, or frontend points at routes that do not exist in `routes.go`.

## Summary

| Area | Status | Notes |
| --- | --- | --- |
| Environment/API base URL | Done | `NEXT_PUBLIC_API_URL` is wired through `api.ts`; requests include credentials. |
| Auth register/login | Done | Register and login pages call matching backend routes. |
| Auth logout | Pending | Backend route exists, but frontend sidebar logout is just a link to `/login`. |
| Profile | Partial | Own profile load/update/avatar are integrated; public profile by id and following list are not. |
| Follow flows | Pending | Backend follow/request endpoints exist, but frontend has no matching helpers or UI. |
| Feed/posts | Partial | Get feed, create post, get post, upload post image are integrated; update/delete/user posts/group posts are not. |
| Comments | Partial | Fetch/create comments are integrated; delete comment and comment image upload are not. |
| Groups/events | Partial | Group list is integrated; detail page uses same-origin raw fetch, and create/join/events are not fully wired. |
| Chat/messages | Partial | Send and private message history have helpers/UI; conversation/read helpers point to backend routes that do not exist. |
| Notifications | Pending | Frontend has helpers/UI for notification HTTP routes, but those routes are not registered in `routes.go`. Live WebSocket notifications are partially consumed. |
| WebSocket | Partial | Frontend connects to `/api/ws` after checking auth; backend route exists. Message/notification consumers depend on payload shapes and surrounding endpoints. |

## Route Checklist

### Auth

| Backend route | Frontend status | Frontend location | Notes |
| --- | --- | --- | --- |
| `POST /api/auth/register` | Done | `authAPI.register`, `app/(auth)/register/page.tsx` | Registration posts to backend and redirects to feed. |
| `POST /api/auth/login` | Done | `authAPI.login`, `app/(auth)/login/page.tsx` | Login posts to backend and redirects to feed. |
| `POST /api/auth/logout` | Pending | None | Add `authAPI.logout()` and wire the sidebar logout action to call it before navigation. |

### Profile and Users

| Backend route | Frontend status | Frontend location | Notes |
| --- | --- | --- | --- |
| `GET /api/profile` | Done | `profileAPI.getMyProfile`, feed/profile/messages/WebSocket | Used for current user, profile page, and auth check before WebSocket connect. |
| `GET /api/profile/{id}` | Pending | None | Needed for viewing another user's profile. |
| `PUT /api/profile` | Done | `profileAPI.updateProfile`, `app/(main)/profile/page.tsx` | Wrapper updates then refetches `/api/profile` because backend returns `204`. |
| `POST /api/profile/avatar` | Done | `profileAPI.uploadAvatar`, `app/(main)/profile/page.tsx` | Avatar upload uses `avatar` form field. |
| `GET /api/users/{id}/followers` | Done | `profileAPI.getFollowers`, `app/(main)/feed/page.tsx` | Used for private post audience selection. |
| `GET /api/users/{id}/following` | Pending | None | Add helper/UI if following lists are needed. |

### Follow Requests

| Backend route | Frontend status | Frontend location | Notes |
| --- | --- | --- | --- |
| `POST /api/follow/requests` | Pending | None | Route is registered as `POST`, while handler comment says `GET /api/follow/requests`; confirm intended method before wiring. |
| `POST /api/follow/{id}` | Pending | None | Needed for follow/request button on user/profile surfaces. |
| `DELETE /api/follow/{id}` | Pending | None | Needed for unfollow. |
| `POST /api/follow/{id}/accept` | Pending | None | Needed for accepting follow requests. |
| `POST /api/follow/{id}/decline` | Pending | None | Needed for declining follow requests. |

### Posts and Feed

| Backend route | Frontend status | Frontend location | Notes |
| --- | --- | --- | --- |
| `POST /api/posts` | Done | `feedAPI.createPost`, `app/(main)/feed/page.tsx` | Composer supports content, privacy, selected viewers. |
| `GET /api/posts/{id}` | Done | `feedAPI.getPost`, `app/(main)/feed/page.tsx` | Used after create/image upload to fetch completed post. |
| `PUT /api/posts/{id}` | Pending | None | No edit post helper/UI. |
| `DELETE /api/posts/{id}` | Pending | None | No delete post helper/UI. |
| `POST /api/posts/{id}/image` | Done | `feedAPI.uploadPostImage`, `app/(main)/feed/page.tsx` | Uses `image` form field. |
| `GET /api/feed` | Done | `feedAPI.getFeed`, `app/(main)/feed/page.tsx` | Main feed loads with pagination params. |
| `GET /api/users/{id}/posts` | Pending | None | Not wired into profile/user pages. |
| `GET /api/groups/{id}/posts` | Pending | None | Group detail page has a discussion placeholder only. |

### Comments

| Backend route | Frontend status | Frontend location | Notes |
| --- | --- | --- | --- |
| `POST /api/posts/{id}/comments` | Done | `feedAPI.createComment`, `app/(main)/feed/page.tsx` | Top-level text comments are wired. Reply support is possible via `parent_comment_id` type but no UI yet. |
| `GET /api/posts/{id}/comments` | Done | `feedAPI.getComments`, `app/(main)/feed/page.tsx` | Comment tree rendering is wired. |
| `DELETE /api/comments/{id}` | Pending | None | No helper/UI. |
| `POST /api/comments/{id}/image` | Pending | None | No helper/UI for comment image upload. |

### Groups and Events

| Backend route | Frontend status | Frontend location | Notes |
| --- | --- | --- | --- |
| `POST /api/groups` | Pending | Button only in `app/(main)/groups/page.tsx` | Create Group button is not wired. |
| `GET /api/groups` | Done | `groupAPI.getGroups`, `app/(main)/groups/page.tsx` | Group list loads from backend. |
| `GET /api/groups/{id}` | Partial | Raw `fetch` in `app/(main)/groups/[id]/page.tsx` | Uses `/api/groups/${id}` without `NEXT_PUBLIC_API_URL`, so it will hit the frontend origin unless proxied. Move into `groupAPI`. |
| `POST /api/groups/{id}/join` | Pending | Button only in group detail page | Join Group button is not wired. |
| `POST /api/groups/{id}/events` | Pending | Button only in group detail page | Event creation button is present but not wired. |
| `GET /api/groups/{id}/events` | Partial | Raw `fetch` in `app/(main)/groups/[id]/page.tsx` | Same base URL issue as group detail. Move into `groupAPI`. |

### Chat and Messages

| Backend route | Frontend status | Frontend location | Notes |
| --- | --- | --- | --- |
| `POST /api/chat/messages` | Done | `chatAPI.sendMessage`, `app/(main)/messages/page.tsx` | Sending private messages is wired. |
| `GET /api/chat/messages/{userId}` | Done | `chatAPI.getMessages`, `app/(main)/messages/page.tsx` | Private message history is wired. |
| `GET /api/groups/{id}/messages` | Pending | None | No group chat UI/helper. |
| No matching backend route | Pending | `chatAPI.getConversations` | Frontend calls `GET /api/chat/conversations`, but `routes.go` does not register it. |
| No matching backend route | Pending | `chatAPI.markConversationRead` | Frontend calls `POST /api/chat/conversations/{userId}/read`, but `routes.go` does not register it. |

### WebSocket

| Backend route | Frontend status | Frontend location | Notes |
| --- | --- | --- | --- |
| `GET /api/ws` | Partial | `WebSocketContext`, messages/notifications/dropdown | Connection is wired and auth-checked first. Consumers listen for `chat_message`, `notification`, `follow_request`, `group_invite`, and `group_event` payloads. |

### Notifications

`routes.go` does not currently register HTTP notification routes, but `frontend/src/lib/api.ts` and `app/(main)/notifications/page.tsx` assume these endpoints:

| Frontend route | Status | Notes |
| --- | --- | --- |
| `GET /api/notifications` | Pending | No registered backend route. |
| `POST /api/notifications/{notificationId}/read` | Pending | No registered backend route. |
| `POST /api/notifications/read-all` | Pending | No registered backend route. |
| `POST /api/follow-requests/{actorId}/accept` | Pending | No registered backend route; backend uses `POST /api/follow/{id}/accept`. |
| `POST /api/follow-requests/{actorId}/decline` | Pending | No registered backend route; backend uses `POST /api/follow/{id}/decline`. |

The backend does have notification storage helpers and WebSocket notification dispatch, but no HTTP notification handler routes are exposed in `routes.go`.

## Recommended Next Steps

1. Add missing `api.ts` helpers for backend routes that already exist: logout, public profile, following list, follow/unfollow/request accept/decline, post update/delete, user posts, group posts, comment delete/image, group create/detail/join/events, group messages.
2. Fix route mismatches before building UI around them:
   - Decide whether pending follow requests should be `GET` or `POST /api/follow/requests`.
   - Either add backend routes for chat conversations/read receipts, or change the messages page to use only available backend routes.
   - Either add backend notification HTTP routes, or remove/disable frontend notification HTTP calls.
3. Replace raw `fetch` calls in `app/(main)/groups/[id]/page.tsx` with `groupAPI` methods so requests use `NEXT_PUBLIC_API_URL` and credentials.
4. Wire the visible but inactive UI controls: logout, create group, view group navigation, join group, create event, follow/unfollow, accept/decline requests, edit/delete post, delete/image comments.
5. After each area is wired, verify with `npx tsc --noEmit` and targeted browser/API testing against the backend on `http://localhost:8080`.

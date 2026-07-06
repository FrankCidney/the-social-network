# Frontend Integration Checklist

This is the current frontend-to-backend integration snapshot, using `backend/internal/routes/routes.go` as backend route source of truth.

Main frontend integration points:

- `frontend/src/lib/api.ts`: main API wrapper, shared request helper, frontend response/payload types, asset URL resolution.
- `frontend/src/contexts/WebSocketContext.tsx`: WebSocket connection setup for chat/notifications.
- `frontend/src/app/(main)/*`: main app pages using the shared API layer.
- `frontend/.env.local` and `frontend/.env.example`: `NEXT_PUBLIC_API_URL=http://localhost:8080`.

Status meanings:

- Done: frontend helper exists and is wired to a matching backend route and visible UI.
- Partial: helper or UI exists, but there are known backend gaps, bugs, or incomplete flows.
- Pending: backend route exists but frontend is not using it yet, or frontend still needs new backend support.

## Summary

| Area | Status | Notes |
| --- | --- | --- |
| Environment/API base URL | Done | `NEXT_PUBLIC_API_URL` is wired through `api.ts`; requests include credentials. |
| Auth register/login | Done | Register and login pages call matching backend routes. |
| Auth logout | Pending | Backend route exists, but frontend sidebar logout is still just a link to `/login`. |
| Profile | Partial | Own profile load/update/avatar plus public profile by id are integrated. Following list still has no dedicated UI surface beyond stats modal. |
| Follow flows | Partial | Follow/unfollow/request is wired on profile pages; accept/decline follow requests and broader discovery follow CTAs are still pending. |
| Feed/posts | Partial | Feed, create post, get post, upload post image, comments, and replies are integrated; edit/delete post still missing. |
| Comments | Partial | Fetch/create comments and comment image upload are integrated in feed; delete comment is still missing. |
| Groups/events | Partial | Group browse, create, detail, join, members, invite picker, events, creator/member/request state, and RSVP wiring are in place. |
| Chat/messages | Partial | Private messages are wired, including opening a direct thread from profile/discovery. Group chat remains wired inside group detail. Presence/online status is not implemented. |
| Notifications | Pending | Frontend has notification HTTP helpers/UI assumptions, but those routes are still not registered in `routes.go`. |
| WebSocket | Partial | Frontend connects to `/api/ws`. Group and follow notifications depend on backend payload/route coverage beyond the current HTTP work. |

## Route Checklist

### Auth

| Backend route | Frontend status | Frontend location | Notes |
| --- | --- | --- | --- |
| `POST /api/auth/register` | Done | `authAPI.register`, `app/(auth)/register/page.tsx` | Registration posts to backend and redirects to feed. |
| `POST /api/auth/login` | Done | `authAPI.login`, `app/(auth)/login/page.tsx` | Login posts to backend and redirects to feed. |
| `POST /api/auth/logout` | Pending | None | Add `authAPI.logout()` and wire the sidebar logout action. |

### Profile and Users

| Backend route | Frontend status | Frontend location | Notes |
| --- | --- | --- | --- |
| `GET /api/profile` | Done | `profileAPI.getMyProfile`, feed/profile/messages/groups/WebSocket | Used for current user and auth-aware UI state. |
| `GET /api/profile/{id}` | Done | `profileAPI.getProfile`, `app/(main)/profile/[id]/page.tsx` | Public/other-user profile route is now wired, including follow and message entry points. |
| `PUT /api/profile` | Done | `profileAPI.updateProfile`, `app/(main)/profile/page.tsx` | Wrapper updates then refetches `/api/profile` because backend returns `204`. |
| `POST /api/profile/avatar` | Done | `profileAPI.uploadAvatar`, `app/(main)/profile/page.tsx` | Avatar upload uses `avatar` form field. |
| `GET /api/users/{id}/followers` | Done | `profileAPI.getFollowers`, `app/(main)/feed/page.tsx` | Used for private post audience selection. |
| `GET /api/users/{id}/following` | Partial | `profileAPI.getFollowing`, profile stats modal | Helper is wired in the profile stats modal; no standalone following page/list flow yet. |
| `GET /api/users/search` | Done | `profileAPI.searchUsers`, layout discovery panel and group manage tab | Used for people discovery and invite picking. |

### Follow Requests

| Backend route | Frontend status | Frontend location | Notes |
| --- | --- | --- | --- |
| `POST /api/follow/requests` | Pending | None | Route method still needs confirmation against intended behavior. |
| `POST /api/follow/{id}` | Done | `followAPI.follow`, `app/(main)/profile/[id]/page.tsx` | Profile CTA respects backend public/private follow behavior. |
| `DELETE /api/follow/{id}` | Done | `followAPI.unfollow`, `app/(main)/profile/[id]/page.tsx` | Profile unfollow CTA is wired. |
| `POST /api/follow/{id}/accept` | Pending | None | Needed for accepting follow requests. |
| `POST /api/follow/{id}/decline` | Pending | None | Needed for declining follow requests. |

### Posts and Feed

| Backend route | Frontend status | Frontend location | Notes |
| --- | --- | --- | --- |
| `POST /api/posts` | Done | `feedAPI.createPost`, `app/(main)/feed/page.tsx`, `app/(main)/groups/[id]/page.tsx` | Feed and group post creation are wired. |
| `GET /api/posts/{id}` | Done | `feedAPI.getPost`, feed and group detail pages | Used after create/image upload to fetch completed post. |
| `PUT /api/posts/{id}` | Pending | None | No edit post helper/UI. |
| `DELETE /api/posts/{id}` | Pending | None | No delete post helper/UI. |
| `POST /api/posts/{id}/image` | Done | `feedAPI.uploadPostImage`, `app/(main)/feed/page.tsx` | Uses `image` form field. |
| `GET /api/feed` | Done | `feedAPI.getFeed`, `app/(main)/feed/page.tsx` | Main feed loads with pagination params. |
| `GET /api/users/{id}/posts` | Done | `feedAPI.getUserPosts`, profile stats modal | User posts are wired into the profile posts modal and use feed-style cards. |
| `GET /api/groups/{id}/posts` | Done | `feedAPI.getGroupPosts`, `app/(main)/groups/[id]/page.tsx` | Group posts load in the main feed area of the group page. |

### Comments

| Backend route | Frontend status | Frontend location | Notes |
| --- | --- | --- | --- |
| `POST /api/posts/{id}/comments` | Done | `feedAPI.createComment`, `app/(main)/feed/page.tsx` | Top-level comments and replies are wired in feed. |
| `GET /api/posts/{id}/comments` | Done | `feedAPI.getComments`, `app/(main)/feed/page.tsx` | Comment tree rendering is wired in feed. |
| `DELETE /api/comments/{id}` | Pending | None | No helper/UI. |
| `POST /api/comments/{id}/image` | Done | `feedAPI.uploadCommentImage`, `app/(main)/feed/page.tsx` | Comment image upload is wired in feed. |

### Groups and Events

| Backend route | Frontend status | Frontend location | Notes |
| --- | --- | --- | --- |
| `POST /api/groups` | Done | `groupAPI.createGroup`, `app/(main)/groups/page.tsx` | Create group is wired and redirects to group detail. |
| `GET /api/groups` | Done | `groupAPI.getGroups`, `app/(main)/groups/page.tsx` | Group list loads from backend. Handles `null` response safely. |
| `GET /api/groups/{id}` | Done | `groupAPI.getGroup`, `app/(main)/groups/[id]/page.tsx` | Uses enriched detail response including creator and membership status. |
| `GET /api/groups/{id}/members` | Done | `groupAPI.getMembers`, `app/(main)/groups/[id]/page.tsx` | Members tab is wired and backend scan bug is fixed. |
| `POST /api/groups/{id}/join` | Done | `groupAPI.requestJoin`, `app/(main)/groups/[id]/page.tsx` | Join request button is wired. |
| `POST /api/groups/{id}/invite` | Done | `groupAPI.inviteUser`, `app/(main)/groups/[id]/page.tsx` | Wired in Manage tab with search-based invite picker. |
| `POST /api/groups/{id}/invite/accept` | Done | `groupAPI.acceptInvite`, `app/(main)/groups/[id]/page.tsx` | Accept invite is wired from group detail header. |
| `POST /api/groups/{id}/invite/decline` | Done | `groupAPI.declineInvite`, `app/(main)/groups/[id]/page.tsx` | Decline invite is wired from group detail header. |
| `GET /api/groups/{id}/requests` | Done | `groupAPI.getJoinRequests`, `app/(main)/groups/[id]/page.tsx` | Creator-only request list is wired in Manage tab. |
| `POST /api/groups/{id}/requests/{userId}/accept` | Done | `groupAPI.acceptJoinRequest`, `app/(main)/groups/[id]/page.tsx` | Creator accept action is wired. |
| `POST /api/groups/{id}/requests/{userId}/decline` | Done | `groupAPI.declineJoinRequest`, `app/(main)/groups/[id]/page.tsx` | Creator decline action is wired. |
| `POST /api/groups/{id}/events` | Done | `groupAPI.createEvent`, `app/(main)/groups/[id]/page.tsx` | Event creation is wired in Events tab. |
| `GET /api/groups/{id}/events` | Done | `groupAPI.getEvents`, `app/(main)/groups/[id]/page.tsx` | Events list loads in Events tab. |
| `POST /api/events/{id}/rsvp` | Done | `groupAPI.rsvpEvent`, `app/(main)/groups/[id]/page.tsx` | RSVP buttons are wired. |

### Chat and Messages

| Backend route | Frontend status | Frontend location | Notes |
| --- | --- | --- | --- |
| `POST /api/chat/messages` | Done | `chatAPI.sendMessage`, messages page and group detail page | Used for private messages and group chat send. Message send now dedupes optimistic UI against WebSocket echo. |
| `GET /api/chat/messages/{userId}` | Done | `chatAPI.getMessages`, `app/(main)/messages/page.tsx` | Private message history is wired. Backend now clamps pagination so previous messages reload correctly. |
| `GET /api/groups/{id}/messages` | Done | `chatAPI.getGroupMessages`, `app/(main)/groups/[id]/page.tsx` | Group chat history is wired in Chat tab. |
| `GET /api/chat/conversations` | Done | `chatAPI.getConversations`, `app/(main)/messages/page.tsx` | Route is registered and drives the private conversation sidebar/previews. |
| `POST /api/chat/conversations/{userId}/read` | Done | `chatAPI.markConversationRead`, `app/(main)/messages/page.tsx` | Route is now registered in backend. |

### WebSocket

| Backend route | Frontend status | Frontend location | Notes |
| --- | --- | --- | --- |
| `GET /api/ws` | Partial | `WebSocketContext`, messages/notifications/dropdown | Connection is wired and auth-checked first. Consumers listen for `chat_message`, `notification`, `follow_request`, `group_invite`, and `group_event`. |

### Notifications

`routes.go` still does not register HTTP notification routes, but `frontend/src/lib/api.ts` and `app/(main)/notifications/page.tsx` assume these endpoints:

| Frontend route | Status | Notes |
| --- | --- | --- |
| `GET /api/notifications` | Pending | No registered backend route. |
| `POST /api/notifications/{notificationId}/read` | Pending | No registered backend route. |
| `POST /api/notifications/read-all` | Pending | No registered backend route. |
| `POST /api/follow-requests/{actorId}/accept` | Pending | No registered backend route; backend uses `POST /api/follow/{id}/accept`. |
| `POST /api/follow-requests/{actorId}/decline` | Pending | No registered backend route; backend uses `POST /api/follow/{id}/decline`. |

## Known Issues

### Messages Read State Migration

Private-message conversations now expect a `read_at` column on `messages`.

Impact:

- Existing local databases need the new migration applied by restarting the backend so unread/read-aware conversation queries work.
- Without the migration, conversation list queries can fail with `no such column: read_at`.

## UI and Code Pattern Notes

- Keep the current UI style and visual language.
- Keep copy minimal. It is an app, not a blog.
- Prefer functional UI over explanatory text blocks.
- Group detail layout direction is now:
  - posts as the main feed column
  - side panel with tabs for chat, events, members, and management
- Group posts should continue to look like app feed posts rather than a totally separate visual system.

## Recommended Next Step

### Finish Follow Request and Notification Flows

Next backend/frontend work should focus on the remaining follow-request and notification gaps.

Why next:

- Profile follow/unfollow is wired, but accepting/declining follow requests is still not surfaced cleanly.
- Notification UI still assumes HTTP routes that are not registered in `routes.go`.
- Logout is still a placeholder link and is one of the remaining auth integration gaps.

Recommended scope:

- add/wire follow request accept and decline flows
- either implement notification HTTP routes or remove/disable those frontend assumptions
- add real logout action wiring

## Recommended Next Steps After That

1. Finish follow request accept/decline UX and notification-driven actions.
2. Add logout action wiring.
3. Revisit notification HTTP routes or remove/disable the current notification HTTP assumptions.
4. Add edit/delete post flows.
5. Add delete comment flow.
6. Decide whether people discovery should gain direct follow CTAs again after shared-state syncing is in place.

## Verification

After each integration area:

- run `npx tsc --noEmit` in `frontend/`
- run targeted browser testing against backend on `http://localhost:8080`
- where backend changes were made, run targeted Go package tests/build checks

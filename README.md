doing things

Problem Statement
    - Design and build a scalable backend system that supports the following:
    User Features
        Browse a list of upcoming events with details (name, venue, time, capacity).
        Book and cancel tickets, ensuring seat availability is updated correctly.
        View booking history.
    Admin Features
        Create, update, and manage events.
        View booking analytics (total bookings, most popular events, capacity utilization).

**Functional Requirements**

Users should be able to:
- Browse of a list of upcoming events
- Book tickets
- Cancel tickets

Admin Features:
- CRUD Events
- Analytics

Non Functional Requirements

- CAP Theorem
    - Consistency >> Availiability
    - Strong Consistency for tickets & high availiability for viewing events

- Reads >> Write
- scalability to handle popular events

**Core Entities**
- Users
- Admins
- Events
- Tickets

Thursday
- [x] postgres DB config
- [x] generate sql files with sqlc
- [x] Connect to postgres DB with sqlc
- [x] Migrations with goose
- [x] logging (just using golang's logging package for now)
- [x] user table setup
- [x] air toml
- [x] user repo
- [x] error domain
- [x] we are not doing payments
- [x] utils for http delivery layer

Friday
- [x] Download postman
- [x] http delivery layer
- [x] Figure out middlewares
- [x] Register User
- [x] Implement Login
- [ ] Implement Logout (do later)
- [ ] create make file (not happening)
- [x] Setup Health router
- [x] Logging
- [x] Events Table
- [x] Bookings Table
- [ ] Reserve Tickets
- [ ] Book Tickets
- [ ] Redis
- [ ] Caching?
- [x] Get Events Paginated
- [x] Delete Events
- [ ] Update Events (do this later)
- [ ] Admin Check
- [x] Return Event ID as well
- [ ] Graceful Shutdown of server

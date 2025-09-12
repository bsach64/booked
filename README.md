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
- [ ] Download postman
- [x] http delivery layer
- [x] Figure out middlewares
- [x] Register User
- [x] Implement Login
- [ ] Implement Logout (maybe later)
- [ ] create make file (not happening)
- [ ] Figure out most of the user features
- [x] Setup Health router
- [x] Logging

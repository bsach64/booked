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
- [ ] Auth
- [ ] create make file (not happening)
- [x] air toml
- [ ] Register User
- [ ] Implement Login
- [x] user repo
- [x] error domain
- [x] we are not doing payments

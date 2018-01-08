## Background

This is a simple 'chan' style image board that I thought would be fun to try over winter break.
A live demo is currently [here](http://board.dexterhaslem.com) but exercise caution.

It has most basic functionality in place. It's by no means a good reference as it's missing 
key security and maintainability.

TLDR: golang + postgres + server templates

### Persistence
Data persistence is done in a postgresql database 
versioned by flyway (see `db/sql for migrations`). To spin up a new db modify `db/conf/flyway.conf` to point at a valid db.
then run `flyway migrate` from the db folder.


### Front end / client
Originally the front end was written in angular5 w/ material, using redux to maintain state.
I decided writing lots of boilerplate and fraomework for a 90% readonly app was not fun and shot it all.
Now the entire client is basic html served from the server using go templates (see `templates/`)


### Backend 
The backend is simple golang, using vanilla http for most things. Thumbnail resizing is done with the
handy [resize](https://github.com/nfnt/resize/) library.


### database details
While I don't have an ERD handy, the schema is simple. Boards are the root. Threads have a 
board_id foreign key. Posts have a thread_id foreign key.

Originally I used join tables using composite keys for board_thread, thread_post 
(and even attachment join table, so there is some slight duplication between thread/post)
 but scrapped it for a more straightforward approach as there was no real benefits given
 the current simple design. If this was a giant system with sharding on the horizon that would be a consideration.

the SERIAL keys probably need to be bumped to BIGSERIAL, or if to be a true 'chan' experience,
keys should be unix timestamps by board



## Future improvements / TODOS

Low hanging fruit: templates have a lot of duplication that could be pulled into subtemplates.

The majority of missing features is statistics like 'R# I#' (replies/images) for threads,
listing threads by most recent bump and recent threads/posts which all could be done with
a couple SQL functions and views.



### license
GPL3, see `LICENSE`
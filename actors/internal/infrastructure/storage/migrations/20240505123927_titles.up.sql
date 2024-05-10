BEGIN;
CREATE TABLE titles (
	id int PRIMARY KEY
);

TRUNCATE actors;

CREATE TABLE actors_titles (
	actor_id int REFERENCES actors(id),
	title_id int REFERENCES titles(id),
	CONSTRAINT actors_titles_pk PRIMARY KEY(actor_id, title_id)
);
COMMIT;



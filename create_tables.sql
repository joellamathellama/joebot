/*
	persons(discord user) table
*/
CREATE TABLE IF NOT EXISTS persons
(
    id serial PRIMARY KEY,
    created_time timestamptz,
    modified_time timestamptz,
);

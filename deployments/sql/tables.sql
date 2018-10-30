CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS payments (
 ID uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
 info json NOT NULL
);

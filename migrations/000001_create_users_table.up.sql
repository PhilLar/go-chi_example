CREATE TABLE IF NOT EXISTS pets
(
  id        BIGSERIAL PRIMARY KEY,
  pet_name   TEXT NOT NULL,
  pet_kind   TEXT NOT NULL
);
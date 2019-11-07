BEGIN;

CREATE TYPE pet_kind_enum AS ENUM (
	'dog',
	'cat'
);
ALTER TABLE pets ALTER COLUMN pet_kind TYPE pet_kind_enum USING pet_kind::TEXT::pet_kind_enum;
ALTER TABLE pets ADD COLUMN pet_age SMALLINT NOT NULL DEFAULT -1;

COMMENT ON COLUMN pets.pet_age IS 'If a pet''s age is -1, then it was not set';

COMMIT;
BEGIN;

UPDATE status_page_groups SET type = 'current_status_indicator' WHERE type = 'none';
UPDATE status_page_monitors SET type = 'current_status_indicator' WHERE type = 'none';

CREATE TYPE status_page_element_type_new AS ENUM (
	'historical_timeline',
	'current_status_indicator'
);

ALTER TABLE status_page_groups
	ALTER COLUMN type TYPE status_page_element_type_new
	USING type::text::status_page_element_type_new;

ALTER TABLE status_page_monitors
	ALTER COLUMN type TYPE status_page_element_type_new
	USING type::text::status_page_element_type_new;

DROP TYPE status_page_element_type;
ALTER TYPE status_page_element_type_new RENAME TO status_page_element_type;

ALTER TABLE status_pages ADD COLUMN title text;
UPDATE status_pages SET title = slug WHERE title IS NULL;
ALTER TABLE status_pages ALTER COLUMN title SET NOT NULL;

COMMIT;

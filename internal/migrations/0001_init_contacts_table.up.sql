CREATE TABLE IF NOT EXISTS contacts (
  "id" TEXT NOT NULL,
  "name_prefix" TEXT DEFAULT NULL,
  "first_name" TEXT DEFAULT NULL,
  "middle_name" TEXT DEFAULT NULL,
  "surname" TEXT DEFAULT NULL,
  "name_suffix" TEXT DEFAULT NULL,
  "created_at" TEXT NOT NULL
) STRICT;

CREATE UNIQUE INDEX IF NOT EXISTS idx_contacts_id ON contacts(id);

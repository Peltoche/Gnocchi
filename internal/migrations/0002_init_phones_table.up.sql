CREATE TABLE IF NOT EXISTS phones (
  "id" TEXT NOT NULL,
  "type" TEXT NOT NULL,
  "iso2_region_code" TEXT NOT NULL,
  "international_formatted" TEXT NOT NULL,
  "national_formatted" TEXT NOT NULL,
  "normalized" TEXT NOT NULL,
  "contact_id" TEXT NOT NULL,
  "created_at" TEXT NOT NULL,
  FOREIGN KEY(contact_id) REFERENCES contacts(id) ON UPDATE RESTRICT ON DELETE CASCADE
) STRICT;

CREATE UNIQUE INDEX IF NOT EXISTS idx_phones_id ON phones(id);
CREATE INDEX IF NOT EXISTS idx_phones_international_id ON phones(international_formatted);
CREATE INDEX IF NOT EXISTS idx_phones_national_id ON phones(national_formatted);

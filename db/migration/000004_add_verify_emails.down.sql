DROP TABLE IF EXISTS "verify_emails" CASCADE;
-- CASCADE = memastikan jika ada record pada table lain yang REFERENCES ke record di table ini maka akan dihapus

ALTER TABLE "users" DROP COLUMN "is_email_verified";
CREATE TABLE "users" (
  "username" varchar(100) PRIMARY KEY,
  "hashed_password" varchar NOT NULL,
  "full_name" varchar(100) NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "password_changed_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00+00',
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE "accounts" ADD FOREIGN KEY ("owner") REFERENCES "users" ("username");

-- CREATE UNIQUE INDEX ON "accounts" ("owner", "currency");
-- sama saja seperti yang diatas, pilih mana aja boleh sama saja
ALTER TABLE "accounts" ADD CONSTRAINT "owner_currency_key" UNIQUE ("owner","currency");
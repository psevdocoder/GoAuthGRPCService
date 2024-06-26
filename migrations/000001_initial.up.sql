CREATE TABLE IF NOT EXISTS "users" (
    "id" INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    "username" TEXT NOT NULL UNIQUE,
    "password_hash" BLOB NOT NULL,
    "role" SMALLINT NOT NULL DEFAULT 0
);

CREATE INDEX IF NOT EXISTS "idx_username" ON "users" ("username");

CREATE TABLE IF NOT EXISTS "apps" (
    "id" INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    "name" TEXT NOT NULL UNIQUE,
    "secret" TEXT NOT NULL
)
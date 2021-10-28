CREATE TABLE IF NOT EXISTS "user" (
  "id" BIGSERIAL PRIMARY KEY,
  "email" VARCHAR(128) NOT NULL,
  "password" TEXT NOT NULL,
  "verified" BOOLEAN NOT NULL DEFAULT FALSE,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now()),
  "deleted_at" timestamptz
);

CREATE INDEX ON "user" ("email");

CREATE TABLE IF NOT EXISTS "password" (
  "id" BIGSERIAL PRIMARY KEY,
  "user_id" BIGINT,
  "password" TEXT NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE "password" ADD FOREIGN KEY ("user_id") REFERENCES "user" ("id");

CREATE INDEX ON "password" ("user_id");

CREATE TABLE IF NOT EXISTS "profile" (
  "id" BIGSERIAL PRIMARY KEY,
  "user_id" BIGINT,
  "fullname" VARCHAR(128) NOT NULL,
  "location" VARCHAR(128) NOT NULL DEFAULT '',
  "bio" TEXT NOT NULL DEFAULT '',
  "web" VARCHAR(128) NOT NULL DEFAULT '',
  "picture" VARCHAR(128) NOT NULL DEFAULT '',
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE "profile" ADD FOREIGN KEY ("user_id") REFERENCES "user" ("id");

CREATE INDEX ON "profile" ("user_id");

CREATE TABLE IF NOT EXISTS "tfa" (
  "id" BIGSERIAL PRIMARY KEY,
  "user_id" BIGINT,
  "secret" VARCHAR(128),
  "enable" BOOLEAN NOT NULL DEFAULT FALSE,
  "enable_at" timestamptz
);

ALTER TABLE "tfa" ADD FOREIGN KEY ("user_id") REFERENCES "user" ("id");

CREATE INDEX ON "tfa" ("user_id");

CREATE TABLE IF NOT EXISTS "tfa_codes" (
  "id" BIGSERIAL PRIMARY KEY,
  "tfa_id" BIGINT,
  "code" VARCHAR(128) NOT NULL
);

ALTER TABLE "tfa_codes" ADD FOREIGN KEY ("tfa_id") REFERENCES "tfa" ("id");

CREATE INDEX ON "tfa_codes" ("tfa_id");

CREATE TABLE IF NOT EXISTS "sessions" (
  "id" BIGSERIAL PRIMARY KEY,
  "user_id" BIGINT,
  "is_current" BOOLEAN NOT NULL DEFAULT TRUE
);

ALTER TABLE "sessions" ADD FOREIGN KEY ("user_id") REFERENCES "user" ("id");

CREATE INDEX ON "sessions" ("user_id");

CREATE TABLE IF NOT EXISTS "client" (
  "id" BIGSERIAL PRIMARY KEY,
  "name" VARCHAR(128) NOT NULL
);

CREATE INDEX ON "client" ("id");

CREATE TABLE IF NOT EXISTS "events" (
  "id" BIGSERIAL PRIMARY KEY,
  "session_id" BIGINT,
  "client_id" BIGINT,
  "event" VARCHAR(255) NOT NULL,
  "user_agent" TEXT NOT NULL,
  "ip" TEXT NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE "events" ADD FOREIGN KEY ("session_id") REFERENCES "sessions" ("id");

ALTER TABLE "events" ADD FOREIGN KEY ("client_id") REFERENCES "client" ("id");

CREATE INDEX ON "events" ("session_id");

CREATE INDEX ON "events" ("client_id");

COMMENT ON COLUMN "user"."password" IS 'bcrypt';

COMMENT ON COLUMN "user"."verified" IS 'true verify, false not verify';

COMMENT ON COLUMN "user"."created_at" IS 'full RFC3339 format';

COMMENT ON COLUMN "user"."updated_at" IS 'full RFC3339 format';

COMMENT ON COLUMN "user"."deleted_at" IS 'full RFC3339 format';

COMMENT ON COLUMN "password"."password" IS 'bcrypt';

COMMENT ON COLUMN "password"."created_at" IS 'full RFC3339 format';

COMMENT ON COLUMN "profile"."updated_at" IS 'full RFC3339 format';

COMMENT ON COLUMN "tfa"."secret" IS 'secret code tfa';

COMMENT ON COLUMN "tfa"."enable" IS 'true require, false not require';

COMMENT ON COLUMN "tfa"."enable_at" IS 'full RFC3339 format';

COMMENT ON COLUMN "sessions"."is_current" IS 'true login, false not login';

COMMENT ON COLUMN "events"."created_at" IS 'full RFC3339 format';

COMMENT ON COLUMN "events"."updated_at" IS 'full RFC3339 format';

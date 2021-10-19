CREATE TABLE "auth" (
  "id" BIGSERIAL PRIMARY KEY,
  "email" VARCHAR(128) UNIQUE NOT NULL,
  "password" TEXT NOT NULL,
  "verified" BOOLEAN NOT NULL DEFAULT FALSE,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "users" (
  "id" BIGSERIAL PRIMARY KEY,
  "auth_id" BIGINT,
  "fullname" VARCHAR(128) NOT NULL,
  "location" VARCHAR(128),
  "bio" TEXT,
  "web" VARCHAR(128),
  "picture" VARCHAR(128),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "tfa" (
  "id" BIGSERIAL PRIMARY KEY,
  "auth_id" BIGINT,
  "secret" VARCHAR(128),
  "enable" BOOLEAN NOT NULL DEFAULT FALSE,
  "enable_at" timestamptz
);

CREATE TABLE "tfa_codes" (
  "id" BIGSERIAL PRIMARY KEY,
  "tfa_id" BIGINT,
  "code" VARCHAR(128) NOT NULL
);

CREATE TABLE "sessions" (
  "id" BIGSERIAL PRIMARY KEY,
  "auth_id" BIGINT,
  "is_current" BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE TABLE "client" (
  "id" BIGSERIAL PRIMARY KEY,
  "name" VARCHAR(128) NOT NULL
);

CREATE TABLE "events" (
  "id" BIGSERIAL PRIMARY KEY,
  "session_id" BIGINT,
  "client_id" BIGINT,
  "event" VARCHAR(255) NOT NULL,
  "user_agent" TEXT NOT NULL,
  "ip" TEXT NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE "users" ADD FOREIGN KEY ("auth_id") REFERENCES "auth" ("id");

ALTER TABLE "tfa" ADD FOREIGN KEY ("auth_id") REFERENCES "auth" ("id");

ALTER TABLE "tfa_codes" ADD FOREIGN KEY ("tfa_id") REFERENCES "tfa" ("id");

ALTER TABLE "sessions" ADD FOREIGN KEY ("auth_id") REFERENCES "auth" ("id");

ALTER TABLE "events" ADD FOREIGN KEY ("session_id") REFERENCES "sessions" ("id");

ALTER TABLE "events" ADD FOREIGN KEY ("client_id") REFERENCES "client" ("id");

CREATE INDEX ON "auth" ("email");

CREATE INDEX ON "users" ("auth_id");

CREATE INDEX ON "tfa" ("auth_id");

CREATE INDEX ON "tfa_codes" ("tfa_id");

CREATE INDEX ON "sessions" ("auth_id");

CREATE INDEX ON "client" ("id");

CREATE INDEX ON "events" ("session_id");

COMMENT ON COLUMN "auth"."password" IS 'bcrypt';

COMMENT ON COLUMN "auth"."verified" IS 'true verify, false not verify';

COMMENT ON COLUMN "auth"."created_at" IS 'full RFC3339 format';

COMMENT ON COLUMN "auth"."updated_at" IS 'full RFC3339 format';

COMMENT ON COLUMN "users"."updated_at" IS 'full RFC3339 format';

COMMENT ON COLUMN "tfa"."secret" IS 'secret code tfa';

COMMENT ON COLUMN "tfa"."enable" IS 'true require, false not require';

COMMENT ON COLUMN "tfa"."enable_at" IS 'full RFC3339 format';

COMMENT ON COLUMN "sessions"."is_current" IS 'true login, false not login';

COMMENT ON COLUMN "events"."created_at" IS 'full RFC3339 format';

COMMENT ON COLUMN "events"."updated_at" IS 'full RFC3339 format';

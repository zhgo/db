DROP TABLE IF EXISTS "passport_login";
CREATE TABLE "passport_login" (
    "LoginID" SERIAL PRIMARY KEY,
    "UserID" integer NOT NULL,
    "CreationTime" text NOT NULL,
    "Source" integer NOT NULL,
    "LoginIp" integer NOT NULL,
    "AnonymousID" text NOT NULL,
    "AuthCode" text NOT NULL,
    "UserAgent" text NOT NULL
);

DROP TABLE IF EXISTS "passport_user";
CREATE TABLE "passport_user" (
    "UserID" SERIAL PRIMARY KEY,
    "CreationTime" text NOT NULL,
    "BirthYear" integer NOT NULL,
    "Gender" text NOT NULL,
    "Nickname" text NOT NULL
);

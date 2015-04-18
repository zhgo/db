DROP TABLE IF EXISTS "passport_login";
CREATE TABLE "passport_login" (
  "LoginID" INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
  "UserID" INTEGER NOT NULL,
  "CreationTime" TEXT NOT NULL,
  "Source" INTEGER NOT NULL DEFAULT "1",
  "LoginIp" INTEGER NOT NULL,
  "AnonymousID" TEXT NOT NULL,
  "AuthCode" TEXT NOT NULL,
  "UserAgent" TEXT NOT NULL
);

UPDATE "sqlite_sequence" SET "seq" = 1000000 WHERE "name" = "passport_login";

DROP TABLE IF EXISTS "passport_user";
CREATE TABLE "passport_user" (
  "UserID" INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
  "CreationTime" TEXT NOT NULL,
  "BirthYear" INTEGER NOT NULL,
  "Gender" TEXT NOT NULL DEFAULT "Secret",
  "Nickname" TEXT NOT NULL
);

UPDATE "sqlite_sequence" SET "seq" = 1000000 WHERE "name" = "passport_user";

CREATE TABLE IF NOT EXISTS "app_version" (
    "id"                SERIAL NOT NULL,
    "android_version"   VARCHAR(255) NOT NULL,
    "ios_version"       VARCHAR(255) NOT NULL,
    "is_force_update"   BOOLEAN NOT NULL,
    "created_at"        TIMESTAMP(0) WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    "updated_at"        TIMESTAMP(0) WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

insert into app_version (android_version, ios_version, is_force_update) values ('1.0.0', '1.0.0', false);

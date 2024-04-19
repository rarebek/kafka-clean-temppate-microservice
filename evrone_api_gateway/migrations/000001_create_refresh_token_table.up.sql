CREATE TABLE IF NOT EXISTS "refresh_tokens" (
    "guid"          UUID NOT NULL,
    "refresh_token" CHARACTER VARYING(200) DEFAULT '',
    "expiry_date"   TIMESTAMP(0) WITH TIME ZONE,
    "created_at"    TIMESTAMP(0) WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT refresh_token_guid_pkey PRIMARY KEY (guid)
);

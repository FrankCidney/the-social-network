-- UP
CREATE TYPE privacy_level AS ENUM ('public', 'almost_private', 'private', 'group');
CREATE TYPE follower_status AS ENUM ('pending', 'accepted');
CREATE TYPE group_member_status AS ENUM ('invited', 'requested', 'accepted');
CREATE TYPE rsvp_status AS ENUM ('going', 'not_going');
CREATE TYPE notification_type AS ENUM ('follow_request', 'group_invite', 'group_request', 'event_created');

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    dob DATE NOT NULL,
    avatar_url VARCHAR(500),
    nickname VARCHAR(100),
    about_me TEXT,
    is_public BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

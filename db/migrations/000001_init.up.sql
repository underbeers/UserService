CREATE TABLE user_profile
(
    id          UUID         NOT NULL PRIMARY KEY,
    first_name  VARCHAR(255) NOT NULL,
    second_name VARCHAR(255),
    sur_name    VARCHAR(255),
    status      INT          NOT NULL
);

CREATE TABLE user_data
(
    id               SERIAL PRIMARY KEY,
    id_profile       UUID UNIQUE REFERENCES user_profile (id),
    password_encoded VARCHAR(255),
    password_salt    VARCHAR(255)
);

CREATE TABLE user_contacts
(
    id                 SERIAL PRIMARY KEY,
    id_profile         UUID UNIQUE REFERENCES user_profile (id),
    push_notifications BOOLEAN,
    email              VARCHAR(255),
    email_subscription BOOLEAN,
    mobile_phone       VARCHAR(255) NOT NULL,
    show_phone         BOOLEAN
);

CREATE TABLE verification
(
    id                    SERIAL PRIMARY KEY,
    id_contacts           INT UNIQUE,
    sms_verification_code CHAR(6),
    sms_code_expiration   TIMESTAMP(0) WITHOUT TIME ZONE,
    block_expiration      TIMESTAMP(0) WITHOUT TIME ZONE
);
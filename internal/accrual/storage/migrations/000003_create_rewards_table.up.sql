create table if not exists rewards(
    match text unique not null,
    reward decimal not null,
    reward_type varchar not null
)
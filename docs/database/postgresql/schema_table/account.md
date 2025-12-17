# account schema

__*TABLES:*__

`account.user`
```sql
/*
LAST UPDATED: YYYY-MM-DD

this table represent core account user data

---

note:
- password_hash:
    - using argon2id

---

after creation:
    - n/a

*/
create table if not exists account.user(
    id              uuid        unique not null primary key default uuidv7(),
    email           text        unique not null,
    password_hash   text        not null,
    dt_created      timestamp   null default now(),
    dt_updated      timestamp   null
);

-- alter table

-- indexes
create index if not exists idx_account_user_email on account.user(email);

-- functions
create or replace function account.user_dt_updated()
    returns trigger as $$
    begin
        new.dt_updated = now();

        return new;
    end;
    $$ language plpgsql;

-- triggers
create or replace trigger account_user_dt_updated_trigger
    before update on account.user
    for each row
    execute function account.user_dt_updated();
```

<br>

---

###### end of account


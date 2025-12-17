# game1 schema

__*TABLES:*__

`game1.stash`
```sql
/*
LAST UPDATED: YYYY-MM-DD

assume there is a game called `game1`
and this table represent that in-game stash from game1

---

note:
- name_norm:
    - normalized name:
        - all lower case
- items:
    - format:
    [
        // ...
        {
            // could be add item_id if has unique identifier
            "item": "Foo Bar Item",
            "quantity": 1
        }
        // ...
    ]

---

after creation:
    - n/a
*/
create table if not exists game1.stash(
    id       	uuid        unique not null primary key default uuidv7(),
	uid			uuid		not null,
    name        text        not null,
    name_norm   text        not null,
    items       jsonb       null,
    dt_created  timestamp   null default now(),
    dt_updated  timestamp   null
);

-- alter table
do $$
begin
    if not exists (
        select 1 from information_schema.table_constraints
        where table_schema = 'game1'
        and table_name = 'stash'
        and constraint_name = 'fk_id_game1_stash_id'
    ) then
        alter table game1.stash
            add constraint fk_id_game1_stash_id
            foreign key (uid)
            references account.user (id)
            on delete no action
            on update cascade;
    end if;
end $$;

-- indexes
create index if not exists idx_game1_stash_uid on game1.stash(uid);
create index if not exists idx_game1_stash_name_norm on game1.stash(name_norm);

-- functions
create or replace function game1.stash_dt_updated()
    returns trigger as $$
    begin
        new.dt_updated = now();

        return new;
    end;
    $$ language plpgsql;

-- triggers
create or replace trigger game1_stash_dt_updated_trigger
    before update on game1.stash
    for each row
    execute function game1.stash_dt_updated();
```

<br>

---

###### end of game1 schema


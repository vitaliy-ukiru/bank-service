BEGIN;
create table accounts
(
    id      integer generated always as identity,
    balance double precision
);
COMMIT;
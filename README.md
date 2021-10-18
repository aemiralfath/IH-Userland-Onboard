# ![Ice House Logo](assets/logo-white.png)

## IH-Userland-Onboard

Userland is account self-management system for Ice House Onboarding Project

## API Contract

[https://userland.docs.apiary.io/#introduction/http-status-codes](https://userland.docs.apiary.io/#introduction/http-status-codes)

## Data Modeling

### Postgre Schema

1. Auth Table

   Fields:
   - email - string (<=128 chars, primary key, unique)
   - password - string (bcrypt, 6-128 chars)
   - verify - boolean (true->verify, false->not verify)
   - updated_at - string (full RFC3339 format)

2. User Table

   Fields:
   - id - int (primary key)
   - fullname - string (3-128 chars)
   - location - string (<=128 chars)
   - bio - string (<=256 chars)
   - web - string (<=128 chars)
   - picture - string (<=128 chars)
   - created_at - string (full RFC3339 format)
   - updated_at - string (full RFC3339 format)
   - auth_id - string (foreign key, auth table primary key)

3. TFA Table

   Fields:
   - id - int (primary key)
   - enable - boolean (true->tfa require, false->tfa not require)
   - ? secret - string (<=128 chars, secret code TFA)
   - enabled_at - string (full RFC3339 format)
   - auth_id - string (foreign key, auth table primary key)

4. TFA Backup Code Table

   Fields:
   - id - int (primary key)
   - code - string (<=128 chars)
   - tfa_id - int (foreign key, tfa table primary key)

5. Events Table

   Fields:
   - id - int (primary key)
   - event - string (<=128 chars)
   - user_agent - string (<=256 chars)
   - ip - string (<=128 chars)
   - created_at - string (full RFC3339 format)
   - updated_at - string (full RFC3339 format)
   - client_id - int (foreign key, client table primary key)

6. Client Table

   Fields:
   - id - int (primary key)
   - name - string (<=128 chars)

7. Session Table

   Fields:
   - id - int (primary key)
   - is_current - boolean (true->login, false->not login)
   - event_id - int (foreign key, event table primary key)
   - auth_id - string (foreign key, auth table primary key)

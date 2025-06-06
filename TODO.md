## TODO LIST

- Auth
    - MH user/pass credentials are not needed for the login, consider removing them
    - JWT is being logged when hitting login endpoint, this is a security concern

- Encryption at rest
    - Passwords are not encrypted: users.password_pri
    - Should encrypt PI data as well?
    - Branch offices api_key and api_secret are not encrypted

- DB
    - Check if the search criteria columns have proper indexes
# API Documentation

## Account APIs

| API                         | Method | Token Required | Request                   | Response                |
|-----------------------------|--------|----------------|---------------------------|--------------------------|
| `/v1/account/login/admin`   | POST   | ❌             | `username`, `password`    | `token`                  |
| `/v1/account/login/user`    | POST   | ❌             | `username`, `password`    | `token`                  |
| `/v1/account/signup`        | POST   | ❌             | `username`, `password`, `name` | `message`             |
| `/v1/account/change-password` | POST | ✅             | `(new) password`          | `message` and `account data` |

## User APIs

| API                         | Method | Token Required | Request                   | Response                |
|-----------------------------|--------|----------------|---------------------------|--------------------------|
| `/v1/user/profile`          | GET    | ✅             | -                         | `user data`             |
| `/v1/user/mutation/transaction` | GET | ✅             | -                         | `all transactions of user` |
| `/v1/user/mutation/deposit` | GET    | ✅             | -                         | `list deposito`         |
| `/v1/user/edit/profile`     | POST   | ✅             | `address`, `id_card`, `mothers_name`, `date_of_birth`, `gender` | `message` and `user data` |
| `/v1/user/register/deposit` | POST   | ✅             | `deposit_id`, `account_id`, `name`, `amount`, `min_amount` | `message`             |

## Admin APIs

| API                         | Method | Token Required | Request                   | Response                |
|-----------------------------|--------|----------------|---------------------------|--------------------------|
| `/v1/admin/list/user`       | GET    | ❌             | -                         | `list all users`        |
| `/v1/admin/list/user/:id`   | GET    | ❌             | -                         | `user data based on ID` |
| `/v1/admin/list/deposit/mutation` | GET | ❌          | -                         | `list all deposit mutations` |
| `/v1/admin/topup`           | POST   | ❌             | `username`, `amount`      | `message` and `user balance` |

---

- **Token Required**: Indicates if a token is required for the API endpoint. 
  - ✅: Token required.
  - ❌: Token not required.
- **Request**: Parameters to be sent in the request body or URL.
- **Response**: Expected response from the server.

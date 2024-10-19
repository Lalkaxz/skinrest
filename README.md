# Rest Api Documentation

## Usage
### Endpoints:

### /api/v1:
- [`GET: /`](#get--health-check)
- [`POST: /user/register`](#post-userregister-register-new-user)
- [`POST: /user/login`](#post-userlogin-login-as-user)
- [`GET: /user/me`](#get-userme-get-info-about-current-user)
- [`POST: /skins/add`](#post-skinsadd-add-skin-in-collection)
- [`GET: /skins`](#get-skins-get-user-skins-collection)
- [`GET: /skins/:id`](#get-skinsid-get-skin-information)
- [`DELETEs: /skins/:id`](#delete-skinsid-delete-skin)


## `GET: /`: Health check

### Response
```json
{
    "status": "ok"
}
```

## `POST: /user/register`: Register New User

### Request Headers:
```
    Content-Type: application/json
```

### Request Body:
```json
{
    "login": "John",
    "password": "123"
}
```
### Response:
### With status code 200
```json
    "token": "token-here"
```


## `POST: /user/login`: Login as user 

### Request Headers:
```
    Content-Type: application/json
```

### Request Body:
```json
{
    "login": "John",
    "password": "123"
}
```
### Response:
### With status code 200
```json
    "token": "token-here"
```

### if token is expired:
### With status code 200
```json
{
    "token": "new-token-here"
}
```



## `GET: /user/me`: Get info about current user

### Request Headers:
```
    Authorization: Bearer (ur-token-here)
```

### Response Body:
```json
{
    "Login": "john",
    "Skins": [
        {
            "Id": 1,
            "Name": "Aid",
            "Type": "Slim",
            "Src": "mojang-nickname-or-url"
        }
    ]
}
```


## `GET: /skins`: Get user skins collection

### Request Headers:
```
    Authorization: Bearer (ur-token-here)
```

### Response Body:
```json
[
    {
            "Id": 1,
            "Name": "Aid",
            "Type": "Slim",
            "Src": "mojang-nickname-or-url"
        },
    {
        "Id": 2,
        "Name": "Notch",
        "Type": "Classic",
        "Src": "mojang-nickname-or-url"
        }
]
```

## `POST: /skins/add`: Add skin in collection

### Request Headers:
```
    Authorization: Bearer (ur-token-here)
```

### Request Body:
```json
{
    "skinname": "Aid",
    "skintype": "Slim",
    "skinsrc": "mojang-nickname-or-url"
}
```

### Response Body:

### With status 201 Created:
```json
{
    "Id": 1,
    "Name": "Aid",
    "Type": "Slim",
    "Src": "mojang-nickname-or-url"
}
```


## `GET: /skins/:id`: Get skin information

### Request Headers:

```
    Authorization: Bearer (ur-token-here)
```

### Response Body:

### With status 200 Ok:
```json
{
    "Id": 1,
    "Name": "Aid",
    "Type": "Slim",
    "Src": "mojang-nickname-or-url"
}
```


## `DELETE: /skins/:id`: Delete skin 

### Request Headers:

```
    Authorization: Bearer (ur-token-here)
```

### With status 200 Ok:
```json
{
    "status": "Success"
}
```





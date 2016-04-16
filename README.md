# Prix

Introduction to the API

# API

#### [GET] /
Return the index.html dub page.


#### [GET] /api/ping
This route is used to check the service status.

- Return:
```javascript
{
  "message":"pong"
}
```


#### [POST] /api/login
This route is used to authenticate an user, it returns an valid token to make authenticated api calls.

- Request:
```javascript
{
  "email":"user@prix.plus",
  "password":"pass"
}
```

- Return:
```javascript
{
  "expire":"2016-05-15T01:43:07-04:00",
  "token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE0NjMyOTA5ODcsImlkIjoxfQ.ROZ9l2I41QE3Mz9jhJdLmqHAQpQr5SazzCU7q-8WSnk"
}
```


#### [GET] /api/refresh_token
This route is used to refresh a valid token.

- Request Header: `Authorization: Bearer TOKEN`

- Return:
```javascript
{
  "expire":"2016-05-15T02:17:12-03:00",
  "token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE0NjMyODk0MzIsImlkIjoxfQ.8Rz7x1s7CJ4xZ-PDomuV8bAvgmhIp6nSoPDjfJ2Bha0"
}
```


#### [GET] /api/me
This route is used to retrieve the current user in the session.

- Request Header: `Authorization: Bearer TOKEN`

- Return:
```javascript
{
  "results":[
    {
      "id":1,
      "password":"pass",
      "email":"user@prix.plus"
    }
  ]
}
```


#### [POST] /api/users
This route is used to create a new user.

- Request:
```javascript
{
  "email":"newuser@prix.plus",
  "password":"pass"
}
```

- Return Header: `Location: /api/users/2`
- Return:
```javascript
{
  "location":"/api/users/2",
  "results":[
    {
      "id":2,
      "password":"pass",
      "email":"newuser@prix.plus"
    }
  ]
}
```


#### [PUT] /api/users/:id
This route is used to update users info. You can't use an `id` differente from yours.

- Request Header: `Authorization: Bearer TOKEN`

- Request:
```javascript
{
  "id":1, // Optional
  "password":"pass",
  "email":"newuser@prix.plus"
}
```

- Return:
```javascript
{
  "results":[
    {
      "id":1,
      "password":"pass",
      "email":"newuser@prix.plus"
    }
  ]
}
```
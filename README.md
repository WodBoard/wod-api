# wod-api
The API of Wodboard is the main server of our working of the day app.

## Routes

Publicly exposes the REST API routes used by the app. Responses are formatted in JSON.

## Database

Contains various methods to handle the storage of our different models in a mongodb database.

# Getting started

Provided examples make use of the command line tool `http` (https://httpie.org/) instead of the commonly used `curl`.

## Setup

First run `docker-compose up` or `docker-compose up -d` in order to run the server.

## Authentication

Now you will need to register a new user on WodBoard, for you to interact with all endpoints.
After your registration you will have to provide an auth token to access all authenticated routes, such as `/hello` which gives basic informations about the connected user.
You can get a token by logging in.

### Signup

To signup, you need to do as shown below:
```
http --json POST http://localhost:4242/signup email=patrice@gmail.com password=mdp123 firstname=patrice lastname=michel weight=3.14 height=6.7 picture_url=prout.png birthday=1994-12-31T00:00:00Z
```

The server should return you en empty response with an ok status (200), which means the user `patrice@gmail.com` is now successfully registered.
```
HTTP/1.1 200 OK
Content-Length: 0
Date: Thu, 30 May 2019 19:11:58 GMT
```

### Login

To login, you need to do as shown below:
```
http --json POST http://localhost:4242/login email=patrice@gmail.com password=mdp123
```

If you successfully logged in, the server should reply you with a status 200 and a json response containing your token.
```
HTTP/1.1 200 OK
Content-Length: 223
Content-Type: application/json; charset=utf-8
Date: Thu, 30 May 2019 19:22:25 GMT

{
    "code": 200,
    "expire": "2019-05-31T19:22:25Z",
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NTkzMzA1NDUsImlkIjoicGF0cmljZUBnbWFpbC5jb20iLCJvcmlnX2lhdCI6MTU1OTI0NDE0NX0.CCP0IzIHZEjhdQso6KW-Z_kad0V0otCpNAjVAlUhztw"
}
```

### Hello

To test if everything worked fine, you can try calling the `/hello` endpoint which is one of the authenticated route that requires a token.
To do so you need to include the token prefixed by `"Bearer"` (`"Bearer <your-token>"`) in the `Authorization` header:

```
http GET "http://localhost:4242/hello" "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NTkzMjg5MjksImlkIjoicGF0cmljZUBnbWFpbC5jb20iLCJvcmlnX2lhdCI6MTU1OTI0MjUyOX0.iQEZbeD6GyZ9EsTeXFq574UhW4mBPGq-JRut0s4QvG4"
```

If you specified a correct token, the http response should have a 200 status and the json content should be like the following:
```
HTTP/1.1 200 OK
Content-Length: 185
Content-Type: application/json; charset=utf-8
Date: Sun, 16 Jun 2019 17:30:52 GMT

{
    "birthday": {
        "seconds": 788832000
    },
    "email": "patrice@gmail.com",
    "firstname": "patrice",
    "height": 6.699999809265137,
    "lastname": "michel",
    "picture_url": "prout.png",
    "weight": 3.140000104904175
}

```
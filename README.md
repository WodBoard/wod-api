# wod-api
The API of Wodboard is the main server of our working of the day app.

## Routes

Publicly exposes the REST API routes used by the app. Responses are formatted in JSON.

*Public:*
* `POST` "/login" - Login route, takes `email` and `password` as arguments.
* `POST` "/signup" - Signup route, takes `email`, `password` and user informations.

*Authenticated:*
* `GET` "/profile" - Profile is an endpoint that returns the currently logged in user's informations.
* `PUT` "/profile" - Profile is an endpoint that returns the currently logged in user's informations.
* `GET` "/trainings" - List user's trainings.
* `PUT` "/trainings" - Add a new training to the users' trainings (`name` of the training is unique to the user).
* `POST` "/trainings" - Edit and updates a existing training by its `name`.
* `DELETE` "/trainings" - Delete an existing training in the user training list.

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

### Profile

#### Get your profile

To test if everything worked fine, you can try calling the `/profile` endpoint which is one of the authenticated route that requires a token.
To do so you need to include the token prefixed by `"Bearer"` (`"Bearer <your-token>"`) in the `Authorization` header:

```
http GET "http://localhost:4242/profile" "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NTkzMjg5MjksImlkIjoicGF0cmljZUBnbWFpbC5jb20iLCJvcmlnX2lhdCI6MTU1OTI0MjUyOX0.iQEZbeD6GyZ9EsTeXFq574UhW4mBPGq-JRut0s4QvG4"
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
    "picture_url": "test.png",
    "weight": 3.140000104904175,
    "affiliated_box": "La salle de Patrice"
}
```

#### Edit your profile

To edit your profile, you need to call `PUT /profile` and insert a `User` structure, it will return you a `200` if everything worked fine.

### Trainings

The following endpoints will let you create a new personal training, and also fetch a list of them.

#### Add a new training

This is a `POST` request with a json body containing your new training characteristics.
Json body (./testdata/add_training.json):
```
cat testdata/add_training.json 
{
   "name":"MyFirstTrainingEver",
   "type":3,
   "exercises":[
      {
         "movement":1,
         "name":"Yoga"
      },
      {
         "movement":2,
         "name":"Velo"
      },
      {
         "movement":2,
         "name":"Tapis"
      }
   ],
   "time_cap":6000
}
```

The actual request to insert it in our database:
```
http --json POST "http://localhost:4242/trainings" "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NjA4ODkxMjAsImlkIjoicGF0cmljaW9AZ21haWwuY29tIiwib3JpZ19pYXQiOjE1NjA4MDI3MjB9.LYJW3Oy1kaG2-GoH2UXCF1Xk2AGv4O0dx-j4MsFlt1Q" < ./testdata/add_training.json
HTTP/1.1 200 OK
Content-Length: 157
Content-Type: application/json; charset=utf-8
Date: Mon, 17 Jun 2019 21:17:56 GMT

{
    "exercises": [
        {
            "movement": 1,
            "name": "Yoga"
        },
        {
            "movement": 2,
            "name": "Velo"
        },
        {
            "movement": 2,
            "name": "Tapis"
        }
    ],
    "name": "MyFirstTrainingEver",
    "time_cap": 6000,
    "type": 3
}
```

- If everything went perfectly, you should be returned an http status of `200` and a json object corresponding to the one you just sent.
- If not, it is probably that you tried to add an already existing training, and you should be returned a status `400 - Bad request`.

#### Edit an existing training

Same as the previous one but lets you update an existing training (can also create it).
This is a `PUT` request with a json body containing your new training characteristics.
Json body (./testdata/add_training.json):
```
cat testdata/add_training.json 
{
   "name":"MyFirstTrainingEver",
   "type":3,
   "exercises":[
      {
         "movement":1,
         "name":"Yoga"
      },
      {
         "movement":2,
         "name":"Velo"
      },
      {
         "movement":2,
         "name":"Tapis"
      }
   ],
   "time_cap":6000
}
```

The actual request to edit it in our database:
```
http --json PUT "http://localhost:4242/trainings" "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NjA4ODkxMjAsImlkIjoicGF0cmljaW9AZ21haWwuY29tIiwib3JpZ19pYXQiOjE1NjA4MDI3MjB9.LYJW3Oy1kaG2-GoH2UXCF1Xk2AGv4O0dx-j4MsFlt1Q" < ./testdata/add_training.json
HTTP/1.1 200 OK
Content-Length: 157
Content-Type: application/json; charset=utf-8
Date: Mon, 17 Jun 2019 21:17:56 GMT

{
    "exercises": [
        {
            "movement": 1,
            "name": "Yoga"
        },
        {
            "movement": 2,
            "name": "Velo"
        },
        {
            "movement": 2,
            "name": "Tapis"
        }
    ],
    "name": "MyFirstTrainingEver",
    "time_cap": 6000,
    "type": 3
}
```

If everything went perfectly, you should be returned an http status of `200` and a json object corresponding to the one you just sent.

#### List your trainings

List your trainings is a simple `GET` endpoints to retrieve all trainings bound to your account:
```
http GET "http://localhost:4242/trainings" "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NjA4ODkxMjAsImlkIjoicGF0cmljaW9AZ21haWwuY29tIiwib3JpZ19pYXQiOjE1NjA4MDI3MjB9.LYJW3Oy1kaG2-GoH2UXCF1Xk2AGv4O0dx-j4MsFlt1Q"
HTTP/1.1 200 OK
Content-Length: 633
Content-Type: application/json; charset=utf-8
Date: Mon, 17 Jun 2019 21:30:51 GMT

[
    {
        "exercises": [
            {
                "movement": 1,
                "name": "Yoga"
            },
            {
                "movement": 2,
                "name": "Velo"
            },
            {
                "movement": 2,
                "name": "Tapis"
            }
        ],
        "name": "MyFirstTrainingEver",
        "time_cap": 6000,
        "type": 3
    }
]
```

#### Delete an existing Training

To delete an existing training, you have to do a `DELETE` request to the "/trainings" endpoint.
You need to provide a json body containing the name of the training.

```
http --json DELETE "http://localhost:4242/trainings" "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NjA5NzgwMDQsImlkIjoicGF0cmljaW9AZ21haWwuY29tIiwib3JpZ19pYXQiOjE1NjA4OTE2MDR9.YtAtLG_vZjqlTZalIeDGJGYx5ULmK6wLH0wsQLVeJBk" name=MyFirstTrainingEver    
HTTP/1.1 200 OK
Content-Length: 0
Date: Tue, 18 Jun 2019 21:20:24 GMT
```

If everything worked fine, you should receive a status `200` and the training should be removed from the database.
## Implementation details
Libraries that I've used in the implementation
- github.com/gorilla/mux - Allows easy setup for API routes for a webapp.
- github.com/jmoiron/sqlx - Lightweight ORM
- github.com/go-sql-driver/mysql - mysql driver

## Components

The follow sections detail the requirements for the three portions of the
application.

### Backend API

The api is located under the `api` directory and was generated with [cookie
cuttter django](https://github.com/pydanny/cookiecutter-django).

For this component, we're going to build out the api routes to interact with
the profiles using the Django Rest Framework. Some code is included to get you
started. That code is found in the `api/api/profilesapp` directory.

#### MUST HAVES

- Profile model/migration (see the [docs](https://docs.djangoproject.com/en/3.0/topics/migrations/) for info)
    * a model has been included but it does not contain all of the fields that we need to return in our api
- Viewset providing the following api routes (see the [docs](https://www.django-rest-framework.org/api-guide/viewsets/#modelviewset) for info)
    * list
        + working pagination
        + includes the following fields:
            - Full name
            - Current title
            - Total years of work experience (the sum of all positions excluding certain overlaps)
            - Has CS degree?
            - Is currently employed? (i.e. their most recent position has no end date)
        + supports the following filters:
            - Text search (over reasonable subset of model fields)
    * detail (route for getting a single profile based on a given identifier)
- Serializer to translate the database model representation into what's returned by the api (see the [docs](https://www.django-rest-framework.org/api-guide/serializers/#modelserializer) for more info)

#### NICE TO HAVES

- Tested
- Include support for the following filters
    * Range search over "Total years of work experience"
    * Has CS degree?
    * Is currently employed?
- Include support for the sortable columns
    * Full name
    * Current title
    * Total years of work experience
- Add a notes field to the profile and add an update api route
- Model includes indices to speed up common query patterns (performance won't
  be an issue with the test data but imagine if we had hundreds of millions of
  profiles)

### Backend Ingest

The ingest code is located with the api under the
`api/api/profiles/app/management/commands` directory.

For this component, we're going to write an ingest command that reads a
newline-delimited json file containing "raw" profiles, processes the profiles,
and persists them to the database, using the model you wrote in the last
section. You'll need to extract the features required for sorting or filtering
(unless you plan to do that on the fly in the API (not recommended)).

#### MUST HAVES

- Supports creating new profiles in the db
- Supports updating existing profiles in the db

#### NICE TO HAVES

- Tested
- Writes in a performant way to the database (e.g. maybe in bulk)


### Ingesting profiles

```
$ docker-compose -f local.yml run --rm django python manage.py ingestprofiles ./test_data/test_profiles.json
```

### View API docs

Go to [http://localhost:8000/redoc/](http://localhost:8000/redoc/)

### View API Response

List view [http://localhost:8000/api/profiles](http://localhost:3000/api/profiles)

Detail view [http://localhost:8000/api/profiles/<id>](http://localhost:3000/api/profiles/<id>)

## Evaluation Criteria

First and foremost, we are looking for a clean, well-designed solution that
accomplishes the "MUST HAVE" items listed above. Some things to keep in mind
while building:

- Extensible — can my be code be easily modified to do something slightly different?
- Testable — does it work? How can we be sure? If changes are made, how would we detect a regression?
- Maintainable — could someone that isn't you come in and make improvements? Are more complicated sections commented?
- Deployment — how do we get this thing into production? Do we need to run any db migrations? Update configurations? Ect.

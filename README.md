## Currency Conversion Service

### About Project
That project is responsible for authenticate users and manage 
currency conversions between authenticated users currency wallets.

### How To Build

To build the project as an executable file you can simply run
the following make command;
```shell
make build
```

or you can build the project as docker image with the following
command;
```shell
docker build -t ${docker_username}/${registry_name}:${image_tag}
```

### How To Test

Before you test the application you need to generate needed mocks
by executing following make command;

````shell
make generate-mocks
````

After you generate mocks, you can simply run following make command
to run all tests;

````shell
make unit-tests
````

### How To Run
Before you run the project, you need to stand the database up.
You can simply use the following docker-compose command;
````shell
docker-compose up -d
````

After you stand up the database, you can either run the project
as docker container or an executable file.

To run the project as executable file, you can use following 
make command combinations;

````shell
make build
make run
````
to first build then run. Or directly use following make command 
to build and run;
````shell
make start
````

To run the project as docker container, you can use following docker
commands;
````shell
docker build -t ${docker_username}/${registry_name}:${image_tag}
docker run -p ${host_port}:4001 ${docker_username}/${registry_name}:${image_tag}
````
### Cleanup The Project

You can use following make command to get rid of unnecessary
files generated during build or test.
````shell
make clean
````
### Handle Migrations

I wrote special seed commands both supports drop all tables
from database by following make command;
````shell
make drop-tables
````
and creating all tables on database with static and example data
````shell
make migrate-tables
````

### Swagger

You can see the swagger documentation of current endpoints after
you start the application in the route;
`/swagger/index.html`


### Linter

I use default golangci-lint linter specifications for linting.
You can use the following make command to check any lint issues;
````shell
make linter
````
## UserService
1. [First run](#first-run)
2. [Regular run](#regular-run)
3. [Run locally](#run-locally)
4. [Before push](#before-push)

[Technical requirements/UserService](https://wiki.andersenlab.com/display/GOR/1+USER-Service "Technical requirements/USER-Service")

[Diagrams & Schemas/UserService](https://dbdiagram.io/d/6396187bbae3ed7c454619b9 "Diagrams & Schemas/USER-Service")
 
### First run
1) Clone repo

```git clone https://github.com/underbeers/UserService.git```

2) Go to the project folder

```cd UserService```

3) Create an env file. yourpassword replace with password you would like to use.

```
UNIX/MAC: 
	echo 'POSTGRES_PASSWORD=yourpassword'>db.env
Windows:
	echo POSTGRES_PASSWORD=yourpassword >db.env
```
Add in db.env ```SECRET_JWT=SomeSecretStringforGeneratingJwt!!!!```

4) Create DB container with ```make db_unix``` If you are on Windows and have problems with this command try ```make db_win```

5) Check for pgsql is running  ```docker ps``` Look in column NAMES for pgsql (If not rey to start it with ```docker start pgsql```

6) Connect to pgsql container ```docker exec -it pgsql sh```

7) Change current user ```su postgres```

8) Start psql ```psql```

9) Create new DB ```CREATE DATABASE user_service;```

10) List all DB's look for user_service ```\l```

11) Return to the terminal type ```exit``` few times

12) Migrate data ```make migrate```

13) Open your favorite DB Manager, and connect to DB localhost:5430 (I am using DBeaver)

14) Insert dummy data in to the user_service DB


15) Add to the userservice/conf/db/local.yaml password from db.env (if you are start service locally)
16) Build user_service image ```make build_image```

17) Check for pgsql and api_gateway is running  ```docker ps``` if not, start them ```docker start pgsql```
    (For the api_gateway go to <a href="https://git.andersenlab.com/Andersen/repo-gor/apigateway#regular-run">documentation</a>)

18) Start user_service ```make run```

### Regular run
1) Start pgsql and api_gateway containers if not. ```docker start pgsql```
   (For the api_gateway go to <a href="https://git.andersenlab.com/Andersen/repo-gor/apigateway#regular-run">documentation</a>)
2) Start user_service ```make run```


### Run locally
1) In conf/db/local.yml replace string ```password: yourpassword``` -> ```password: password from the db.env file```

2) Start user_service ```make local```

### Before push
1) Run linters. In root directory run ```golangci-lint run```

2) Replace password from userservice/conf/db/local.yaml with ```yourpassword```

# Golang-API-Assessment

## Table of Contents  
- [About](#about) <a name="about"/>  
- [Deployed API](#deployed-api) <a name="deployed-api"/>
- [Project Structure](#project-structure) <a name="project-structure"/>
- [API Endpoints](#api-endpoints) <a name="api-endpoints"/>
- [Design Decisions](#design-decisions) <a name="design-decisions"/>
- [Setup (Local & Production)](#setup-(local-&-production)) <a name="setup-(local-&-production)"/>
- [Unit Tests](#unit-tests) <a name="unit-tests"/>
- [Proposed Testing Sequence](#proposed-testing-sequence) <a name="proposed-testing-sequence"/>
- [Git Workflow Practices](#git-workflow-practices) <a name="git-workflow-practices"/>

<br>

## About
Backend application that will be part of a system which teachers can use to perform administrative functions for their students. Teachers and students are identified by their email addresses. Assessment instructions are located [here](https://docs.google.com/document/d/1X0DwX8pSb4XnVwUMRPb9KRMdO61N0ots/edit?usp=sharing&ouid=112691497186686761815&rtpof=true&sd=true).

<br>

## Deployed API
URL: https://golang-api-assessment-hueiiming.onrender.com

<br>

## Project Structure
```bash
Golang-API-Assessment/
│
│── .github/
│   └── workflows
│      └── ci.yml                          # Run unit tests in Github CI
│
├── cmd/                   
│   └── main.go                            # main app to run
│
├── pkg/                       
│   ├── api/                   
│   │   ├── handlers.go                    # Handles api endpoints
│   │   └── handlers_test.go               # Handlers unit test
│   │   └── server.go                      # Start the server
│   │
│   └── repository/                        # Database
│   │   ├── mocks/
│   │   │   └── Repository.go              # mockery file for unit test
│   │   │       
│   │   ├── postgresql_repository.go       # Database methods
│   │   └── postgresql_repository_test.go  # Database test
│   │   └── repository.go                  # Interface for database methods
│   │
│   └── types/                             # Common structs
│   │   ├── requests.go                    
│   │   └── response.go
│   │
│   └── utils/                             # Helper methods    
│       ├── utils.go   
│
├── postman/                   
│   ├── local.postman_collection.json       # postman collection to import for local test
│   └── production.postman_collection.json  # postman collection to import for prod test
│
├── .env                                   # Will be provided by me
├── .gitignore
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

<br>

## API Endpoints
- POST `/api/register`
  - Headers: Content-Type: application/json
  - Success Response status: HTTP 204
  - Request body example:
  ```
  {
    "teacher": "teacherken@gmail.com"
    "students":
      [
        "studentjon@gmail.com",
        "studenthon@gmail.com"
      ]
  }
  ```
- GET `/api/commonstudents`
  - Success response status: HTTP 200
  - Request example 1: GET `/api/commonstudents?teacher=teacherken%40gmail.com`
  - Success body example:
  ```
  {
    "students" :
      [
        "commonstudent1@gmail.com", 
        "commonstudent2@gmail.com",
        "student_only_under_teacher_ken@gmail.com"
      ]
  }
  ```
  - Request example 2: GET `/api/commonstudents?teacher=teacherken%40gmail.com&teacher=teacherjoe%40gmail.com`
  - Success body example:
  ```
  {
    "students" :
      [
        "commonstudent1@gmail.com", 
        "commonstudent2@gmail.com"
      ]
  }
  ```
- POST `/api/suspend`
  - Headers: Content-Type: application/json
  - Success Response status: HTTP 204
  - Request body example:
  ```
  {
    "student": "studentmary@gmail.com"
  }
  
  ```
- POST `/api/retrievefornotifications`
  - Headers: Content-Type: application/json
  - Success response status: HTTP 200
  - Request body example:
  ```
  {
    "teacher":  "teacherken@gmail.com",
    "notification": "Hello students! @studentagnes@gmail.com @studentmiche@gmail.com"
  }
  ```
   - Success response body:
  ```
  {
    "recipients":
      [
        "studentbob@gmail.com",
        "studentagnes@gmail.com", 
        "studentmiche@gmail.com"
      ]   
  }
  ```
- POST `/api/populatestudentsandteachers`
  - Headers: Content-Type: application/json
  - Success response status: HTTP 204
  - Request body example:
  ```
  {
    "teachers":
    [
        "teacherken@gmail.com",
        "teacherjoe@gmail.com",
        "teachermax@gmail.com"
    ],
    "students":
    [
        "studentjon@gmail.com",
        "studenthon@gmail.com",
        "studentmay@gmail.com",
        "studentagnes@gmail.com",
        "studentmiche@gmail.com",
        "studentbob@gmail.com",
        "studentbad@gmail.com",
        "studentmary@gmail.com"
    ]
  }
  ```
- POST `/api/cleardatabase`
  - Headers: Content-Type: application/json
  - Success response status: HTTP 204
- Error handling
  - `"message": "error: status method not allowed"`
  - `"message": "error: invalid teacher email"`
  - `"message": "error: invalid student email"`
  - `"message": "error: missing student request"`
  - `"message": "error: invalid teacher or notification request"`
  - `"message": "error: teacher with email teacherkesn@gmail.com not found"`

<br>

## Design Decisions

- ### Database Design
  - #### Entity Relationship Diagram
    <img width="717" alt="image" src="https://github.com/hueiiming/Golang-API-Assessment/assets/61011188/c700bb27-afda-4e9e-bff3-2f4d8e539272">
  
    - A teacher registers zero or many students
    - A student is registered by zero or many teachers
    - A teacher suspends zero or many students
    - A student is suspended by zero or 1 teacher
    
    ##### Assumptions:
    - A student can only be suspended by 1 teacher as teacher_email is not recorded when suspending a student
    - A teacher can suspend multiple different students

- ### Design Principles
  - #### SOLID
    - Single Responsibility Principle (SRP)
      <br>
      Each structs only has 1 responsibility
      
      ![image](https://github.com/hueiiming/Golang-API-Assessment/assets/61011188/4e08d7e7-e792-477a-ad0c-3ab1c2d03d79)
    
    - Open-Closed Principle (OCP)
      <br>
      I have implemented the `Repository` interface that allows for easy extension by implementing new methods without modifying existing ones.

      ![image](https://github.com/hueiiming/Golang-API-Assessment/assets/61011188/b03b3fff-cdaf-4e11-9854-1f3b254b8fca)

    - Interface Segregation Principle (ISP)
      <br>
      Repository interface follows ISP with methods that are specific to the needs of the consumer. Each method serves a distinct purpose and is small and specific.
    - Dependency Inversion Principle (DIP)
      <br>
      I have implemented Server struct to depend on the Repository interface instead of a specific database implementation. This allows for flexibility in swapping out
different database implementations without affecting the higher-level application logic where they rely on abstractions.

      ![image](https://github.com/hueiiming/Golang-API-Assessment/assets/61011188/32ac2927-14e4-4800-86c0-22243113367e)

  - #### Don't Repeat Yourself (DRY)
    I have implemented my code based on the DRY principle, for example,
    - `MakeHTTPHandle` in `handlers.go` is wrapped around all of the handlers method, ensuring consistent error handling without repeating the error response logic in each handler.
    - `WriteToJSON` in `handlers.go` is used in all of the handlers method, ensuring consistent handling of writing JSON responses to the user.

- ### Design Patterns
  - #### Factory
    I've tried to design my code based on the Factory pattern to allow more flexibility and scalability. For example:
    <br>
    `repository.go` provides a common interface for different databases to implement its methods. In my code, I have used `postgresql_repository.go` to implement these `Repository` methods but if I wish to use other databases, I can simply make use of this interface to implement the common methods and also write new methods specific to the new instance.
    
    ![image](https://github.com/hueiiming/Golang-API-Assessment/assets/61011188/61d3fb6d-f6ea-46f4-9c96-d618e2eef981)
    
    This can be further seen with `main.go` initializing a new `PostgreSQLRepository` and passing it into `NewServer` at `server.go`, where `PostgreSQLRepository` is seen implementing the `Repository` interface.
    <br>
    
    **main.go**
    
    ![image](https://github.com/hueiiming/Golang-API-Assessment/assets/61011188/9dff6f15-e5cf-4bcb-9210-256a5cfe683b)
    
    **server.go**
    
    ![image](https://github.com/hueiiming/Golang-API-Assessment/assets/61011188/aec5e8ab-acba-4083-8bf1-c01dfbff501e)

    
    Therefore, this fulfills the Factory pattern that defines providing an interface for creating objects in a superclass, but allows subclasses to alter the type of objects that will be created.
  
<br>

## Setup (Local & Production)

### Prerequisites
- Install Homebrew for macOS or Linux 
```
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
```
- Install go 
```brew install go``` or visit https://go.dev/doc/install
- Install mockery 
```brew install mockery``` or visit https://go.dev/doc/install
- Install DBeaver https://dbeaver.io/download/ or any other database that supports PostgreSQL
- Include `.env` file in the root folder **Note**: DB_PASSWORD is for production connection to PostgreSQL, the password will be provided by me
  - `ENV=local` for running application in local machine
  - `ENV=prod` for running application in production
  ```
  ENV=local
  PORT=8080
  DB_PASSWORD=<PASSWORD>
  ```

### Run Locally
  #### Step 1: Clone repo
  ```
  git clone https://github.com/hueiiming/Golang-API-Assessment.git
  ```
  
  #### Step 2: Run go project
  ```
  make run
  ```
  
  #### Step 3: Connect to PostgreSQL and view data
  - Open up DBeaver
  - Click file -> new -> expand DBeaver -> click Database Connection
  - Select PostgreSQL and click next
  - Connection settings
    - **Host**: localhost
    - **Database**: postgres
    - **Port**: 5432
    - **username**: postgres
    - **password**: root
  - Click on Finish
  - You have now successfully connected to postgres
    
  <img width="674" alt="image" src="https://github.com/hueiiming/Golang-API-Assessment/assets/61011188/caa935f8-6493-4ab5-a53d-aa17c7d574a4">


### View production Database
  #### Connect to PostgreSQL and view data
  - Open up DBeaver
  - Click file -> new -> expand DBeaver -> click Database Connection
  - Select PostgreSQL and click next
  - Connection settings
    - **Host**: aws-0-ap-southeast-1.pooler.supabase.com
    - **Database**: postgres
    - **Port**: 5432
    - **username**: postgres.wxmkhkkcxatyzukbfqtw
    - **password**: Will be provided by me
  - Click on Finish
  - You have now successfully connected to postgres in production environment
    
  <img width="519" alt="image" src="https://github.com/hueiiming/Golang-API-Assessment/assets/61011188/e57e1ae1-9d94-4e93-913d-f3bb83e50e1e">

<br><br>

## Unit Tests
Unit tests are being executed on every Pull Request or Push to main branch using Github Actions
  - handlers_test.go (Follows a table-driven test to be more descriptive and easy to understand. mockery was used to mock database interface)
  - postgresql_repository_test.go (Uses sqlmock library `"github.com/DATA-DOG/go-sqlmock"`, without needing a real database connection)
  ### Run tests locally:
  ```
  make test
  ```

<br><br>

## Proposed Testing Sequence
Recommended Postman collections are included in `postman` folder to import
1. Clear all database tables to start a new testing scenario using POST `/api/cleardatabase`
2. Populate students and teachers database tables using POST `/api/populatestudentsandteachers`
3. Feel free to test any of the other 4 main endpoints (Note: Registration table and Suspension table are currently empty at this step)

<br><br>

## Git Workflow Practices
Throughout this project, I have been adhering to the git workflow best practices by:
- Branching out for every code change such as new features/bug fixes
- Creating Pull Requests (PR) before merging into the main branch
- Squash and Merge PRs to keep the main branch commits clean
- Keeping the main branch stable at all times
- Committing frequently with descriptive messages
- Integrated CI testing using GitHub Actions to ensure every PR passes before merging into the main branch

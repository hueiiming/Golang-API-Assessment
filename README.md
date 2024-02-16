# Golang-API-Assessment

## Table of Contents  
- [About](#about) <a name="about"/>  
- [Deployed API](#deployed-api) <a name="deployed-api"/>  
- [API Endpoints](#api-endpoints) <a name="api-endpoints"/>  
- [Run Locally](#run-locally) <a name="run-locally"/>   
- [User Stories](#user-stories) <a name="user-stories"/>

<br>

## About
Backend application that will be part of a system which teachers can use to perform administrative functions for their students. Teachers and students are identified by their email addresses.

<br>

## Deployed API
URL: https://golang-api-assessment-hueiiming.onrender.com

<br>

## API Endpoints
- `/api/register`
- `/api/commonstudents?teacher=teacherken%40gmail.com`
- `/api/commonstudents?teacher=teacherken%40gmail.com&teacher=teacherjoe%40gmail.com`
- `/api/suspend`
- `/api/retrievefornotifications`

<br>

## Run Locally

### Prerequisites
- Install Homebrew for macOS or Linux 
```
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
```
- Install go 
```brew install go``` or visit https://go.dev/doc/install
- Install DBeaver https://dbeaver.io/download/ or any other database that supports PostgreSQL
- Include `.env` file in the root folder
```
ENV=local
PORT=8080
DB_PASSWORD=<PASSWORD>
```
<small>Note: DB_PASSWORD is for production connection to PostgreSQL, the password will be provided by me</small>

### Step 1: Clone repo
```
git clone https://github.com/hueiiming/Golang-API-Assessment.git
```

### Step 2: Run go project
```
make run
```

### Step 3: Connect to PostgreSQL and view data
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

<br><br>

## User Stories
1. As a teacher, I want to register one or more students to a specified teacher.
A teacher can register multiple students. A student can also be registered to multiple teachers.
- Endpoint: POST /api/register
- Headers: Content-Type: application/json
- Success response status: HTTP 204
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

2. As a teacher, I want to retrieve a list of students common to a given list of teachers (i.e. retrieve students who are registered to ALL of the given teachers).
-	Endpoint: GET /api/commonstudents
-	Success response status: HTTP 200
-	Request example 1: GET /api/commonstudents?teacher=teacherken%40gmail.com
-	Success response body 1:
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
-	Request example 2: GET /api/commonstudents?teacher=teacherken%40gmail.com&teacher=teacherjoe%40gmail.com
-	Success response body 2:
```
{
  "students" :
    [
      "commonstudent1@gmail.com", 
      "commonstudent2@gmail.com"
    ]
}
```

3. As a teacher, I want to suspend a specified student.
-	Endpoint: POST /api/suspend
-	Headers: Content-Type: application/json
-	Success response status: HTTP 204
-	Request body example:
```
{
  "student" : "studentmary@gmail.com"
}
```

4. As a teacher, I want to retrieve a list of students who can receive a given notification.
- A notification consists of:
  -	the teacher who is sending the notification, and
  -	the text of the notification itself.
- To receive notifications from e.g. 'teacherken@gmail.com', a student:
  -	MUST NOT be suspended,
  -	AND MUST fulfill AT LEAST ONE of the following:
    - is registered with “teacherken@gmail.com"
    -	has been @mentioned in the notification
- The list of students retrieved should not contain any duplicates/repetitions.
  -	Endpoint: POST /api/retrievefornotifications
  -	Headers: Content-Type: application/json
  -	Success response status: HTTP 200
  -	Request body example 1:
```
{
  "teacher":  "teacherken@gmail.com",
  "notification": "Hello students! @studentagnes@gmail.com @studentmiche@gmail.com"
}
```
  -	Success response body 1:
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
In the example above, studentagnes@gmail.com and studentmiche@gmail.com can receive the notification from teacherken@gmail.com, regardless whether they are registered to him, because they are @mentioned in the notification text. studentbob@gmail.com however, has to be registered to teacherken@gmail.com.
  -	Request body example 2:
```
{
  "teacher":  "teacherken@gmail.com",
  "notification": "Hey everybody"
}
```
  -	Success response body 2:
```
{
  "recipients":
    [
      "studentbob@gmail.com"
    ]   
}
```

### Error Responses
For all the above API endpoints, error responses should:
-	have an appropriate HTTP response code
-	have a JSON response body containing a meaningful error message:
```
{ "message": "Some meaningful error message" }
```

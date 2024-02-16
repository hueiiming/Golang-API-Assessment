<h1>Golang-API-Assessment</h1>

Backend application that will be part of a system which teachers can use to perform administrative functions for their students. Teachers and students are identified by their email addresses.

<h2>Deployed API</h2>

URL: https://golang-api-assessment-hueiiming.onrender.com
<h2>API Endpoints</h2>

- `/api/register`

- `/api/commonstudents?teacher=teacherken%40gmail.com`

- `/api/commonstudents?teacher=teacherken%40gmail.com&teacher=teacherjoe%40gmail.com`

- `/api/suspend`

- `/api/retrievefornotifications`


<h2>Run Locally</h2>

<h3>Prerequisites</h3>
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

<h3>Step 1: Clone repo</h3>

```
git clone https://github.com/hueiiming/Golang-API-Assessment.git
```

<h3>Step 2: Run go project</h3>

```
make run
```

<h3>Step 3: Connect to PostgreSQL and view data</h3>
- Open up DBeaver

- Click file -> new -> expand DBeaver -> click Database Connection

- Select PostgreSQL and click next

- Connection settings
  - **Host**: localhost

  - **Database**: postgres
  
  - **username**: postgres
  
  - **password**: root
  
- 
- click on Finish

- You have now successfully connected to postgres
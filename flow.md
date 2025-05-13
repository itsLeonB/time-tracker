1. Login
2. Find a task to work on
3. Insert UserTask
4. When log, insert UserTaskLog


GET /projects
Get list of projects

GET /projects/:id
Get project detail
Queries: 
include=tasks,logs

GET /user-tasks
Get list of user tasks

POST /user-tasks/:id/logs

Setup service client auth

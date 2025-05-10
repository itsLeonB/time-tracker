1. Login
2. Find a task to work on
3. Insert UserTask
4. When log, insert UserTaskLog

store external_id in projects
store external_id in tasks

GET /projects
Get list of projects

GET /projects/:id
Get project detail
Queries: 
include=tasks,logs

GET /user-tasks
Get list of user tasks

GET /tasks?number=US-1
Find task with externalId = number
if not exists
Find Youtrack Task with idReadable = number
Upsert Project
Upsert Task

POST /user-tasks/:id/logs
# Forum

This project had already been started before but I decided to start from scratch.

Previous versions did not meet my expectations, hence this new approach.

To use this code, follow these steps:

## 1. Download SQlite

##  Click [Here](https://www.sqlite.org/download.html) for download SQlite

For the moment I choose SQlite but in the future I change to a another database (postgreSQL, MySQL)

## 2 You also need Docker to build the project in the dockerFile

In this project you can post everithing you can (messages, likes dislikes not images for the moment but ça viendra avec le futur),

I use several import like gorrilla session, crypto wbesocket securecookie.

Gorrilla session can make a session for the user and it's facilite la tache pour cela, crypto hash the password in the database and is more protection if the database is contourned, websocket is for speaking between users, securecookie is for a token if the user forget his password.


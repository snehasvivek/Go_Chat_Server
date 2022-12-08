# Project Title: Chat Server in Go  

# Project Description
Written in Go, this project creates a local server that allows users to connect to a server and create rooms. Users are able to create rooms and send messages to the 
rest of the users in the room. 

# How to Install and Run the Project 
Go is required to run this project. First, install Go if it is not already installed. 

To set up the server:
1. Open Command Prompt/Terminal.
2. Change Directory to location of the Files 
3. Enter command: go build
4. Enter command: chat_server.exe

To connect to the server as a user:
1. Open Command Prompt.
2. Enter command: telnet localhost 8888

# Command Descriptions: How to use the server 
comand        description
-------       --------------
*name {NAME}  sets the name of the user 
*join {ROOM}  if ROOM exists, puts user in ROOM; if ROOM does not exist, creates user and puts user in ROOM 
*rooms        lists all current rooms 
*msg {MSG}    broadcasts MSG to all members in the room 

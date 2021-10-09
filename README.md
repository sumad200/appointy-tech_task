# appointy-tech_task

**This repository was created as a submission to Appointy October 2021 Internship recruitment Tech Task round**
   
 Basic HTTP REST API written in Golang, Responses are in JSON, uses MongoDB for storage 
 
 ### Data Attributes: ###
  -**User**
   * ID
   * Name
   * Email
   * Password
   
  -**Posts**
   * ID
   * Caption
   * Image URL
   * Posted Timestamp 
   
 ### Endpoints Implemented: ###
 
 -**Create an User**
   * POST request
   * Use JSON request body
   * URL : ‘/users'
   
 -**Get an user using id**
   * GET request
   * URL : ‘/users/{id}'
   
 -**Create a Post**
   * POST request
   * Use JSON request body
   * URL : ‘/posts'  
   
 -**Get a post using id**
   * GET request
   * URL : ‘/posts/{id}'
   
  -**Get all posts in the database**
   * GET request
   * URL : ‘/allposts'
   
  -**Get all users in the database**
   * GET request
   * URL : ‘/allusers' 


 
 ### Additional Features: ###
  * Passwords securely stored with SHA256 such they can't be reverse engineered

 ### Installation  ###
  * Clone this repo or download code
  * Place all files except readme in $GOPATH/src
  * Ensure that go.mongodb.org/mongo-driver is present
  * go run main.go

## Installation:

Secret_Santa requires go 1.8 or later

`go get -u github.com/sachinagada/Secret_Santa`

## Usage

After cloning the project, go to the resources directory and change the username and password for the email address from which the emails should be sent for letting people know the name of the person they are assigned.

 Go the cmd directory and run the main file to start the application
  ```
  cd Path/to/Secret_Santa/cmd
  go run main.go
  ```
 Go to the following [link](http://localhost:8090/secretSanta) to bring up the
 front end of the application. In the forms, insert the name and the
 corresponding email address and hit the submit button when done. The application will randomly assign the Secret Santa to each name and email everyone their paired match. It will make sure that no one gets themselves as their Secret Santa.

## Purpose
This application makes it easier for friends and family to continue the tradition of Secret Santa near the holiday season despite living far apart. By just typing in the names of the people who want to participate in the tradition, the assigning of people is done within seconds and prevents the coordination required to match a large group without having some get themselves as their own Secret Santa.

### Personal Purpose

I am using this application to try out angular2 front end with a go application. I also wanted to learn docker and how to containerize an app so using this application as a means of learning this technology.

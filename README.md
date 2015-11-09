#Issue-Tracker

Full tracking system written in Go/HTML


###Introduction
I have been working with [Jira](https://www.atlassian.com/software/jira) for years. In spite how much I like Jira, I find it a bit complex to both ease of use and ability to integrate with other programs.

Lately I've had to develop separately different versions of the issue-tracker for another projects, and thought it would make sense to have one unified open-source version of it. And this is the reason I decided to create this project.

Issue-Tracker is written in Go and Javascript. The backend is a golang API and the storage is into a Postgres database. In the future I would like to make an implementation also in MongoDB, but I'm not sure I have time available. Postgres is now my bet and gives me the assurance that everything is working within transactions with referential integrity. Even attachments are stored in Postgres, so backup/restore operations are simple to achieve.

I hope you find useful this project, either you do what I do: integrate with other applications; or simply out of curiosity if you are someone who is learning to develop golang.

The frontend is done in HTML / JS with angularjs. I'm not a FrontEnd guy, so possibly the appearance is awful for an expert. Please have mercy.

I wish more people to contribute to the project to make it more robust, incorporate additional features and watch it grow. You can work as possible, I'm open to suggestions from detecting a bug or expose an improvement. Any feedback is welcome.

Below are some sample screens and the instructions to configure the application at the bottom.

####Dashboard
![demo1](https://cloud.githubusercontent.com/assets/1648558/9989053/161dbc80-601b-11e5-81a2-b7e3dd063932.png)

####Projects
![demo2](https://cloud.githubusercontent.com/assets/1648558/9989111/83dc6956-601b-11e5-9ffe-c10fd0de6748.png)

####Issues
![demo3](https://cloud.githubusercontent.com/assets/1648558/9989133/a4314b7c-601b-11e5-80f7-7c84d37794c4.png)

### Installation
*NOTE: This has been tested to work in Ubuntu 14.04.3 but should be easy to port these instructions to MacOS or Windows*
```
go get github.com/mauleyzaola/issue-tracker
```
##### Requirements
Make sure to prepare the setup for the UI before continuing. The instructions are [here](/static)

You will need [Postgres](http://www.postgresql.org/download/) installed and change ```postgres``` user's password to: ```nevermind```. So it can match the sample configuration files.

##### Setup
Run these commands to copy the sample configuration files to real ones:
```
cd $GOPATH/src/github.com/mauleyzaola/issue-tracker
cp dbmigrations/pg/db/dbconf.yml.sample dbmigrations/pg/db/dbconf.yml
cp server/config.json.sample server/config.json
cp test/config.json.sample test/config.json
```
Next, create two databases in postgres; one is for running the unit tests, the other for the application.
```
sudo -u postgres psql -c 'create database tracker;'
sudo -u postgres psql -c 'create database tracker_test;'
```
When databases have been created, execute the [check.sh](check.sh) script. It will download the go dependencies, create the database objects and run the unit tests.
```
./check.sh
```
Change to ```server/``` directory, build the application and execute
```
cd server/
go build
./server
```
When you're done with [UI setup](/static/), browse to your [http://localhost/](http://localhost/) and the application should be ready. The credentials to access are:

email: ```admin@admin.com```
password: ```admin```

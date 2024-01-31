# Quiz WebApp
## Brief Desctription
The project is a web application developed in Go (Golang) that provides a Quiz Game
with user authentication, registration, and leaderboard
management. It allows users to register with unique usernames and
securely log in using bcrypt-hashed passwords. Once authenticated,
users can participate in a quiz, where questions are presented one by
one from a database, and their progress is tracked. Upon completing the
quiz, users' scores are displayed on a leaderboard. Real-time communication
is developed using WebSocket support for features such as adding users
to a game queue and updating the leaderboard queue. The project
aims to offer an interactive and engaging web experience while demonstrating
the capabilities of Go and the Gorilla toolkit.

## Installation (v1)
To try the Quiz game all you need to do is:
1. Dowload the `/Templates`, `/Data` and `QuizWebApp` files (or `QuizWebAppWin.exe` if you are on Windows).
2. go in the project folder and execute the `QuizWebApp.exe` if you are on windows,
   else you can execute the unix file `QuizWebApp` if you are on macOS
4. go to your browser and search for `localhost:8080` in the search bar
5. Use the app with your friends and multiple people by connecting all to the same ip

## Installation (v2)
If the previous installation doesn't work it's because sometimes running a file by
double-clicking it may change the absolute path of the execution.
1. Dowload the `/Templates`, `/Data` and `QuizWebApp` files (or `QuizWebAppWin.exe` if you are on Windows).
2. open the terminal and navigate to the folder where you cloned the project using
   `cd path-to-folder`
3. write `./QuizWebApp` in the terminal and press enter (or`./QuizWebApp.exe` if you are on windows
4. the Quiz should now start
5. go to your browser and search for `localhost:8080` in the search bar
6. Use the app with your friends and multiple people by connecting all to the same ip

## Installation (v3)
If the previous 2 installation didn't work you will need to build or run the code in the repository.
You'll need Golang to successfully run the program so first install Golang, then:
1. Clone the repository to your local machine using
`git clone https://github.com/pippodima/DP_Project`
2. open the project with your IDE or go with the terminal to the folder where you dowloaded the repo
3. build or run the program with `go run .` or `go build .` command in the terminal
4. go to your browser and search for `localhost:8080` in the search bar
5. Use the app with your friends and multiple people by connecting all to the same ip

## Flags
There are 4 flags you can use by running the App from your terminal:
1. `./QuizWebApp -help` shows the possible flags and their description
2. `./QuizWebApp -p N` set the number of players of the Quiz to N (default is 3)
3. `./QuizWebApp -q N` set the number of Question per round (default 5)
4. `./QuizWebApp -verbose` it will print in the console all errors and information messages

## Note
If the database is not working just navigate to the folder of the project then click on
the folder `/Data`, delete the database.db file and run CreateDB.exe, this script will create a
new DataBase with no users and the 50 standard question


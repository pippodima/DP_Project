package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

type Question struct {
	ID         int    `json:"id"`
	Text       string `json:"text"`
	Options    string `json:"options"`
	CorrectIdx int    `json:"correct_idx"`
}

var db *sql.DB

func addQuestion(question Question) error {
	_, err := db.Exec(`
		INSERT INTO questions (text, options, correct_idx)
		VALUES (?, ?, ?)
	`, question.Text, question.Options, question.CorrectIdx)
	return err
}

func main() {
	var err error
	db, err = sql.Open("sqlite3", "database.db")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT UNIQUE NOT NULL,
			password TEXT NOT NULL,
		    totalPoints INTEGER,  
		    gamesPlayed INTEGER               
		)
	`)
	if err != nil {
		fmt.Println("error creating users table: ", err)
	}
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS questions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			text TEXT,
			options TEXT,
			correct_idx INTEGER
		)
	`)
	if err != nil {
		fmt.Println("error creating questions table: ", err)
	}

	defer db.Close()

	questionsToAdd := []Question{
		{
			Text:       "What is the capital of France?",
			Options:    "Berlin,Paris,Madrid,Rome",
			CorrectIdx: 1,
		},
		{
			Text:       "Which planet is known as the Red Planet?",
			Options:    "Mars,Venus,Jupiter,Saturn",
			CorrectIdx: 0,
		},
		{
			Text:       "In which year did World War II end?",
			Options:    "1943,1945,1947,1950",
			CorrectIdx: 1,
		},
		{
			Text:       "Who wrote 'Romeo and Juliet'?",
			Options:    "Charles Dickens,William Shakespeare,Emily Brontë,Jane Austen",
			CorrectIdx: 1,
		},
		{
			Text:       "What is the largest mammal on Earth?",
			Options:    "Elephant,Blue Whale,Giraffe,Gorilla",
			CorrectIdx: 1,
		},
		{
			Text:       "Which chemical element has the symbol 'H'?",
			Options:    "Helium,Hydrogen,Mercury,Neon",
			CorrectIdx: 1,
		},
		{
			Text:       "What is the main ingredient in guacamole?",
			Options:    "Tomato,Avocado,Onion,Cilantro",
			CorrectIdx: 1,
		},
		{
			Text:       "Who painted the Mona Lisa?",
			Options:    "Vincent van Gogh,Leonardo da Vinci,Claude Monet,Pablo Picasso",
			CorrectIdx: 1,
		},
		{
			Text:       "Which country is known as the Land of the Rising Sun?",
			Options:    "China,Japan,Korea,Vietnam",
			CorrectIdx: 1,
		},
		{
			Text:       "What is the speed of light in a vacuum?",
			Options:    "299.792 km/s,150.000 km/s,450.000 km/s,600.000 km/s",
			CorrectIdx: 0,
		},
		{
			Text:       "Who wrote 'To Kill a Mockingbird'?",
			Options:    "Harper Lee,J.K. Rowling,George Orwell,Ernest Hemingway",
			CorrectIdx: 0,
		},
		{
			Text:       "What is the capital of Brazil?",
			Options:    "Buenos Aires,Bogotá,Lima,Brasília",
			CorrectIdx: 3,
		},
		{
			Text:       "Which gas makes up the majority of Earth's atmosphere?",
			Options:    "Oxygen,Carbon Dioxide,Nitrogen,Helium",
			CorrectIdx: 2,
		},
		{
			Text:       "What is the currency of Japan?",
			Options:    "Won,Yen,Ringgit,Peso",
			CorrectIdx: 1,
		},
		{
			Text:       "Who is known as the 'Father of Computer Science'?",
			Options:    "Alan Turing,Charles Babbage,Ada Lovelace,Bill Gates",
			CorrectIdx: 0,
		},
		{
			Text:       "What is the longest river in the world?",
			Options:    "Amazon,Nile,Yangtze,Mississippi",
			CorrectIdx: 0,
		},
		{
			Text:       "Which element has the chemical symbol 'Au'?",
			Options:    "Silver,Gold,Platinum,Iron",
			CorrectIdx: 1,
		},
		{
			Text:       "Who discovered penicillin?",
			Options:    "Marie Curie,Alexander Fleming,Louis Pasteur,Joseph Lister",
			CorrectIdx: 1,
		},
		{
			Text:       "What is the largest planet in our solar system?",
			Options:    "Earth,Jupiter,Mars,Saturn",
			CorrectIdx: 1,
		},
		{
			Text:       "Which ocean is the largest?",
			Options:    "Indian,Atlantic,Arctic,Pacific",
			CorrectIdx: 3,
		},
		{
			Text:       "Who played Neo in 'The Matrix'?",
			Options:    "Keanu Reeves,Brad Pitt,Leonardo DiCaprio,Johnny Depp",
			CorrectIdx: 0,
		},
		{
			Text:       "What is the square root of 144?",
			Options:    "10,12,14,16",
			CorrectIdx: 1,
		},
		{
			Text:       "What is the chemical symbol for water?",
			Options:    "H2O,CO2,O2,N2",
			CorrectIdx: 0,
		},
		{
			Text:       "Who wrote the 'Harry Potter' series?",
			Options:    "J.K. Rowling,George R.R. Martin,Stephen King,J.R.R. Tolkien",
			CorrectIdx: 0,
		},
		{
			Text:       "What is the capital of Australia?",
			Options:    "Sydney,Melbourne,Canberra,Brisbane",
			CorrectIdx: 2,
		},
		{
			Text:       "Which famous scientist developed the theory of general relativity?",
			Options:    "Isaac Newton,Albert Einstein,Galileo Galilei,Nikola Tesla",
			CorrectIdx: 1,
		},
		{
			Text:       "What is the main component of Earth's core?",
			Options:    "Iron,Nickel,Lead,Copper",
			CorrectIdx: 0,
		},
		{
			Text:       "Who wrote '1984'?",
			Options:    "George Orwell,Aldous Huxley,Ray Bradbury,Arthur C. Clarke",
			CorrectIdx: 0,
		},
		{
			Text:       "What is the largest desert in the world?",
			Options:    "Gobi Desert,Sahara Desert,Arabian Desert,Karakum Desert",
			CorrectIdx: 1,
		},
		{
			Text:       "In what year did the first manned moon landing occur?",
			Options:    "1965,1969,1973,1981",
			CorrectIdx: 1,
		},
		{
			Text:       "Which animal is known as the 'King of the Jungle'?",
			Options:    "Giraffe,Elephant,Lion,Tiger",
			CorrectIdx: 2,
		},
		{
			Text:       "Who painted 'The Starry Night'?",
			Options:    "Claude Monet,Pablo Picasso,Vincent van Gogh,Salvador Dalí",
			CorrectIdx: 2,
		},
		{
			Text:       "What is the currency of China?",
			Options:    "Won,Yuan,Ruble,Peso",
			CorrectIdx: 1,
		},
		{
			Text:       "Who is the author of 'The Great Gatsby'?",
			Options:    "F. Scott Fitzgerald,Ernest Hemingway,J.D. Salinger,Mark Twain",
			CorrectIdx: 0,
		},
		{
			Text:       "What is the largest island in the world?",
			Options:    "Greenland,Australia,Java,Madagascar",
			CorrectIdx: 0,
		},
		{
			Text:       "Who is known as the 'Queen of Pop'?",
			Options:    "Madonna,Beyoncé,Lady Gaga,Taylor Swift",
			CorrectIdx: 0,
		},
		{
			Text:       "Which planet is known as the 'Morning Star' or 'Evening Star'?",
			Options:    "Mercury,Venus,Mars,Jupiter",
			CorrectIdx: 1,
		},
		{
			Text:       "What is the smallest prime number?",
			Options:    "1,2,3,5",
			CorrectIdx: 1,
		},
		{
			Text:       "Who is the famous physicist known for his theory of relativity?",
			Options:    "Niels Bohr,Werner Heisenberg,Max Planck,Albert Einstein",
			CorrectIdx: 3,
		},
		{
			Text:       "What is the capital of South Africa?",
			Options:    "Cape Town,Pretoria,Johannesburg,Durban",
			CorrectIdx: 1,
		},
		{
			Text:       "Which famous scientist formulated the laws of motion and universal gravitation?",
			Options:    "Galileo Galilei,Isaac Newton,Nikola Tesla,Marie Curie",
			CorrectIdx: 1,
		},
		{
			Text:       "What is the chemical symbol for gold?",
			Options:    "Ag,Au,Hg,Pt",
			CorrectIdx: 1,
		},
		{
			Text:       "Who is the author of 'Pride and Prejudice'?",
			Options:    "Charlotte Brontë,Emily Brontë,Jane Austen,Charles Dickens",
			CorrectIdx: 2,
		},
		{
			Text:       "Which planet is known as the 'Blue Planet'?",
			Options:    "Earth,Mars,Venus,Uranus",
			CorrectIdx: 0,
		},
		{
			Text:       "What is the largest moon in our solar system?",
			Options:    "Io,Enceladus,Ganymede,Titan",
			CorrectIdx: 2,
		},
		{
			Text:       "In what year did Christopher Columbus reach the Americas?",
			Options:    "1492,1501,1510,1525",
			CorrectIdx: 0,
		},
		{
			Text:       "Who is the famous Greek philosopher known for his teachings on ethics and virtue?",
			Options:    "Socrates,Aristotle,Plato,Pythagoras",
			CorrectIdx: 1,
		},
		{
			Text:       "What is the largest bird in the world?",
			Options:    "Ostrich,Emu,Condor,Albatross",
			CorrectIdx: 0,
		},
		{
			Text:       "Which element is essential for human bones and teeth?",
			Options:    "Calcium,Iron,Potassium,Sodium",
			CorrectIdx: 0,
		},
		{
			Text:       "What is the currency of Russia?",
			Options:    "Ruble,Yen,Euro,Pound",
			CorrectIdx: 0,
		},
		// Add more questions as needed
	}

	for _, q := range questionsToAdd {
		err := addQuestion(q)
		if err != nil {
			log.Fatal(err)
		}
	}

}

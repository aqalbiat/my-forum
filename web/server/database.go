package server

import (
	"database/sql"
	"errors"
	"log"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const INITIATE_DB = `CREATE TABLE IF NOT EXISTS users (
				userID INTEGER PRIMARY KEY,
				userName TEXT,
				userEmail TEXT,
				passwordHash TEXT);`

const TEST_USERS_IN_DB = `1|Doda      |loperamid@simes.ok       |ff45jkl356
						  2|Clipper   |chromsfield@sin.ok       |efd893oww
						  3|Flokster  |christopher@slopesh.poz  |sidef56789asd`

const INITIATE_SESSIONS = `CREATE TABLE IF NOT EXISTS sessions (
						id TEXT,
						userName TEXT,
						expires INT);`

const INITIATE_POSTS = `CREATE TABLE IF NOT EXISTS posts (
						postID INTEGER PRIMARY KEY,
						userName TEXT,
						title TEXT,
						content TEXT,
						parentID INTEGER,
						categories TEXT,
						likes INTEGER,
						dislikes INTEGER);`

const INITIATE_LIKES_AND_DISLIKES = `CREATE TABLE IF NOT EXISTS likesndislikes (
				postID INT,				
				userName TEXT,
				typeText TEXT);`

const INITIATE_COMMENTS = `CREATE TABLE IF NOT EXISTS comments (
				commentID INTEGER PRIMARY KEY,
				content TEXT,
				postID INT,
				author TEXT,
				likes INT,
				dislikes INT);`

const INITIATE_COMMENTS_FEEDBACKS = `CREATE TABLE IF NOT EXISTS feedbacks (
				commentID INT,
				author TEXT,
				typeText TEXT);`

var (
	db       *sql.DB
	posts    *sql.DB
	comms    *sql.DB
	initErr  error
	postsErr error
	commsErr error
)

func InitiateDatabase() {
	db, initErr = sql.Open("sqlite3", "./web/database/userlist.db")
	if initErr != nil {
		log.Println(initErr)
		return
	}
	if _, err := db.Exec("PRAGMA journal_mod=WAL"); err != nil {
		log.Println(err)
		return
	}
	if _, err := db.Exec(INITIATE_DB); err != nil {
		log.Println(err)
		return
	}
	if _, err := db.Exec(INITIATE_SESSIONS); err != nil {
		log.Println(err)
		return
	}
	posts, postsErr = sql.Open("sqlite3", "./web/database/posts.db")
	if postsErr != nil {
		log.Println(postsErr)
		return
	}
	if _, err := posts.Exec("PRAGMA journal_mod=WAL"); err != nil {
		log.Println(err)
		return
	}
	if _, err := posts.Exec(INITIATE_POSTS); err != nil {
		log.Println(err)
		return
	}
	if _, err := posts.Exec(INITIATE_LIKES_AND_DISLIKES); err != nil {
		log.Println(err)
		return
	}
	comms, commsErr = sql.Open("sqlite3", "./web/database/comments.db")
	if commsErr != nil {
		log.Println(commsErr)
		return
	}
	if _, err := comms.Exec("PRAGMA journal_mod=WAL"); err != nil {
		log.Println(err)
		return
	}
	if _, err := comms.Exec(INITIATE_COMMENTS); err != nil {
		log.Println(err)
		return
	}
	if _, err := comms.Exec(INITIATE_COMMENTS_FEEDBACKS); err != nil {
		log.Println(err)
		return
	}
	log.Println("database initiation succesfully ended!")
}

func AddNewUser(NewUser UserCredentials) (bool, error) {
	db.Exec("BEGIN;")
	defer db.Exec("COMMIT;")
	log.Println("Adding new user!")
	log.Println(NewUser)
	_, err := db.Exec("INSERT INTO users VALUES (NULL,?,?,?);", NewUser.NameOfUser, NewUser.EmailOfUser, NewUser.PassHash)
	if err != nil {
		log.Println(err)
		return false, err
	}
	return true, nil
}

func CheckForFreeName(username string) (bool, error) {
	db.Exec("BEGIN;")
	defer db.Exec("COMMIT;")
	rows, err := db.Query("SELECT userID FROM users WHERE userName=?;", username)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	count := 0
	for rows.Next() {
		var curr int = -1
		scanErr := rows.Scan(&curr)
		if scanErr != nil {
			return false, scanErr
		}
		count++
	}
	return (count < 1), nil
}

func CheckForFreeEmail(email string) (bool, error) {
	db.Exec("BEGIN;")
	defer db.Exec("COMMIT;")
	rows, err := db.Query("SELECT userID FROM users WHERE userEmail=?;", email)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	count := 0
	for rows.Next() {
		var curr int = -1
		scanErr := rows.Scan(&curr)
		if scanErr != nil {
			return false, scanErr
		}
		count++
	}
	return (count < 1), nil
}

func CheckForGoodEntry(user UserCredentials) (bool, error) {
	db.Exec("BEGIN;")
	defer db.Exec("COMMIT;")
	log.Println("Try to log-in!")
	log.Println(user)
	rows, err := db.Query("SELECT passwordHash FROM users WHERE userEmail=?;", user.EmailOfUser)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	count := 0
	ok := false
	for rows.Next() {
		var pass string
		scanErr := rows.Scan(&pass)
		if scanErr != nil {
			return false, scanErr
		}
		count++
		if pass == user.PassHash {
			ok = true
		}
	}
	return (ok && count == 1), nil
}

func GetUsernameFromEmail(emailOfUser string) (string, error) {
	db.Exec("BEGIN;")
	defer db.Exec("COMMIT;")
	rows, err := db.Query("SELECT userName from users WHERE userEmail=?;", emailOfUser)
	if err != nil {
		return ".", err
	}
	defer rows.Close()
	counts := 0
	lastOccurence := "#"
	for rows.Next() {
		var currentName string
		scanErr := rows.Scan(&currentName)
		if scanErr != nil {
			return ".", err
		}
		lastOccurence = currentName
		counts++
	}
	if counts != 1 {
		return "", errors.New("database query: not an uniqe email or it does not exist")
	}
	return lastOccurence, nil
}

func GetUsernameBySessionID(sid string) (string, error) {
	db.Exec("BEGIN;")
	defer db.Exec("COMMIT;")
	rows, err := db.Query("SELECT userName FROM sessions WHERE id=?;", sid)
	if err != nil {
		return "#", err
	}
	defer rows.Close()
	counts := 0
	lastOccurence := "r"
	for rows.Next() {
		var buffer string
		scanErr := rows.Scan(&buffer)
		if scanErr != nil {
			return "#", scanErr
		}
		counts++
		lastOccurence = buffer
	}
	if counts != 1 {
		return "#", errors.New("database query: not an uniqe session ID or it does not exist")
	}
	return lastOccurence, nil
}

func CheckExistanceOfSessionID(sid string) (bool, error) {
	db.Exec("BEGIN;")
	defer db.Exec("COMMIT;")
	currentTime := time.Now().Unix()
	rows, err := db.Query("SELECT expires FROM sessions WHERE id=?;", sid)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	counts := 0
	for rows.Next() {
		var expiringTime int64
		scanErr := rows.Scan(&expiringTime)
		if scanErr != nil {
			return false, scanErr
		}
		if expiringTime > currentTime {
			counts++
		}
	}
	if counts > 1 {
		return false, errors.New("database query: not an uniqe session ID")
	}
	return (counts == 1), nil
}

func AddNewSessionID(sessionID string, userName string) (bool, error) {
	db.Exec("BEGIN;")
	defer db.Exec("COMMIT;")
	expireTime := time.Now().Unix() + 2400
	_, insertionErr := db.Exec("INSERT INTO sessions VALUES (?, ?, ?);", sessionID, userName, expireTime)
	if insertionErr != nil {
		return false, insertionErr
	}
	return true, nil
}

func DeleteSessionWithID(sessionID string) (bool, error) {
	db.Exec("BEGIN;")
	defer db.Exec("COMMIT;")
	_, updateErr := db.Exec("UPDATE sessions SET expires=1 WHERE id=?", sessionID)
	if updateErr != nil {
		return false, updateErr
	}
	return true, nil
}

func AddNewPost(username string, title string, content string, parentID int, categories string) error {
	posts.Exec("BEGIN;")
	defer posts.Exec("COMMIT;")
	_, insertionErr := posts.Exec("INSERT INTO posts VALUES (NULL, ?, ?, ?, ?, ?, 0, 0);", username, title, content, parentID, categories)
	if insertionErr != nil {
		log.Println("[AddNewPost: insertionErr] -> " + insertionErr.Error())
		return insertionErr
	}
	return nil
}

func GetNumberOfPosts() (int, error) {
	defer posts.Exec("COMMIT;")
	posts.Exec("BEGIN;")
	rows, execErr := posts.Query("SELECT COUNT(*) FROM posts;")
	if execErr != nil {
		return 0, execErr
	}
	defer rows.Close()
	var SIZE int
	for rows.Next() {
		scanErr := rows.Scan(&SIZE)
		if scanErr != nil {
			return 0, scanErr
		}
	}
	return SIZE, nil
}

func GetSubsequenceOfPosts(start int, finish int) ([]ArticlePrinter, error) {
	defer posts.Exec("COMMIT;")
	posts.Exec("BEGIN;")
	const ask = "SELECT postID, userName, title, content, categories FROM posts WHERE postID >= ? AND postID <= ?;"
	rows, queryErr := posts.Query(ask, start, finish)
	if queryErr != nil {
		log.Println("[GetSubsequenceOfPosts: queryErr] -> " + queryErr.Error())
		return nil, queryErr
	}
	defer rows.Close()
	result := []ArticlePrinter{}
	for rows.Next() {
		buff := ArticlePrinter{PostId: 0, UserNick: "DNE", Title: "DNE", Categories: []string{"DNE"}, Content: "DNE"}
		var catline string
		scanErr := rows.Scan(&buff.PostId, &buff.UserNick, &buff.Title, &buff.Content, &catline)
		if scanErr != nil {
			log.Println("[GetSubsequenceOfPosts: scanErr] ->" + scanErr.Error())
			return nil, scanErr
		}
		buff.Categories = strings.Split(catline, ";")
		result = append(result, buff)
	}
	return result, nil
}

func CheckLikesDislikesAction(postID int, userName string) (bool, error) {
	defer posts.Exec("COMMIT;")
	posts.Exec("BEGIN;")
	rows, rowErr := posts.Query("SELECT userName FROM likesndislikes WHERE postID=?;", postID)
	if rowErr != nil {
		return false, rowErr
	}
	defer rows.Close()
	for rows.Next() {
		var nick string
		scanErr := rows.Scan(&nick)
		if scanErr != nil {
			return false, scanErr
		}
		if nick == userName {
			return true, nil
		}
	}
	return false, nil
}

func GetLikesDislikesAction(postID int, userName string) (string, error) {
	posts.Exec("BEGIN;")
	defer posts.Exec("COMMIT;")
	rows, rowErr := posts.Query("SELECT typeText FROM likesndislikes WHERE postID=? AND userName=?;", postID, userName)
	if rowErr != nil {
		return "", rowErr
	}
	defer rows.Close()
	var action string
	counter := 0
	for rows.Next() {
		scanErr := rows.Scan(&action)
		if scanErr != nil {
			return "", scanErr
		}
		counter++
	}
	if counter != 1 {
		return "", errors.New("multiple occurance of like updates")
	}
	return action, nil
}

func UpdateLikesDislikes(postId int, userName string, typeText string) error {
	defer posts.Exec("COMMIT;")
	posts.Exec("BEGIN;")
	_, execErr := posts.Exec("UPDATE likesndislikes SET typeText=? WHERE postID=? AND userName=?;", typeText, postId, userName)
	if execErr != nil {
		return execErr
	}
	return nil
}

func AddLikesDislikes(postId int, userName string, typeText string) error {
	defer posts.Exec("COMMIT;")
	posts.Exec("BEGIN;")
	_, execErr := posts.Exec("INSERT INTO likesndislikes VALUES (?,?,?);", postId, userName, typeText)
	if execErr != nil {
		return execErr
	}
	return nil
}

func CountLikesDislikes(postId int) (int, int, error) {
	posts.Exec("BEGIN;")
	defer posts.Exec("COMMIT;")
	rows, rowsErr := posts.Query("SELECT likes, dislikes FROM posts WHERE postID=?;", postId)
	if rowsErr != nil {
		return -1, -1, rowsErr
	}
	defer rows.Close()
	likes := 0
	dislikes := 0
	counts := 0
	for rows.Next() {
		scanErr := rows.Scan(&likes, &dislikes)
		if scanErr != nil {
			return -1, -1, scanErr
		}
		counts++
	}
	if counts != 1 {
		return -1, -1, errors.New("duplicate postid error")
	}
	return likes, dislikes, nil
}

func PutLikesDislikes(postId int, likes int, dislikes int) error {
	defer posts.Exec("COMMIT;")
	posts.Exec("BEGIN;")
	_, execErr := posts.Exec("UPDATE posts SET likes=?, dislikes=? WHERE postID=?;", likes, dislikes, postId)
	if execErr != nil {
		return execErr
	}
	return nil
}

func CountCommentsOfThePost(postID int) (int, error) {
	defer comms.Exec("COMMIT;")
	comms.Exec("BEGIN;")
	res, qErr := comms.Query("SELECT COUNT(*) FROM comments WHERE postID=?;", postID)
	if qErr != nil {
		return -1, qErr
	}
	defer res.Close()
	var count int = 0
	for res.Next() {
		scanErr := res.Scan(&count)
		if scanErr != nil {
			return -1, scanErr
		}
	}
	return count, nil
}

func AddNewCommentToPost(postID int, comm CommentContainer) (int, error) {
	defer comms.Exec("COMMIT;")
	comms.Exec("BEGIN;")
	const run = "INSERT INTO comments VALUES (NULL, ?, ?, ?, 0, 0);"
	_, execErr := comms.Exec(run, comm.Contain, comm.PostId, comm.Author)
	if execErr != nil {
		return -1, execErr
	}
	row, qrErr := comms.Query("SELECT commentID FROM comments ORDER BY commentID DESC LIMIT 1;")
	if qrErr != nil {
		return -1, qrErr
	}
	defer row.Close()
	var last int = -1
	for row.Next() {
		scanErr := row.Scan(&last)
		if scanErr != nil {
			return -1, scanErr
		}
	}
	return last, nil
}

func GetCommentsWithPostId(postID int) ([]CommentContainer, error) {
	defer comms.Exec("COMMIT;")
	comms.Exec("BEGIN;")
	rows, qErr := comms.Query("SELECT * FROM comments;")
	if qErr != nil {
		return nil, qErr
	}
	defer rows.Close()
	result := []CommentContainer{}
	for rows.Next() {
		var buffer CommentContainer
		scanErr := rows.Scan(&buffer.CommentId, &buffer.Contain, &buffer.PostId, &buffer.Author, &buffer.Likes, &buffer.Dislikes)
		if scanErr != nil {
			return nil, scanErr
		}
		if buffer.PostId == postID {
			result = append(result, buffer)
		}
	}
	return result, nil
}

func CountCommentsLikesAndDislikes(commentID int) (int, int, error) {
	comms.Exec("BEGIN;")
	defer comms.Exec("COMMIT;")
	res, qErr := comms.Query("SELECT likes, dislikes FROM comments WHERE commentID=?;", commentID)
	if qErr != nil {
		return -1, -1, qErr
	}
	defer res.Close()
	var likes, dislikes int
	for res.Next() {
		scanErr := res.Scan(&likes, &dislikes)
		if scanErr != nil {
			return -1, -1, scanErr
		}
	}
	return likes, dislikes, nil
}

func UpdateCommentsLikesAndDislikes(commentID int, likes int, dislikes int) error {
	comms.Exec("BEGIN;")
	defer comms.Exec("COMMIT;")
	_, execErr := comms.Exec("UPDATE comments SET likes=?, dislikes=? WHERE commentID=?;", likes, dislikes, commentID)
	if execErr != nil {
		return execErr
	}
	return nil
}

func CheckCommentFeedback(commentID int, author string) (bool, error) {
	comms.Exec("BEGIN;")
	defer comms.Exec("COMMIT;")
	rows, qErr := comms.Query("SELECT typeText FROM feedbacks WHERE commentID=? AND author=?;", commentID, author)
	if qErr != nil {
		return false, qErr
	}
	defer rows.Close()
	var occurance int = 0
	var feedback string = "#"
	for rows.Next() {
		scanErr := rows.Scan(&feedback)
		if scanErr != nil {
			return false, scanErr
		}
		occurance++
	}
	if occurance == 0 {
		return false, nil
	}
	if occurance > 1 {
		return false, errors.New("multiple occurances of feedback's like/dislike action")
	}
	return true, nil
}

func AddCommentFeedback(commentID int, author string, grade string) error {
	comms.Exec("BEGIN;")
	defer comms.Exec("COMMIT;")
	_, execErr := comms.Exec("INSERT INTO feedbacks VALUES (?, ?, ?);", commentID, author, grade)
	if execErr != nil {
		return execErr
	}
	return nil
}

func UpdateCommentFeedback(commentID int, author string, grade string) error {
	comms.Exec("BEGIN;")
	defer comms.Exec("COMMIT;")
	_, execErr := comms.Exec("UPDATE feedbacks SET typeText=? WHERE commentID=? AND author=?;", grade, commentID, author)
	if execErr != nil {
		return execErr
	}
	return nil
}

func GetCommentFeedback(commentID int, author string) (string, error) {
	comms.Exec("BEGIN;")
	defer comms.Exec("COMMIT;")
	rows, qErr := comms.Query("SELECT typeText from feedbacks WHERE commentID=? AND author=?;", commentID, author)
	if qErr != nil {
		return "", qErr
	}
	defer rows.Close()
	var occurances int = 0
	var result string = ""
	for rows.Next() {
		scanErr := rows.Scan(&result)
		if scanErr != nil {
			return "", scanErr
		}
		occurances++
	}
	if occurances != 1 {
		return "", errors.New("feedback log did not existed")
	}
	return result, nil
}

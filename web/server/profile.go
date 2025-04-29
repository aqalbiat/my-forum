package server

import (
	"errors"
	"html/template"
	"log"
	"net/http"
)

/*
	Files                               Tables                               Golang objects

	comments.db                         comments, feedbacks                  comms

	posts.db                            posts, likesndislikes                posts

	userlist.db                         users, sessions                      db

*/

func GetUserProfileInfo(userName string) (*UserProfilePrinter, error) {
	var result *UserProfilePrinter = &UserProfilePrinter{}
	(*result).NotLogged = false
	(*result).Logged = true
	(*result).UserHandler = userName
	emailRow, emailRowErr := db.Query("SELECT userEmail FROM users WHERE userName=?;", userName)
	if emailRowErr != nil {
		log.Println("[GetUserProfileInfo: emailRowsErr] -> " + emailRowErr.Error())
		return nil, emailRowErr
	}
	defer emailRow.Close()
	for emailRow.Next() {
		scanErr := emailRow.Scan(&(*result).UserEmail)
		if scanErr != nil {
			log.Println("[GetUserProfileInfo: scanErr] -> " + scanErr.Error())
			return nil, scanErr
		}
	}
	countPosts, countingPostsErr := posts.Query("SELECT COUNT(postID) FROM posts WHERE userName=?;", userName)
	if countingPostsErr != nil {
		log.Println("[GetUserProfileInfo: countingPostsErr] -> " + countingPostsErr.Error())
		return nil, countingPostsErr
	}
	defer countPosts.Close()
	for countPosts.Next() {
		scanPostsCountErr := countPosts.Scan(&(*result).NumberOfCreatedPosts)
		if scanPostsCountErr != nil {
			return nil, scanPostsCountErr
		}
	}
	const countLikesQuery = "SELECT COUNT (postID) FROM likesndislikes WHERE userName=? AND typeText='upgrade';"
	countLikesPosts, countLikesErr := posts.Query(countLikesQuery, userName)
	if countLikesErr != nil {
		log.Println("[GetUserProfileInfo: countLikesErr] -> " + countLikesErr.Error())
		return nil, countLikesErr
	}
	defer countLikesPosts.Close()
	if countLikesPosts.Next() {
		scanPostLikesErr := countLikesPosts.Scan(&(*result).NumberOfLikedPosts)
		if scanPostLikesErr != nil {
			log.Println("[GetUserProfileInfo: scanPostLikesErr] -> " + scanPostLikesErr.Error())
			return nil, scanPostLikesErr
		}
	}
	const countDislikesQuery = "SELECT COUNT (postID) FROM likesndislikes WHERE userName=? AND typeText='downgrade';"
	countDislikesPosts, countDislikesErr := posts.Query(countDislikesQuery, userName)
	if countDislikesErr != nil {
		log.Println("[GetUserProfileInfo: countDislikesErr] -> " + countDislikesErr.Error())
		return nil, countDislikesErr
	}
	defer countDislikesPosts.Close()
	for countDislikesPosts.Next() {
		scanPostDislikesErr := countDislikesPosts.Scan(&(*result).NumberOfDislikedPosts)
		if scanPostDislikesErr != nil {
			log.Println("[GetUserProfileInfo: scanPostDislikesErr] -> " + scanPostDislikesErr.Error())
			return nil, scanPostDislikesErr
		}
	}
	countComments, commentsErr := comms.Query("SELECT COUNT (commentID) FROM comments WHERE author=?;", userName)
	if commentsErr != nil {
		log.Println("[GetUserProfileInfo: commentsErr] -> " + commentsErr.Error())
		return nil, commentsErr
	}
	defer countComments.Close()
	for countComments.Next() {
		scanCountCommentsErr := countComments.Scan(&(*result).NumberOfCreatedComments)
		if scanCountCommentsErr != nil {
			log.Println("[GetUserProfileInfo: scanCountCommentsErr] -> " + scanCountCommentsErr.Error())
			return nil, scanCountCommentsErr
		}
	}
	const commentLikesQuery = "SELECT COUNT (commentID) FROM feedbacks WHERE author=? AND typeText='upgrade';"
	countCommentLikes, commentLikesErr := comms.Query(commentLikesQuery, userName)
	if commentLikesErr != nil {
		log.Println("[GetUserProfileInfo: commentLikesErr] -> " + commentLikesErr.Error())
		return nil, commentLikesErr
	}
	defer countCommentLikes.Close()
	for countCommentLikes.Next() {
		scanCommentLikesErr := countCommentLikes.Scan(&(*result).NumberOfLikedComments)
		if scanCommentLikesErr != nil {
			log.Println("[GetUserProfileInfo: scanCommentLikesErr] -> " + scanCommentLikesErr.Error())
			return nil, scanCommentLikesErr
		}
	}
	const commentDislikesQuery = "SELECT COUNT (commentID) FROM feedbacks WHERE author=? AND typeText='downgrade';"
	countCommentDislikes, commentDislikesErr := comms.Query(commentDislikesQuery, userName)
	if commentDislikesErr != nil {
		log.Println("[GetUserProfileInfo: commentDislikesErr] -> " + commentDislikesErr.Error())
		return nil, commentDislikesErr
	}
	defer countCommentDislikes.Close()
	for countCommentDislikes.Next() {
		scanCommentDislikesErr := countCommentDislikes.Scan(&(*result).NumberOfDislikedComments)
		if scanCommentDislikesErr != nil {
			log.Println("[GetUserProfileInfo: scanCommentDislikesErr] -> " + scanCommentDislikesErr.Error())
			return nil, scanCommentDislikesErr
		}
	}
	return result, nil
}

func UserProfileHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/profile/" {
		log.Println("[UserProfileHandler: 404]")
		ErrorHandler(w, http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		log.Println("[UserProfileHandler: 405]")
		ErrorHandler(w, http.StatusMethodNotAllowed)
		return
	}
	cookieJar, loadSessionErr := r.Cookie("sid")
	if loadSessionErr != nil {
		if errors.Is(loadSessionErr, http.ErrNoCookie) {
			tmpl, tmplLoadErr := template.ParseFiles("./web/UI/templates/profile.html")
			if tmplLoadErr != nil {
				log.Println("[UserProfileHandler: tmplLoadErr] -> " + tmplLoadErr.Error())
				ErrorHandler(w, http.StatusInternalServerError)
				return
			}
			tmplExecErr := tmpl.Execute(w, UserProfilePrinter{NotLogged: true, Logged: false})
			if tmplExecErr != nil {
				log.Println("[UserProfileHandler: tmplExecErr] -> " + tmplExecErr.Error())
				ErrorHandler(w, http.StatusInternalServerError)
			}
			return
		}
		log.Println("[UserProfileHandler: loadSessionErr] -> " + loadSessionErr.Error())
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}
	sessionId := cookieJar.Value
	isLogged, getLogErr := CheckExistanceOfSessionID(sessionId)
	if getLogErr != nil {
		log.Println("[UserProfileHandler: getLogErr] -> " + getLogErr.Error())
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}
	if !isLogged {
		tmpl, tmplLoadErr := template.ParseFiles("./web/UI/templates/profile.html")
		if tmplLoadErr != nil {
			log.Println("[UserProfileHandler: tmplLoadErr] -> " + tmplLoadErr.Error())
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}
		tmplExecErr := tmpl.Execute(w, UserProfilePrinter{NotLogged: true, Logged: false})
		if tmplExecErr != nil {
			log.Println("[UserProfileHandler: tmplExecErr] -> " + tmplExecErr.Error())
			ErrorHandler(w, http.StatusInternalServerError)
		}
		return
	}
	userHandle, getHandleErr := GetUsernameBySessionID(sessionId)
	if getHandleErr != nil {
		log.Println("[UserProfileHandler: userHandle] -> " + getHandleErr.Error())
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}
	userNick := r.URL.Query().Get("userhandle")
	if userHandle != userNick {
		log.Println("[UserProfileHandler: usernames not equal]")
		ErrorHandler(w, http.StatusNotAcceptable)
		return
	}
	info, getInfoErr := GetUserProfileInfo(userHandle)
	if getInfoErr != nil {
		log.Println("[UserProfileHandler: getInfoErr] -> " + getInfoErr.Error())
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}
	profTmpl, profTmplErr := template.ParseFiles("./web/UI/templates/profile.html")
	if profTmplErr != nil {
		log.Println("[UserProfileHandler: profTmplErr] -> " + profTmplErr.Error())
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}
	execProfTmplErr := profTmpl.Execute(w, (*info))
	if execProfTmplErr != nil {
		log.Println("[UserProfileHandler: execProfTmplErr] -> " + execProfTmplErr.Error())
		ErrorHandler(w, http.StatusInternalServerError)
	}
}

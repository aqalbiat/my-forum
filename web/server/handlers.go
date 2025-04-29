package server

import (
	"encoding/json"
	"errors"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

func MainHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		ErrorHandler(w, 404)
		return
	}
	if r.Method != http.MethodGet {
		ErrorHandler(w, 405)
		return
	}
	tmpl, err := template.ParseFiles("./web/UI/templates/index.html")
	if err != nil {
		ErrorHandler(w, 501)
		return
	}
	cookie, cookieErr := r.Cookie("sid")
	if errors.Is(cookieErr, http.ErrNoCookie) {
		tmpl.Execute(w, PageLoader{IsEmpty: true, Username: "#"})
		return
	}
	sessionID := string(cookie.Value)
	log.Println("Got the sessionID = " + sessionID)
	status, getStatErr := sessionList.IfSessionIdExists(sessionID)
	if getStatErr != nil {
		log.Println(getStatErr)
		ErrorHandler(w, 500)
		return
	}
	if status {
		userName, getNickErr := sessionList.RetrieveUsernameWithSID(sessionID)
		if getNickErr != nil {
			log.Println(getNickErr)
			ErrorHandler(w, 500)
			return
		}
		tmpl.Execute(w, PageLoader{IsEmpty: false, Username: userName})
		return
	}
	loadData := PageLoader{IsEmpty: true, Username: "#"}
	tmpl.Execute(w, loadData)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/logout/" {
		ErrorHandler(w, http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		ErrorHandler(w, http.StatusMethodNotAllowed)
		return
	}
	cookie, readErr := r.Cookie("sid")
	if errors.Is(readErr, http.ErrNoCookie) {
		tmpl, err := template.ParseFiles("./web/UI/templates/index.html")
		if err != nil {
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, PageLoader{IsEmpty: true, Username: "#"})
		cookie.Value = ""
		return
	}
	if readErr != nil {
		log.Println(readErr)
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}
	_, deletionErr := DeleteSessionWithID(cookie.Value)
	if deletionErr != nil {
		log.Println(deletionErr)
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}
	tmpl, err := template.ParseFiles("./web/UI/templates/index.html")
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, PageLoader{IsEmpty: true, Username: "#"})
}

func MenuHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/menu/" {
		ErrorHandler(w, http.StatusNotFound)
		return
	}
	if r.Method == http.MethodPost {
		ErrorHandler(w, http.StatusNotImplemented)
		return
	}
	if r.Method != http.MethodGet {
		ErrorHandler(w, http.StatusMethodNotAllowed)
		return
	}
	button := r.URL.Query().Get("button")
	if button == "about_us" {
		tmpl, err := template.ParseFiles("./web/UI/templates/aboutus.html")
		if err != nil {
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
		return
	}
	if button == "posts" {
		biscuits, biscErr := r.Cookie("sid")
		if errors.Is(http.ErrNoCookie, biscErr) {
			tmpl, tmplErr := template.ParseFiles("./web/UI/templates/message.html")
			if tmplErr != nil {
				log.Println("[MenuHandler] -> message template error (no cookie)")
				ErrorHandler(w, http.StatusInternalServerError)
				return
			}
			tmpl.Execute(w, UserMessage{Message: "PLEASE, LOGIN FIRST!"})
			return
		}
		if biscErr != nil {
			log.Println("cookie load error")
			ErrorHandler(w, http.StatusBadRequest)
			return
		}
		status, statErr := sessionList.IfSessionIdExists(biscuits.Value)
		if statErr != nil {
			log.Println("[MenuHandler: statErr] -> " + statErr.Error())
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}
		if status {
			tmpl, tmplErr := template.ParseFiles("./web/UI/templates/posting.html")
			if tmplErr != nil {
				log.Println("[MenuHandler: tmplErr (loading posting.html)] -> " + tmplErr.Error())
				ErrorHandler(w, http.StatusInternalServerError)
				return
			}
			tmpl.Execute(w, categoryJar)
		} else {
			tmpl, tmplErr := template.ParseFiles("./web/UI/templates/message.html")
			if tmplErr != nil {
				ErrorHandler(w, http.StatusInternalServerError)
				return
			}
			tmpl.Execute(w, UserMessage{Message: "PLEASE, LOGIN FIRST!"})
		}
		return
	}
	if button == "categories" {
		letter := r.URL.Query().Get("letter")
		id := int([]rune(letter)[0] - 'A')
		tmpl, tmplErr := template.ParseFiles("./web/UI/templates/categories.html")
		if tmplErr != nil {
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}
		log.Println("Category list with letter " + letter + " was sended!")
		tmpl.Execute(w, categoryJar.Categories[id])
		return
	}
	if button == "main" {
		page := r.URL.Query().Get("pg")
		nbr, convErr := strconv.Atoi(page)
		if convErr != nil {
			ErrorHandler(w, http.StatusBadRequest)
			return
		}
		log.Println("User sended a query to get page #" + page + " of the articles.")
		countOfArt, dbErr := GetNumberOfPosts()
		if dbErr != nil {
			log.Println("[MenuHandler: dbErr] -> " + dbErr.Error())
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}
		l := (nbr-1)*5 + 1
		r := nbr * 5
		if l > countOfArt {
			tmpl, tmplErr := template.ParseFiles("./web/UI/templates/message.html")
			if tmplErr != nil {
				log.Println("[MenuHandler: tmplErr] -> " + tmplErr.Error())
				ErrorHandler(w, http.StatusInternalServerError)
				return
			}
			tmpl.Execute(w, UserMessage{Message: "NO SUCH NUMBER OF PAGES"})
			return
		}
		printers, printErr := GetSubsequenceOfPosts(l, r)
		if printErr != nil {
			log.Println("[MenuHandler: printErr] -> " + printErr.Error())
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}
		for i := 0; i < len(printers); i++ {
			printers[i].Content = printers[i].Content[:100]
		}
		tmplMain, tmplMainErr := template.ParseFiles("./web/UI/templates/main.html")
		if tmplMainErr != nil {
			log.Println("[MenuHandler: tmplMainErr] -> " + tmplMainErr.Error())
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}
		bufferedSender := MainPageCollection{}
		bufferedSender.ArticlesFace = printers
		bufferedSender.PageNumber = nbr
		tmplMain.Execute(w, bufferedSender)
		return
	}
	if button == "search" {
		tmpl, err := template.ParseFiles("./web/UI/templates/search.html")
		if err != nil {
			log.Println("[MenuHandler: (search template error)] -> " + err.Error())
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}
		execTmplErr := tmpl.Execute(w, categoryJar)
		if execTmplErr != nil {
			log.Println("[MenuHandler: execTmplErr] -> " + execTmplErr.Error())
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}
		return
	}
	ErrorHandler(w, http.StatusNotFound)
}

func ArticlesHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/api/articles/" {
		ErrorHandler(w, http.StatusNotFound)
		return
	}
	if r.Method == http.MethodPost {
		jar, cookieErr := r.Cookie("sid")
		if errors.Is(http.ErrNoCookie, cookieErr) {
			answer, answerErr := json.Marshal(UserMessage{Message: "NOT REGISTERED"})
			if answerErr != nil {
				ErrorHandler(w, http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write(answer)
			return
		}
		if cookieErr != nil {
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}
		sessionId := jar.Value
		status, sidErr := sessionList.IfSessionIdExists(sessionId)
		if sidErr != nil {
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}
		if !status {
			tmpl, tmplErr := template.ParseFiles("./web/UI/templates/message.html")
			if tmplErr != nil {
				ErrorHandler(w, http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusNotAcceptable)
			tmpl.Execute(w, UserMessage{Message: "LOGIN PLEASE!!!"})
			return
		}
		userNick, userErr := GetUsernameBySessionID(sessionId)
		if userErr != nil {
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}
		act := r.URL.Query().Get("action")
		if act == "create" {
			log.Println("Recieved create call from user")
			var buffer ArticleSender
			jsonErr := json.NewDecoder(r.Body).Decode(&buffer)
			if jsonErr != nil {
				ErrorHandler(w, http.StatusInternalServerError)
				return
			}
			log.Println("Recieved object")
			log.Println(buffer)
			postErr := AddNewPost(userNick, buffer.Title, buffer.Content, 0, buffer.Categories)
			if postErr != nil {
				ErrorHandler(w, http.StatusInternalServerError)
				return
			}
			tmpl, tmplErr := template.ParseFiles("./web/UI/templates/message.html")
			if tmplErr != nil {
				ErrorHandler(w, http.StatusInternalServerError)
				return
			}
			tmpl.Execute(w, UserMessage{Message: "YOUR ARTICLE WAS ADDED"})
			return
		}
		if act == "likes" {
			log.Println("Recieved like/dislike call from user")
			id := r.URL.Query().Get("id")
			grade := r.URL.Query().Get("type")
			nbr, err := strconv.Atoi(id)
			if err != nil {
				log.Println("[ArticleHandler: err (string conversion)] -> " + err.Error())
				ErrorHandler(w, http.StatusInternalServerError)
				return
			}
			status, statusErr := CheckLikesDislikesAction(nbr, userNick)
			if statusErr != nil {
				log.Println("[ArticleHandler: statusErr] -> " + statusErr.Error())
				ErrorHandler(w, http.StatusInternalServerError)
				return
			}
			likes, dislikes, ldErr := CountLikesDislikes(nbr)
			if ldErr != nil {
				log.Println("[ArticleHandler: ldErr] -> " + ldErr.Error())
				ErrorHandler(w, http.StatusInternalServerError)
				return
			}
			if status {
				dbGrade, gradeErr := GetLikesDislikesAction(nbr, userNick)
				if gradeErr != nil {
					log.Println("[ArticleHandler: gradeErr] -> " + gradeErr.Error())
					ErrorHandler(w, http.StatusInternalServerError)
					return
				}
				log.Printf("Recieved from %v to %v the post\n", userNick, grade)
				if dbGrade != grade {
					if dbGrade == "upgrade" {
						if grade == "neutral" {
							likes--
						} else {
							likes--
							dislikes++
						}
					} else {
						if grade == "neutral" {
							dislikes--
						} else {
							likes++
							dislikes--
						}
					}
					dbErr := UpdateLikesDislikes(nbr, userNick, grade)
					if dbErr != nil {
						log.Println("[ArticleHandler: dbErr] -> " + dbErr.Error())
						ErrorHandler(w, http.StatusInternalServerError)
						return
					}
					putErr := PutLikesDislikes(nbr, likes, dislikes)
					if putErr != nil {
						log.Println("[ArticleHandler: putErr] -> " + putErr.Error())
						ErrorHandler(w, http.StatusInternalServerError)
						return
					}
					buffer, buffErr := json.Marshal(UserMessage{Message: "Accepted"})
					if buffErr != nil {
						log.Println("[ArticleHandler: buffErr] -> " + buffErr.Error())
						ErrorHandler(w, http.StatusInternalServerError)
						return
					}
					w.Write(buffer)
				} else {
					buffer, buffErr := json.Marshal(UserMessage{Message: "EXISTS"})
					if buffErr != nil {
						log.Println("[ArticleHandler: buffErr] -> " + buffErr.Error())
						ErrorHandler(w, http.StatusInternalServerError)
						return
					}
					w.WriteHeader(http.StatusOK)
					w.Write(buffer)
				}
			} else {
				if grade == "upgrade" {
					likes++
				} else {
					dislikes++
				}
				putErr := PutLikesDislikes(nbr, likes, dislikes)
				if putErr != nil {
					log.Println("[ArticleHandler: putErr] -> " + putErr.Error())
					ErrorHandler(w, http.StatusInternalServerError)
					return
				}
				addErr := AddLikesDislikes(nbr, userNick, grade)
				if addErr != nil {
					log.Println("[ArticleHandler: addErr] -> " + addErr.Error())
					ErrorHandler(w, http.StatusInternalServerError)
					return
				}
				buffer, buffErr := json.Marshal(UserMessage{Message: "Accepted"})
				if buffErr != nil {
					log.Println("[ArticleHandler: buffErr] -> " + buffErr.Error())
					ErrorHandler(w, http.StatusInternalServerError)
					return
				}
				w.WriteHeader(http.StatusOK)
				w.Write(buffer)
			}
			return
		}
		if act == "comments" {
			log.Println("Recieved a call for attachment of a comment")
			var userData CommentContainer
			jsonErr := json.NewDecoder(r.Body).Decode(&userData)
			if jsonErr != nil {
				log.Println("[ArticlesHandler: jsonErr] -> " + jsonErr.Error())
				ErrorHandler(w, http.StatusInternalServerError)
				return
			}
			log.Println("Recieved a comment from user: ", userData)
			userData.Author = userNick
			commentId, dbErr := AddNewCommentToPost(userData.PostId, userData)
			if dbErr != nil {
				log.Println("[ArticlesHandler: dbErr] -> " + dbErr.Error())
				ErrorHandler(w, http.StatusInternalServerError)
				return
			}
			buffer, convErr := json.Marshal(UserMessage{Message: "OK" + "#" + userNick + "#" + strconv.Itoa(commentId)})
			if convErr != nil {
				log.Println("[ArticleHandler: convErr] -> " + convErr.Error())
				ErrorHandler(w, http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write(buffer)
			return
		}
		if act == "feedback" {
			log.Println("Recieved a call to tag a like/dislike on comment")
			bufferedCommentId := r.URL.Query().Get("id")
			commentId, convertErr := strconv.Atoi(bufferedCommentId)
			if convertErr != nil {
				log.Println("[ArticlesHandler: convertErr] -> " + convertErr.Error())
				ErrorHandler(w, http.StatusInternalServerError)
				return
			}
			status, feedbackErr := CheckCommentFeedback(commentId, userNick)
			if feedbackErr != nil {
				log.Println("[ArticlesHandler: feedbackErr] -> " + feedbackErr.Error())
				ErrorHandler(w, http.StatusInternalServerError)
				return
			}
			likes, dislikes, countErr := CountCommentsLikesAndDislikes(commentId)
			if countErr != nil {
				log.Println("[ArticlesHandler: countErr] -> " + countErr.Error())
				ErrorHandler(w, http.StatusInternalServerError)
				return
			}
			grade := r.URL.Query().Get("grade")
			if status {
				preGrade, dbLoadErr := GetCommentFeedback(commentId, userNick)
				if dbLoadErr != nil {
					log.Println("[ArticlesHandler: dbLoadErr (comment feedback part)] -> " + dbLoadErr.Error())
					ErrorHandler(w, http.StatusInternalServerError)
					return
				}
				dbRefactorErr := UpdateCommentFeedback(commentId, userNick, grade)
				if dbRefactorErr != nil {
					log.Println("[ArticlesHandler: dbRefactorErr] -> " + dbRefactorErr.Error())
					ErrorHandler(w, http.StatusInternalServerError)
					return
				}
				if preGrade != grade {
					if grade == "upgrade" {
						if preGrade == "neutral" {
							likes++
						} else {
							likes--
							dislikes++
						}
					} else if grade == "downgrade" {
						if preGrade == "neutral" {
							dislikes++
						} else {
							dislikes++
							likes--
						}
					} else {
						if preGrade == "upgrade" {
							likes--
						} else {
							dislikes++
						}
					}
					updateErr := UpdateCommentsLikesAndDislikes(commentId, likes, dislikes)
					if updateErr != nil {
						log.Println("[ArticlesHandler: updateErr (comment feedback part)] -> " + updateErr.Error())
						ErrorHandler(w, http.StatusInternalServerError)
						return
					}
					bufferedSender, jsonMarshalError := json.Marshal(UserMessage{Message: "OK"})
					if jsonMarshalError != nil {
						log.Println("[Articleshandler: jsonMarhallError] -> " + jsonMarshalError.Error())
						ErrorHandler(w, http.StatusInternalServerError)
						return
					}
					w.WriteHeader(http.StatusOK)
					w.Write(bufferedSender)
				} else {
					bufferedAnswer, jsonErr := json.Marshal(UserMessage{Message: "EXISTS"})
					if jsonErr != nil {
						log.Println("[ArticlesHandler: jsonErr (comment feedback part)] -> " + jsonErr.Error())
						ErrorHandler(w, http.StatusInternalServerError)
						return
					}
					w.WriteHeader(http.StatusOK)
					w.Write(bufferedAnswer)
				}
			} else {
				if grade == "upgrade" {
					likes++
				} else {
					dislikes++
				}
				updateErr := UpdateCommentsLikesAndDislikes(commentId, likes, dislikes)
				if updateErr != nil {
					log.Println("[ArticlesHandler: updateErr] -> " + updateErr.Error())
					ErrorHandler(w, http.StatusInternalServerError)
					return
				}
				addFeedbackError := AddCommentFeedback(commentId, userNick, grade)
				if addFeedbackError != nil {
					log.Println("[ArticlesHandler: addFeedbackErr] -> " + addFeedbackError.Error())
					ErrorHandler(w, http.StatusInternalServerError)
					return
				}
				bufferedSender, jsonErr := json.Marshal(UserMessage{Message: "OK"})
				if jsonErr != nil {
					log.Println("[ArticlesHandler: jsonErr (feedback comment part)] -> " + jsonErr.Error())
					ErrorHandler(w, http.StatusInternalServerError)
					return
				}
				w.WriteHeader(http.StatusOK)
				w.Write(bufferedSender)
			}
			return
		}
		ErrorHandler(w, http.StatusNotFound)
		return
	}
	if r.Method == http.MethodGet {
		act := r.URL.Query().Get("action")
		if act == "page_number" {
			nbr, nbrErr := GetNumberOfPosts()
			response := MaxNumberOfPages{}
			if nbrErr != nil {
				response.NumberOfPages = -1
				response.ServerMessage = nbrErr.Error()
			} else {
				const div = 5
				response.NumberOfPages = (nbr + div - 1) / div
				response.ServerMessage = "OK"
			}
			sender, Err := json.Marshal(response)
			if Err != nil {
				ErrorHandler(w, http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusAccepted)
			w.Write(sender)
			return
		}
		if act == "textpage" {
			id := r.URL.Query().Get("id")
			log.Printf("user sended request for article #%v\n", id)
			idx, err := strconv.Atoi(id)
			if err != nil {
				ErrorHandler(w, http.StatusInternalServerError)
				return
			}
			textPage, textErr := GetSubsequenceOfPosts(idx, idx)
			if textErr != nil {
				ErrorHandler(w, http.StatusInternalServerError)
				return
			}
			likes, dislikes, countErr := CountLikesDislikes(idx)
			if countErr != nil {
				ErrorHandler(w, http.StatusInternalServerError)
				return
			}
			textPage[0].Likes = likes
			textPage[0].Dislikes = dislikes
			cookieJar, cookieErr := r.Cookie("sid")
			if errors.Is(cookieErr, http.ErrNoCookie) {
				textPage[0].UpdateGradeStatus("neutral")
			} else {
				if cookieErr != nil {
					log.Println("[ArticlesHandler: cookieErr]" + cookieErr.Error())
					ErrorHandler(w, http.StatusInternalServerError)
					return
				}
				sessionId := cookieJar.Value
				userStat, respErr := CheckExistanceOfSessionID(sessionId)
				if respErr != nil {
					log.Println("[ArticlesHandler: respErr] -> " + respErr.Error())
					ErrorHandler(w, http.StatusInternalServerError)
					return
				}
				if userStat {
					userName, retriveFailErr := GetUsernameBySessionID(sessionId)
					if retriveFailErr != nil {
						log.Println("[ArticlesHandler: retriveFailErr] -> " + retriveFailErr.Error())
						ErrorHandler(w, http.StatusInternalServerError)
						return
					}
					isActionExists, callFailErr := CheckLikesDislikesAction(idx, userName)
					if callFailErr != nil {
						log.Println("[ArticlesHandler: callFailErr] -> " + callFailErr.Error())
						ErrorHandler(w, http.StatusInternalServerError)
						return
					}
					if isActionExists {
						grading, loadGradeFailErr := GetLikesDislikesAction(idx, userName)
						if loadGradeFailErr != nil {
							log.Println("[ArticlesHandler: loadGradeFailErr] -> " + loadGradeFailErr.Error())
							ErrorHandler(w, http.StatusInternalServerError)
							return
						}
						textPage[0].UpdateGradeStatus(grading)
					} else {
						textPage[0].UpdateGradeStatus("neutral")
					}
				} else {
					textPage[0].UpdateGradeStatus("neutral")
				}
			}
			// start loading the articles comments
			numberOfComms, commsCountErr := CountCommentsOfThePost(idx)
			if commsCountErr != nil {
				log.Println("[ArticlesHandler: commsCountErr] -> " + commsCountErr.Error())
				ErrorHandler(w, http.StatusInternalServerError)
				return
			}
			if numberOfComms > 0 {
				collectionOfComments, loadCommsErr := GetCommentsWithPostId(idx)
				if loadCommsErr != nil {
					log.Println("[ArticlesHandler: loadCommsErr] -> " + loadCommsErr.Error())
					ErrorHandler(w, http.StatusInternalServerError)
					return
				}
				textPage[0].ListOfComments = collectionOfComments
				if errors.Is(cookieErr, http.ErrNoCookie) {
					for i := 0; i < len(textPage[0].ListOfComments); i++ {
						textPage[0].ListOfComments[i].SetGrade("neutral")
					}
				} else {
					sessionId := cookieJar.Value
					userStat, loadStatErr := CheckExistanceOfSessionID(sessionId)
					if loadStatErr != nil {
						log.Println("[ArticlesHandler: loadStatErr] -> " + loadStatErr.Error())
						ErrorHandler(w, http.StatusInternalServerError)
						return
					}
					if !userStat {
						for i := 0; i < len(textPage[0].ListOfComments); i++ {
							textPage[0].ListOfComments[i].SetGrade("neutral")
						}
					} else {
						userNick, loadNickErr := GetUsernameBySessionID(sessionId)
						if loadNickErr != nil {
							log.Println("[ArticlesHandler: loadNickErr] -> " + loadNickErr.Error())
							ErrorHandler(w, http.StatusInternalServerError)
							return
						}
						for i := 0; i < len(textPage[0].ListOfComments); i++ {
							feedbackStat, fbStatErr := CheckCommentFeedback(textPage[0].ListOfComments[i].CommentId, userNick)
							if fbStatErr != nil {
								log.Println("[ArticlesHandler: fbStatErr] -> " + fbStatErr.Error())
								ErrorHandler(w, http.StatusInternalServerError)
								return
							}
							if feedbackStat {
								grade, feedbackGradeErr := GetCommentFeedback(textPage[0].ListOfComments[i].CommentId, userNick)
								if feedbackGradeErr != nil {
									log.Println("[ArticlesHandler: feedbackGradeErr] -> " + feedbackGradeErr.Error())
									ErrorHandler(w, http.StatusInternalServerError)
									return
								}
								textPage[0].ListOfComments[i].SetGrade(grade)
							} else {
								textPage[0].ListOfComments[i].SetGrade("neutral")
							}
						}
					}
				}
			}
			// start rendering it to the client
			tmpl, errTmpl := template.ParseFiles("./web/UI/templates/article.html")
			if errTmpl != nil {
				log.Println("[ArticlesHandler: errTmpl] ->" + errTmpl.Error())
				ErrorHandler(w, http.StatusInternalServerError)
				return
			}
			tmpl.Execute(w, textPage)
			return
		}
		ErrorHandler(w, http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet && r.Method != http.MethodPost && r.Method != http.MethodHead {
		ErrorHandler(w, http.StatusMethodNotAllowed)
		return
	}
	ErrorHandler(w, http.StatusNotFound)
}

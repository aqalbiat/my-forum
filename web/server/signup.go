package server

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

// merged from 2024-07-13

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	const currentPath string = "/login/"
	if r.URL.Path[:len(currentPath)] != currentPath {
		ErrorHandler(w, 404)
		return
	}
	if r.Method != http.MethodGet {
		ErrorHandler(w, 405)
		return
	}
	tmpl, err := template.ParseFiles("./web/UI/templates/login.html")
	if err != nil {
		ErrorHandler(w, 500)
		return
	}
	tmpl.Execute(w, nil)
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	const currentPath string = "/register/"
	if r.URL.Path[:len(currentPath)] != currentPath {
		ErrorHandler(w, 404)
		return
	}
	if r.Method != http.MethodGet {
		ErrorHandler(w, 405)
		return
	}
	tmpl, err := template.ParseFiles("./web/UI/templates/register.html")
	if err != nil {
		ErrorHandler(w, 500)
		return
	}
	tmpl.Execute(w, nil)
}

func CheckCredentials(w http.ResponseWriter, r *http.Request) {
	const currentPath string = "/api/credentials/"
	if r.URL.Path[:len(currentPath)] != currentPath {
		ErrorHandler(w, 404)
		return
	}
	if r.Method != http.MethodGet {
		ErrorHandler(w, 405)
		return
	}
	userName := r.URL.Query().Get("username")
	userEmail := r.URL.Query().Get("useremail")
	okName, nameErr := CheckForFreeName(userName)
	okEmail, emailErr := CheckForFreeEmail(userEmail)
	result := ResponseDB{UserNameResponse: "OK", UserEmailResponse: "OK"}
	if !okName {
		if nameErr != nil {
			result.UserNameResponse = nameErr.Error()
		} else {
			result.UserNameResponse = "EXISTS"
		}
	}
	if !okEmail {
		if emailErr != nil {
			result.UserEmailResponse = emailErr.Error()
		} else {
			result.UserEmailResponse = "EXISTS"
		}
	}
	responseJson, jsonErr := json.Marshal(result)
	if jsonErr != nil {
		ErrorHandler(w, 500)
		log.Println(jsonErr)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseJson)
}

func UserLoadHandler(w http.ResponseWriter, r *http.Request) {
	var currentPath string = "/loader/"
	if r.URL.Path[:len(currentPath)] != currentPath {
		ErrorHandler(w, 404)
		return
	}
	if r.Method != http.MethodPost {
		ErrorHandler(w, 405)
		log.Println("Problem here " + r.Method)
		return
	}
	requiredAction := r.URL.Query().Get("action")
	if requiredAction == "register" {
		var buff UserCredentials
		formErr := r.ParseForm()
		if formErr != nil {
			ErrorHandler(w, 403)
			log.Println(formErr)
			return
		}
		buff.NameOfUser = r.FormValue("NameOfUser")
		buff.EmailOfUser = r.FormValue("EmailOfUser")
		cookie, cookieErr := r.Cookie("password")
		if cookieErr != nil {
			ErrorHandler(w, 501)
			log.Println(cookieErr)
			return
		}
		buff.PassHash = string(cookie.Value)
		log.Println("Credentials recieved")
		log.Println(buff)
		h := md5.New()
		_, md5Err := h.Write([]byte(buff.PassHash))
		if md5Err != nil {
			ErrorHandler(w, 501)
			log.Println(md5Err)
			return
		}
		wh := sha256.New()
		wh.Write([]byte(buff.PassHash))
		buff.PassHash = fmt.Sprintf("%x", h.Sum(nil)) + fmt.Sprintf("%x", wh.Sum(nil))
		statusWrite, writeErr := AddNewUser(buff)
		if writeErr != nil {
			ErrorHandler(w, 501)
			log.Println(writeErr)
			return
		}
		log.Printf("Recieved status:%v\n", statusWrite)
		var message UserMessage
		if statusWrite {
			message.Message = "New user added!"
		} else {
			message.Message = "Something went wrong!"
		}
		resultTmpl, loadResultTmplErr := template.ParseFiles("./web/UI/templates/signed.html")
		if loadResultTmplErr != nil {
			log.Println("[UserLoadHandler: loadResultTmplErr] -> " + loadResultTmplErr.Error())
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}
		execResultTmplErr := resultTmpl.Execute(w, message)
		if execResultTmplErr != nil {
			log.Println("[UserLoadHandler: execResultTmplErr] -> " + execResultTmplErr.Error())
			ErrorHandler(w, http.StatusInternalServerError)
		}
		return
	}
	if requiredAction == "login" {
		var buff UserCredentials
		formErr := r.ParseForm()
		if formErr != nil {
			ErrorHandler(w, 403)
			log.Println(formErr)
			return
		}
		buff.EmailOfUser = r.FormValue("EmailOfUser")
		cookie, cookieErr := r.Cookie("password")
		if cookieErr != nil {
			ErrorHandler(w, 501)
			log.Println(cookieErr)
			return
		}
		buff.PassHash = string(cookie.Value)
		log.Println("Recieved login info")
		log.Println("Login credentials")
		log.Println(buff)
		h := md5.New()
		_, md5Err := h.Write([]byte(buff.PassHash))
		if md5Err != nil {
			ErrorHandler(w, 501)
			log.Println(md5Err)
			return
		}
		wh := sha256.New()
		_, sha256Err := wh.Write([]byte(buff.PassHash))
		if sha256Err != nil {
			ErrorHandler(w, 501)
			log.Println(sha256Err)
			return
		}
		buff.PassHash = fmt.Sprintf("%x", h.Sum(nil)) + fmt.Sprintf("%x", wh.Sum(nil))
		logRes, logErr := CheckForGoodEntry(buff)
		if logErr != nil {
			ErrorHandler(w, 501)
			log.Println(logErr)
			return
		}
		log.Printf("Recieved status: %v\n", logRes)
		if logRes {
			nickName, parseErr := GetUsernameFromEmail(buff.EmailOfUser)
			if parseErr != nil {
				ErrorHandler(w, 500)
				return
			}
			buff.NameOfUser = nickName
			sid, loadErr := sessionList.InsertUserSession(buff)
			if loadErr != nil {
				log.Println("[UserLoadHandler: loadErr] -> " + loadErr.Error())
				ErrorHandler(w, 500)
				return
			}
			settledCookie := http.Cookie{
				Name:     "sid",
				Value:    sid,
				Path:     "/",
				MaxAge:   2400,
				HttpOnly: true,
				Secure:   true,
				SameSite: http.SameSiteLaxMode,
			}
			http.SetCookie(w, &settledCookie)
			tmpl, loadTmplErr := template.ParseFiles("./web/UI/templates/index.html")
			if loadTmplErr != nil {
				ErrorHandler(w, 500)
				log.Println(loadTmplErr)
				return
			}
			loadFlow := PageLoader{IsEmpty: false, Username: buff.NameOfUser}
			tmpl.Execute(w, loadFlow)
		} else {
			message := UserMessage{Message: "No such user!"}
			resultTmpl, loadResultTmplErr := template.ParseFiles("./web/UI/templates/signed.html")
			if loadResultTmplErr != nil {
				log.Println("[UserLoadHandler: loadResultTmplErr] -> " + loadResultTmplErr.Error())
				ErrorHandler(w, http.StatusInternalServerError)
				return
			}
			execResultTmplErr := resultTmpl.Execute(w, message)
			if execResultTmplErr != nil {
				log.Println("[UserLoadHandler: execResultTmplErr] -> " + execResultTmplErr.Error())
				ErrorHandler(w, http.StatusInternalServerError)
			}
		}
		return
	}
	ErrorHandler(w, 404)
}

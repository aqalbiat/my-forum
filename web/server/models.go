package server

import (
	"crypto/md5"
	"crypto/sha256"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type ErrorBar struct {
	ErrorCode int
	ErrorMsg  string
}

type UserCredentials struct {
	NameOfUser  string `json:"nameOfUser"`
	EmailOfUser string `json:"emailOfUser"`
	PassHash    string
}

type ResponseDB struct {
	UserNameResponse  string
	UserEmailResponse string
}

type PageLoader struct {
	IsEmpty    bool
	Username   string
	LeftRoute  string
	RightRoute string
}

type UserMessage struct {
	Message string
}

type UID struct {
	QueryCount       int64
	QueryHistory     []string
	CurrentUserQueue []string
	InsertionTime    []int64
}

func (userID *UID) InsertUserSession(user UserCredentials) (string, error) {
	currentTime := time.Now().Unix()
	entropySource := user.NameOfUser + user.EmailOfUser + time.Now().String() + user.PassHash
	hashMD5 := md5.New()
	_, md5Err := hashMD5.Write([]byte(entropySource))
	if md5Err != nil {
		return "", md5Err
	}
	hashSHA256 := sha256.New()
	_, sha256Err := hashSHA256.Write([]byte(entropySource))
	if sha256Err != nil {
		return "", sha256Err
	}
	sid := fmt.Sprintf("%x", hashMD5.Sum(nil)) + fmt.Sprintf("%x", hashSHA256.Sum(nil))
	addStatus, addErr := AddNewSessionID(sid, user.NameOfUser)
	if addErr != nil {
		log.Println("[InserUserSession: addErr] -> " + addErr.Error())
		return "", addErr
	}
	if addStatus {
		userID.QueryCount++
		userID.QueryHistory = append(userID.QueryHistory, "ADD USER")
		userID.CurrentUserQueue = append(userID.CurrentUserQueue, user.NameOfUser)
		userID.InsertionTime = append(userID.InsertionTime, currentTime)
	}
	return sid, nil
}

func (userID *UID) RetrieveUsernameWithSID(sid string) (string, error) {
	userName, getErr := GetUsernameBySessionID(sid)
	if getErr != nil {
		return "", getErr
	}
	userID.QueryCount++
	userID.QueryHistory = append(userID.QueryHistory, "GET USERNAME WITH SID")
	userID.CurrentUserQueue = append(userID.CurrentUserQueue, userName)
	userID.InsertionTime = append(userID.InsertionTime, time.Now().Unix())
	return userName, nil
}

func (userID *UID) IfSessionIdExists(sid string) (bool, error) {
	status, getErr := CheckExistanceOfSessionID(sid)
	if getErr != nil {
		return false, getErr
	}
	userID.QueryCount++
	userID.QueryHistory = append(userID.QueryHistory, "CHECK SESSION ID")
	userID.CurrentUserQueue = append(userID.CurrentUserQueue, "N/A")
	userID.InsertionTime = append(userID.InsertionTime, time.Now().Unix())
	return status, nil
}

type Theme struct {
	ThemeId   int
	TopicName string
}

type LetterThemes struct {
	Letter string
	Topics []Theme
}

type CategoriesList struct {
	Categories []LetterThemes
}

const MAXN = (1 << 16)

func (cats *CategoriesList) LoadAllData(path string) error {
	file, openErr := os.Open(path)
	if openErr != nil {
		return openErr
	}
	var inputBuffer [MAXN]byte
	len, readErr := file.Read(inputBuffer[:])
	if readErr != nil {
		return readErr
	}
	current := string(inputBuffer[:len])
	lines := strings.Split(current, "\r\n")
	pos := 0
	for i := 0; i < 26; i++ {
		currentRune := lines[pos][4]
		currentSlice := lines[pos+1 : pos+11]
		(*cats).Categories = append((*cats).Categories, LetterThemes{Letter: string(currentRune)})
		for j := 0; j < 10; j++ {
			parts := strings.Split(currentSlice[j], ".")
			retId, retErr := strconv.Atoi(parts[0])
			if retErr != nil {
				return retErr
			}
			(*cats).Categories[i].Topics = append((*cats).Categories[i].Topics, Theme{ThemeId: retId, TopicName: parts[1][1:]})
		}
		pos += 12
	}
	return nil
}

type CommentContainer struct {
	CommentId int
	Contain   string `json:"content"`
	PostId    int    `json:"postId"`
	Author    string
	Likes     int
	Dislikes  int
	IsNeutral bool
	IsRed     bool
	IsGreen   bool
}

func (cc *CommentContainer) SetGrade(grade string) {
	if grade == "neutral" {
		cc.IsNeutral = true
		cc.IsGreen = false
		cc.IsRed = false
	}
	if grade == "upgrade" {
		cc.IsGreen = true
		cc.IsNeutral = false
		cc.IsRed = false
	}
	if grade == "downgrade" {
		cc.IsRed = true
		cc.IsNeutral = false
		cc.IsRed = false
	}
}

type ArticleSender struct {
	Title      string `json:"title"`
	Categories string `json:"category"`
	Content    string `json:"article"`
}

type ArticlePrinter struct {
	PostId         int
	UserNick       string
	Title          string
	Categories     []string
	Content        string
	Likes          int
	Dislikes       int
	IsRed          bool
	IsGreen        bool
	IsNeutral      bool
	ListOfComments []CommentContainer
}

func (ap *ArticlePrinter) UpdateGradeStatus(stat string) {
	if stat == "neutral" {
		(*ap).IsNeutral = true
		(*ap).IsGreen = false
		(*ap).IsRed = false
	} else if stat == "upgrade" {
		(*ap).IsNeutral = false
		(*ap).IsGreen = true
		(*ap).IsRed = false
	} else {
		(*ap).IsNeutral = false
		(*ap).IsGreen = false
		(*ap).IsRed = true
	}
}

type MainPageCollection struct {
	ArticlesFace []ArticlePrinter
	PageNumber   int
}

type MaxNumberOfPages struct {
	NumberOfPages int
	ServerMessage string
}

type PostFilter struct {
	Categories  string
	TitleHints  string
	Keywords    string
	MinLikes    int
	MaxLikes    int
	MinDiss     int
	MaxDiss     int
	HasComments bool
}

type FilteredOutput struct {
	HaveResults     bool
	NoResults       bool
	NumberOfResults int
	Articles        []ArticlePrinter
}

type UserProfilePrinter struct {
	NotLogged                bool
	Logged                   bool
	UserHandler              string
	UserEmail                string
	NumberOfCreatedPosts     int
	NumberOfLikedPosts       int
	NumberOfDislikedPosts    int
	NumberOfCreatedComments  int
	NumberOfLikedComments    int
	NumberOfDislikedComments int
}

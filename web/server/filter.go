package server

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func FilterPosts(filters PostFilter) ([]ArticlePrinter, error) {
	catsFilter := "%"
	titleFilter := "%"
	contentFIlter := "%"
	if len(filters.Categories) > 0 {
		chunks := strings.Split(filters.Categories, ",")
		catsFilter += strings.Join(chunks, "%")
		catsFilter += "%"
	}
	if len(filters.TitleHints) > 0 {
		chunks := strings.Split(filters.TitleHints, " ")
		titleFilter += strings.Join(chunks, "%")
		titleFilter += "%"
	}
	if len(filters.Keywords) > 0 {
		chunks := strings.Split(filters.Keywords, " ")
		contentFIlter += strings.Join(chunks, "%")
		contentFIlter += "%"
	}
	log.Println("DB recieved filters: ", filters)
	ask := "SELECT * FROM posts WHERE categories LIKE ? AND title LIKE ? AND content LIKE ?;"
	rows, filterErr := posts.Query(ask, catsFilter, titleFilter, contentFIlter)
	if filterErr != nil {
		log.Println("[FilterPosts: filterErr] -> " + filterErr.Error())
		return nil, filterErr
	}
	defer rows.Close()
	result := []ArticlePrinter{}
	for rows.Next() {
		var buf ArticlePrinter
		var x int
		var cats string
		scanErr := rows.Scan(&buf.PostId, &buf.UserNick, &buf.Title, &buf.Content, &x, &cats, &buf.Likes, &buf.Dislikes)
		if scanErr != nil {
			log.Println("[FilterPosts: scanErr] -> " + scanErr.Error())
			return nil, scanErr
		}
		buf.Categories = strings.Split(cats, ";")
		if filters.HasComments {
			cnt, countErr := CountCommentsOfThePost(buf.PostId)
			if countErr != nil {
				log.Println("[FilterPosts: countErr] -> " + countErr.Error())
				return nil, countErr
			}
			if cnt == 0 {
				continue
			}
		}
		var IsLikesOk bool = (filters.MinLikes <= buf.Likes && buf.Likes <= filters.MaxLikes)
		var IsDissOk bool = (filters.MinDiss <= buf.Dislikes && buf.Dislikes <= filters.MaxDiss)
		if IsLikesOk && IsDissOk {
			result = append(result, buf)
		}
	}
	return result, nil
}

func FilterPostsApi(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/api/filters/" {
		log.Println("[FilterPostsApi: 404]")
		ErrorHandler(w, http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		log.Println("[FilterPostsApi: 405]")
		ErrorHandler(w, http.StatusMethodNotAllowed)
		return
	}
	var filters PostFilter
	filters.Categories = r.URL.Query().Get("categories")
	filters.TitleHints = r.URL.Query().Get("title")
	filters.Keywords = r.URL.Query().Get("keywords")
	maxl, maxlConvErr := strconv.Atoi(r.URL.Query().Get("maxl"))
	if maxlConvErr != nil {
		log.Println("[FilterPostApi: maxlConvErr] -> " + maxlConvErr.Error())
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}
	filters.MaxLikes = maxl
	minl, minlConvErr := strconv.Atoi(r.URL.Query().Get("minl"))
	if minlConvErr != nil {
		log.Println("[FilterPostsApi: minlConvErr] -> " + minlConvErr.Error())
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}
	filters.MinLikes = minl
	maxd, maxdConvErr := strconv.Atoi(r.URL.Query().Get("maxd"))
	if maxdConvErr != nil {
		log.Println("[FilterPostsApi: maxdConvErr] -> " + maxdConvErr.Error())
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}
	filters.MaxDiss = maxd
	mind, mindConvErr := strconv.Atoi(r.URL.Query().Get("mind"))
	if mindConvErr != nil {
		log.Println("[FilterPostsApi: mindConvErr] -> " + mindConvErr.Error())
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}
	filters.MinDiss = mind
	hasComments := r.URL.Query().Get("commented")
	filters.HasComments = (hasComments == "ok")
	log.Println("Recievied a call to filter some posts: ", filters)
	articles, loadArticlesErr := FilterPosts(filters)
	if loadArticlesErr != nil {
		log.Println("[FilterPostsApi: loadArticlesErr] -> " + loadArticlesErr.Error())
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}
	var output FilteredOutput
	for i := 0; i < len(articles); i++ {
		articles[i].Content = (articles[i].Content[:100] + "....")
	}
	output.Articles = articles
	output.NumberOfResults = len(articles)
	output.HaveResults = (len(articles) > 0)
	output.NoResults = !output.HaveResults
	tmpl, loadTmplErr := template.ParseFiles("./web/UI/templates/filtered.html")
	if loadTmplErr != nil {
		log.Println("[FilterPostsApi: loadTmplErr] -> " + loadTmplErr.Error())
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}
	log.Println("Output to filtering of posts: ", output)
	loadResponseErr := tmpl.Execute(w, output)
	if loadResponseErr != nil {
		log.Println("[FilterPostsApi: loadResponseErr] -> " + loadResponseErr.Error())
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}
}

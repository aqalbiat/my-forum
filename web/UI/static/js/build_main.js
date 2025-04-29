
async function LoadAboutUs() {
    let resp = await fetch("http://localhost:8680/menu/?button=about_us");
    let text = await resp.text();
    if (resp.status >= 300) {
        document.open();
        document.write(text);
        document.close();
        return;
    }
    let neededBox = document.getElementById("loading_box");
    neededBox.innerHTML = text;
}

async function LoadSearchings() {
    let resp = await fetch("http://localhost:8680/menu/?button=search");
    let text = await resp.text();
    if (resp.status >= 300) {
        document.open();
        document.write(text);
        document.close();
        return;
    }
    let neededBox = document.getElementById("loading_box");
    neededBox.innerHTML = text;
}

async function LoadCreatePost() {
    let resp = await fetch("http://localhost:8680/menu/?button=posts");
    let text = await resp.text();
    if (resp.status >= 300) {
        document.open();
        document.write(text);
        document.close();
        return;
    }
    let neededBox = document.getElementById("loading_box");
    neededBox.innerHTML = text;
}

async function LoadCategories() {
    let resp = await fetch("http://localhost:8680/menu/?button=categories&letter=A");
    let text = await resp.text();
    if (resp.status >= 300) {
        document.open();
        document.write(text);
        document.close();
        return;
    }
    let neededBox = document.getElementById("loading_box");
    neededBox.innerHTML = text;
}


function LoadPosts() {
    const buttons = document.getElementsByClassName("button_page");
    for (let i = 0; i < buttons.length; i++) {
        if (buttons[i].getAttribute("listener") === "true") {
            continue;
        }
        let postId = buttons[i].id; 
        console.log(postId);
        buttons[i].setAttribute("listener", "true");
        buttons[i].addEventListener("click", async function() {
            let response = await fetch("http://localhost:8680/api/articles/?action=textpage&id=" + postId);
            let bodyText = await response.text();
            if (response.status >= 300) {
                document.open();
                document.write(bodyText);
                document.close();
                return;
            }
            let neededBox = document.getElementById("loading_box");
            neededBox.innerHTML = bodyText;
        });
    }
}

function SubmitArticle() {
    let ArticleTitle = document.getElementById("title").value;
    if (ArticleTitle.length < 4) {
        alert("Title of the article should be at least 4");
        return;
    }
    let Checkboxes = document.getElementsByClassName("checks");
    console.log("Number of checks: " + Checkboxes.length);
    let Categories = "";
    for (let i = 0; i < Checkboxes.length; i++) {
        if (Checkboxes[i].checked === true) {
            Categories += (Checkboxes[i].name + ";");
        }
    }
    if (Categories.length === 0) {
        alert("You should choose at least one category");
        return;
    }
    let Article = document.getElementById("article").value;
    if (Article.length < 100) {
        alert("One article should have at least 100 characters");
        return;
    }
    fetch("http://localhost:8680/api/articles/?action=create", {
        method: "POST",
        body: JSON.stringify({
            title: ArticleTitle,
            category: Categories,
            article: Article,
        }),
        headers: {
            "Content-Type" : "application/json",
        },
    }).then((resp) => resp.text()).then((txt) => {
        let NeededRectangle = document.getElementById("loading_box");
        NeededRectangle.innerHTML = txt;
    });
}

function GotoLeft() {
    let currentLetter = document.getElementsByClassName("letter")[0].id;
    if (currentLetter == "A") {
        alert("It is already in leftmost page");
        return;
    }
    let nbr = currentLetter.charCodeAt(0);
    nbr--;
    let prevLetter = String.fromCharCode(nbr);
    fetch("http://localhost:8680/menu/?button=categories&letter=" + prevLetter).then((resp) => resp.text()).then((text) => {
        let neededBox = document.getElementById("loading_box");
        neededBox.innerHTML = text;
    });
}

function GotoRight() {
    let currentLetter = document.getElementsByClassName("letter")[0].id;
    if (currentLetter == "Z") {
        alert("It is already in rightmost page");
        return;
    }
    let nbr = currentLetter.charCodeAt(0);
    nbr++;
    let nextLetter = String.fromCharCode(nbr);
    console.log(nextLetter);
    fetch("http://localhost:8680/menu/?button=categories&letter=" + nextLetter).then((resp) => resp.text()).then((text) => {
        let neededBox = document.getElementById("loading_box");
        neededBox.innerHTML = text;
    });
}

function LoadMainPage() {
    fetch("http://localhost:8680/menu/?button=main&pg=1").then((resp) => resp.text()).then((text) => {
        let NeededBox = document.getElementById("loading_box");
        NeededBox.innerHTML = text;
    });
}

function LoadLeftForMainPage() {
    let pg = 0 + 1 * (document.getElementById("page_number").innerHTML);
    if (pg === 1) {
        alert("Already at leftmost page");
        return;
    }
    pg--;
    fetch("http://localhost:8680/menu/?button=main&pg=" + pg).then((resp) => resp.text()).then((text) => {
        let NeededBox = document.getElementById("loading_box");
        NeededBox.innerHTML = text;
    });
}

function LoadRightForMainPage() {
    fetch("http://localhost:8680/api/articles/?action=page_number").then((resp) => resp.text()).then((text) => {
        let ans = JSON.parse(text)
        if (ans.ServerMessage === "OK") {
            let pg = 0 + 1 * (document.getElementById("page_number").innerHTML);
            let last = 0 + 1 * (ans.NumberOfPages);
            if (pg === last) {
                alert("It is already rightmost page!");
            } else {
                pg++;
                fetch("http://localhost:8680/menu/?button=main&pg=" + pg).then((resp) => resp.text()).then((text) => {
                    let NeededBox = document.getElementById("loading_box");
                    NeededBox.innerHTML = text;
                });
            }
        } else {
            alert("Some error occured!");
        }
    });
}

function HistoryBack(){
    let id = 0+1*(document.getElementsByClassName("article_wrapper")[0].id);
    let pg = Math.floor((id+4)/5);
    fetch("http://localhost:8680/menu/?button=main&pg=" + pg).then((resp) => resp.text()).then((text) => {
        let NeededBox = document.getElementById("loading_box");
        NeededBox.innerHTML = text;
    });
}

function LikeButton() {
    let id = 0+1*(document.getElementsByClassName("article_wrapper")[0].id);
    fetch("http://localhost:8680/api/articles/?action=likes&id="+id+"&type=upgrade", {
        method: "POST",
        body: "Hello world!",
        headers: {
            "Content-Type" : "plain/text",
        },
    }).then((resp) => resp.text()).then((text) => {
        console.log(text);
        let ans = JSON.parse(text);
        if (ans.Message === "EXISTS") {
            return;
        }
        if (ans.Message === "NOT REGISTERED") {
            alert("Please sign in first!");
            return;
        }
        if (ans.Message === "Accepted") {
            console.log("liked");
            let box = document.getElementsByClassName("likes")[0];
            box.style = "background-color: green";
            let otherBox=document.getElementsByClassName("dislikes")[0];
            otherBox.style = "background-color: white";
            let likecount = document.getElementById("likescount");
            likecount.innerHTML = 1+1*(likecount.innerHTML);
            box.setAttribute("data-is-checked", "true");
            if (otherBox.getAttribute("data-is-checked") === "true") {
                let dislikescount = document.getElementById("dislikescount");
                dislikescount = 1*(dislikescount.innerHTML)-1;
            }
            otherBox.setAttribute("data-is-checked", "false");
        }
    });
}

function DislikeButton() {
    let id = 0+1*(document.getElementsByClassName("article_wrapper")[0].id);
    fetch("http://localhost:8680/api/articles/?action=likes&id="+id+"&type=downgrade", {
        method: "POST",
        body: "Hello world!",
        headers: {
            "Content-Type" : "plain/text",
        },
    }).then((resp) => resp.text()).then((text) => {
        console.log(text);
        let ans = JSON.parse(text);
        if (ans.Message === "EXISTS") {
            return;
        }
        if (ans.Message === "NOT REGISTERED") {
            alert("Please sign in first!");
            return;
        }
        if (ans.Message === "Accepted") {
            console.log("disliked");
            let box = document.getElementsByClassName("dislikes")[0];
            box.style = "background-color: red";
            let otherBox=document.getElementsByClassName("likes")[0];
            otherBox.style = "background-color: white";
            let dislikecount = document.getElementById("dislikescount");
            dislikecount.innerHTML = 1+1*(dislikecount.innerHTML);
            box.setAttribute("data-is-checked", "true");
            if (otherBox.getAttribute("data-is-checked") === "true") {
                let likescount = document.getElementById("likescount");
                likescount.innerHTML = 1*(likescount.innerHTML)-1;
            }
            otherBox.setAttribute("data-is-checked", "false");
        }
    });
}

async function SendComment() {
    let comment = document.getElementById("user_comment").value;
    let articleId = 0+1*document.getElementsByClassName("article_wrapper")[0].id;
    let response = await fetch("http://localhost:8680/api/articles/?action=comments", {
        method: "POST",
        body: JSON.stringify({
            content: comment,
            postId: articleId,
        }),
        headers: {
            "Content-Type" : "application/json",
        },
    });
    if (response.status >= 300) {
        let result = await response.text();
        document.open();
        document.write(result);
        document.close();
        return;
    }
    let text = await response.text();
    const textBlock = document.getElementById("user_comment");
    textBlock.innerHTML = "";
    let result = JSON.parse(text);
    if (result.Message === "NOT REGISTERED") {
        alert("You are not logged in!");
        return;
    }
    if (result.Message.slice(0, 2) !== "OK") {
        alert("Ooops! Some error occured!");
        return;
    }
    let parts = result.Message.split("#");
    let userName = parts[1];
    let commentId = parts[2];
    let commentBox = document.createElement("div");
    commentBox.classList.add("comment_printer");
    commentBox.id = commentId;
    let commentAuthour = document.createElement("div");
    commentAuthour.classList.add("comment_author_block");
    commentAuthour.innerHTML = userName;
    commentBox.appendChild(commentAuthour);
    let commentText = document.createElement("div");
    commentText.classList.add("comment_text_block");
    commentText.innerHTML = comment;
    commentBox.appendChild(commentText);
    let feedbackBox = document.createElement("div");
    feedbackBox.classList.add("likes_n_dislikes_container");
    let countLikes = document.createElement("div");
    countLikes.classList.add("counting_numbers");
    countLikes.innerHTML = "0";
    feedbackBox.appendChild(countLikes);
    let likeButton = document.createElement("button");
    likeButton.classList.add("likes");
    likeButton.setAttribute("touched", "false");
    likeButton.innerHTML = "&#128077;";
    feedbackBox.appendChild(likeButton);
    let countDislikes = document.createElement("div");
    countDislikes.classList.add("counting_numbers");
    countDislikes.innerHTML = "0";
    feedbackBox.appendChild(countDislikes);
    let dislikeButton = document.createElement("button");
    dislikeButton.classList.add("dislikes");
    dislikeButton.setAttribute("touched", "false");
    dislikeButton.innerHTML = "&#128078;";
    feedbackBox.appendChild(dislikeButton);
    commentBox.appendChild(feedbackBox);
    let commentsList = document.getElementById("comments_list");
    commentsList.append(commentBox);
}

function AddCommentFeedbackListeners() {
    let commentListBox = document.getElementById("comments_list").children;
    for (let i = 1; i < commentListBox.length; i++) {
        if (commentListBox[i].hasAttribute("listener")) {
            continue;
        }
        let feedbackBox = commentListBox[i].children[2];
        let likeCount = feedbackBox.children[0];
        let likeButton = feedbackBox.children[1];
        let dislikeCount = feedbackBox.children[2];
        let dislikeButton = feedbackBox.children[3];
        const commentId = commentListBox[i].id;
        likeButton.addEventListener("click", async function(e) {
            let resp = await fetch("http://localhost:8680/api/articles/?action=feedback&id=" + commentId + "&grade=upgrade", {
                method: "POST",
                body: JSON.stringify({
                    "request": "need-for-feedback"
                }),
                headers: {
                    "Content-Type": "application/json"
                },
            });
            if (resp.status >= 300) {
                let text = await resp.text();
                document.open();
                document.write(text);
                document.close();
                return;
            }
            let text = await resp.text();
            let msg = JSON.parse(text);
            if (msg.Message === "EXISTS") {
                return;
            }
            if (msg.Message === "NOT REGISTERED") {
                alert("Login first, please!");
                return;
            }
            if (msg.Message !== "OK") {
                alert("Some error occured!");
                return;
            }
            likeButton.style = "background-color: green;";
            likeButton.setAttribute("touched", "true");
            likeCount.innerHTML = (0+1*likeCount.innerHTML) + 1;
            dislikeButton.style = "background-color: white;";
            if (dislikeButton.getAttribute("touched") === "true") {
                dislikeCount.innerHTML = (0+1*dislikeCount.innerHTML) - 1;
            }
            dislikeButton.setAttribute("touched", "false");
        });
        dislikeButton.addEventListener("click", async function (e) {
            let resp = await fetch("http://localhost:8680/api/articles/?action=feedback&id=" + commentId + "&grade=downgrade", {
                method: "POST",
                body: JSON.stringify({
                    "request": "need-for-feedback"
                }),
                headers: {
                    "Content-Type": "application/json"
                },
            });
            if (resp.status >= 300) {
                let text = await resp.text();
                document.open();
                document.write(text);
                document.close();
                return;
            }
            let text = await resp.text();
            let msg = JSON.parse(text);
            if (msg.Message === "EXISTS") {
                return;
            }
            if (msg.Message === "NOT REGISTERED") {
                alert("Login first, please!");
                return;
            }
            if (msg.Message !== "OK") {
                alert("Some error occured!");
                return;
            }
            dislikeButton.style = "background-color: red;";
            dislikeButton.setAttribute("touched", "true");
            dislikeCount.innerHTML = (0+1*dislikeCount.innerHTML) + 1;
            likeButton.style = "background-color: white;";
            if (likeButton.getAttribute("touched") === "true") {
                likeCount.innerHTML = (0+1*likeCount.innerHTML) - 1;
            }
            likeButton.setAttribute("touched", "false");
        });
        commentListBox[i].setAttribute("listener", "true");
    }
}

function UploadReverseFilterButtons() {
    let buttons = document.getElementsByClassName("button_page");
    for (let i = 0; i < buttons.length; i++) {
        if (buttons[i].hasAttribute("listener")) {
            continue;
        }
        buttons[i].setAttribute("listener", "true");
        let articleId = buttons[i].id;
        buttons[i].addEventListener("click", async function() {
            let response = await fetch("http://localhost:8680/api/articles/?action=textpage&id=" + articleId);
            let txt = await response.text();
            if (response.status >= 300) {
                document.open();
                document.write(txt);
                document.close();
                return;
            }
            let neededBox = document.getElementById("loading_box");
            neededBox.innerHTML = txt;
            let backButton = document.getElementsByClassName("back_button")[0];
            backButton.removeAttribute("onclick");
            backButton.addEventListener("click", async function() {
                let shell = document.getElementById("shell");
                let query = shell.getAttribute("data-prev-query");
                console.log("Query retrieved form shell " + query);
                let response = await fetch("http://localhost:8680/api/filters/?" + query);
                let txt = await response.text();
                if (response.status >= 300) {
                    document.open();
                    document.write(txt);
                    document.close();
                    return;
                }
                let neededBox = document.getElementById("loading_box");
                neededBox.innerHTML = txt;
            });
        });
    }
}

function CreateCategoryLoaders() {
    let buttons = document.getElementsByClassName("class_category_names");
    const MAX_VAL = 1000*1000*1000+167;
    for (let i = 0; i < buttons.length; i++) {
        if (buttons[i].hasAttribute("listener")) {
            continue;
        }
        let category = buttons[i].innerHTML;
        let query = "categories=" + category;
        query += "&title=%25&keywords=%25";
        query += "&maxl=" + MAX_VAL;
        query += "&minl=0";
        query += "&maxd=" + MAX_VAL;
        query += "&mind=0";
        query += "&commented=no";
        buttons[i].setAttribute("listener", "true");
        buttons[i].addEventListener("click", async function() {
            let response = await fetch("http://localhost:8680/api/filters/?" + query);
            let txt = await response.text();
            if (response.status >= 300) {
                document.open();
                document.write(txt);
                document.close();
                return;
            }
            let neededBox = document.getElementById("loading_box");
            neededBox.innerHTML = txt;
            console.log("Settled query values: " + query);
            let shell = document.getElementById("shell");
            shell.setAttribute("data-prev-query", query);
        });
    }
}

async function LoadSearchedPosts() {
    let Checkboxes = document.getElementsByClassName("checks");
    let categories = "";
    for (let i = 0; i < Checkboxes.length; i++) {
        if (Checkboxes[i].checked === true) {
            categories += (Checkboxes[i].name + ",");
        }
    }
    if (categories.length === 0) {
        alert("Must have at least one category in filter!");
        return;
    }
    let query = ("categories=" + categories);
    let titleKeywords = document.getElementById("title-keywords").value;
    let contentKeywords = document.getElementById("content-keywords").value;
    let maxLikes = document.getElementById("max-likes").value;
    let minLikes = document.getElementById("min-likes").value;
    let maxDisses = document.getElementById("max-disses").value;
    let minDisses = document.getElementById("min-disses").value;
    let hasComments = document.getElementById("has-comments").value;
    if (titleKeywords.length > 0) {
        query += ("&title=" + titleKeywords);
    } else {
        query += "&title=%25";
    }
    if (contentKeywords.length > 0) {
        query += ("&keywords=" + contentKeywords);
    } else {
        query += "&keywords=%25";
    }
    query += ("&maxl=" + maxLikes);
    query += ("&minl=" + minLikes);
    query += ("&maxd=" + maxDisses);
    query += ("&mind=" + minDisses);
    if (hasComments.checked === true) {
        query += "&commented=ok";
    } else {
        query += "&commented=no";
    }
    let response = await fetch("http://localhost:8680/api/filters/?" + query);
    let txt = await response.text();
    if (response.status >= 300) {
        document.open();
        document.write(txt);
        document.close();
        return;
    }
    let neededBox = document.getElementById("loading_box");
    neededBox.innerHTML = txt;
    let shell = document.getElementById("shell");
    console.log("User's query for search: " + query);
    shell.setAttribute("data-prev-query", query);
}

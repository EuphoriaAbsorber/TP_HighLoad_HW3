package Config

var Port = ":5000"
var BasePath = "/api"
var PathDocs = BasePath + "/docs"
var PathCreateUser = BasePath + "/user/{nickname}/create"
var PathProfile = BasePath + "/user/{nickname}/profile"

var PathCreateForum = BasePath + "/forum/create"
var PathForumInfo = BasePath + "/forum/{slug}/details"
var PathCreateThread = BasePath + "/forum/{slug}/create"
var PathGetForumUsers = BasePath + "/forum/{slug}/users"
var PathGetForumThreads = BasePath + "/forum/{slug}/threads"

var PathGetServiceStatus = BasePath + "/service/status"
var PathServiceClear = BasePath + "/service/clear"

var PathCreatePosts = BasePath + "/thread/{slug_or_id}/create"
var PathThreadInfo = BasePath + "/thread/{slug_or_id}/details"
var PathThreadVote = BasePath + "/thread/{slug_or_id}/vote"
var PathGetThreadPosts = BasePath + "/thread/{slug_or_id}/posts"

var PathPost = BasePath + "/post/{id}/details"

var Headers = map[string]string{
	//"Access-Control-Allow-Credentials": "true",
	//"Access-Control-Allow-Headers":     "Origin, Content-Type, accept",
	//"Access-Control-Allow-Methods":     "GET, POST, DELETE, OPTIONS",
	"Content-Type": "application/json; charset=utf-8",
}

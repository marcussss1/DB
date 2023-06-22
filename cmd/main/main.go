package main

import (
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"net/http"

	handlForum "project/internal/forum/delivery/http"
	handlPost "project/internal/post/delivery/http"
	handlService "project/internal/service/delivery/http"
	handlThread "project/internal/thread/delivery/http"
	handlUser "project/internal/user/delivery/http"
	handlVote "project/internal/vote/delivery/http"

	usecaseForum "project/internal/forum/usecase"
	usecasePost "project/internal/post/usecase"
	usecaseSerivce "project/internal/service/usecase"
	usecaseThread "project/internal/thread/usecase"
	usecaseUser "project/internal/user/usecase"
	usecaseVote "project/internal/vote/usecase"

	repoForum "project/internal/forum/repository"
	repoPost "project/internal/post/repository"
	repoService "project/internal/service/repository"
	repoThread "project/internal/thread/repository"
	repoUser "project/internal/user/repository"
	repoVote "project/internal/vote/repository"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"project/internal/pkg"
)

func main() {
	logger := logrus.Logger{}

	router := mux.NewRouter()

	db, errr := sqlx.Open("pgx", "user=brabra password=brabra dbname=brabra host=localhost port=5432 sslmode=disable")
	if errr != nil {
		logrus.Fatal(errr)
	}

	forumStorage := repoForum.NewForumPostgres(db)
	userStorage := repoUser.NewUserPostgres(db)
	postStorage := repoPost.NewPostPostgres(db)
	threadStorage := repoThread.NewThreadPostgres(db)
	voteStorage := repoVote.NewVotePostgres(db)
	serviceStorage := repoService.NewServicePostgres(db)

	forumService := usecaseForum.NewForumService(forumStorage, userStorage)
	userService := usecaseUser.NewUserService(userStorage)
	postService := usecasePost.NewPostService(postStorage)
	threadService := usecaseThread.NewThreadService(threadStorage, forumStorage, userStorage, postStorage)
	voteService := usecaseVote.NewVoteService(voteStorage, threadStorage, userStorage)
	serivceService := usecaseSerivce.NewService(serviceStorage)

	forumHandler := handlForum.NewForumHandler(forumService, router)
	router.HandleFunc("/api/forum/create", forumHandler.CreateForumHandler).Methods(http.MethodPost)
	router.HandleFunc("/api/forum/{slug}/details", forumHandler.GetForumHandler).Methods(http.MethodGet)
	router.HandleFunc("/api/forum/{slug}/threads", forumHandler.GetForumThreads).Methods(http.MethodGet)
	router.HandleFunc("/api/forum/{slug}/users", forumHandler.GetForumUsersHandler).Methods(http.MethodGet)

	postHandler := handlPost.NewPostHandler(postService, router)
	router.HandleFunc("/api/post/{id}/details", postHandler.GetPostHandler).Methods(http.MethodGet)
	router.HandleFunc("/api/post/{id}/details", postHandler.UpdatePostHandler).Methods(http.MethodPost)

	serviceHandler := handlService.NewServiceHandler(serivceService, router)
	router.HandleFunc("/api/service/clear", serviceHandler.ServiceClearHandler).Methods(http.MethodPost)
	router.HandleFunc("/api/service/status", serviceHandler.ServiceStatusHandler).Methods(http.MethodGet)

	threadHandler := handlThread.NewThreadHandler(threadService, router)
	router.HandleFunc("/api/thread/{slug_or_id}/create", threadHandler.CreatePostsHandler).Methods(http.MethodPost)
	router.HandleFunc("/api/forum/{slug}/create", threadHandler.CreateThreadHandler).Methods(http.MethodPost)
	router.HandleFunc("/api/thread/{slug_or_id}/details", threadHandler.GetThreadHandler).Methods(http.MethodGet)
	router.HandleFunc("/api/thread/{slug_or_id}/posts", threadHandler.GetPostsHandler).Methods(http.MethodGet)
	router.HandleFunc("/api/thread/{slug_or_id}/details", threadHandler.UpdateThreadHandler).Methods(http.MethodPost)

	voteHandler := handlVote.NewVoteHandler(voteService, router)
	router.HandleFunc("/api/thread/{slug_or_id}/vote", voteHandler.VoteHandler).Methods(http.MethodPost)

	userHandler := handlUser.NewUserHandler(userService, router)
	router.HandleFunc("/api/user/{nickname}/create", userHandler.CreateUserHandler).Methods(http.MethodPost)
	router.HandleFunc("/api/user/{nickname}/profile", userHandler.GetProfileHandler).Methods(http.MethodGet)
	router.HandleFunc("/api/user/{nickname}/profile", userHandler.UpdateProfileHandler).Methods(http.MethodPost)

	logrus.Info("server started :5000")

	server := pkg.NewServerHTTP(&logger)

	err := server.Launch(router)
	if err != nil {
		logrus.Fatal(err)
	}
}

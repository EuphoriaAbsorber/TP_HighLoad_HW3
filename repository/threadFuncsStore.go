package repository

import (
	"dbproject/model"
	"fmt"
	"time"

	"github.com/jackc/pgx"
)

func (s *Store) CheckAllPostParentIds(threadId int, in []int) error {
	dbcount := 0
	err := s.db.QueryRow(`SELECT count(id) FROM (SELECT id FROM posts WHERE thread = $1 GROUP BY id HAVING $2 @> array_agg(id)) AS S;`, threadId, in).Scan(&dbcount)
	if err != nil {
		return err
	}
	if dbcount < len(in) {
		return model.ErrConflict409
	}
	return nil
}

func (s *Store) GetThreadById(id int) (*model.Thread, error) {
	rows, err := s.db.Query(`SELECT id, title, author, forum, message, votes, slug, created FROM threads WHERE id = $1;`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	dat := model.Thread{}
	for rows.Next() {
		err := rows.Scan(&dat.Id, &dat.Title, &dat.Author, &dat.Forum, &dat.Message, &dat.Votes, &dat.Slug, &dat.Created)
		if err != nil {
			return nil, err
		}
	}
	if dat.Id != 0 {
		return &dat, nil
	}
	return nil, model.ErrNotFound404
}

func (s *Store) GetThreadBySlug(slug string) (*model.Thread, error) {
	if slug == "" {
		return nil, nil
	}
	rows, err := s.db.Query(`SELECT id, title, author, forum, message, votes, slug, created FROM threads WHERE slug = $1;`, slug)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	dat := model.Thread{}
	for rows.Next() {
		err := rows.Scan(&dat.Id, &dat.Title, &dat.Author, &dat.Forum, &dat.Message, &dat.Votes, &dat.Slug, &dat.Created)
		if err != nil {
			return nil, err
		}
	}
	if dat.Id != 0 {
		return &dat, nil
	}
	return nil, model.ErrNotFound404
}

func (s *Store) InsertNPostsDB(in *model.Posts, position int, Ncount int, createTime time.Time, threadId int, forumSlug string) error {
	query := "INSERT INTO posts (parent, author, message, forum, isedited, thread, created) VALUES "
	args := make([]interface{}, 0)
	j := 0
	for i := position; i < position+Ncount; i++ {
		user, err := s.GetProfile((*in)[i].Author)
		if err != nil {
			return err
		}
		(*in)[i].Author = user.Nickname
		(*in)[i].Forum = forumSlug
		(*in)[i].Thread = threadId
		(*in)[i].Created = createTime.Format(time.RFC3339)
		query += fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d),", j*7+1, j*7+2, j*7+3, j*7+4, j*7+5, j*7+6, j*7+7)
		if (*in)[i].Parent != 0 {
			args = append(args, (*in)[i].Parent, (*in)[i].Author, (*in)[i].Message, (*in)[i].Forum, (*in)[i].IsEdited, (*in)[i].Thread, (*in)[i].Created)
		} else {
			args = append(args, nil, (*in)[i].Author, (*in)[i].Message, (*in)[i].Forum, (*in)[i].IsEdited, (*in)[i].Thread, (*in)[i].Created)
		}
		j++
	}

	query = query[:len(query)-1]
	query += " RETURNING id;"

	resultRows, err := s.db.Query(query, args...)
	if err != nil {
		return model.ErrConflict409
	}
	defer resultRows.Close()
	for i := position; resultRows.Next(); i++ {
		var id int
		if err = resultRows.Scan(&id); err != nil {
			return err
		}
		(*in)[i].Id = id
	}
	return nil
}

func (s *Store) CreatePosts(in *model.Posts, threadId int, forumSlug string) (*model.Posts, error) {
	var err error
	createTime := time.Now()
	postsForOneInsert := 20
	parts := len(*in) / postsForOneInsert
	for i := 0; i < parts+1; i++ {
		if i == parts {
			if i*postsForOneInsert != len(*in) {
				err = s.InsertNPostsDB(in, i*postsForOneInsert, len(*in)-i*postsForOneInsert, createTime, threadId, forumSlug)
				if err != nil {
					return nil, err
				}
			}
		} else {
			err = s.InsertNPostsDB(in, i*postsForOneInsert, postsForOneInsert, createTime, threadId, forumSlug)
			if err != nil {
				return nil, err
			}
		}
	}
	_, err = s.db.Exec(`UPDATE forums SET posts = posts + $1 WHERE slug = $2;`, len(*in), forumSlug)
	if err != nil {
		return nil, err
	}
	return in, nil
}

// func (s *Store) CreatePosts(in *model.Posts, threadId int, forumSlug string) ([]*model.Post, error) {

// 	posts := []*model.Post{}
// 	createTime := time.Now()
// 	createdFormatted := createTime.Format(time.RFC3339)
// 	var err error
// 	for _, post := range *in {
// 		id := 0
// 		insertModel := model.Post{Parent: post.Parent, Author: post.Author, Message: post.Message, IsEdited: false, Thread: threadId, Forum: forumSlug, Created: createdFormatted}
// 		if post.Parent == 0 {
// 			err = s.db.QueryRow(`INSERT INTO posts (parent, author, message, forum, thread, isedited, created) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id;`, nil, insertModel.Author, insertModel.Message, insertModel.Forum, insertModel.Thread, insertModel.IsEdited, createdFormatted).Scan(&id)
// 		} else {
// 			err = s.db.QueryRow(`INSERT INTO posts (parent, author, message, forum, thread, isedited, created) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id;`, insertModel.Parent, insertModel.Author, insertModel.Message, insertModel.Forum, insertModel.Thread, insertModel.IsEdited, createdFormatted).Scan(&id)
// 		}
// 		if err != nil {
// 			return nil, model.ErrConflict409
// 		}
// 		insertModel.Id = id
// 		insertModel.Created = createdFormatted
// 		posts = append(posts, &insertModel)
// 	}
// 	_, err = s.db.Exec(`UPDATE forums SET posts = posts + $1 WHERE slug = $2;`, len(*in), forumSlug)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return posts, nil
// }

func (s *Store) UpdateThreadInfo(in *model.ThreadUpdate, id int) error {
	_, err := s.db.Exec(`UPDATE threads SET message = $1, title = $2 WHERE id = $3;`, in.Message, in.Title, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) VoteForThread(in *model.Vote, threadID int, threadVotes int) (int, error) {
	oldVote := 0
	rows, err := s.db.Query(`SELECT voice FROM votes WHERE thread = $1 AND nickname = $2;`, threadID, in.Nickname)
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&oldVote)
		if err != nil {
			return 0, err
		}
	}
	newVote := 0
	err = s.db.QueryRow(`INSERT INTO votes (nickname, thread, voice) VALUES ($1, $2, $3) ON CONFLICT (nickname, thread) DO UPDATE SET voice = EXCLUDED.voice RETURNING voice;`, in.Nickname, threadID, in.Voice).Scan(&newVote)
	if err != nil {
		return 0, err
	}
	if oldVote != newVote {
		err = s.db.QueryRow(`UPDATE threads SET votes = votes - $1 + $2 WHERE id = $3 RETURNING votes;`, oldVote, newVote, threadID).Scan(&threadVotes)
		if err != nil {
			return 0, err
		}
	}
	return threadVotes, nil
}

func (s *Store) GetThreadPostsFlatSort(threadId int, limit int, since int, desc bool) ([]*model.Post, error) {
	posts := []*model.Post{}
	var rows *pgx.Rows
	var err error
	if !desc {
		rows, err = s.db.Query(`SELECT id, parent, author, message, forum, thread, isedited, created FROM posts WHERE thread = $1 AND id > $2 ORDER BY (created, id) LIMIT $3;`, threadId, since, limit)
	}
	if desc {
		if since == 0 {
			since = 1e9
		}
		rows, err = s.db.Query(`SELECT id, parent, author, message, forum, thread, isedited, created FROM posts WHERE thread = $1 AND id < $2 ORDER BY (created, id) DESC LIMIT $3;`, threadId, since, limit)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		dat := model.Post{}
		date := time.Now()
		var parent *int
		err := rows.Scan(&dat.Id, &parent, &dat.Author, &dat.Message, &dat.Forum, &dat.Thread, &dat.IsEdited, &date)
		dat.Created = date.Format(time.RFC3339)
		if parent != nil {
			dat.Parent = *parent
		}
		if err != nil {
			return nil, err
		}
		posts = append(posts, &dat)
	}
	return posts, nil
}

func (s *Store) GetThreadPostsTreeSort(threadId int, limit int, since int, desc bool) ([]*model.Post, error) {
	posts := []*model.Post{}
	var rows *pgx.Rows
	var err error

	if !desc {
		if since == 0 {
			rows, err = s.db.Query(`SELECT id, parent, author, message, forum, thread, isedited, created FROM posts WHERE thread = $1 ORDER BY path LIMIT $2;`, threadId, limit)
		} else {
			rows, err = s.db.Query(`SELECT id, parent, author, message, forum, thread, isedited, created FROM posts WHERE thread = $1 AND path > (SELECT path FROM posts WHERE id = $2) ORDER BY path LIMIT $3;`, threadId, since, limit)
		}
	}
	if desc {
		if since == 0 {
			rows, err = s.db.Query(`SELECT id, parent, author, message, forum, thread, isedited, created FROM posts WHERE thread = $1 ORDER BY path DESC LIMIT $2;`, threadId, limit)
		} else {
			rows, err = s.db.Query(`SELECT id, parent, author, message, forum, thread, isedited, created FROM posts WHERE thread = $1 AND path < (SELECT path FROM posts WHERE id = $2) ORDER BY path DESC LIMIT $3;`, threadId, since, limit)
		}
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		dat := model.Post{}
		date := time.Now()
		var parent *int
		err = rows.Scan(&dat.Id, &parent, &dat.Author, &dat.Message, &dat.Forum, &dat.Thread, &dat.IsEdited, &date)
		dat.Created = date.Format(time.RFC3339)
		if parent != nil {
			dat.Parent = *parent
		}
		posts = append(posts, &dat)
	}
	return posts, nil
}

func (s *Store) GetThreadPostsTreeParentSort(threadId int, limit int, since int, desc bool) ([]*model.Post, error) {
	posts := []*model.Post{}
	var rows *pgx.Rows
	var err error
	if !desc {
		if since == 0 {
			rows, err = s.db.Query(`SELECT id, parent, author, message, forum, thread, isedited, created FROM posts WHERE path[1] IN (SELECT id FROM posts WHERE thread = $1 AND parent IS NULL ORDER BY id LIMIT $2) ORDER BY path;`, threadId, limit)
		} else {
			rows, err = s.db.Query(`SELECT id, parent, author, message, forum, thread, isedited, created FROM posts WHERE path[1] IN (SELECT id FROM posts WHERE thread = $1 AND parent IS NULL AND path[1] > (SELECT path[1] FROM posts WHERE id = $2) ORDER BY id LIMIT $3) ORDER BY path;`, threadId, since, limit)
		}
	}
	if desc {
		if since == 0 {
			rows, err = s.db.Query(`SELECT id, parent, author, message, forum, thread, isedited, created FROM posts WHERE path[1] IN (SELECT id FROM posts WHERE thread = $1 AND parent IS NULL ORDER BY id DESC LIMIT $2) ORDER BY path[1] DESC, path ASC, id ASC;`, threadId, limit)
		} else {
			rows, err = s.db.Query(`SELECT id, parent, author, message, forum, thread, isedited, created FROM posts WHERE path[1] IN (SELECT id FROM posts WHERE thread = $1 AND parent IS NULL AND path[1] < (SELECT path[1] FROM posts WHERE id = $2) ORDER BY id DESC LIMIT $3) ORDER BY path[1] DESC, path ASC, id ASC;`, threadId, since, limit)
		}
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		dat := model.Post{}
		date := time.Now()
		var parent *int
		err = rows.Scan(&dat.Id, &parent, &dat.Author, &dat.Message, &dat.Forum, &dat.Thread, &dat.IsEdited, &date)
		dat.Created = date.Format(time.RFC3339)
		if parent != nil {
			dat.Parent = *parent
		}
		posts = append(posts, &dat)
	}
	return posts, nil
}

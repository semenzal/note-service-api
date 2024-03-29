package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/semenzal/note-service-api/internal/model"
	"github.com/semenzal/note-service-api/internal/pkg/db"
	def "github.com/semenzal/note-service-api/internal/repository"
	desc "github.com/semenzal/note-service-api/pkg/note_v1"

	sq "github.com/Masterminds/squirrel"
)

const (
	tableName         = "note"
	idColumn          = "id"
	titleColumn       = "title"
	textColumn        = "text"
	authorColumn      = "author"
	emailColumn       = "email"
	created_at_Column = "created_at"
	updated_at_Column = "updated_at"
)

type repository struct {
	client db.Client
}

func NewRepository(client db.Client) def.NoteRepository {
	return &repository{
		client: client,
	}
}

func (r *repository) Create(ctx context.Context, noteInfo *model.NoteInfo) (int64, error) {
	builder := sq.Insert(tableName).
		PlaceholderFormat(sq.Dollar).
		Columns(titleColumn, textColumn, authorColumn, emailColumn).
		Values(noteInfo.Title, noteInfo.Text, noteInfo.Author, noteInfo.Email).
		Suffix("returning id")

	query, args, err := builder.ToSql()
	if err != nil {
		return 0, err
	}

	q := db.Query{
		Name:     "Create",
		QueryRaw: query,
	}

	row, err := r.client.DB().QueryContext(ctx, q, args...)
	if err != nil {
		return 0, err
	}
	defer row.Close()

	row.Next()
	var id int64
	err = row.Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *repository) Get(ctx context.Context, id int64) (*model.Note, error) {
	builder := sq.Select(idColumn, titleColumn, textColumn, authorColumn, created_at_Column, updated_at_Column, emailColumn).
		From(tableName).
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).
		Limit(1)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	q := db.Query{
		Name:     "Get",
		QueryRaw: query,
	}

	note := new(model.Note)
	err = r.client.DB().GetContext(ctx, note, q, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("note not found")
		}
		return nil, err
	}

	return note, nil
}

func (r *repository) GetList(ctx context.Context, filter *model.Filter) ([]*model.Note, error) {
	builder := sq.Select(idColumn, titleColumn, textColumn, authorColumn, created_at_Column, updated_at_Column, emailColumn).
		From(tableName).
		PlaceholderFormat(sq.Dollar)

	if filter.Title.Valid {
		builder = builder.Where(sq.Eq{titleColumn: filter.Title.String})
	}

	if filter.Text.Valid {
		builder = builder.Where(sq.Eq{textColumn: filter.Text.String})
	}

	if filter.Author.Valid {
		builder = builder.Where(sq.Eq{authorColumn: filter.Author.String})
	}

	if filter.Email.Valid {
		builder = builder.Where(sq.Eq{emailColumn: filter.Email.String})
	}

	if filter.Limit.Valid {
		builder = builder.Limit(uint64(filter.Limit.Int64))
	}

	if filter.Offset.Valid {
		builder = builder.Offset(uint64(filter.Offset.Int64))
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	q := db.Query{
		Name:     "GetList",
		QueryRaw: query,
	}

	var notes []*model.Note
	err = r.client.DB().SelectContext(ctx, &notes, q, args...)
	if err != nil {
		return nil, err
	}

	return notes, nil
}

func (r *repository) Update(ctx context.Context, id int64, updateInfo *model.UpdateNoteInfo) error {
	builder := sq.Update(tableName).
		PlaceholderFormat(sq.Dollar).
		Set(updated_at_Column, time.Now()).
		Where(sq.Eq{"id": id})

	if updateInfo.Title.Valid {
		builder = builder.Set(titleColumn, updateInfo.Title)
	}

	if updateInfo.Text.Valid {
		builder = builder.Set(textColumn, updateInfo.Text)
	}

	if updateInfo.Author.Valid {
		builder = builder.Set(authorColumn, updateInfo.Author)
	}

	if updateInfo.Email.Valid {
		builder = builder.Set(emailColumn, updateInfo.Email)
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	q := db.Query{
		Name:     "Update",
		QueryRaw: query,
	}

	_, err = r.client.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) Delete(ctx context.Context, req *desc.DeleteRequest) error {
	builder := sq.Delete(tableName).
		Where(sq.Eq{idColumn: req.Id}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	q := db.Query{
		Name:     "Delete",
		QueryRaw: query,
	}
	_, err = r.client.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return err
	}

	return nil
}

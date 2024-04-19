package postgresql

import (
	"context"

	"evrone_api_gateway/internal/entity"
	"evrone_api_gateway/internal/infrastructure/repository/postgresql/repo"
	"evrone_api_gateway/internal/pkg/otlp"
	"evrone_api_gateway/internal/pkg/postgres"
)

type appVersionRepo struct {
	tableName string
	db        *postgres.PostgresDB
}

func NewAppVersionRepo(db *postgres.PostgresDB) repo.AppVersionRepo {
	return &appVersionRepo{
		tableName: "app_version",
		db:        db,
	}
}

func (r *appVersionRepo) Get(ctx context.Context) (*entity.AppVersion, error) {
	// tracing
	ctx, span := otlp.Start(ctx, "appVersionService", "appVersionRepoGet")
	defer span.End()

	query := r.db.Sq.Builder.
		Select(
			"android_version",
			"ios_version",
			"is_force_update",
		).
		From(r.tableName)

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return nil, r.db.ErrSQLBuild(err, r.tableName+" read")
	}

	var res entity.AppVersion
	err = r.db.QueryRow(ctx, sqlStr, args...).Scan(
		&res.AndroidVersion,
		&res.IOSVersion,
		&res.IsForceUpdate,
	)
	if err != nil {
		return nil, r.db.Error(err)
	}

	return &res, nil
}

func (r *appVersionRepo) Create(ctx context.Context, m *entity.AppVersion) error {
	// tracing
	ctx, span := otlp.Start(ctx, "appVersionService", "appVersionRepoCreate")
	defer span.End()

	clauses := map[string]interface{}{
		"android_version": m.AndroidVersion,
		"ios_version":     m.IOSVersion,
		"is_force_update": m.IsForceUpdate,
		"created_at":      m.CreatedAt,
		"updated_at":      m.UpdatedAt,
	}

	sqlStr, args, err := r.db.Sq.Builder.Insert(r.tableName).SetMap(clauses).ToSql()
	if err != nil {
		return r.db.ErrSQLBuild(err, r.tableName+" create")
	}

	if _, err = r.db.Exec(ctx, sqlStr, args...); err != nil {
		return r.db.Error(err)
	}
	return nil
}

func (r *appVersionRepo) Update(ctx context.Context, m *entity.AppVersion) error {
	// tracing
	ctx, span := otlp.Start(ctx, "appVersionService", "appVersionRepoUpdate")
	defer span.End()

	clauses := map[string]interface{}{
		"android_version": m.AndroidVersion,
		"ios_version":     m.IOSVersion,
		"is_force_update": m.IsForceUpdate,
		"updated_at":      m.UpdatedAt,
	}

	sqlStr, args, err := r.db.Sq.Builder.Update(r.tableName).SetMap(clauses).ToSql()
	if err != nil {
		return r.db.ErrSQLBuild(err, r.tableName+" update")
	}

	if _, err = r.db.Exec(ctx, sqlStr, args...); err != nil {
		return r.db.Error(err)
	}
	return nil
}

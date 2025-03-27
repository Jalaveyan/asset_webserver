package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"go-asset-service/internal/models"
)

// AssetRepository отвечает за выполнение операций с таблицей assets в базе данных.
type AssetRepository struct {
	db *pgxpool.Pool // Пул соединений с базой данных
}

// NewAssetRepository создаёт новый экземпляр AssetRepository.
func NewAssetRepository(db *pgxpool.Pool) *AssetRepository {
	return &AssetRepository{db: db}
}

// CreateAsset сохраняет новый asset (файл/данные) в базе данных.
func (r *AssetRepository) CreateAsset(ctx context.Context, asset *models.Asset) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO assets (name, uid, data, created_at)
		 VALUES ($1, $2, $3, $4)`,
		asset.Name, asset.UID, asset.Data, asset.CreatedAt,
	)
	return err
}

// GetAsset извлекает asset по имени и идентификатору пользователя.
func (r *AssetRepository) GetAsset(ctx context.Context, name string, uid int64) (*models.Asset, error) {
	row := r.db.QueryRow(ctx,
		`SELECT name, uid, data, created_at
		 FROM assets
		 WHERE name = $1 AND uid = $2`,
		name, uid,
	)
	var a models.Asset
	err := row.Scan(&a.Name, &a.UID, &a.Data, &a.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

// ListAssets возвращает список всех assets (файлов), загруженных пользователем с заданным uid.
func (r *AssetRepository) ListAssets(ctx context.Context, uid int64) ([]models.Asset, error) {
	rows, err := r.db.Query(ctx,
		`SELECT name, uid, created_at FROM assets WHERE uid = $1`,
		uid,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var assets []models.Asset
	for rows.Next() {
		var a models.Asset
		err := rows.Scan(&a.Name, &a.UID, &a.CreatedAt)
		if err != nil {
			return nil, err
		}
		assets = append(assets, a)
	}
	return assets, nil
}

// DeleteAsset удаляет asset с указанным именем и идентификатором пользователя из базы данных.
func (r *AssetRepository) DeleteAsset(ctx context.Context, name string, uid int64) error {
	_, err := r.db.Exec(ctx, `DELETE FROM assets WHERE name = $1 AND uid = $2`, name, uid)
	return err
}

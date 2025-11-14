package repositories

import (
	"database/sql"
	"fmt"

	"github.com/stiffis/UniCLI/internal/models"
)

// CategoryRepository handles category data operations
type CategoryRepository struct {
	*BaseRepository
}

// NewCategoryRepository creates a new category repository
func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

// FindAll retrieves all categories from the database
func (r *CategoryRepository) FindAll() ([]models.Category, error) {
	query := `SELECT id, name, color FROM categories ORDER BY name ASC`

	rows, err := r.DB().Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query categories: %w", err)
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var category models.Category
		if err := rows.Scan(&category.ID, &category.Name, &category.Color); err != nil {
			return nil, fmt.Errorf("failed to scan category: %w", err)
		}
		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating categories: %w", err)
	}

	return categories, nil
}

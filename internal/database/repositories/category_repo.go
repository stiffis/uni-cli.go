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

func (r *CategoryRepository) Create(category *models.Category) error {
	query := `INSERT INTO categories (id, name, color) VALUES (?, ?, ?)`
	_, err := r.DB().Exec(query, category.ID, category.Name, category.Color)
	if err != nil {
		return fmt.Errorf("failed to insert category: %w", err)
	}
	return nil
}

func (r *CategoryRepository) Update(category *models.Category) error {
	query := `UPDATE categories SET name = ?, color = ? WHERE id = ?`
	_, err := r.DB().Exec(query, category.Name, category.Color, category.ID)
	if err != nil {
		return fmt.Errorf("failed to update category: %w", err)
	}
	return nil
}

func (r *CategoryRepository) Delete(id string) error {
	query := `DELETE FROM categories WHERE id = ?`
	_, err := r.DB().Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
	}
	return nil
}

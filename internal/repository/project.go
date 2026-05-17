package repository

import (
	"context"
	"fmt"

	"showcase/internal/database"
	"showcase/internal/model"
)

func GetAllProjects(ctx context.Context) ([]model.Project, error) {
	rows, err := database.Pool.Query(ctx, `
		SELECT id, slug, title, description, content, tech_stack, status,
		       features, challenges, learnings, future_plans,
		       github_url, live_url, image_url, featured, sort_order,
		       created_at, updated_at
		FROM projects
		ORDER BY sort_order ASC, created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []model.Project
	for rows.Next() {
		var p model.Project
		err := rows.Scan(
			&p.ID, &p.Slug, &p.Title, &p.Description, &p.Content,
			&p.TechStack, &p.Status, &p.Features, &p.Challenges, &p.Learnings, &p.FuturePlans,
			&p.GitHubURL, &p.LiveURL, &p.ImageURL,
			&p.Featured, &p.SortOrder, &p.CreatedAt, &p.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		projects = append(projects, p)
	}

	return projects, rows.Err()
}

func GetFeaturedProjects(ctx context.Context) ([]model.Project, error) {
	rows, err := database.Pool.Query(ctx, `
		SELECT id, slug, title, description, content, tech_stack, status,
		       features, challenges, learnings, future_plans,
		       github_url, live_url, image_url, featured, sort_order,
		       created_at, updated_at
		FROM projects
		WHERE featured = true
		ORDER BY sort_order ASC, created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []model.Project
	for rows.Next() {
		var p model.Project
		err := rows.Scan(
			&p.ID, &p.Slug, &p.Title, &p.Description, &p.Content,
			&p.TechStack, &p.Status, &p.Features, &p.Challenges, &p.Learnings, &p.FuturePlans,
			&p.GitHubURL, &p.LiveURL, &p.ImageURL,
			&p.Featured, &p.SortOrder, &p.CreatedAt, &p.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		projects = append(projects, p)
	}

	return projects, rows.Err()
}

func GetProjectBySlug(ctx context.Context, slug string) (*model.Project, error) {
	var p model.Project
	err := database.Pool.QueryRow(ctx, `
		SELECT id, slug, title, description, content, tech_stack, status,
		       features, challenges, learnings, future_plans,
		       github_url, live_url, image_url, featured, sort_order,
		       created_at, updated_at
		FROM projects
		WHERE slug = $1
	`, slug).Scan(
		&p.ID, &p.Slug, &p.Title, &p.Description, &p.Content,
		&p.TechStack, &p.Status, &p.Features, &p.Challenges, &p.Learnings, &p.FuturePlans,
		&p.GitHubURL, &p.LiveURL, &p.ImageURL,
		&p.Featured, &p.SortOrder, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func GetProjectByID(ctx context.Context, id int) (*model.Project, error) {
	var p model.Project
	err := database.Pool.QueryRow(ctx, `
		SELECT id, slug, title, description, content, tech_stack, status,
		       features, challenges, learnings, future_plans,
		       github_url, live_url, image_url, featured, sort_order,
		       created_at, updated_at
		FROM projects
		WHERE id = $1
	`, id).Scan(
		&p.ID, &p.Slug, &p.Title, &p.Description, &p.Content,
		&p.TechStack, &p.Status, &p.Features, &p.Challenges, &p.Learnings, &p.FuturePlans,
		&p.GitHubURL, &p.LiveURL, &p.ImageURL,
		&p.Featured, &p.SortOrder, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func GetFilteredProjects(ctx context.Context, status, tech string) ([]model.Project, error) {
	query := `
		SELECT id, slug, title, description, content, tech_stack, status,
		       features, challenges, learnings, future_plans,
		       github_url, live_url, image_url, featured, sort_order,
		       created_at, updated_at
		FROM projects
		WHERE 1=1
	`
	args := []interface{}{}
	argNum := 1

	if status != "" {
		query += fmt.Sprintf(" AND status = $%d", argNum)
		args = append(args, status)
		argNum++
	}

	if tech != "" {
		query += fmt.Sprintf(" AND $%d = ANY(tech_stack)", argNum)
		args = append(args, tech)
		argNum++
	}

	query += ` ORDER BY sort_order ASC, created_at DESC`

	rows, err := database.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []model.Project
	for rows.Next() {
		var p model.Project
		err := rows.Scan(
			&p.ID, &p.Slug, &p.Title, &p.Description, &p.Content,
			&p.TechStack, &p.Status, &p.Features, &p.Challenges, &p.Learnings, &p.FuturePlans,
			&p.GitHubURL, &p.LiveURL, &p.ImageURL,
			&p.Featured, &p.SortOrder, &p.CreatedAt, &p.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		projects = append(projects, p)
	}

	return projects, rows.Err()
}

func GetAllStatuses(ctx context.Context) ([]string, error) {
	rows, err := database.Pool.Query(ctx, `
		SELECT DISTINCT status FROM projects ORDER BY status
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var statuses []string
	for rows.Next() {
		var s string
		if err := rows.Scan(&s); err != nil {
			return nil, err
		}
		statuses = append(statuses, s)
	}
	return statuses, rows.Err()
}

func GetAllTechnologies(ctx context.Context) ([]string, error) {
	rows, err := database.Pool.Query(ctx, `
		SELECT DISTINCT unnest(tech_stack) AS tech FROM projects ORDER BY tech
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var techs []string
	for rows.Next() {
		var t string
		if err := rows.Scan(&t); err != nil {
			return nil, err
		}
		techs = append(techs, t)
	}
	return techs, rows.Err()
}

func CreateProject(ctx context.Context, p *model.Project) error {
	return database.Pool.QueryRow(ctx, `
		INSERT INTO projects (
			slug, title, description, content, tech_stack, status,
			features, challenges, learnings, future_plans,
			github_url, live_url, image_url, featured, sort_order
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
		RETURNING id, created_at, updated_at
	`,
		p.Slug, p.Title, p.Description, p.Content, p.TechStack, p.Status,
		p.Features, p.Challenges, p.Learnings, p.FuturePlans,
		p.GitHubURL, p.LiveURL, p.ImageURL, p.Featured, p.SortOrder,
	).Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)
}

func UpdateProject(ctx context.Context, p *model.Project) error {
	_, err := database.Pool.Exec(ctx, `
		UPDATE projects SET
			slug = $2, title = $3, description = $4, content = $5,
			tech_stack = $6, status = $7, features = $8,
			challenges = $9, learnings = $10, future_plans = $11,
			github_url = $12, live_url = $13, image_url = $14,
			featured = $15, sort_order = $16, updated_at = NOW()
		WHERE id = $1
	`,
		p.ID, p.Slug, p.Title, p.Description, p.Content,
		p.TechStack, p.Status, p.Features,
		p.Challenges, p.Learnings, p.FuturePlans,
		p.GitHubURL, p.LiveURL, p.ImageURL,
		p.Featured, p.SortOrder,
	)
	return err
}

func DeleteProject(ctx context.Context, id int) error {
	_, err := database.Pool.Exec(ctx, `DELETE FROM projects WHERE id = $1`, id)
	return err
}

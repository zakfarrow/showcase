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

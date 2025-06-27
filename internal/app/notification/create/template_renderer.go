package create

import (
	"context"
	"fmt"
	"strings"
	"text/template"

	layoutPorts "notification/internal/domain/layout/ports"
	"notification/internal/domain/notification"
	"notification/pkg/logger"
)

type TemplateRenderer struct {
	layoutRepo layoutPorts.LayoutRepository
	logger     logger.Logger
}

func NewTemplateRenderer(layoutRepo layoutPorts.LayoutRepository, logger logger.Logger) *TemplateRenderer {
	return &TemplateRenderer{
		layoutRepo: layoutRepo,
		logger:     logger,
	}
}

func (r *TemplateRenderer) Render(ctx context.Context, n *notification.Notification) (string, string, error) {
	layout, err := r.layoutRepo.GetByID(ctx, n.LayoutID)
	if err != nil {
		return "", "", fmt.Errorf("failed to get layout: %w", err)
	}

	if !layout.IsActive {
		return "", "", fmt.Errorf("layout %s is not active", layout.ID)
	}

	title, err := r.renderTemplate(layout.Subject, n.Data)
	if err != nil {
		return "", "", fmt.Errorf("failed to render title: %w", err)
	}

	content, err := r.renderTemplate(layout.Body, n.Data)
	if err != nil {
		return "", "", fmt.Errorf("failed to render content: %w", err)
	}

	return title, content, nil
}

func (r *TemplateRenderer) renderTemplate(tmplText string, data map[string]interface{}) (string, error) {
	tmpl, err := template.New("notification").Parse(tmplText)
	if err != nil {
		return "", fmt.Errorf("template parsing error: %w", err)
	}

	var buf strings.Builder
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("template execution error: %w", err)
	}

	return buf.String(), nil
}

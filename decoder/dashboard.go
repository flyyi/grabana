package decoder

import (
	"fmt"

	"github.com/K-Phoen/grabana/dashboard"
	"github.com/K-Phoen/grabana/row"
)

var ErrPanelNotConfigured = fmt.Errorf("panel not configured")

type dashboardModel struct {
	Title           string
	Editable        bool
	SharedCrosshair bool `yaml:"shared_crosshair"`
	Tags            []string
	AutoRefresh     string `yaml:"auto_refresh"`

	TagsAnnotation []dashboard.TagAnnotation `yaml:"tags_annotations"`
	Variables      []dashboardVariable

	Rows []dashboardRow
}

func (d *dashboardModel) toDashboardBuilder() (dashboard.Builder, error) {
	emptyDashboard := dashboard.Builder{}
	opts := []dashboard.Option{
		d.editable(),
		d.sharedCrossHair(),
	}

	if len(d.Tags) != 0 {
		opts = append(opts, dashboard.Tags(d.Tags))
	}

	if d.AutoRefresh != "" {
		opts = append(opts, dashboard.AutoRefresh(d.AutoRefresh))
	}

	for _, tagAnnotation := range d.TagsAnnotation {
		opts = append(opts, dashboard.TagsAnnotation(tagAnnotation))
	}

	for _, variable := range d.Variables {
		opt, err := variable.toOption()
		if err != nil {
			return emptyDashboard, err
		}

		opts = append(opts, opt)
	}

	for _, r := range d.Rows {
		opt, err := r.toOption()
		if err != nil {
			return emptyDashboard, err
		}

		opts = append(opts, opt)
	}

	return dashboard.New(d.Title, opts...), nil
}

func (d *dashboardModel) sharedCrossHair() dashboard.Option {
	if d.SharedCrosshair {
		return dashboard.SharedCrossHair()
	}

	return dashboard.DefaultTooltip()
}

func (d *dashboardModel) editable() dashboard.Option {
	if d.Editable {
		return dashboard.Editable()
	}

	return dashboard.ReadOnly()
}

type dashboardPanel struct {
	Graph      *dashboardGraph
	Table      *dashboardTable
	SingleStat *dashboardSingleStat `yaml:"single_stat"`
	Text       *dashboardText
}

func (panel dashboardPanel) toOption() (row.Option, error) {
	if panel.Graph != nil {
		return panel.Graph.toOption()
	}
	if panel.Table != nil {
		return panel.Table.toOption()
	}
	if panel.SingleStat != nil {
		return panel.SingleStat.toOption()
	}
	if panel.Text != nil {
		return panel.Text.toOption(), nil
	}

	return nil, ErrPanelNotConfigured
}
package app

import (
	"github.com/mum4k/termdash/align"
	"github.com/mum4k/termdash/container"
	"github.com/mum4k/termdash/container/grid"
	"github.com/mum4k/termdash/linestyle"
	"github.com/mum4k/termdash/terminal/terminalapi"
)

func newGrid(
	b *buttonSet,
	w *widgets,
	s *summary,
	t terminalapi.Terminal,
) (*container.Container, error) {

	builder := grid.New()
	// Add first row
	builder.Add(
		grid.RowHeightPerc(30,
			// Add first column
			grid.ColWidthPercWithOpts(30,
				[]container.Option{
					container.Border(linestyle.Light),
					container.BorderTitle("Press Q to Quit"),
				},
				// Add donut timer
				grid.RowHeightPerc(80,
					grid.Widget(w.donTimer)),
				// Add text timer under donut timer
				grid.RowHeightPercWithOpts(20,
					[]container.Option{
						container.AlignHorizontal(align.HorizontalCenter),
					},
					grid.Widget(w.txtTimer,
						container.AlignHorizontal(align.HorizontalCenter),
						container.AlignVertical(align.VerticalMiddle),
						container.PaddingLeftPercent(49),
					),
				),
			),
			// Add second column
			grid.ColWidthPerc(70,
				// Display type
				grid.RowHeightPerc(80,
					grid.Widget(w.disType, container.Border(linestyle.Light)),
				),
				// Display info
				grid.RowHeightPerc(20,
					grid.Widget(w.txtInfo, container.Border(linestyle.Light)),
				),
			),
		),
	)

	// Add second row
	builder.Add(
		grid.RowHeightPerc(10,
			// Add button 'start'
			grid.ColWidthPerc(50,
				grid.Widget(b.btStart),
			),
			// Add button 'pause'
			grid.ColWidthPerc(50,
				grid.Widget(b.btPause),
			),
		),
	)

	// Add third row
	builder.Add(
		grid.RowHeightPerc(
			60,
			grid.ColWidthPerc(
				30,
				grid.Widget(
					s.bcDay,
					container.Border(linestyle.Light),
					container.BorderTitle("Daily Summary (minutes)"),
				),
			),
			grid.ColWidthPerc(
				70,
				grid.Widget(
					s.lcWeekly,
					container.Border(linestyle.Light),
					container.BorderTitle("Weekly Summary"),
				),
			),
		),
	)

	gridOpts, err := builder.Build()
	if err != nil {
		return nil, err
	}

	c, err := container.New(t, gridOpts...)
	if err != nil {
		return nil, err
	}

	return c, nil
}

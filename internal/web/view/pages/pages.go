package pages

import (
	"JoaoRafa19/myhousetask/internal/web/view/components"
	"JoaoRafa19/myhousetask/store"

	"github.com/a-h/templ"
)

func FamiliesTableComponent(families []store.ListRecentFamiliesRow) templ.Component {
	return components.FamiliesTable(families)
}

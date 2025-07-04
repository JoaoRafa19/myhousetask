package pages

import (
	db "JoaoRafa19/myhousetask/db/gen"
	"JoaoRafa19/myhousetask/internal/web/view/components"

	"github.com/a-h/templ"
)

func FamiliesTableComponent(families []db.ListRecentFamiliesRow) templ.Component {
	return components.FamiliesTable(families)
}

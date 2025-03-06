package bookmark

import (
	"database/sql"

	"github.com/go-gorp/gorp"

	"github.com/ovh/cds/sdk"
)

// LoadAll returns all bookmarks with icons and their description
func LoadAll(db gorp.SqlExecutor, userID string) ([]sdk.Bookmark, error) {
	var data []sdk.Bookmark
	query := `
		WITH results AS (
			(
				SELECT DISTINCT 'project' AS type, project.projectkey AS id, project.name AS label
				FROM project
				JOIN project_favorite ON project.id = project_favorite.project_id AND project_favorite.authentified_user_id = $1
			)
			UNION
			(
				SELECT 'workflow-legacy' AS type, CONCAT(project.projectkey, '/', workflow.name) AS id, workflow.name AS label
				FROM project
				JOIN workflow ON workflow.project_id = project.id
				JOIN workflow_favorite ON workflow.id = workflow_favorite.workflow_id AND workflow_favorite.authentified_user_id = $1
			)
		)		
		SELECT *
		FROM results
		ORDER BY type ASC, label ASC
	`

	if _, err := db.Select(&data, query, userID); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, sdk.WrapError(err, "cannot load bookmarks as admin")
	}

	return data, nil
}

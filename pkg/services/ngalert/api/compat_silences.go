package api

import (
	"fmt"

	"github.com/grafana/grafana/pkg/services/ngalert/api/tooling/definitions"
	"github.com/grafana/grafana/pkg/services/ngalert/models"
)

// Silence-specific compat functions to convert between API and model types.

func SilenceToGettableGrafanaSilence(s *models.SilenceWithMetadata) definitions.GettableGrafanaSilence {
	gettable := definitions.GettableGrafanaSilence{
		GettableSilence: (*definitions.GettableSilence)(s.Silence),
	}
	if s.Metadata != nil {
		gettable.Metadata = &definitions.SilenceMetadata{
			RuleUID:     s.Metadata.RuleUID,
			RuleTitle:   s.Metadata.RuleTitle,
			FolderUID:   s.Metadata.FolderUID,
			Permissions: make(map[definitions.SilencePermission]struct{}, len(s.Metadata.Permissions)),
		}
		for _, permission := range models.SilencePermissions() {
			if s.Metadata.Permissions.Has(permission) {
				p, err := SilencePermissionToAPI(permission)
				if err != nil {
					// Skip unknown permissions in response.
					continue
				}
				gettable.Metadata.Permissions[p] = struct{}{}
			}
		}
	}
	return gettable
}

func SilencesToGettableGrafanaSilences(silences []*models.SilenceWithMetadata) definitions.GettableGrafanaSilences {
	res := make(definitions.GettableGrafanaSilences, 0, len(silences))
	for _, sil := range silences {
		apiSil := SilenceToGettableGrafanaSilence(sil)
		res = append(res, &apiSil)
	}
	return res
}

func PostableSilenceToSilence(s definitions.PostableSilence) models.Silence {
	return models.Silence{
		ID:        &s.ID,
		Status:    nil,
		UpdatedAt: nil,
		Silence:   s.Silence,
	}
}

func SilencePermissionToAPI(p models.SilencePermission) (definitions.SilencePermission, error) {
	switch p {
	case models.SilencePermissionRead:
		return definitions.SilencePermissionRead, nil
	case models.SilencePermissionCreate:
		return definitions.SilencePermissionCreate, nil
	case models.SilencePermissionWrite:
		return definitions.SilencePermissionWrite, nil
	default:
		return "", fmt.Errorf("unknown permission: %s", p)
	}
}

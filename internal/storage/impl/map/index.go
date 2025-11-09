package maps

import (
	"vtb-apihack-2025/client-pilot/pe"
	"vtb-apihack-2025/internal/storage/interfaces"
)

type Cache struct {
	consents map[string]map[pe.PermissionsType]map[string]interfaces.Consent
}

func NewCache() *Cache {
	return &Cache{
		consents: make(map[string]map[pe.PermissionsType]map[string]interfaces.Consent, 0),
	}
}

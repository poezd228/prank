package dependencies

import (
	"github.com/mvd-inc/anibliss/pkg/time_manager"
)

func (d *Dependencies) TimeManager() time_manager.TimeManager {
	if d.timeManager == nil {
		d.timeManager = time_manager.New(3 * 3600) // TODO брать из конфига
	}
	return d.timeManager
}

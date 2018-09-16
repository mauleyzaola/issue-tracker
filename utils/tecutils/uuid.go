package tecutils

import (
	"github.com/nu7hatch/gouuid"
)

func UUID() string {
	uuid, _ := uuid.NewV4()
	return uuid.String()
}

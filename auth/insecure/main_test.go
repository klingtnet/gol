package insecure

import (
	tu "../../util/testing"
	"testing"
)

// TODO: TestOpen

func TestLogin(t *testing.T) {
	auth := &Auth{
		mapping: map[string]string{
			"joe":  "doe",
			"jane": "sane",
		},
	}

	tu.RequireNil(t, auth.Login("joe", "doe"))
	tu.RequireNil(t, auth.Login("jane", "sane"))
	tu.RequireNotNil(t, auth.Login("joe", "sane"))
	tu.RequireNotNil(t, auth.Login("joe", "oops"))
	tu.RequireNotNil(t, auth.Login("jane", "doe"))
	tu.RequireNotNil(t, auth.Login("jane", ""))
	tu.RequireNotNil(t, auth.Login("mo", "joe"))
}

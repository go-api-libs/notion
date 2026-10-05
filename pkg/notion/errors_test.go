package notion

import "testing"

func TestErrorError(t *testing.T) {
	for _, tc := range []struct {
		err  *Error
		want string
	}{
		{
			&Error{Status: 404, Code: ErrorCodeObjectNotFound, Message: "Could not find page."},
			"notion: 404 object_not_found: Could not find page.",
		},
		{
			&Error{
				Status: 403, Code: ErrorCodeAgentCreditLimitReached, Message: "Limit reached.",
				AdditionalData: map[string]PublicAPICommonErrorAdditionalDataValue{
					"recovery_kind": {String: "none"},
					"agents":        {String2: []string{"a", "b"}},
				},
			},
			"notion: 403 agent_credit_limit_reached: Limit reached. (agents: a, b; recovery_kind: none)",
		},
	} {
		if got := tc.err.Error(); got != tc.want {
			t.Errorf("got %q, want %q", got, tc.want)
		}
	}
}

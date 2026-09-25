package app

import (
	"github.com/velocitykode/velocity/contract"
)

// Errors adds the app's rules to the error handler. main.go passes this
// to v.Errors(...); rules registered here outrank the framework defaults.
//
// The framework already maps its own errors (unmatched routes 404/405,
// validation 422, unauthenticated 401, forbidden 403, orm not-found 404,
// panics 500) and answers JSON clients with application/problem+json.
// The one line below makes that the answer for every client. Add rules
// once the app has errors of its own, using
// github.com/velocitykode/velocity/problem.
func Errors(h contract.ErrorHandler) {
	// API-only app: every error answers application/problem+json, whatever
	// the client's Accept header says. Routes under r.API already do; this
	// covers everything else (unmatched paths outside the prefix included).
	h.SetAPIMode(true)

	// Answer a domain sentinel as a 410 problem:
	//
	//	problem.MapIs(h, ErrInviteExpired, func(err error) error {
	//		return problem.Gone("this invite has expired").WithCause(err)
	//	})
	//
	// Write your own JSON body for a domain error type:
	//
	//	problem.RenderFor[*QuotaError](h, func(rc problem.RenderContext, err *QuotaError, _ *problem.ErrorContext) bool {
	//		rc.SetHeader("Content-Type", "application/json")
	//		rc.WriteHeader(http.StatusTooManyRequests)
	//		_, _ = rc.Write([]byte(`{"error":"quota exceeded"}`))
	//		return true
	//	})
}

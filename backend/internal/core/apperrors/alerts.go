package apperrors

// Alert are common static strings used to indicate an error requiring an alert
// in our observability platform. This allows easily filter common errors vs
// errors that may require manual intervention.
// ENUM(
//
//	failed-subscription // Used for events where a user attempts to subscribe but fails after payment (stripe) and may require developer intervention.
//
// )
type Alert string

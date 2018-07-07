// Package statsd creates a simple statistics server that accepts UDP
// data and peridically pushes the stats to a statistics end point.
//
// UDP data packets each contain a single stat. If the packet begins with
// a `#` character, then the stat is a floating point value to be set. Otherwise
// the stat is an integer count that increments the stat's total. The data format
// is the type character (if a value), a number, followed by a `=` (equals sign),
// followed by the name of the statistic (max 255 characters in length).
//
// The statsd server aggregates statistics and pushes them to a final "end point"
// at certain intervals. Currently, the supported end points are:
//
// * StatHat
//
// Future end points we hope to support include:
//
// * AWS Kinesis
// * AWS Lambda
// * Graphite
package statsd

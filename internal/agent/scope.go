package agent

import "strconv"

func formatFloat(f float64) string { return strconv.FormatFloat(f, 'f', -1, 64) }

func formatInt(i int64) string { return strconv.FormatInt(i, 10) }

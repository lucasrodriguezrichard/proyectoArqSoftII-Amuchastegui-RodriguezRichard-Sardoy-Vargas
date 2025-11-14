package service

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/blassardoy/restaurant-reservas/search-api/internal/repository"
	"github.com/blassardoy/restaurant-reservas/search-api/internal/solr"
)

var (
	isoDateRegex = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)

	// Characters that need to be escaped for Solr/Lucene queries.
	solrEscaper = strings.NewReplacer(
		`+`, `\+`,
		`-`, `\-`,
		`&&`, `\&&`,
		`||`, `\||`,
		`!`, `\!`,
		`(`, `\(`,
		`)`, `\)`,
		`{`, `\{`,
		`}`, `\}`,
		`[`, `\[`,
		`]`, `\]`,
		`^`, `\^`,
		`"`, `\"`,
		`~`, `\~`,
		`?`, `\?`,
		`:`, `\:`,
		`\\`, `\\\\`,
		`/`, `\/`,
	)

	sortFieldMap = map[string]string{
		"date":         solr.FieldDate,
		"table":        solr.FieldTableNumber,
		"table_number": solr.FieldTableNumber,
		"capacity":     solr.FieldCapacity,
		"created":      solr.FieldCreatedAt,
		"created_at":   solr.FieldCreatedAt,
		"updated":      solr.FieldUpdatedAt,
		"updated_at":   solr.FieldUpdatedAt,
	}
)

func normalizeQuery(q repository.SearchQuery) repository.SearchQuery {
	normalized := q
	normalized.Q = buildFullTextQuery(q.Q)
	normalized.Filters = normalizeFilters(q.Filters)
	normalized.Sort, normalized.Order = sanitizeSort(q.Sort, q.Order)
	return normalized
}

func buildFullTextQuery(term string) string {
	trimmed := strings.TrimSpace(term)
	if trimmed == "" || trimmed == "*:*" {
		return "*:*"
	}
	// Allow advanced queries using Solr syntax (e.g. table_number:5)
	if strings.Contains(trimmed, ":") {
		return trimmed
	}

	tokens := strings.Fields(trimmed)
	if len(tokens) == 0 {
		return "*:*"
	}

	clauses := make([]string, 0, len(tokens))
	for _, token := range tokens {
		if clause := buildTokenClause(token); clause != "" {
			clauses = append(clauses, clause)
		}
	}
	if len(clauses) == 0 {
		return "*:*"
	}
	if len(clauses) == 1 {
		return clauses[0]
	}
	return strings.Join(clauses, " AND ")
}

func buildTokenClause(token string) string {
	token = strings.TrimSpace(token)
	if token == "" {
		return ""
	}
	lower := strings.ToLower(token)
	escaped := escapeValue(lower)
	wildcard := wildcardValue(lower)

	var clauseParts []string
	clauseParts = append(clauseParts,
		fmt.Sprintf("%s:%s", solr.FieldID, wildcard),
		fmt.Sprintf("%s:%s", solr.FieldReservationID, wildcard),
		fmt.Sprintf("%s:%s", solr.FieldMealType, wildcard),
		fmt.Sprintf("%s:%s", solr.FieldDate, wildcard),
	)

	if isoDateRegex.MatchString(lower) {
		clauseParts = append(clauseParts, fmt.Sprintf("%s:\"%s\"", solr.FieldDate, escaped))
	}

	if num, err := strconv.Atoi(lower); err == nil {
		clauseParts = append(clauseParts,
			fmt.Sprintf("%s:%d", solr.FieldTableNumber, num),
			fmt.Sprintf("%s:%d", solr.FieldCapacity, num),
		)
	}

	return fmt.Sprintf("(%s)", strings.Join(dedupe(clauseParts), " OR "))
}

func normalizeFilters(filters map[string]string) map[string]string {
	if len(filters) == 0 {
		return nil
	}

	out := make(map[string]string)
	for key, raw := range filters {
		value := strings.TrimSpace(raw)
		if value == "" {
			continue
		}
		switch normalizeFilterKey(key) {
		case solr.FieldMealType:
			out[solr.FieldMealType] = strings.ToLower(value)
		case solr.FieldIsAvailable:
			val := strings.ToLower(value)
			if val == "true" || val == "false" {
				out[solr.FieldIsAvailable] = val
			}
		case solr.FieldCapacity:
			if _, err := strconv.Atoi(value); err == nil {
				out[solr.FieldCapacity] = fmt.Sprintf("[%s TO *]", value)
			}
		case solr.FieldDate:
			if rangeQuery, ok := buildDateRangeQuery(value); ok {
				out[solr.FieldDate] = rangeQuery
			}
		case solr.FieldTableNumber:
			if _, err := strconv.Atoi(value); err == nil {
				out[solr.FieldTableNumber] = value
			}
		default:
			out[key] = value
		}
	}

	if len(out) == 0 {
		return nil
	}
	return out
}

func normalizeFilterKey(key string) string {
	switch strings.ToLower(strings.TrimSpace(key)) {
	case "meal_type":
		return solr.FieldMealType
	case "is_available":
		return solr.FieldIsAvailable
	case "capacity":
		return solr.FieldCapacity
	case "date":
		return solr.FieldDate
	case "table", "table_number":
		return solr.FieldTableNumber
	default:
		return key
	}
}

func sanitizeSort(field, order string) (string, string) {
	field = strings.ToLower(strings.TrimSpace(field))
	if field == "" {
		return "", ""
	}
	if mapped, ok := sortFieldMap[field]; ok {
		field = mapped
	}

	order = strings.ToLower(strings.TrimSpace(order))
	if order != "desc" {
		order = "asc"
	}
	return field, order
}

func escapeValue(value string) string {
	if value == "" {
		return ""
	}
	return solrEscaper.Replace(value)
}

func wildcardValue(value string) string {
	if value == "" {
		return ""
	}
	return fmt.Sprintf("*%s*", escapeValue(value))
}

func dedupe(values []string) []string {
	if len(values) <= 1 {
		return values
	}
	seen := make(map[string]struct{}, len(values))
	result := make([]string, 0, len(values))
	for _, v := range values {
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		result = append(result, v)
	}
	sort.Strings(result)
	return result
}

func normalizeDateValue(value string) (string, bool) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", false
	}
	if isoDateRegex.MatchString(value) {
		return value, true
	}

	layouts := []string{
		"02/01/2006",
		"02-01-2006",
		"2006/01/02",
		"2006.01.02",
	}

	for _, layout := range layouts {
		if t, err := time.Parse(layout, value); err == nil {
			return t.Format("2006-01-02"), true
		}
	}
	return "", false
}

func buildDateRangeQuery(value string) (string, bool) {
	normalized, ok := normalizeDateValue(value)
	if !ok {
		return "", false
	}

	date, err := time.Parse("2006-01-02", normalized)
	if err != nil {
		return "", false
	}

	start := date.UTC()
	end := start.Add(24*time.Hour - time.Nanosecond)

	return fmt.Sprintf("[%s TO %s]", start.Format(time.RFC3339), end.Format(time.RFC3339)), true
}

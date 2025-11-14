package domain

// TableConfig represents the fixed configuration for a table
type TableConfig struct {
	TableNumber int    `json:"table_number"`
	Capacity    int    `json:"capacity"`
	MealType    string `json:"meal_type"`
}

// GetPredefinedTables returns all predefined tables for each meal type
func GetPredefinedTables() []TableConfig {
	tables := []TableConfig{}

	// Breakfast tables (10 tables: 2-6 capacity)
	breakfastCapacities := []int{2, 2, 4, 4, 4, 6, 6, 6, 8, 8}
	for i, capacity := range breakfastCapacities {
		tables = append(tables, TableConfig{
			TableNumber: i + 1,
			Capacity:    capacity,
			MealType:    MealTypeBreakfast,
		})
	}

	// Lunch tables (10 tables: 2-8 capacity)
	lunchCapacities := []int{2, 2, 4, 4, 4, 6, 6, 6, 8, 8}
	for i, capacity := range lunchCapacities {
		tables = append(tables, TableConfig{
			TableNumber: i + 1,
			Capacity:    capacity,
			MealType:    MealTypeLunch,
		})
	}

	// Dinner tables (10 tables: 2-8 capacity)
	dinnerCapacities := []int{2, 2, 4, 4, 4, 6, 6, 6, 8, 8}
	for i, capacity := range dinnerCapacities {
		tables = append(tables, TableConfig{
			TableNumber: i + 1,
			Capacity:    capacity,
			MealType:    MealTypeDinner,
		})
	}

	// Event tables (10 tables: 8-20 capacity, larger groups)
	eventCapacities := []int{8, 10, 10, 12, 12, 15, 15, 18, 20, 20}
	for i, capacity := range eventCapacities {
		tables = append(tables, TableConfig{
			TableNumber: i + 1,
			Capacity:    capacity,
			MealType:    MealTypeEvent,
		})
	}

	return tables
}

// GetTableCapacity returns the capacity for a specific table number and meal type
func GetTableCapacity(tableNumber int, mealType string) (int, bool) {
	tables := GetPredefinedTables()
	for _, table := range tables {
		if table.TableNumber == tableNumber && table.MealType == mealType {
			return table.Capacity, true
		}
	}
	return 0, false
}

// GetTablesForMealType returns all tables for a specific meal type
func GetTablesForMealType(mealType string) []TableConfig {
	allTables := GetPredefinedTables()
	filtered := []TableConfig{}
	for _, table := range allTables {
		if table.MealType == mealType {
			filtered = append(filtered, table)
		}
	}
	return filtered
}

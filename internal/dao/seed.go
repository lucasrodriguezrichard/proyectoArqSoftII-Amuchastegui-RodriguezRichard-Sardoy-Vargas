package dao

import (
	"fmt"
	"log"

	"github.com/google/uuid"
)

// SeedDatabase inserts initial test data
func SeedDatabase(db *DB) error {
	log.Println("Starting database seeding...")

	// Create restaurant
	restaurantID := uuid.New().String()
	restaurantQuery := `
		INSERT INTO restaurants (id, name, address, phone, email, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
		ON CONFLICT (id) DO NOTHING
	`
	_, err := db.Exec(restaurantQuery, restaurantID, "Restaurant El Buen Sabor",
		"Av. Colón 123, Córdoba", "+54 351 123-4567",
		"contacto@elbuensabor.com",
		"Restaurante especializado en comida argentina con ambiente familiar")
	if err != nil {
		return fmt.Errorf("failed to seed restaurant: %w", err)
	}
	log.Println("✓ Restaurant seeded")

	// Create tables
	tables := []struct {
		number   int
		capacity int
	}{
		{1, 2},
		{2, 2},
		{3, 4},
		{4, 4},
		{5, 6},
		{6, 6},
		{7, 8},
		{8, 4},
		{9, 2},
		{10, 10},
	}

	tableQuery := `
		INSERT INTO tables (id, restaurant_id, number, capacity, is_available)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT DO NOTHING
	`

	for _, table := range tables {
		tableID := uuid.New().String()
		_, err := db.Exec(tableQuery, tableID, restaurantID, table.number, table.capacity, true)
		if err != nil {
			return fmt.Errorf("failed to seed table %d: %w", table.number, err)
		}
	}
	log.Println("✓ Tables seeded")

	// Create menu items
	menuItems := []struct {
		name        string
		description string
		price       float64
		category    string
	}{
		// Entradas
		{"Empanadas", "Empanadas de carne, jamón y queso, o pollo", 8.50, "Entradas"},
		{"Tabla de Quesos", "Selección de quesos argentinos con frutos secos", 15.00, "Entradas"},
		{"Provoleta", "Queso provolone a la parrilla con orégano", 12.00, "Entradas"},

		// Platos Principales
		{"Asado Completo", "Asado de tira, chorizo, morcilla y ensalada", 28.00, "Platos Principales"},
		{"Bife de Chorizo", "Bife de chorizo 400g con guarnición", 25.00, "Platos Principales"},
		{"Milanesa Napolitana", "Milanesa con jamón, queso y salsa", 18.50, "Platos Principales"},
		{"Ravioles", "Ravioles caseros con salsa a elección", 16.00, "Platos Principales"},
		{"Pollo Grillé", "Pechuga de pollo con verduras grilladas", 17.00, "Platos Principales"},

		// Postres
		{"Flan Casero", "Flan con dulce de leche y crema", 6.50, "Postres"},
		{"Tiramisu", "Tiramisu italiano tradicional", 7.00, "Postres"},
		{"Helado", "2 bochas de helado artesanal", 5.50, "Postres"},

		// Bebidas
		{"Vino Tinto", "Copa de vino tinto Malbec", 8.00, "Bebidas"},
		{"Vino Blanco", "Copa de vino blanco Torrontés", 7.50, "Bebidas"},
		{"Cerveza Artesanal", "Pinta de cerveza artesanal", 6.00, "Bebidas"},
		{"Gaseosa", "Gaseosa línea Coca-Cola", 3.50, "Bebidas"},
		{"Agua Mineral", "Agua mineral con o sin gas", 2.50, "Bebidas"},
		{"Café", "Café expreso", 3.00, "Bebidas"},
	}

	menuQuery := `
		INSERT INTO menu_items (id, restaurant_id, name, description, price, category, is_available)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT DO NOTHING
	`

	for _, item := range menuItems {
		itemID := uuid.New().String()
		_, err := db.Exec(menuQuery, itemID, restaurantID, item.name,
			item.description, item.price, item.category, true)
		if err != nil {
			return fmt.Errorf("failed to seed menu item %s: %w", item.name, err)
		}
	}
	log.Println("✓ Menu items seeded")

	// Create sample user (admin)
	userQuery := `
		INSERT INTO users (id, username, email, password, first_name, last_name, phone, role, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), NOW())
		ON CONFLICT DO NOTHING
	`

	adminID := uuid.New().String()
	// Password: admin123 (should be hashed in production)
	_, err = db.Exec(userQuery, adminID, "admin", "admin@elbuensabor.com",
		"$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy", // bcrypt hash of "admin123"
		"Admin", "Sistema", "+54 351 000-0000", "admin", true)
	if err != nil {
		return fmt.Errorf("failed to seed admin user: %w", err)
	}
	log.Println("✓ Admin user seeded (username: admin, password: admin123)")

	log.Println("Database seeding completed successfully!")
	log.Println("=====================================================")
	log.Println("Restaurant ID:", restaurantID)
	log.Println("Admin User: admin / admin123")
	log.Println("=====================================================")

	return nil
}

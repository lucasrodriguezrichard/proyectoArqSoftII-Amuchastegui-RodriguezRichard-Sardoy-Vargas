#!/bin/bash

# Script para crear el primer usuario administrador

echo "=== Creando usuario administrador ==="
echo ""

# Solicitar datos del admin
read -p "Username: " ADMIN_USERNAME
read -p "Email: " ADMIN_EMAIL
read -p "First Name: " ADMIN_FIRSTNAME
read -p "Last Name: " ADMIN_LASTNAME
read -sp "Password: " ADMIN_PASSWORD
echo ""

# Datos del admin
ADMIN_DATA=$(cat <<EOF
{
  "username": "$ADMIN_USERNAME",
  "email": "$ADMIN_EMAIL",
  "password": "$ADMIN_PASSWORD",
  "first_name": "$ADMIN_FIRSTNAME",
  "last_name": "$ADMIN_LASTNAME"
}
EOF
)

echo ""
echo "Creando usuario administrador..."

# Crear el usuario admin usando el endpoint de registro y luego actualizar el rol en la BD
# Primero registramos el usuario normalmente
RESPONSE=$(curl -s -X POST http://localhost:8080/api/users/register \
  -H "Content-Type: application/json" \
  -d "$ADMIN_DATA")

echo "Usuario creado. Respuesta del servidor:"
echo "$RESPONSE"

echo ""
echo "=== IMPORTANTE ==="
echo "Ahora necesitas actualizar el rol a 'admin' en la base de datos MySQL."
echo ""
echo "Ejecuta estos comandos:"
echo ""
echo "docker-compose exec users-db mysql -u root -pbla2ucc users -e \"UPDATE users SET role='admin' WHERE username='$ADMIN_USERNAME';\""
echo ""
echo "Luego podrás iniciar sesión en http://localhost:3000/login"
echo "Username/Email: $ADMIN_USERNAME o $ADMIN_EMAIL"
echo "Password: [el que ingresaste]"

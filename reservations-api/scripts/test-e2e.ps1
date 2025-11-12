Param(
  [string]$UsersApiUrl,
  [string]$ResApiUrl,
  [string]$OwnerId,
  [int]$TableNumber,
  [int]$Guests,
  [string]$MealType,
  [string]$DateTime
)

$ErrorActionPreference = 'Stop'

# Defaults compatible with Windows PowerShell 5.1
if (-not $UsersApiUrl -or $UsersApiUrl -eq '') { $UsersApiUrl = if ($env:USERS_API_URL) { $env:USERS_API_URL } else { 'http://localhost:8080' } }
if (-not $ResApiUrl -or $ResApiUrl -eq '')   { $ResApiUrl   = if ($env:RES_API_URL)   { $env:RES_API_URL }   else { 'http://localhost:8081' } }
if (-not $OwnerId -or $OwnerId -eq '')       { $OwnerId     = if ($env:OWNER_ID)      { $env:OWNER_ID }      else { '1' } }
if (-not $TableNumber)                        { $TableNumber = if ($env:TABLE_NUMBER)  { [int]$env:TABLE_NUMBER } else { 5 } }
if (-not $Guests)                             { $Guests      = if ($env:GUESTS)       { [int]$env:GUESTS }       else { 4 } }
if (-not $MealType -or $MealType -eq '')     { $MealType    = if ($env:MEAL_TYPE)    { $env:MEAL_TYPE }    else { 'dinner' } }
if (-not $DateTime -or $DateTime -eq '')     { $DateTime    = if ($env:DATETIME)     { $env:DATETIME }     else { '2025-11-20T20:00:00Z' } }

Write-Host "[1/7] Health checks"
Invoke-RestMethod -Method GET -Uri "$UsersApiUrl/health" | Out-Null
Invoke-RestMethod -Method GET -Uri "$ResApiUrl/health" | Out-Null

Write-Host "[2/7] Ensure user exists: $UsersApiUrl/api/users/$OwnerId"
try {
  $status = (Invoke-WebRequest -Method GET -Uri "$UsersApiUrl/api/users/$OwnerId" -UseBasicParsing).StatusCode
  Write-Host "status: $status"
} catch {
  Write-Warning "User check failed: $($_.Exception.Message)"
}

Write-Host "[3/7] Create reservation"
$payload = @{ 
  owner_id = $OwnerId
  table_number = $TableNumber
  guests = $Guests
  date_time = $DateTime
  meal_type = $MealType
  special_requests = 'Mesa junto a la ventana'
} | ConvertTo-Json -Depth 5

$create = Invoke-RestMethod -Method POST -Uri "$ResApiUrl/api/reservations" -ContentType 'application/json' -Body $payload
Write-Host "[created]" ($create | ConvertTo-Json -Depth 5)

$resId = $create.id
if (-not $resId) { throw "Failed to parse reservation id" }

Write-Host "[4/7] Get reservation by id: $resId"
Invoke-RestMethod -Method GET -Uri "$ResApiUrl/api/reservations/$resId" | Out-Null
Write-Host "[got]"

Write-Host "[5/7] List reservations (limit=5)"
Invoke-RestMethod -Method GET -Uri "$ResApiUrl/api/reservations?limit=5&offset=0" | Out-Null
Write-Host "[list ok]"

Write-Host "[6/7] Confirm reservation"
Invoke-RestMethod -Method POST -Uri "$ResApiUrl/api/reservations/$resId/confirm" -ContentType 'application/json' -Body '{}' | Out-Null
Write-Host "[confirm ok]"

Write-Host "[7/7] Update reservation (guests=6)"
Invoke-RestMethod -Method PUT -Uri "$ResApiUrl/api/reservations/$resId" -ContentType 'application/json' -Body '{"guests":6}' | Out-Null
Write-Host "[update ok]"

Write-Host "E2E OK - reservation id: $resId"

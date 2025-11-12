Param(
  [string]$SolrUrl = $env:SOLR_URL,
  [string]$Core = $env:SOLR_CORE
)

if (-not $SolrUrl -or $SolrUrl -eq '') { $SolrUrl = 'http://localhost:8983/solr' }
if (-not $Core -or $Core -eq '') { $Core = 'reservations' }

$schema = "$SolrUrl/$Core/schema"

function Add-Field($name, $type) {
  $body = @{ 'add-field' = @{ name = $name; type = $type; stored = $true; indexed = $true } } | ConvertTo-Json -Depth 5
  try { Invoke-RestMethod -Method POST -Uri $schema -ContentType 'application/json' -Body $body | Out-Null } catch { }
}

Write-Host "Init Solr schema for $Core"
Add-Field 'id' 'string'
Add-Field 'owner_id' 'string'
Add-Field 'table_number' 'pint'
Add-Field 'guests' 'pint'
Add-Field 'date_time' 'pdate'
Add-Field 'meal_type' 'string'
Add-Field 'status' 'string'
Add-Field 'total_price' 'pfloat'
Add-Field 'created_at' 'pdate'
Add-Field 'updated_at' 'pdate'
Write-Host 'Schema init done.'


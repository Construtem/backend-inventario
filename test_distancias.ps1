# Script para probar la funcionalidad de cálculo de distancias
# Ejecutar desde PowerShell en la carpeta del proyecto

Write-Host "🚀 Probando API de cálculo de distancias" -ForegroundColor Green
Write-Host "=" * 50

# Función para hacer peticiones HTTP
function Test-Endpoint {
    param(
        [string]$Url,
        [string]$Method = "GET",
        [hashtable]$Body = $null,
        [string]$Description
    )
    
    Write-Host "`n📍 $Description" -ForegroundColor Yellow
    
    try {
        $params = @{
            Uri = $Url
            Method = $Method
            ContentType = "application/json"
        }
        
        if ($Body) {
            $params.Body = ($Body | ConvertTo-Json)
        }
        
        $response = Invoke-RestMethod @params
        Write-Host "✅ Éxito:" -ForegroundColor Green
        $response | ConvertTo-Json -Depth 3
    }
    catch {
        Write-Host "❌ Error:" -ForegroundColor Red
        Write-Host $_.Exception.Message
        if ($_.Exception.Response) {
            $reader = [System.IO.StreamReader]::new($_.Exception.Response.GetResponseStream())
            $errorBody = $reader.ReadToEnd()
            Write-Host "Detalles: $errorBody"
        }
    }
}

# URL base del API
$baseUrl = "http://localhost:8080/api"

# Test 1: Verificar servidor
Test-Endpoint -Url "$baseUrl/despachos" -Description "Verificando conexión al servidor"

# Test 2: Calcular distancia simple
$distanciaBody = @{
    origen = "Santiago Centro, Santiago, Chile"
    destino = "Las Condes, Santiago, Chile"
}

Test-Endpoint -Url "$baseUrl/calcular-distancia" -Method "POST" -Body $distanciaBody -Description "Calculando distancia básica"

# Test 3: Obtener despachos con distancias
Test-Endpoint -Url "$baseUrl/despachos-distancia" -Description "Obteniendo despachos con distancias"

# Test 4: Probar con direcciones específicas chilenas
Write-Host "`n🇨🇱 Probando con direcciones chilenas específicas:" -ForegroundColor Cyan

$rutasChilenas = @(
    @{
        origen = "Av. Providencia 1234, Providencia, Santiago"
        destino = "Av. Las Condes 567, Las Condes, Santiago"
    },
    @{
        origen = "Plaza de Armas, Santiago Centro, Santiago"
        destino = "Costanera Center, Providencia, Santiago"
    },
    @{
        origen = "Universidad de Chile, Santiago"
        destino = "Aeropuerto Internacional Arturo Merino Benítez, Pudahuel"
    }
)

foreach ($ruta in $rutasChilenas) {
    $descripcion = "Ruta: $($ruta.origen) → $($ruta.destino)"
    Test-Endpoint -Url "$baseUrl/calcular-distancia" -Method "POST" -Body $ruta -Description $descripcion
    Start-Sleep -Milliseconds 500  # Pausa para no sobrecargar la API
}

Write-Host "`n📋 Resumen de pruebas completado" -ForegroundColor Green
Write-Host "Si hay errores de API key, asegúrate de:" -ForegroundColor Yellow
Write-Host "1. Configurar GOOGLE_MAPS_API_KEY en tu archivo .env" -ForegroundColor White
Write-Host "2. Habilitar Distance Matrix API en Google Cloud Console" -ForegroundColor White
Write-Host "3. Verificar que la API key no tenga restricciones de dominio" -ForegroundColor White

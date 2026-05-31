<?php
require '../vendor/autoload.php'; // Load Composer dependencies

use GeoIp2\Database\Reader;

// Path to GeoLite2-City.mmdb
$databasePath = __DIR__ . '/GeoLite2-City.mmdb';

// AbuseIPDB API key (replace with your key)
$abuseIPDBApiKey = 'df79338205938c4342a9567851526e05f87845fe483d2a324f323a72333a1f97f411e7984118cab5';

// Check if database file exists
if (!file_exists($databasePath)) {
    http_response_code(500);
    echo json_encode(['error' => 'GeoLite2-City.mmdb file not found']);
    exit;
}

try {
    // Initialize GeoIP2 Reader
    $reader = new Reader($databasePath);

    // Get IP from query parameter or client
    $ip = isset($_GET['ip']) ? $_GET['ip'] : $_SERVER['REMOTE_ADDR'];

    // Validate IP
    if (!filter_var($ip, FILTER_VALIDATE_IP)) {
        http_response_code(400);
        echo json_encode(['error' => 'Invalid IP address']);
        exit;
    }

    // Fetch geolocation data
    $record = $reader->city($ip);

    // Extract continent, country, and city
    $continent = $record->continent->names['en'] ?? 'Unknown';
    $countryName = $record->country->names['en'] ?? 'Unknown';
    $countryCode = $record->country->isoCode ?? 'Unknown';
    $city = $record->city->names['en'] ?? 'Unknown';

    // Check IP reputation with AbuseIPDB
    $reputation = 'unknown';
    $url = "https://api.abuseipdb.com/api/v2/check?ipAddress=" . urlencode($ip) . "&maxAgeInDays=90";
    $options = [
        'http' => [
            'header' => "Key: $abuseIPDBApiKey\r\nAccept: application/json\r\n",
            'method' => 'GET',
        ],
    ];
    $context = stream_context_create($options);
    $response = @file_get_contents($url, false, $context);

    if ($response !== false) {
        $data = json_decode($response, true);
        if (isset($data['data']['abuseConfidenceScore'])) {
            $score = $data['data']['abuseConfidenceScore'];
            // Simple reputation logic: score < 25 is "good", otherwise "bad"
            $reputation = $score < 25 ? 'good' : 'bad';
        }
    } else {
        // Fallback if API call fails
        $reputation = 'unknown';
    }

    // Output in requested JSON format
    header('Content-Type: application/json; charset=utf-8');
    echo json_encode([
        'ip' => $ip,
        'continent' => $continent,
        'country' => [
            'name' => $countryName,
            'code' => $countryCode,
        ],
        'city' => $city,
        'reputation' => $reputation
    ], JSON_UNESCAPED_UNICODE | JSON_PRETTY_PRINT);

} catch (\GeoIp2\Exception\AddressNotFoundException $e) {
    http_response_code(404);
    echo json_encode(['error' => 'IP not found in database']);
} catch (Exception $e) {
    http_response_code(500);
    echo json_encode(['error' => 'Error processing request: ' . $e->getMessage()]);
}

// Close Reader
$reader->close();
?>
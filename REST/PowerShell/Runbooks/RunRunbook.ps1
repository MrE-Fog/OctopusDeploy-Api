$ErrorActionPreference = "Stop";

# Define working variables
$octopusURL = "https://youroctourl"
$octopusAPIKey = "API-YOURAPIKEY"
$header = @{ "X-Octopus-ApiKey" = $octopusAPIKey }
$spaceName = "default"
$projectName = "MyProject"
$runbookName = "MyRunbook"
$environmentNames = @("Development", "Staging")
$environmentIds = @()

# Optional Tenant
$tenantName = ""
$tenantId = $null

# Get space
$space = (Invoke-RestMethod -Method Get -Uri "$octopusURL/api/spaces/all" -Headers $header) | Where-Object {$_.Name -eq $spaceName}

# Get project
$project = (Invoke-RestMethod -Method Get -Uri "$octopusURL/api/$($space.Id)/projects/all" -Headers $header) | Where-Object {$_.Name -eq $projectName}

# Get runbook
$runbook = (Invoke-RestMethod -Method Get -Uri "$octopusURL/api/$($space.Id)/runbooks/all" -Headers $header) | Where-Object {$_.Name -eq $runbookName -and $_.ProjectId -eq $project.Id}

# Get environments
$environments = (Invoke-RestMethod -Method Get -Uri "$octopusURL/api/$($space.Id)/environments/all" -Headers $header) | Where-Object {$environmentNames -contains $_.Name}
foreach ($environment in $environments)
{
    $environmentIds += $environment.Id
}

# Optionally get tenant
if (![string]::IsNullOrEmpty($tenantName)) {
    $tenant = (Invoke-RestMethod -Method Get -Uri "$octopusURL/api/$($space.Id)/tenants/all" -Headers $header) | Where-Object {$_.Name -eq $tenantName} | Select-Object -First 1
    $tenantId = $tenant.Id
}

$runbook = (Invoke-RestMethod -Method Get -Uri "$octopusURL/api/$($space.Id)/runbooks/all" -Headers $header) | Where-Object {$_.Name -eq $runbookName -and $_.ProjectId -eq $project.Id}

# Run runbook per selected environment
foreach ($environmentId in $environmentIds)
{
    # Create json payload
    $jsonPayload = @{
        RunbookId = $runbook.Id
        RunbookSnapshotId = $runbook.PublishedRunbookSnapshotId
        EnvironmentId = $environmentId
        TenantId = $tenantId
        SkipActions = @()
        SpecificMachineIds = @()
        ExcludedMachineIds = @()
    }

    # Run runbook
    Invoke-RestMethod -Method Post -Uri "$octopusURL/api/$($space.Id)/runbookRuns" -Body ($jsonPayload | ConvertTo-Json -Depth 10) -Headers $header
}
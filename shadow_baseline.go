package main

// generateShadowBaseline создает сбалансированную дефолтную конфигурацию
func generateShadowBaseline(budget BudgetRange, useCase UseCase) BuildConfig {
	cfg := BuildConfig{
		Budget:             budget,
		UseCase:            useCase,
		GameType:           "all",
		MonitorResolution:  "1440p",
		MonitorRefreshRate: "144-180",
		CPUVendor:          "amd",
		GPUVendor:          "nvidia",
		RAMSize:            "16",
		StorageType:        "nvme",
		StorageSize:        "1",
		Color:              "black",
		Cooling:            "air",
		Case:               "mid-tower",
		UpgradeHorizon:     "upgrade",
	}

	// Корректировки под бюджет
	switch budget {
	case BudgetLow:
		cfg.MonitorResolution = "1080p"
		cfg.MonitorRefreshRate = "60-75"
		cfg.RAMSize = "16"
		cfg.StorageSize = "512"
		cfg.Cooling = "air"
		cfg.Case = "mid-tower"
		cfg.UpgradeHorizon = "5years"
	case BudgetMidLow:
		cfg.MonitorResolution = "1440p"
		cfg.MonitorRefreshRate = "144-180"
		cfg.RAMSize = "32"
		cfg.StorageSize = "1"
		cfg.Cooling = "air"
		cfg.Case = "mid-tower"
		cfg.UpgradeHorizon = "upgrade"
	case BudgetMidHi:
		cfg.MonitorResolution = "1440p"
		cfg.MonitorRefreshRate = "144-180"
		cfg.RAMSize = "32"
		cfg.StorageSize = "2"
		cfg.Cooling = "water"
		cfg.Case = "aquarium"
		cfg.UpgradeHorizon = "upgrade"
	case BudgetHigh:
		cfg.MonitorResolution = "4k"
		cfg.MonitorRefreshRate = "144-180"
		cfg.RAMSize = "64"
		cfg.StorageSize = "2"
		cfg.Cooling = "water"
		cfg.Case = "aquarium"
		cfg.UpgradeHorizon = "upgrade"
	}

	// Корректировки под use case
	switch useCase {
	case UseCaseOffice:
		cfg.GameType = "none"
		cfg.MonitorResolution = "1080p"
		cfg.MonitorRefreshRate = "60-75"
		cfg.RAMSize = "16"
		cfg.StorageSize = "512"
		cfg.Cooling = "air"
		cfg.Case = "mid-tower"
		cfg.UpgradeHorizon = "5years"
	case UseCaseWork:
		cfg.RAMSize = "64"
		cfg.StorageSize = "2"
		cfg.Cooling = "water"
	case UseCaseStreaming:
		cfg.RAMSize = "32"
		cfg.CPUVendor = "amd" // для NVENC можно и Intel
	}

	return cfg
}

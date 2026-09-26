package main

import "time"

type Step int

const (
	StepIdle Step = iota
	StepBudget
	StepUseCase
	StepGameType
	StepMonitorResolution
	StepMonitorRefreshRate
	StepCPUVendor
	StepGPUVendor
	StepRAMSize
	StepStorageType
	StepStorageSize
	StepColor
	StepCooling
	StepCase
	StepUpgradeHorizon
	StepComplete
)

type BudgetRange string

const (
	BudgetLow    BudgetRange = "low"
	BudgetMidLow BudgetRange = "mid_low"
	BudgetMidHi  BudgetRange = "mid_hi"
	BudgetHigh   BudgetRange = "high"
)

type UseCase string

const (
	UseCaseGaming    UseCase = "gaming"
	UseCaseStreaming UseCase = "streaming"
	UseCaseWork      UseCase = "work"
	UseCaseOffice    UseCase = "office"
	UseCaseUniversal UseCase = "universal"
)

type Component struct {
	ID         int
	Name       string
	Price      int
	Specs      string
	Socket     string
	FormFactor string
	RAMType    string
	TDP        int
	Power      int
	HasIGPU    bool
}

type BuildResult struct {
	CPU         Component
	GPU         Component
	Motherboard Component
	RAM         Component
	Storage     Component
	Case        Component
	Cooling     Component
	PSU         Component
	TotalPrice  int
}

type BuildConfig struct {
	Budget             BudgetRange
	UseCase            UseCase
	GameType           string
	MonitorResolution  string
	MonitorRefreshRate string
	CPUVendor          string
	GPUVendor          string
	RAMSize            string
	StorageType        string
	StorageSize        string
	Color              string
	Cooling            string
	Case               string
	UpgradeHorizon     string
}

type UserSession struct {
	UserID      int64
	CurrentStep Step
	Default     BuildConfig
	User        BuildConfig
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

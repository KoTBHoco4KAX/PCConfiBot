package main

import (
	"log"

	tele "gopkg.in/telebot.v3"
)

type BotHandlers struct {
	bot       *tele.Bot
	sessions  *SessionManager
	keyboards *KeyboardBuilder
}

func NewBotHandlers(bot *tele.Bot, sm *SessionManager) *BotHandlers {
	return &BotHandlers{
		bot:       bot,
		sessions:  sm,
		keyboards: NewKeyboardBuilder(),
	}
}

func (h *BotHandlers) RegisterHandlers() {
	h.bot.Handle("/start", h.handleStart)
	h.bot.Handle("/reset", h.handleReset)
	h.bot.Handle(&tele.Btn{Unique: "restart"}, h.handleReset)

	h.registerBudgetHandlers()
	h.registerUseCaseHandlers()
	h.registerGameTypeHandlers()
	h.registerMonitorHandlers()
	h.registerVendorHandlers()
	h.registerRAMHandlers()
	h.registerStorageHandlers()
	h.registerOtherHandlers()
}

func (h *BotHandlers) handleStart(c tele.Context) error {
	session := h.sessions.Mutate(c.Sender().ID, func(s *UserSession) {
		s.CurrentStep = StepBudget
	})

	text := "Привет! Я бот-конфигуратор ПК.\n\n" +
		"Помогу подобрать сбалансированный сетап под твои задачи.\n" +
		"Для начала укажи планируемый бюджет:"

	return c.EditOrSend(text, h.keyboards.buildBudgetKeyboard(&session))
}

func (h *BotHandlers) handleReset(c tele.Context) error {
	h.sessions.Reset(c.Sender().ID)
	return h.handleStart(c)
}

func (h *BotHandlers) registerBudgetHandlers() {
	budgets := []struct {
		data   string
		budget BudgetRange
	}{
		{"budget_low", BudgetLow},
		{"budget_mid_low", BudgetMidLow},
		{"budget_mid_hi", BudgetMidHi},
		{"budget_high", BudgetHigh},
	}

	for _, item := range budgets {
		b := item.budget
		h.bot.Handle(&tele.Btn{Unique: item.data}, func(c tele.Context) error {
			session := h.sessions.Mutate(c.Sender().ID, func(s *UserSession) {
				s.User.Budget = b
				s.CurrentStep = StepUseCase
				s.Default = generateShadowBaseline(b, s.Default.UseCase)
			})

			text := "Отлично! Теперь выбери основной сценарий использования:"
			return c.EditOrSend(text, h.keyboards.buildUseCaseKeyboard(&session))
		})
	}
}

func (h *BotHandlers) registerUseCaseHandlers() {
	useCases := []struct {
		data    string
		useCase UseCase
	}{
		{"usecase_gaming", UseCaseGaming},
		{"usecase_streaming", UseCaseStreaming},
		{"usecase_work", UseCaseWork},
		{"usecase_office", UseCaseOffice},
		{"usecase_universal", UseCaseUniversal},
	}

	for _, item := range useCases {
		u := item.useCase
		h.bot.Handle(&tele.Btn{Unique: item.data}, func(c tele.Context) error {
			var session UserSession

			if u == UseCaseOffice {
				session = h.sessions.Mutate(c.Sender().ID, func(s *UserSession) {
					s.User.UseCase = u
					s.CurrentStep = StepMonitorResolution
					s.Default = generateShadowBaseline(s.User.Budget, u)
				})
				text := "Для офисного ПК игры не требуются. Выбери разрешение монитора:"
				return c.EditOrSend(text, h.keyboards.buildMonitorResolutionKeyboard(&session))
			}

			session = h.sessions.Mutate(c.Sender().ID, func(s *UserSession) {
				s.User.UseCase = u
				s.CurrentStep = StepGameType
				s.Default = generateShadowBaseline(s.User.Budget, u)
			})

			text := "Какие игры планируешь запускать?"
			return c.EditOrSend(text, h.keyboards.buildGameTypeKeyboard(&session))
		})
	}
}

func (h *BotHandlers) registerGameTypeHandlers() {
	gameTypes := []struct {
		data  string
		value string
	}{
		{"game_cybersport", "cybersport"},
		{"game_aaa", "aaa"},
		{"game_all", "all"},
	}

	for _, item := range gameTypes {
		g := item.value
		h.bot.Handle(&tele.Btn{Unique: item.data}, func(c tele.Context) error {
			session := h.sessions.Mutate(c.Sender().ID, func(s *UserSession) {
				s.User.GameType = g
				s.CurrentStep = StepMonitorResolution
			})

			text := "Выбери разрешение монитора:"
			return c.EditOrSend(text, h.keyboards.buildMonitorResolutionKeyboard(&session))
		})
	}
}

func (h *BotHandlers) registerMonitorHandlers() {
	resolutions := []struct {
		data  string
		value string
	}{
		{"res_1080p", "1080p"},
		{"res_1440p", "1440p"},
		{"res_4k", "4k"},
		{"res_auto", "auto"},
	}

	for _, item := range resolutions {
		res := item.value
		h.bot.Handle(&tele.Btn{Unique: item.data}, func(c tele.Context) error {
			session := h.sessions.Mutate(c.Sender().ID, func(s *UserSession) {
				if res == "auto" {
					s.User.MonitorResolution = s.Default.MonitorResolution
				} else {
					s.User.MonitorResolution = res
				}
				s.CurrentStep = StepMonitorRefreshRate
			})

			text := "Выбери частоту обновления монитора:"
			return c.EditOrSend(text, h.keyboards.buildMonitorRefreshRateKeyboard(&session))
		})
	}

	refreshRates := []struct {
		data  string
		value string
	}{
		{"refresh_60", "60-75"},
		{"refresh_144", "144-180"},
		{"refresh_240", "240+"},
		{"refresh_auto", "auto"},
	}

	for _, item := range refreshRates {
		rate := item.value
		h.bot.Handle(&tele.Btn{Unique: item.data}, func(c tele.Context) error {
			session := h.sessions.Mutate(c.Sender().ID, func(s *UserSession) {
				if rate == "auto" {
					s.User.MonitorRefreshRate = s.Default.MonitorRefreshRate
				} else {
					s.User.MonitorRefreshRate = rate
				}
				s.CurrentStep = StepCPUVendor
			})

			text := "Выбери производителя процессора:"
			return c.EditOrSend(text, h.keyboards.buildCPUVendorKeyboard(&session))
		})
	}
}

func (h *BotHandlers) registerVendorHandlers() {
	cpuVendors := []struct {
		data  string
		value string
	}{
		{"cpu_amd", "amd"},
		{"cpu_intel", "intel"},
		{"cpu_auto", "auto"},
	}

	for _, item := range cpuVendors {
		vendor := item.value
		h.bot.Handle(&tele.Btn{Unique: item.data}, func(c tele.Context) error {
			session := h.sessions.Mutate(c.Sender().ID, func(s *UserSession) {
				if vendor == "auto" {
					s.User.CPUVendor = s.Default.CPUVendor
				} else {
					s.User.CPUVendor = vendor
				}
				s.CurrentStep = StepGPUVendor
			})

			text := "Выбери производителя видеокарты:"
			return c.EditOrSend(text, h.keyboards.buildGPUVendorKeyboard(&session))
		})
	}

	gpuVendors := []struct {
		data  string
		value string
	}{
		{"gpu_nvidia", "nvidia"},
		{"gpu_amd", "amd"},
		{"gpu_auto", "auto"},
	}

	for _, item := range gpuVendors {
		vendor := item.value
		h.bot.Handle(&tele.Btn{Unique: item.data}, func(c tele.Context) error {
			session := h.sessions.Mutate(c.Sender().ID, func(s *UserSession) {
				if vendor == "auto" {
					s.User.GPUVendor = s.Default.GPUVendor
				} else {
					s.User.GPUVendor = vendor
				}
				s.CurrentStep = StepRAMSize
			})

			text := "Выбери объем оперативной памяти:"
			return c.EditOrSend(text, h.keyboards.buildRAMSizeKeyboard(&session))
		})
	}
}

func (h *BotHandlers) registerRAMHandlers() {
	ramOptions := []struct {
		data  string
		value string
	}{
		{"ram_16", "16"},
		{"ram_32", "32"},
		{"ram_64", "64"},
		{"ram_auto", "auto"},
	}

	for _, item := range ramOptions {
		ram := item.value
		h.bot.Handle(&tele.Btn{Unique: item.data}, func(c tele.Context) error {
			session := h.sessions.Mutate(c.Sender().ID, func(s *UserSession) {
				if ram == "auto" {
					s.User.RAMSize = s.Default.RAMSize
				} else {
					s.User.RAMSize = ram
				}
				s.CurrentStep = StepStorageType
			})

			text := "Выбери конфигурацию накопителей:"
			return c.EditOrSend(text, h.keyboards.buildStorageTypeKeyboard(&session))
		})
	}
}

func (h *BotHandlers) registerStorageHandlers() {
	storageTypes := []struct {
		data  string
		value string
	}{
		{"storage_nvme", "nvme"},
		{"storage_combo", "combo"},
		{"storage_auto", "auto"},
	}

	for _, item := range storageTypes {
		sType := item.value
		h.bot.Handle(&tele.Btn{Unique: item.data}, func(c tele.Context) error {
			session := h.sessions.Mutate(c.Sender().ID, func(s *UserSession) {
				if sType == "auto" {
					s.User.StorageType = s.Default.StorageType
				} else {
					s.User.StorageType = sType
				}
				s.CurrentStep = StepStorageSize
			})

			text := "Выбери объем накопителя:"
			return c.EditOrSend(text, h.keyboards.buildStorageSizeKeyboard(&session))
		})
	}

	storageSizes := []struct {
		data  string
		value string
	}{
		{"storage_512", "512"},
		{"storage_1tb", "1"},
		{"storage_2tb", "2"},
		{"storage_auto", "auto"},
	}

	for _, item := range storageSizes {
		sSize := item.value
		h.bot.Handle(&tele.Btn{Unique: item.data}, func(c tele.Context) error {
			session := h.sessions.Mutate(c.Sender().ID, func(s *UserSession) {
				if sSize == "auto" {
					s.User.StorageSize = s.Default.StorageSize
				} else {
					s.User.StorageSize = sSize
				}
				s.CurrentStep = StepColor
			})

			text := "Выбери цвет компонентов сборки:"
			return c.EditOrSend(text, h.keyboards.buildColorKeyboard(&session))
		})
	}
}

func (h *BotHandlers) registerOtherHandlers() {
	colors := []struct {
		data  string
		value string
	}{
		{"color_black", "black"},
		{"color_white", "white"},
		{"color_auto", "auto"},
	}

	for _, item := range colors {
		clr := item.value
		h.bot.Handle(&tele.Btn{Unique: item.data}, func(c tele.Context) error {
			session := h.sessions.Mutate(c.Sender().ID, func(s *UserSession) {
				if clr == "auto" {
					s.User.Color = s.Default.Color
				} else {
					s.User.Color = clr
				}
				s.CurrentStep = StepCooling
			})

			text := "Выбери тип системы охлаждения:"
			return c.EditOrSend(text, h.keyboards.buildCoolingKeyboard(&session))
		})
	}

	coolingOptions := []struct {
		data  string
		value string
	}{
		{"cooling_air", "air"},
		{"cooling_water", "water"},
		{"cooling_auto", "auto"},
	}

	for _, item := range coolingOptions {
		cool := item.value
		h.bot.Handle(&tele.Btn{Unique: item.data}, func(c tele.Context) error {
			session := h.sessions.Mutate(c.Sender().ID, func(s *UserSession) {
				if cool == "auto" {
					s.User.Cooling = s.Default.Cooling
				} else {
					s.User.Cooling = cool
				}
				s.CurrentStep = StepCase
			})

			text := "Выбери форм-фактор корпуса:"
			return c.EditOrSend(text, h.keyboards.buildCaseKeyboard(&session))
		})
	}

	caseOptions := []struct {
		data  string
		value string
	}{
		{"case_mid", "mid-tower"},
		{"case_aquarium", "aquarium"},
		{"case_mini", "mini-itx"},
		{"case_auto", "auto"},
	}

	for _, item := range caseOptions {
		cs := item.value
		h.bot.Handle(&tele.Btn{Unique: item.data}, func(c tele.Context) error {
			session := h.sessions.Mutate(c.Sender().ID, func(s *UserSession) {
				if cs == "auto" {
					s.User.Case = s.Default.Case
				} else {
					s.User.Case = cs
				}
				s.CurrentStep = StepUpgradeHorizon
			})

			text := "Выбери горизонт эксплуатации:"
			return c.EditOrSend(text, h.keyboards.buildUpgradeHorizonKeyboard(&session))
		})
	}

	horizons := []struct {
		data  string
		value string
	}{
		{"horizon_5years", "5years"},
		{"horizon_upgrade", "upgrade"},
		{"horizon_auto", "auto"},
	}

	for _, item := range horizons {
		hrz := item.value
		h.bot.Handle(&tele.Btn{Unique: item.data}, func(c tele.Context) error {
			h.sessions.Mutate(c.Sender().ID, func(s *UserSession) {
				if hrz == "auto" {
					s.User.UpgradeHorizon = s.Default.UpgradeHorizon
				} else {
					s.User.UpgradeHorizon = hrz
				}
				s.CurrentStep = StepComplete
			})

			return h.handleComplete(c)
		})
	}
}

func (h *BotHandlers) handleComplete(c tele.Context) error {
	session := h.sessions.GetOrCreate(c.Sender().ID)

	build, err := buildPC(&session)
	if err != nil {
		log.Printf("Ошибка сборки конфигурации: %v", err)
		return c.EditOrSend("Не удалось подобрать сбалансированные комплектующие под выбранные критерии. Нажмите /reset для сброса.")
	}

	text := formatBuild(build, &session)
	return c.EditOrSend(text, h.keyboards.buildFinalKeyboard())
}

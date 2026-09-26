package main

import (
	tele "gopkg.in/telebot.v3"
)

// KeyboardBuilder создает динамическую клавиатуру для конкретного шага
type KeyboardBuilder struct {
	menu *tele.ReplyMarkup
}

func NewKeyboardBuilder() *KeyboardBuilder {
	return &KeyboardBuilder{
		menu: &tele.ReplyMarkup{},
	}
}

// Добавь эти методы в keyboards.go

func (kb *KeyboardBuilder) buildMonitorRefreshRateKeyboard(session *UserSession) *tele.ReplyMarkup {
	kb.menu = &tele.ReplyMarkup{}

	rates := []struct {
		label string
		data  string
	}{
		{"60–75 Гц", "refresh_60"},
		{"144–180 Гц", "refresh_144"},
		{"240+ Гц", "refresh_240"},
		{"🤷‍♂️ Без разницы", "refresh_auto"},
	}

	rows := make([]tele.Row, 0, len(rates))
	for _, rate := range rates {
		btn := kb.menu.Data(rate.label, rate.data)
		rows = append(rows, kb.menu.Row(btn))
	}

	kb.menu.Inline(rows...)
	return kb.menu
}

func (kb *KeyboardBuilder) buildCPUVendorKeyboard(session *UserSession) *tele.ReplyMarkup {
	kb.menu = &tele.ReplyMarkup{}

	vendors := []struct {
		label string
		data  string
	}{
		{"🔴 AMD", "cpu_amd"},
		{"🔵 Intel", "cpu_intel"},
		{"🤷‍♂️ Без разницы", "cpu_auto"},
	}

	rows := make([]tele.Row, 0, len(vendors))
	for _, vendor := range vendors {
		btn := kb.menu.Data(vendor.label, vendor.data)
		rows = append(rows, kb.menu.Row(btn))
	}

	kb.menu.Inline(rows...)
	return kb.menu
}

func (kb *KeyboardBuilder) buildGPUVendorKeyboard(session *UserSession) *tele.ReplyMarkup {
	kb.menu = &tele.ReplyMarkup{}

	vendors := []struct {
		label string
		data  string
	}{
		{"🟢 NVIDIA", "gpu_nvidia"},
		{"🔴 AMD", "gpu_amd"},
		{"🔵 Intel", "gpu_intel"},
		{"🤷‍♂️ Без разницы", "gpu_auto"},
	}

	rows := make([]tele.Row, 0, len(vendors))
	for _, vendor := range vendors {
		btn := kb.menu.Data(vendor.label, vendor.data)
		rows = append(rows, kb.menu.Row(btn))
	}

	kb.menu.Inline(rows...)
	return kb.menu
}

func (kb *KeyboardBuilder) buildStorageTypeKeyboard(session *UserSession) *tele.ReplyMarkup {
	kb.menu = &tele.ReplyMarkup{}

	types := []struct {
		label string
		data  string
	}{
		{"⚡ Только NVMe", "storage_nvme"},
		{"💽 Только SATA SSD", "storage_sata"},
		{"💾 Только HDD", "storage_hdd"},
		{"🔀 Комбо NVMe + HDD", "storage_combo"},
		{"🤷‍♂️ Без разницы", "storage_auto"},
	}

	rows := make([]tele.Row, 0, len(types))
	for _, storageType := range types {
		btn := kb.menu.Data(storageType.label, storageType.data)
		rows = append(rows, kb.menu.Row(btn))
	}

	kb.menu.Inline(rows...)
	return kb.menu
}

func (kb *KeyboardBuilder) buildColorKeyboard(session *UserSession) *tele.ReplyMarkup {
	kb.menu = &tele.ReplyMarkup{}

	colors := []struct {
		label string
		data  string
	}{
		{"🖤 Черный", "color_black"},
		{"🤍 Белый", "color_white"},
		{"🤷‍♂️ Без разницы", "color_auto"},
	}

	rows := make([]tele.Row, 0, len(colors))
	for _, color := range colors {
		btn := kb.menu.Data(color.label, color.data)
		rows = append(rows, kb.menu.Row(btn))
	}

	kb.menu.Inline(rows...)
	return kb.menu
}

func (kb *KeyboardBuilder) buildCoolingKeyboard(session *UserSession) *tele.ReplyMarkup {
	kb.menu = &tele.ReplyMarkup{}

	coolingOptions := []struct {
		label string
		data  string
	}{
		{"🌪 Воздушная башня", "cooling_air"},
		{"💧 Водянка (СЖО)", "cooling_water"},
		{"🤷‍♂️ Без разницы", "cooling_auto"},
	}

	rows := make([]tele.Row, 0, len(coolingOptions))
	for _, cooling := range coolingOptions {
		btn := kb.menu.Data(cooling.label, cooling.data)
		rows = append(rows, kb.menu.Row(btn))
	}

	kb.menu.Inline(rows...)
	return kb.menu
}

func (kb *KeyboardBuilder) buildCaseKeyboard(session *UserSession) *tele.ReplyMarkup {
	kb.menu = &tele.ReplyMarkup{}

	caseOptions := []struct {
		label string
		data  string
	}{
		{"📦 Mid-Tower (Сетка)", "case_mid"},
		{"🪟 Аквариум", "case_aquarium"},
		{"💼 Mini-ITX", "case_mini"},
		{"🤷‍♂️ Без разницы", "case_auto"},
	}

	rows := make([]tele.Row, 0, len(caseOptions))
	for _, caseOpt := range caseOptions {
		btn := kb.menu.Data(caseOpt.label, caseOpt.data)
		rows = append(rows, kb.menu.Row(btn))
	}

	kb.menu.Inline(rows...)
	return kb.menu
}

func (kb *KeyboardBuilder) buildUpgradeHorizonKeyboard(session *UserSession) *tele.ReplyMarkup {
	kb.menu = &tele.ReplyMarkup{}

	horizons := []struct {
		label string
		data  string
	}{
		{"🔒 На 5 лет и забыть", "horizon_5years"},
		{"🛠 С запасом под апгрейд", "horizon_upgrade"},
		{"🤷‍♂️ Без разницы", "horizon_auto"},
	}

	rows := make([]tele.Row, 0, len(horizons))
	for _, horizon := range horizons {
		btn := kb.menu.Data(horizon.label, horizon.data)
		rows = append(rows, kb.menu.Row(btn))
	}

	kb.menu.Inline(rows...)
	return kb.menu
}

func (kb *KeyboardBuilder) buildFinalKeyboard() *tele.ReplyMarkup {
	kb.menu = &tele.ReplyMarkup{}

	btnRestart := kb.menu.Data("🔄 Начать заново", "restart")
	kb.menu.Inline(kb.menu.Row(btnRestart))

	return kb.menu
}

// buildBudgetKeyboard — шаг 1
func (kb *KeyboardBuilder) buildBudgetKeyboard(session *UserSession) *tele.ReplyMarkup {
	kb.menu = &tele.ReplyMarkup{}

	budgets := []struct {
		label string
		data  string
		value BudgetRange
	}{
		{"💰 До 65 000 ₽", "budget_low", BudgetLow},
		{"💵 70k – 110k ₽", "budget_mid_low", BudgetMidLow},
		{"💎 120k – 170k ₽", "budget_mid_hi", BudgetMidHi},
		{"🚀 180k+ ₽", "budget_high", BudgetHigh},
	}

	rows := make([]tele.Row, 0, len(budgets))
	for _, b := range budgets {
		btn := kb.menu.Data(b.label, b.data)
		rows = append(rows, kb.menu.Row(btn))
	}

	kb.menu.Inline(rows...)
	return kb.menu
}

// buildUseCaseKeyboard — шаг 2
func (kb *KeyboardBuilder) buildUseCaseKeyboard(session *UserSession) *tele.ReplyMarkup {
	kb.menu = &tele.ReplyMarkup{}

	useCases := []struct {
		label string
		data  string
		value UseCase
	}{
		{"🎮 Гейминг", "usecase_gaming", UseCaseGaming},
		{"🎙 Стриминг", "usecase_streaming", UseCaseStreaming},
		{"💻 Работа (3D/Код)", "usecase_work", UseCaseWork},
		{"📄 Офис", "usecase_office", UseCaseOffice},
		{"⚖️ Универсал", "usecase_universal", UseCaseUniversal},
	}

	rows := make([]tele.Row, 0, len(useCases))
	for _, uc := range useCases {
		btn := kb.menu.Data(uc.label, uc.data)
		rows = append(rows, kb.menu.Row(btn))
	}

	kb.menu.Inline(rows...)
	return kb.menu
}

// buildGameTypeKeyboard — шаг 3
func (kb *KeyboardBuilder) buildGameTypeKeyboard(session *UserSession) *tele.ReplyMarkup {
	kb.menu = &tele.ReplyMarkup{}

	// Если офис — пропускаем этот шаг
	if session.Default.UseCase == UseCaseOffice {
		return nil
	}

	games := []struct {
		label string
		data  string
	}{
		{"🎯 Киберспорт", "game_cybersport"},
		{"🌄 AAA-игры", "game_aaa"},
		{"🎲 Всё подряд", "game_all"},
	}

	rows := make([]tele.Row, 0, len(games))
	for _, g := range games {
		btn := kb.menu.Data(g.label, g.data)
		rows = append(rows, kb.menu.Row(btn))
	}

	kb.menu.Inline(rows...)
	return kb.menu
}

// buildMonitorResolutionKeyboard — шаг 4.1
func (kb *KeyboardBuilder) buildMonitorResolutionKeyboard(session *UserSession) *tele.ReplyMarkup {
	kb.menu = &tele.ReplyMarkup{}

	resolutions := []struct {
		label string
		data  string
	}{
		{"📺 1080p", "res_1080p"},
		{"🖥 1440p (2K)", "res_1440p"},
		{"🎮 4K", "res_4k"},
		{"🤷‍♂️ Без разницы", "res_auto"},
	}

	rows := make([]tele.Row, 0, len(resolutions))
	for _, r := range resolutions {
		// UI Pruning: для низкого бюджета скрываем 4K
		if r.data == "res_4k" && session.Default.Budget == BudgetLow {
			continue
		}
		btn := kb.menu.Data(r.label, r.data)
		rows = append(rows, kb.menu.Row(btn))
	}

	kb.menu.Inline(rows...)
	return kb.menu
}

// buildRAMSizeKeyboard — шаг 6 с UI Pruning
func (kb *KeyboardBuilder) buildRAMSizeKeyboard(session *UserSession) *tele.ReplyMarkup {
	kb.menu = &tele.ReplyMarkup{}

	ramOptions := []struct {
		label string
		data  string
		size  int
	}{
		{"4 GB", "ram_4", 4},
		{"8 GB", "ram_8", 8},
		{"16 GB", "ram_16", 16},
		{"32 GB", "ram_32", 32},
		{"64 GB", "ram_64", 64},
		{"🤷‍♂️ Без разницы", "ram_auto", 0},
	}

	rows := make([]tele.Row, 0, len(ramOptions))
	for _, ram := range ramOptions {
		// UI Pruning
		if ram.size > 0 {
			// Для 4K монитора или высокого бюджета скрываем 4-8 GB
			if (session.User.MonitorResolution == "4k" ||
				session.Default.Budget == BudgetHigh ||
				session.Default.Budget == BudgetMidHi) && ram.size < 16 {
				continue
			}
			// Для офиса скрываем 64 GB
			if session.Default.UseCase == UseCaseOffice && ram.size >= 64 {
				continue
			}
		}

		btn := kb.menu.Data(ram.label, ram.data)
		rows = append(rows, kb.menu.Row(btn))
	}

	kb.menu.Inline(rows...)
	return kb.menu
}

// buildStorageSizeKeyboard — шаг 7.2 с UI Pruning
func (kb *KeyboardBuilder) buildStorageSizeKeyboard(session *UserSession) *tele.ReplyMarkup {
	kb.menu = &tele.ReplyMarkup{}

	storageOptions := []struct {
		label string
		data  string
		size  int
	}{
		{"512 GB", "storage_512", 512},
		{"1 TB", "storage_1tb", 1024},
		{"2 TB", "storage_2tb", 2048},
		{"4 TB", "storage_4tb", 4096},
		{"🤷‍♂️ Без разницы", "storage_auto", 0},
	}

	rows := make([]tele.Row, 0, len(storageOptions))
	for _, storage := range storageOptions {
		// UI Pruning
		if storage.size > 0 {
			// Для низкого бюджета скрываем 4 TB
			if (session.Default.Budget == BudgetLow ||
				session.Default.Budget == BudgetMidLow) && storage.size >= 4096 {
				continue
			}
			// Для высокого бюджета скрываем 512 GB
			if (session.Default.Budget == BudgetHigh ||
				session.Default.Budget == BudgetMidHi) && storage.size <= 512 {
				continue
			}
		}

		btn := kb.menu.Data(storage.label, storage.data)
		rows = append(rows, kb.menu.Row(btn))
	}

	kb.menu.Inline(rows...)
	return kb.menu
}

// ... аналогичные методы для остальных шагов

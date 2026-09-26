package main

import (
	"database/sql"
	"fmt"
	"math"
)

func buildPC(session *UserSession) (*BuildResult, error) {
	build := &BuildResult{}

	cpu, err := selectCPU(session)
	if err != nil {
		return nil, fmt.Errorf("ошибка подбора CPU: %w", err)
	}
	build.CPU = cpu

	gpu, err := selectGPU(session)
	if err != nil {
		return nil, fmt.Errorf("ошибка подбора GPU: %w", err)
	}
	build.GPU = gpu

	mb, err := selectMotherboard(session, cpu)
	if err != nil {
		return nil, fmt.Errorf("ошибка подбора материнской платы: %w", err)
	}
	build.Motherboard = mb

	ram, err := selectRAM(session, mb)
	if err != nil {
		return nil, fmt.Errorf("ошибка подбора RAM: %w", err)
	}
	build.RAM = ram

	storage, err := selectStorage(session)
	if err != nil {
		return nil, fmt.Errorf("ошибка подбора накопителя: %w", err)
	}
	build.Storage = storage

	pcCase, err := selectCase(session)
	if err != nil {
		return nil, fmt.Errorf("ошибка подбора корпуса: %w", err)
	}
	build.Case = pcCase

	cooling, err := selectCooling(session, cpu)
	if err != nil {
		return nil, fmt.Errorf("ошибка подбора охлаждения: %w", err)
	}
	build.Cooling = cooling

	psu, err := selectPSU(session, cpu, gpu)
	if err != nil {
		return nil, fmt.Errorf("ошибка подбора БП: %w", err)
	}
	build.PSU = psu

	build.TotalPrice = cpu.Price + gpu.Price + mb.Price + ram.Price +
		storage.Price + pcCase.Price + cooling.Price + psu.Price

	return build, nil
}

func selectCPU(session *UserSession) (Component, error) {
	vendor := session.User.CPUVendor
	if vendor == "" {
		vendor = session.Default.CPUVendor
	}

	minPerf := 5
	switch session.User.Budget {
	case BudgetLow:
		minPerf = 4
	case BudgetMidLow:
		minPerf = 6
	case BudgetMidHi:
		minPerf = 7
	case BudgetHigh:
		minPerf = 8
	}

	if session.User.UseCase == UseCaseWork || session.User.UseCase == UseCaseStreaming {
		minPerf++
	}

	// Для офисного ПК жестко требуем наличие встроенной графики
	query := `
  SELECT id, name, price, socket, tdp, has_igpu,
         printf('%d ядер, %dW TDP', cores, tdp) as specs
  FROM cpu 
  WHERE vendor = ? AND performance >= ?
 `
	if session.User.UseCase == UseCaseOffice {
		query += " AND has_igpu = 1 "
	}
	query += " ORDER BY price ASC LIMIT 1 "

	var comp Component
	var hasIGPUInt int
	err := db.QueryRow(query, vendor, minPerf).Scan(
		&comp.ID, &comp.Name, &comp.Price, &comp.Socket, &comp.TDP, &hasIGPUInt, &comp.Specs,
	)
	if err != nil {
		return Component{}, err
	}
	comp.HasIGPU = (hasIGPUInt == 1)
	return comp, nil
}

func selectGPU(session *UserSession) (Component, error) {
	if session.User.UseCase == UseCaseOffice {
		return Component{
			Name:    "Встроенная графика (iGPU)",
			Price:   0,
			Specs:   "Интегрированное в процессор видеоядро",
			TDP:     0,
			HasIGPU: true,
		}, nil
	}

	vendor := session.User.GPUVendor
	if vendor == "" {
		vendor = session.Default.GPUVendor
	}

	minPerf := 5
	switch session.User.MonitorResolution {
	case "1080p":
		minPerf = 4
	case "1440p":
		minPerf = 6
	case "4k":
		minPerf = 9
	}

	switch session.User.Budget {
	case BudgetLow:
		if minPerf > 5 {
			minPerf = 5
		}
	case BudgetMidLow:
		if minPerf > 7 {
			minPerf = 7
		}
	}

	query := `
  SELECT id, name, price, tdp,
         printf('%dGB VRAM, %dW TDP', vram, tdp) as specs
  FROM gpu 
  WHERE vendor = ? AND performance >= ?
  ORDER BY price ASC
  LIMIT 1
 `
	var comp Component
	err := db.QueryRow(query, vendor, minPerf).Scan(
		&comp.ID, &comp.Name, &comp.Price, &comp.TDP, &comp.Specs,
	)
	return comp, err
}

func selectMotherboard(session *UserSession, cpu Component) (Component, error) {
	caseForm := session.User.Case
	if caseForm == "" {
		caseForm = session.Default.Case
	}

	var query string
	if caseForm == "mini-itx" {
		query = `
   SELECT id, name, price, ram_type, form_factor,
          printf('Сокет %s, %s, %s', socket, chipset, ram_type) as specs
   FROM motherboard 
   WHERE socket = ? AND form_factor = 'ITX'
   ORDER BY price ASC LIMIT 1`
	} else {
		query = `
   SELECT id, name, price, ram_type, form_factor


,
          printf('Сокет %s, %s, %s', socket, chipset, ram_type) as specs
   FROM motherboard 
   WHERE socket = ? AND form_factor IN ('ATX', 'mATX', 'ITX')
   ORDER BY price ASC LIMIT 1`
	}

	rows, err := db.Query(query, cpu.Socket)
	if err != nil {
		return Component{}, err
	}
	defer rows.Close()

	if !rows.Next() {
		return Component{}, fmt.Errorf("нет совместимой материнской платы под сокет %s", cpu.Socket)
	}

	var comp Component
	err = rows.Scan(&comp.ID, &comp.Name, &comp.Price, &comp.RAMType, &comp.FormFactor, &comp.Specs)
	return comp, err
}

func selectRAM(session *UserSession, mb Component) (Component, error) {
	ramSize := session.User.RAMSize
	if ramSize == "" {
		ramSize = session.Default.RAMSize
	}

	capacity := 16
	switch ramSize {
	case "4":
		capacity = 4
	case "8":
		capacity = 8
	case "16":
		capacity = 16
	case "32":
		capacity = 32
	case "64":
		capacity = 64
	}

	query := `
  SELECT id, name, price, 
         printf('%dGB %s %dMHz', capacity, type, speed) as specs
  FROM ram 
  WHERE type = ? AND capacity = ?
  ORDER BY price ASC
  LIMIT 1
 `
	var comp Component
	err := db.QueryRow(query, mb.RAMType, capacity).Scan(
		&comp.ID, &comp.Name, &comp.Price, &comp.Specs,
	)
	return comp, err
}

func selectStorage(session *UserSession) (Component, error) {
	storageType := session.User.StorageType
	if storageType == "" {
		storageType = session.Default.StorageType
	}

	storageSize := session.User.StorageSize
	if storageSize == "" {
		storageSize = session.Default.StorageSize
	}

	capacity := 1024
	switch storageSize {
	case "512":
		capacity = 512
	case "1":
		capacity = 1024
	case "2":
		capacity = 2048
	case "4":
		capacity = 4096
	}

	var dbType string
	switch storageType {
	case "sata":
		dbType = "SATA SSD"
	case "hdd":
		dbType = "HDD"
	case "combo":
		return selectComboStorage(capacity)
	default:
		dbType = "NVMe SSD"
	}

	query := `
  SELECT id, name, price, 
         printf('%s %dGB', type, capacity) as specs
  FROM storage 
  WHERE type = ? AND capacity = ?
  ORDER BY price ASC
  LIMIT 1
 `
	var comp Component
	err := db.QueryRow(query, dbType, capacity).Scan(
		&comp.ID, &comp.Name, &comp.Price, &comp.Specs,
	)
	return comp, err
}

func selectComboStorage(nvmeCapacity int) (Component, error) {
	var nvmeName string
	var nvmePrice int
	err := db.QueryRow(`
  SELECT name, price FROM storage 
  WHERE type = 'NVMe SSD' AND capacity = ?
  ORDER BY price ASC LIMIT 1
 `, nvmeCapacity).Scan(&nvmeName, &nvmePrice)
	if err != nil {
		return Component{}, err
	}

	var hddName string
	var hddPrice, hddCapacity int
	err = db.QueryRow(`
  SELECT name, price, capacity FROM storage 
  WHERE type = 'HDD' 
  ORDER BY capacity DESC, price ASC LIMIT 1
 `).Scan(&hddName, &hddPrice, &hddCapacity)
	if err != nil {
		return Component{}, err
	}

	return Component{
		Name:  fmt.Sprintf("%s + %s", nvmeName, hddName),
		Price: nvmePrice + hddPrice,
		Specs: fmt.Sprintf("NVMe %dGB + HDD %dGB", nvmeCapacity, hddCapacity),
	}, nil
}

func selectCase(session *UserSession) (Component, error) {
	caseType := session.User.Case
	if caseType == "" {
		caseType = session.Default.Case
	}

	color := session.User.Color
	if color == "" {
		color = session.Default.Color
	}

	var comp Component
	query := `
  SELECT id, name, price, form_factor,
         printf('%s, %s', form_factor, color) as specs
  FROM pc_case 
  WHERE form_factor = ?
  ORDER BY (color = ?) DESC, price ASC
  LIMIT 1
 `
	err := db.QueryRow(query, caseType, color).Scan(
		&comp.ID, &comp.Name, &comp.Price, &comp.FormFactor, &comp.Specs,
	)
	return comp, err
}

func selectCooling(session *UserSession, cpu Component) (Component, error) {
	coolingType := session.User.Cooling
	if coolingType == "" {
		coolingType = session.Default.Cooling
	}

	requiredTDP := int(float64(cpu.TDP) * 1.3)

	query := `
  SELECT id, name, price, 
         printf('%s, %dW TDP', type, tdp) as specs
  FROM cooling 
  WHERE type = ? AND tdp >= ?
  ORDER BY price ASC
  LIMIT 1
 `
	var comp Component
	err := db.QueryRow(query, coolingType, requiredTDP).Scan(
		&comp.ID, &comp.Name, &comp.Price, &comp.Specs,
	)
	if err == sql.ErrNoRows && coolingType == "air" {
		err = db.QueryRow(query, "water", requiredTDP).Scan(
			&comp.ID, &comp.Name, &comp.Price, &comp.Specs,
		)
	}
	return comp, err
}

func selectPSU(session *UserSession, cpu Component, gpu Component) (Component, error) {
	systemBase := 75
	totalEstimated := cpu.TDP + gpu.TDP + systemBase

	margin := 1.25
	if gpu.TDP >= 250 {
		margin = 1.35
	}

	targetPower := int(math.Ceil(float64(totalEstimated)*margin/50.0)) * 50

	query := `
  SELECT id, name, price, power,
         printf('%dW', power) as specs
  FROM psu 
  WHERE power >= ?
  ORDER BY power ASC
  LIMIT 1
 `
	var comp Component
	err := db.QueryRow(query, targetPower).Scan(
		&comp.ID, &comp.Name, &comp.Price, &comp.Power, &comp.Specs,
	)
	return comp, err
}

func formatBuild(build *BuildResult, session *UserSession) string {
	text := "🖥 ВАША СБОРКА\n\n"

	text += fmt.Sprintf("🔲 Процессор: %s\n", build.CPU.Name)
	text += fmt.Sprintf("   %s\n", build.CPU.Specs)
	text += fmt.Sprintf("   Цена: %d ₽\n\n", build.CPU.Price)

	if build.GPU.Price > 0 {
		text += fmt.Sprintf("🎮 Видеокарта: %s\n", build.GPU.Name)
		text += fmt.Sprintf("   %s\n", build.GPU.Specs)
		text += fmt.Sprintf("   Цена: %d ₽\n\n", build.GPU.Price)
	} else {
		text += fmt.Sprintf("🎮 Видеокарта: %s\n", build.GPU.Name)
		text += fmt.Sprintf("   %s\n\n", build.GPU.Specs)
	}

	text += fmt.Sprintf("🔧 Материнская плата: %s\n", build.Motherboard.Name)
	text += fmt.Sprintf("   %s\n", build.Motherboard.Specs)
	text += fmt.Sprintf("   Цена: %d ₽\n\n", build.Motherboard.Price)

	text += fmt.Sprintf("🧠 Оперативная память: %s\n", build.RAM.Name)
	text += fmt.Sprintf("   %s\n", build.RAM.Specs)
	text += fmt.Sprintf("   Цена: %d ₽\n\n", build.RAM.Price)

	text += fmt.Sprintf("💾 Накопитель: %s\n", build.Storage.Name)
	text += fmt.Sprintf("   %s\n", build.Storage.Specs)
	text += fmt.Sprintf("   Цена: %d ₽\n\n", build.Storage.Price)

	text += fmt.Sprintf("📦 Корпус: %s\n", build.Case.Name)
	text += fmt.Sprintf("   %s\n", build.Case.Specs)
	text += fmt.Sprintf("   Цена: %d ₽\n\n", build.Case.Price)

	text += fmt.Sprintf("❄️ Охлаждение: %s\n", build.Cooling.Name)
	text += fmt.Sprintf("   %s\n", build.Cooling.Specs)
	text += fmt.Sprintf("   Цена: %d ₽\n\n", build.Cooling.Price)

	text += fmt.Sprintf("⚡ Блок питания: %s\n", build.PSU.Name)
	text += fmt.Sprintf("   %s\n", build.PSU.Specs)
	text += fmt.Sprintf("   Цена: %d ₽\n\n", build.PSU.Price)

	text += fmt.Sprintf("💰 ИТОГО: %d ₽", build.TotalPrice)

	return text
}

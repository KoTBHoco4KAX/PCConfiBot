package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func initDatabase() {
	var err error
	db, err = sql.Open("sqlite3", "./pc_builder.db")
	if err != nil {
		log.Fatalf("Ошибка открытия базы данных: %v", err)
	}

	createTables()
	seedComponents()
}

func createTables() {
	schemas := []string{
		`CREATE TABLE IF NOT EXISTS cpu (
   id INTEGER PRIMARY KEY AUTOINCREMENT,
   name TEXT NOT NULL,
   vendor TEXT NOT NULL,
   cores INTEGER NOT NULL,
   price INTEGER NOT NULL,
   performance INTEGER NOT NULL,
   socket TEXT NOT NULL,
   tdp INTEGER NOT NULL,
   has_igpu INTEGER NOT NULL DEFAULT 0
  );`,
		`CREATE TABLE IF NOT EXISTS gpu (
   id INTEGER PRIMARY KEY AUTOINCREMENT,
   name TEXT NOT NULL,
   vendor TEXT NOT NULL,
   vram INTEGER NOT NULL,
   price INTEGER NOT NULL,
   performance INTEGER NOT NULL,
   tdp INTEGER NOT NULL
  );`,
		`CREATE TABLE IF NOT EXISTS motherboard (
   id INTEGER PRIMARY KEY AUTOINCREMENT,
   name TEXT NOT NULL,
   socket TEXT NOT NULL,
   chipset TEXT,
   form_factor TEXT NOT NULL,
   ram_type TEXT NOT NULL,
   price INTEGER NOT NULL
  );`,
		`CREATE TABLE IF NOT EXISTS ram (
   id INTEGER PRIMARY KEY AUTOINCREMENT,
   name TEXT NOT NULL,
   capacity INTEGER NOT NULL,
   type TEXT NOT NULL,
   speed INTEGER NOT NULL,
   price INTEGER NOT NULL
  );`,
		`CREATE TABLE IF NOT EXISTS storage (
   id INTEGER PRIMARY KEY AUTOINCREMENT,
   name TEXT NOT NULL,
   type TEXT NOT NULL,
   capacity INTEGER NOT NULL,
   price INTEGER NOT NULL
  );`,
		`CREATE TABLE IF NOT EXISTS pc_case (
   id INTEGER PRIMARY KEY AUTOINCREMENT,
   name TEXT NOT NULL,
   form_factor TEXT NOT NULL,
   color TEXT NOT NULL,
   price INTEGER NOT NULL
  );`,
		`CREATE TABLE IF NOT EXISTS cooling (
   id INTEGER PRIMARY KEY AUTOINCREMENT,
   name TEXT NOT NULL,
   type TEXT NOT NULL,
   tdp INTEGER NOT NULL,
   price INTEGER NOT NULL
  );`,
		`CREATE TABLE IF NOT EXISTS psu (
   id INTEGER PRIMARY KEY AUTOINCREMENT,
   name TEXT NOT NULL,
   power INTEGER NOT NULL,
   price INTEGER NOT NULL
  );`,
	}

	for _, schema := range schemas {
		if _, err := db.Exec(schema); err != nil {
			log.Fatalf("Ошибка применения схемы: %v", err)
		}
	}
}

func seedComponents() {
	var count int
	if err := db.QueryRow("SELECT COUNT(*) FROM cpu").Scan(&count); err != nil || count > 0 {
		return
	}

	// =========================
	// 1. ПРОЦЕССОРЫ
	// =========================
	cpus := []struct {
		name, vendor, socket             string
		cores, price, perf, tdp, hasIGPU int
	}{
		// AM4 — APU / офис / бюджет
		{"AMD Ryzen 3 3200G", "amd", "AM4", 4, 4500, 3, 65, 1},
		{"AMD Ryzen 5 3400G", "amd", "AM4", 4, 6500, 4, 65, 1},
		{"AMD Ryzen 5 4600G", "amd", "AM4", 6, 8500, 5, 65, 1},
		{"AMD Ryzen 5 5600G", "amd", "AM4", 6, 11500, 6, 65, 1},
		{"AMD Ryzen 7 5700G", "amd", "AM4", 8, 14500, 7, 65, 1},
		// AM4 — чистые камни
		{"AMD Ryzen 5 3600", "amd", "AM4", 6, 7500, 5, 65, 0},
		{"AMD Ryzen 5 5600", "amd", "AM4", 6, 10500, 6, 65, 0},
		{"AMD Ryzen 5 5600X", "amd", "AM4", 6, 12500, 6, 65, 0},
		{"AMD Ryzen 7 3700X", "amd", "AM4", 8, 9500, 6, 65, 0},
		{"AMD Ryzen 7 5700X", "amd", "AM4", 8, 14500, 7, 65, 0},
		{"AMD Ryzen 7 5700X3D", "amd", "AM4", 8, 19500, 8, 105, 0},
		{"AMD Ryzen 7 5800X3D", "amd", "AM4", 8, 24500, 9, 105, 0},
		{"AMD Ryzen 9 5900X", "amd", "AM4", 12, 26500, 9, 105, 0},
		{"AMD Ryzen 9 5950X", "amd", "AM4", 16, 38500, 10, 105, 0},
		// AM5 — APU / базовые
		{"AMD Ryzen 5 8500G", "amd", "AM5", 6, 13500, 6, 65, 1},
		{"AMD Ryzen 5 8600G", "amd", "AM5", 6, 17500, 7, 65, 1},
		{"AMD Ryzen 7 8700G", "amd", "AM5", 8, 25500, 8, 65, 1},
		{"AMD Ryzen 5 7500F", "amd", "AM5", 6, 13500, 7, 65, 0},
		{"AMD Ryzen 5 7600", "amd", "AM5", 6, 17000, 7, 65, 1},
		{"AMD Ryzen 5 7600X", "amd", "AM5", 6, 19500, 8, 105, 1},
		{"AMD Ryzen 7 7700", "amd", "AM5", 8, 23500, 8, 65, 1},
		{"AMD Ryzen 7 7700X", "amd", "AM5", 8, 28500, 8, 105, 1},
		{"AMD Ryzen 7 7800X3D", "amd", "AM5", 8, 39000, 10, 120, 1},
		{"AMD Ryzen 9 7900", "amd", "AM5", 12, 34500, 9, 65, 1},
		{"AMD Ryzen 9 7900X", "amd", "AM5", 12, 37500, 9, 170, 1},
		{"AMD Ryzen 9 7950X", "amd", "AM5", 16, 55500, 10, 170, 1},
		{"AMD Ryzen 9 7950X3D", "amd", "AM5", 16, 59000, 10, 120, 1},
		// AM5 — Ryzen 9000 (Zen 5)
		{"AMD Ryzen 5 9600X", "amd", "AM5", 6, 23500, 8, 65, 1},
		{"AMD Ryzen 7 9700X", "amd", "AM5", 8, 31500, 8, 65, 1},
		{"AMD Ryzen 7 9800X3D", "amd", "AM5", 8, 49500, 10, 120, 1},
		{"AMD Ryzen 9 9900X", "amd", "AM5", 12, 43500, 9, 120, 1},
		{"AMD Ryzen 9 9950X", "amd", "AM5", 16, 59500, 10, 170, 1},
		{"AMD Ryzen 9 9950X3D", "amd", "AM5", 16, 69900, 10, 170, 1},
		// LGA1700 — младшие Intel
		{"Intel Core i3-12100", "intel", "LGA1700", 4, 8500, 4, 60, 1},
		{"Intel Core i3-12100F", "intel", "LGA1700", 4, 7200, 4, 58, 0},
		{"Intel Core i3-13100", "intel", "LGA1700", 4, 9500, 4, 60, 1},
		{"Intel Core i3-13100F", "intel", "LGA1700", 4, 9500, 4, 58, 0},
		{"Intel Core i5-12400", "intel", "LGA1700", 6, 13500, 6, 65, 1},
		{"Intel Core i5-12400F", "intel", "LGA1700", 6, 11500, 6, 65, 0},
		{"Intel Core i5-13400", "intel", "LGA1700", 10, 17500, 7, 65, 1},
		{"Intel Core i5-13400F", "intel", "LGA1700", 10, 16500, 7, 65, 0},
		{"Intel Core i5-14400F", "intel", "LGA1700", 10, 16500, 7, 65, 0},
		{"Intel Core i5-13600K", "intel", "LGA1700", 14, 27500, 8, 125, 1},
		{"Intel Core i5-13600KF", "intel", "LGA1700", 14, 25500, 8, 125, 0},
		{"Intel Core i5-14600K", "intel", "LGA1700", 14, 29500, 8, 125, 1},
		{"Intel Core i5-14600KF", "intel", "LGA1700", 14, 27500, 8, 125, 0},
		// LGA1700 — старшие Intel
		{"Intel Core i7-12700K", "intel", "LGA1700", 12, 24500, 8, 190, 1},
		{"Intel Core i7-13700K", "intel", "LGA1700", 16, 35500, 9, 253, 1},
		{"Intel Core i7-13700KF", "intel", "LGA1700", 16, 33500, 9, 253, 0},
		{"Intel Core i7-14700K", "intel", "LGA1700", 20, 42500, 9, 253, 1},
		{"Intel Core i7-14700KF", "intel", "LGA1700", 20, 39500, 9, 253, 0},
		{"Intel Core i9-12900K", "intel", "LGA1700", 16, 32500, 9, 241, 1},
		{"Intel Core i9-13900K", "intel", "LGA1700", 24, 52500, 10, 253, 1},
		{"Intel Core i9-13900KF", "intel", "LGA1700", 24, 49500, 10, 253, 0},
		{"Intel Core i9-14900K", "intel", "LGA1700", 24, 57000, 10, 253, 1},
		{"Intel Core i9-14900KF", "intel", "LGA1700", 24, 49500, 10, 253, 0},
		{"Intel Core i9-14900KS", "intel", "LGA1700", 24, 64500, 10, 150, 1},
	}

	for _, c := range cpus {
		db.Exec("INSERT INTO cpu (name, vendor, socket, cores, price, performance, tdp, has_igpu) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
			c.name, c.vendor, c.socket, c.cores, c.price, c.perf, c.tdp, c.hasIGPU)
	}

	// =========================
	// 2. ВИДЕОКАРТЫ
	// =========================
	gpus := []struct {
		name, vendor           string
		vram, price, perf, tdp int
	}{
		// Бюджет
		{"NVIDIA GTX 1650", "nvidia", 4, 12500, 3, 75},
		{"AMD Radeon RX 6400", "amd", 4, 12500, 3, 53},
		{"Intel Arc A380", "intel", 6, 13500, 4, 75},
		{"AMD Radeon RX 6500 XT", "amd", 4, 15500, 4, 107},
		{"NVIDIA GTX 1660 Super", "nvidia", 6, 18500, 4, 125},
		{"Intel Arc A580", "intel", 8, 19500, 5, 185},
		{"AMD Radeon RX 6600", "amd", 8, 23500, 5, 132},
		{"NVIDIA RTX 3050 8GB", "nvidia", 8, 22500, 4, 130},
		{"Intel Arc A750", "intel", 8, 23500, 6, 225},
		// Средний сегмент
		{"NVIDIA RTX 3060 12GB", "nvidia", 12, 27500, 5, 170},
		{"AMD Radeon RX 6650 XT", "amd", 8, 26500, 6, 180},
		{"AMD Radeon RX 7600", "amd", 8, 28500, 6, 165},
		{"NVIDIA RTX 4060", "nvidia", 8, 32000, 6, 115},
		{"AMD Radeon RX 7600 XT", "amd", 16, 33500, 6, 190},
		{"NVIDIA RTX 3060 Ti", "nvidia", 8, 32500, 6, 200},
		{"NVIDIA RTX 4060 Ti 8GB", "nvidia", 8, 38500, 7, 160},
		{"Intel Arc B580", "intel", 12, 29500, 7, 190},
		// Верхний средний
		{"AMD Radeon RX 6700 XT", "amd", 12, 34500, 7, 230},
		{"NVIDIA RTX 3070", "nvidia", 8, 38500, 7, 220},
		{"NVIDIA RTX 4060 Ti 16GB", "nvidia", 16, 45000, 7, 165},
		{"AMD Radeon RX 7700 XT", "amd", 12, 46500, 7, 245},
		{"NVIDIA RTX 3070 Ti", "nvidia", 8, 44500, 7, 290},
		{"AMD Radeon RX 6800", "amd", 16, 42500, 7, 250},
		// Высокий сегмент
		{"AMD Radeon RX 7800 XT", "amd", 16, 54500, 8, 263},
		{"NVIDIA RTX 3080 10GB", "nvidia", 10, 52500, 8, 320},
		{"NVIDIA RTX 4070", "nvidia", 12, 54500, 8, 200},
		{"NVIDIA RTX 4070 Super", "nvidia", 12, 64500, 8, 220},
		{"AMD Radeon RX 6900 XT", "amd", 16, 58500, 8, 300},
		{"AMD Radeon RX 7900 GRE", "amd", 16, 61500, 8, 260},
		{"NVIDIA RTX 3080 Ti", "nvidia", 12, 62500, 8, 350},
		{"NVIDIA RTX 4070 Ti", "nvidia", 12, 72500, 8, 285},
		{"AMD Radeon RX 6950 XT", "amd", 16, 68500, 9, 335},
		// Топ
		{"NVIDIA RTX 3090", "nvidia", 24, 78500, 9, 350},
		{"AMD Radeon RX 7900 XT", "amd", 20, 82500, 9, 315},
		{"NVIDIA RTX 4070 Ti Super", "nvidia", 16, 86000, 9, 285},
		{"NVIDIA RTX 3090 Ti", "nvidia", 24, 92500, 9, 450},
		{"AMD Radeon RX 7900 XTX", "amd", 24, 104000, 10, 355},
		{"NVIDIA RTX 4080", "nvidia", 16, 105000, 9, 320},
		{"NVIDIA RTX 4080 Super", "nvidia", 16, 112000, 9, 320},
		{"NVIDIA RTX 4090", "nvidia", 24, 215000, 10, 450},
		{"NVIDIA RTX 5090", "nvidia", 32, 285000, 10, 575},
		{"NVIDIA RTX 5080", "nvidia", 16, 145000, 10, 360},
	}

	for _, g := range gpus {
		db.Exec("INSERT INTO gpu (name, vendor, vram, price, performance, tdp) VALUES (?, ?, ?, ?, ?, ?)",
			g.name, g.vendor, g.vram, g.price, g.perf, g.tdp)
	}

	// =========================
	// 3. МАТЕРИНСКИЕ ПЛАТЫ
	// =========================
	mbs := []struct {
		name, socket, chipset, formFactor, ramType string
		price                                      int
	}{
		// AM4
		{"ASUS PRIME A320M-K", "AM4", "A320", "mATX", "DDR4", 4500},
		{"ASUS PRIME B450M-K II", "AM4", "B450", "mATX", "DDR4", 6200},
		{"MSI B450 TOMAHAWK MAX II", "AM4", "B450", "ATX", "DDR4", 8500},
		{"MSI B550M PRO-VDH WIFI", "AM4", "B550", "mATX", "DDR4", 10500},
		{"Gigabyte B550 AORUS ELITE V2", "AM4", "B550", "ATX", "DDR4", 13500},
		{"ASRock B550M-ITX/ac", "AM4", "B550", "ITX", "DDR4", 14200},
		{"ASUS ROG STRIX X570-E GAMING WIFI II", "AM4", "X570", "ATX", "DDR4", 27500},
		{"Gigabyte X570 AORUS MASTER", "AM4", "X570", "ATX", "DDR4", 25500},
		// AM5
		{"MSI PRO A620M-E", "AM5", "A620", "mATX", "DDR5", 8500},
		{"ASRock B650M-HDV/M.2", "AM5", "B650", "mATX", "DDR5", 12500},
		{"MSI PRO B650-S WIFI", "AM5", "B650", "ATX", "DDR5", 16500},
		{"Gigabyte B650 AORUS ELITE AX", "AM5", "B650", "ATX", "DDR5", 22000},
		{"ASRock B650I Lightning WiFi", "AM5", "B650", "ITX", "DDR5", 21500},
		{"ASUS ROG STRIX B650E-F GAMING WIFI", "AM5", "B650E", "ATX", "DDR5", 28500},
		{"MSI MAG B650 TOMAHAWK WIFI", "AM5", "B650", "ATX", "DDR5", 19500},
		{"MSI MPG X670E CARBON WIFI", "AM5", "X670E", "ATX", "DDR5", 42500},
		{"Gigabyte X670E AORUS MASTER", "AM5", "X670E", "ATX", "DDR5", 45500},
		{"ASUS ROG CROSSHAIR X870E HERO", "AM5", "X870E", "ATX", "DDR5", 55500},
		{"MSI MAG X870E TOMAHAWK WIFI", "AM5", "X870E", "ATX", "DDR5", 32500},
		{"ASRock B850M Pro RS WiFi", "AM5", "B850", "mATX", "DDR5", 15500},
		{"SAPPHIRE PULSE B850M", "AM5", "B850", "mATX", "DDR5", 12500},
		// LGA1700
		{"MSI PRO H610M-E DDR4", "LGA1700", "H610", "mATX", "DDR4", 6800},
		{"ASRock B760M Pro RS/D4", "LGA1700", "B760", "mATX", "DDR4", 10500},
		{"MSI PRO B760M-A WIFI DDR4", "LGA1700", "B760", "mATX", "DDR4", 12500},
		{"MSI B760 GAMING PLUS WIFI", "LGA1700", "B760", "ATX", "DDR5", 16500},
		{"ASUS ROG STRIX B760-I GAMING WIFI", "LGA1700", "B760", "ITX", "DDR5", 23500},
		{"ASUS TUF GAMING Z790-PLUS WIFI", "LGA1700", "Z790", "ATX", "DDR5", 27000},
		{"ASUS ROG STRIX Z790-E GAMING WIFI", "LGA1700", "Z790", "ATX", "DDR5", 43500},
		{"MSI MAG Z790 TOMAHAWK WIFI", "LGA1700", "Z790", "ATX", "DDR5", 28500},
		{"Gigabyte Z790 AORUS ELITE AX", "LGA1700", "Z790", "ATX", "DDR5", 24500},
	}

	for _, m := range mbs {
		db.Exec("INSERT INTO motherboard (name, socket, chipset, form_factor, ram_type, price) VALUES (?, ?, ?, ?, ?, ?)",
			m.name, m.socket, m.chipset, m.formFactor, m.ramType, m.price)
	}

	// =========================
	// 4. ОПЕРАТИВНАЯ ПАМЯТЬ
	// =========================
	rams := []struct {
		name, ramType          string
		capacity, speed, price int
	}{
		// DDR4
		{"Crucial 8GB (1x8) 2666", "DDR4", 8, 2666, 1800},
		{"Crucial 16GB (1x16) 3200", "DDR4", 16, 3200, 3200},
		{"G.Skill Aegis 16GB (2x8) 3200", "DDR4", 16, 3200, 3500},
		{"Crucial 16GB (2x8) 3200", "DDR4", 16, 3200, 3800},
		{"Patriot Viper Steel 16GB (2x8) 3600", "DDR4", 16, 3600, 4200},
		{"Kingston FURY Beast 32GB (2x16) 3200", "DDR4", 32, 3200, 6900},
		{"G.Skill Ripjaws V 32GB (2x16) 3600", "DDR4", 32, 3600, 7800},
		{"Kingston FURY Renegade 32GB (2x16) 3600", "DDR4", 32, 3600, 8200},
		{"G.Skill Trident Z RGB 32GB (2x16) 3600", "DDR4", 32, 3600, 8500},
		{"Corsair Vengeance RGB Pro 32GB (2x16) 3600", "DDR4", 32, 3600, 8900},
		{"Corsair Vengeance LPX 64GB (2x32) 3200", "DDR4", 64, 3200, 14500},
		{"G.Skill Ripjaws V 64GB (2x32) 3600", "DDR4", 64, 3600, 16500},
		// DDR5
		{"Crucial 16GB (2x8) 4800", "DDR5", 16, 4800, 5500},
		{"Kingston FURY Beast 16GB (2x8) 5200", "DDR5", 16, 5200, 6500},
		{"Kingston FURY Beast 32GB (2x16) 5600", "DDR5", 32, 5600, 10500},
		{"Kingston FURY Beast 32GB (2x16) 6000", "DDR5", 32, 6000, 11800},
		{"G.Skill Trident Z5 Neo RGB 32GB (2x16) 6000", "DDR5", 32, 6000, 13500},
		{"Corsair Vengeance 32GB (2x16) 6000", "DDR5", 32, 6000, 12500},
		{"G.Skill Ripjaws S5 64GB (2x32) 6000", "DDR5", 64, 6000, 22500},
		{"Corsair Vengeance 64GB (2x32) 6000", "DDR5", 64, 6000, 22500},
		{"Kingston FURY Beast 64GB (2x32) 6000 CL36", "DDR5", 64, 6000, 23500},
		{"G.Skill Trident Z5 Neo 64GB (2x32) 6000", "DDR5", 64, 6000, 24500},
	}

	for _, r := range rams {
		db.Exec("INSERT INTO ram (name, capacity, type, speed, price) VALUES (?, ?, ?, ?, ?)",
			r.name, r.capacity, r.ramType, r.speed, r.price)
	}

	// =========================
	// 5. НАКОПИТЕЛИ
	// =========================
	storages := []struct {
		name, sType string
		cap, price  int
	}{
		// SATA SSD
		{"Kingston A400 240GB", "SATA SSD", 240, 1800},
		{"Kingston A400 480GB", "SATA SSD", 480, 2900},
		{"Crucial BX500 512GB", "SATA SSD", 512, 3200},
		{"Samsung 870 EVO 500GB", "SATA SSD", 512, 4200},
		{"Crucial MX500 1TB", "SATA SSD", 1024, 6200},
		{"WD Blue SA510 1TB", "SATA SSD", 1024, 6200},
		{"Samsung 870 EVO 1TB", "SATA SSD", 1024, 7200},
		{"Samsung 870 EVO 2TB", "SATA SSD", 2048, 14500},
		{"Samsung 870 EVO 4TB", "SATA SSD", 4096, 28500},
		// NVMe SSD
		{"Kingston NV2 500GB", "NVMe SSD", 512, 3200},
		{"Kingston NV2 512GB", "NVMe SSD", 512, 3800},
		{"Crucial P3 Plus 1TB", "NVMe SSD", 1024, 4800},
		{"Kingston NV2 1TB", "NVMe SSD", 1024, 5200},
		{"ADATA XPG GAMMIX S11 Pro 1TB", "NVMe SSD", 1024, 7200},
		{"Samsung 980 1TB", "NVMe SSD", 1024, 7500},
		{"Samsung 980 Pro 1TB", "NVMe SSD", 1024, 8500},
		{"WD Black SN850X 1TB", "NVMe SSD", 1024, 8500},
		{"Kingston KC3000 1TB", "NVMe SSD", 1024, 9800},
		{"Samsung 990 Pro 1TB", "NVMe SSD", 1024, 10500},
		{"Samsung 980 Pro 2TB", "NVMe SSD", 2048, 14500},
		{"WD Black SN850X 2TB", "NVMe SSD", 2048, 14500},
		{"Kingston KC3000 2TB", "NVMe SSD", 2048, 15500},
		{"Samsung 990 Pro 2TB", "NVMe SSD", 2048, 17500},
		{"Lexar NM790 4TB", "NVMe SSD", 4096, 28500},
		{"WD Black SN850X 4TB", "NVMe SSD", 4096, 28500},
		{"Samsung 990 Pro 4TB", "NVMe SSD", 4096, 32500},
		// HDD
		{"Seagate BarraCuda 1TB", "HDD", 1024, 3800},
		{"WD Blue 1TB", "HDD", 1024, 3900},
		{"Seagate BarraCuda 2TB", "HDD", 2048, 6200},
		{"Toshiba P300 3TB", "HDD", 3072, 6800},
		{"Seagate BarraCuda 4TB", "HDD", 4096, 8200},
		{"WD Blue 4TB", "HDD", 4096, 9800},
		{"Seagate IronWolf 4TB", "HDD", 4096, 10500},
		{"WD Red Plus 4TB", "HDD", 4096, 11500},
		{"Seagate IronWolf 8TB", "HDD", 8192, 18500},
		{"WD Red Plus 8TB", "HDD", 8192, 19500},
	}

	for _, s := range storages {
		db.Exec("INSERT INTO storage (name, type, capacity, price) VALUES (?, ?, ?, ?)",
			s.name, s.sType, s.cap, s.price)
	}

	// =========================
	// 6. КОРПУСА
	// =========================
	cases := []struct {
		name, formFactor, color string
		price                   int
	}{
		// Mid-tower
		{"DeepCool CC560 V2", "mid-tower", "black", 4500},
		{"Montech X3 Mesh White", "mid-tower", "white", 5200},
		{"Montech AIR 903 MAX Black", "mid-tower", "black", 7200},
		{"Montech AIR 903 MAX White", "mid-tower", "white", 7500},
		{"Fractal Design Pop Air Black", "mid-tower", "black", 7200},
		{"Fractal Design Pop Air White", "mid-tower", "white", 7200},
		{"Corsair 4000D Airflow Black", "mid-tower", "black", 8200},
		{"Corsair 4000D Airflow White", "mid-tower", "white", 8500},
		{"Lian Li Lancool 216 Black", "mid-tower", "black", 8500},
		{"Lian Li Lancool 216 White", "mid-tower", "white", 8800},
		{"NZXT H5 Flow Black", "mid-tower", "black", 8500},
		{"NZXT H5 Flow White", "mid-tower", "white", 8500},
		{"NZXT H7 Flow Black", "mid-tower", "black", 12500},
		{"NZXT H7 Flow White", "mid-tower", "white", 12500},
		// Aquarium
		{"1stPlayer MEGAVIEW MV6 Black", "aquarium", "black", 6500},
		{"1stPlayer MEGAVIEW MV6 White", "aquarium", "white", 6800},
		{"Lian Li O11 Dynamic EVO Black", "aquarium", "black", 15500},
		{"Lian Li O11 Dynamic EVO White", "aquarium", "white", 16000},
		{"Lian Li O11 Vision Black", "aquarium", "black", 17500},
		{"Lian Li O11 Vision White", "aquarium", "white", 17800},
		{"HYTE Y60 Black", "aquarium", "black", 18500},
		{"HYTE Y60 White", "aquarium", "white", 18500},
		{"HYTE Y70 Black", "aquarium", "black", 22500},
		{"HYTE Y70 White", "aquarium", "white", 22500},
		// Mini-ITX
		{"Cooler Master NR200P Black", "mini-itx", "black", 8500},
		{"Cooler Master NR200P White", "mini-itx", "white", 8500},
		{"Cooler Master NR200P V2 Black", "mini-itx", "black", 12750},
		{"Cooler Master NR200P V2 White", "mini-itx", "white", 12750},
		{"SSUPD Meshlicious Black", "mini-itx", "black", 14500},
		{"SSUPD Meshlicious White", "mini-itx", "white", 14500},
		{"Fractal Design Terra Black", "mini-itx", "black", 18500},
		{"Fractal Design Terra Silver", "mini-itx", "silver", 18500},
	}

	for _, cs := range cases {
		db.Exec("INSERT INTO pc_case (name, form_factor, color, price) VALUES (?, ?, ?, ?)",
			cs.name, cs.formFactor, cs.color, cs.price)
	}

	// =========================
	// 7. ОХЛАЖДЕНИЕ
	// =========================
	coolers := []struct {
		name, cType string
		tdp, price  int
	}{
		// Воздушные — бюджет
		{"DeepCool GAMMAXX 400", "air", 130, 1500},
		{"ID-COOLING SE-214-XT", "air", 180, 1800},
		{"ID-COOLING SE-224-XTS", "air", 220, 2200},
		{"DeepCool AK400", "air", 220, 2800},
		{"DeepCool AK400 Digital", "air", 220, 4200},
		// Воздушные — средние
		{"Thermalright Peerless Assassin 120", "air", 260, 4500},
		{"ID-COOLING FROZN A620", "air", 250, 4500},
		{"Thermalright Phantom Spirit 120 SE", "air", 265, 5200},
		{"DeepCool AK620", "air", 260, 5800},
		{"DeepCool AK620 Digital", "air", 260, 6800},
		// Воздушные — топовые
		{"be quiet! Dark Rock Pro 5", "air", 270, 10500},
		{"be quiet! Dark Rock Elite", "air", 280, 12500},
		{"Noctua NH-U12A chromax.black", "air", 220, 10800},
		{"Noctua NH-D15 chromax.black", "air", 250, 12500},
		// СЖО 240mm
		{"DeepCool LE520 240mm", "water", 250, 6500},
		{"DeepCool LT520 240mm", "water", 280, 8500},
		{"Arctic Liquid Freezer III 240", "water", 280, 9500},
		{"Lian Li Galahad II Trinity 240 Black", "water", 300, 12500},
		{"NZXT Kraken 240 Black", "water", 280, 14500},
		{"Corsair iCUE H100i ELITE CAPELLIX XT", "water", 280, 16500},
		// СЖО 280mm
		{"Arctic Liquid Freezer III 280", "water", 300, 11500},
		{"NZXT Kraken 280 Black", "water", 300, 16500},
		// СЖО 360mm
		{"DeepCool LT720 360mm", "water", 350, 11800},
		{"Arctic Liquid Freezer III 360", "water", 350, 13200},
		{"Lian Li Galahad II Trinity 360 Black", "water", 350, 16500},
		{"NZXT Kraken 360 Black", "water", 350, 19500},
		{"Corsair iCUE H150i ELITE CAPELLIX XT", "water", 350, 21500},
	}

	for _, cl := range coolers {
		db.Exec("INSERT INTO cooling (name, type, tdp, price) VALUES (?, ?, ?, ?)",
			cl.name, cl.cType, cl.tdp, cl.price)
	}

	// =========================
	// 8. БЛОКИ ПИТАНИЯ
	// =========================
	psus := []struct {
		name         string
		power, price int
	}{
		// Бюджет
		{"Aerocool VX Plus 500W", 500, 2700},
		{"Aerocool VX Plus 600W", 600, 3200},
		{"DeepCool PF500 500W", 500, 3400},
		{"DeepCool PK500D 500W", 500, 3800},
		{"DeepCool PK550D 550W", 550, 4200},
		{"DeepCool PF700 700W", 700, 4800},
		{"Chieftec ECO GPE-700S 700W", 700, 5200},
		// Средние
		{"Montech Century 650W Gold", 650, 7200},
		{"Sharkoon Rebel P20 750W Gold", 750, 7500},
		{"DeepCool PM750D 750W Gold", 750, 8400},
		{"Sharkoon Rebel P20 850W Gold", 850, 8500},
		{"Seasonic Core GX-650 ATX 3", 650, 8500},
		{"Sharkoon Rebel P20 1000W Gold", 1000, 9500},
		{"Seasonic Core GX-750 ATX 3", 750, 9500},
		{"Corsair RM750e ATX 3.0", 750, 10500},
		{"Montech TITAN GOLD 750W ATX 3.0", 750, 11200},
		{"Seasonic Core GX-850 ATX 3", 850, 11500},
		{"Corsair RM850e ATX 3.0", 850, 12500},
		{"Montech TITAN GOLD 850W ATX 3.0", 850, 13200},
		{"Seasonic Core GX-1000 ATX 3", 1000, 14500},
		{"Corsair RM1000e ATX 3.0", 1000, 15500},
		{"Chieftec Polaris 3.0 1050W", 1050, 16500},
		// Premium
		{"be quiet! Straight Power 12 850W Platinum", 850, 21500},
		{"be quiet! Dark Power 13 1000W", 1000, 23500},
		{"Seasonic Prime TX-1000 1000W Titanium", 1000, 32500},
		{"Corsair AX1600i 1600W Titanium", 1600, 58500},
	}

	for _, p := range psus {
		db.Exec("INSERT INTO psu (name, power, price) VALUES (?, ?, ?)",
			p.name, p.power, p.price)
	}

	log.Println("База данных успешно инициализирована расширенным каталогом")
}

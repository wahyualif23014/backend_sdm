package controllers

import (
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/wahyualif23014/backendGO/initializers"
	"github.com/xuri/excelize/v2"
)

type RecapResponse struct {
	ID           string  `json:"id"`
	NamaWilayah  string  `json:"nama_wilayah"`
	PotensiLahan float64 `json:"potensi_lahan"`
	TanamLahan   float64 `json:"tanam_lahan"`
	PanenLuas    float64 `json:"panen_luas"`
	PanenTon     float64 `json:"panen_ton"`
	Serapan      float64 `json:"serapan"`
	Level        string  `json:"level"`
	NamaPolsek   string  `json:"nama_polsek,omitempty"`
}

// --- GET DATA UNTUK UI (HIERARKI) ---
func GetRecapData(c *gin.Context) {
	query := `
		SELECT 
			w.kode as id,
			w.nama as nama_wilayah,
			COALESCE(SUM(l.luaslahan), 0) as potensi_lahan,
			COALESCE(SUM(t.luastanam), 0) as tanam_lahan,
			COALESCE(SUM(p.luaspanen), 0) as panen_luas,
			COALESCE(SUM(p.totalpanen), 0) as panen_ton,
			COALESCE(SUM(d.totaldistribusi), 0) as serapan,
			'desa' as level,
			COALESCE(pk.nama, '-') as nama_polsek
		FROM wilayah w
		LEFT JOIN lahan l ON l.idwilayah = w.kode
		LEFT JOIN tanam t ON t.idlahan = l.idlahan
		LEFT JOIN panen p ON p.idlahan = l.idlahan
		LEFT JOIN distribusi d ON d.idlahan = l.idlahan
		LEFT JOIN wilayah pk ON pk.kode = SUBSTR(w.kode, 1, 8)
		WHERE CHAR_LENGTH(w.kode) > 8
		GROUP BY w.kode, w.nama, pk.nama
		ORDER BY w.kode ASC
	`
	rows, err := initializers.DB.Raw(query).Rows()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	polresMap := make(map[string]*RecapResponse)
	polsekMap := make(map[string]*RecapResponse)
	var desaList []RecapResponse

	for rows.Next() {
		var r RecapResponse
		rows.Scan(&r.ID, &r.NamaWilayah, &r.PotensiLahan, &r.TanamLahan, &r.PanenLuas, &r.PanenTon, &r.Serapan, &r.Level, &r.NamaPolsek)

		// Validasi ID agar tidak panic saat slicing
		if len(r.ID) < 8 {
			continue
		}

		pID, sID := r.ID[:5], r.ID[:8]

		// Agregasi Level Polres
		if _, ok := polresMap[pID]; !ok {
			var n string
			initializers.DB.Table("wilayah").Select("nama").Where("kode = ?", pID).Scan(&n)
			polresMap[pID] = &RecapResponse{
				ID:          pID,
				NamaWilayah: strings.TrimPrefix(n, "KAB. "),
				Level:       "polres",
			}
		}
		addSums(polresMap[pID], r)

		// Agregasi Level Polsek
		if _, ok := polsekMap[sID]; !ok {
			polsekMap[sID] = &RecapResponse{
				ID:          sID,
				NamaWilayah: r.NamaPolsek,
				Level:       "polsek",
			}
		}
		addSums(polsekMap[sID], r)

		desaList = append(desaList, r)
	}

	var finalData []RecapResponse

	// Menyusun data Hierarki (Flat List untuk Frontend)
	// Catatan: Iterasi Map di Go itu acak, jadi kita perlu sorting manual nanti
	for pID, pData := range polresMap {
		finalData = append(finalData, *pData)
		for sID, sData := range polsekMap {
			if strings.HasPrefix(sID, pID) {
				finalData = append(finalData, *sData)
				for _, dData := range desaList {
					if strings.HasPrefix(dData.ID, sID) {
						finalData = append(finalData, dData)
					}
				}
			}
		}
	}

	// PENTING: Sort data berdasarkan ID agar urutannya rapi (Polres -> Polsek -> Desa)
	// Karena map di Go iterasinya random, langkah ini wajib agar report tidak berantakan
	sort.Slice(finalData, func(i, j int) bool {
		return finalData[i].ID < finalData[j].ID
	})

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": finalData})
}

// --- FUNGSI EXPORT EXCEL ---
func ExportRecapExcel(c *gin.Context) {
	// 1. Buat File Excel Baru
	f := excelize.NewFile()
	sheet := "Sheet1"

	index, _ := f.NewSheet(sheet)
	f.SetActiveSheet(index)

	// 2. Styling Header
	styleID, err := f.NewStyle(&excelize.Style{
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#1B9E5E"}, Pattern: 1},
		Font:      &excelize.Font{Bold: true, Color: "#FFFFFF", Size: 11},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat style excel"})
		return
	}

	// 3. Tulis Header
	headers := []string{"NO", "WILAYAH", "POTENSI (HA)", "TANAM (HA)", "PANEN (HA)", "PANEN (TON)", "SERAPAN (TON)"}
	for i, h := range headers {
		cell := fmt.Sprintf("%c1", 'A'+i)
		f.SetCellValue(sheet, cell, h)
		f.SetCellStyle(sheet, cell, cell, styleID)
		colName := fmt.Sprintf("%c", 'A'+i)
		f.SetColWidth(sheet, colName, colName, 20)
	}

	// 4. Query Data
	query := `
		SELECT 
			w.nama, 
			COALESCE(SUM(l.luaslahan), 0), 
			COALESCE(SUM(t.luastanam), 0), 
			COALESCE(SUM(p.luaspanen), 0), 
			COALESCE(SUM(p.totalpanen), 0), 
			COALESCE(SUM(d.totaldistribusi), 0)
		FROM wilayah w 
		LEFT JOIN lahan l ON l.idwilayah = w.kode 
		LEFT JOIN tanam t ON t.idlahan = l.idlahan
		LEFT JOIN panen p ON p.idlahan = l.idlahan 
		LEFT JOIN distribusi d ON d.idlahan = l.idlahan
		WHERE CHAR_LENGTH(w.kode) > 8 
		GROUP BY w.kode 
		ORDER BY w.kode ASC
	`

	rows, err := initializers.DB.Raw(query).Rows()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data database"})
		return
	}
	defer rows.Close()

	// 5. Isi Data
	rowIdx := 2
	no := 1
	for rows.Next() {
		var namaWilayah string
		var potensi, tanam, panenLuas, panenTon, serapan float64

		rows.Scan(&namaWilayah, &potensi, &tanam, &panenLuas, &panenTon, &serapan)

		f.SetCellValue(sheet, fmt.Sprintf("A%d", rowIdx), no)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", rowIdx), namaWilayah)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", rowIdx), potensi)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", rowIdx), tanam)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", rowIdx), panenLuas)
		f.SetCellValue(sheet, fmt.Sprintf("F%d", rowIdx), panenTon)
		f.SetCellValue(sheet, fmt.Sprintf("G%d", rowIdx), serapan)

		rowIdx++
		no++
	}

	// 6. Header HTTP untuk Download
	fileName := fmt.Sprintf("Rekap_Presisi_%s.xlsx", time.Now().Format("20060102_150405"))

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Cache-Control", "no-cache")

	if err := f.Write(c.Writer); err != nil {
		fmt.Println("Error writing excel file:", err)
	}
}

// --- HELPER FUNCTION (Wajib Ada) ---
func addSums(t *RecapResponse, s RecapResponse) {
	t.PotensiLahan += s.PotensiLahan
	t.TanamLahan += s.TanamLahan
	t.PanenLuas += s.PanenLuas
	t.PanenTon += s.PanenTon
	t.Serapan += s.Serapan
}

package services

import (
	"errors"

	"gorm.io/gorm"

	"myApp/go-api/models"
)

// Error yang bisa dicek controller lewat errors.Is.
var (
	ErrNotFound      = errors.New("data not found")
	ErrDuplicateCode = errors.New("the code has already been taken")
)

// UnitInput adalah data masukan untuk membuat/mengubah unit.
type UnitInput struct {
	Code        string
	Name        string
	Description *string
	IsActive    *bool // nil -> default true
}

// UnitsService berisi aturan bisnis dan query untuk master units.
type UnitsService struct {
	db *gorm.DB
}

func NewUnitsService(db *gorm.DB) *UnitsService {
	return &UnitsService{db: db}
}
 
// Paginate mengambil data untuk datatable server side.
// Mengembalikan: data, total semua data, total setelah filter.
func (s *UnitsService) Paginate(filter string, start, length int) ([]models.Units, int64, int64, error) {

	var total, filtered int64

	if err := s.db.Model(&models.Units{}).Count(&total).Error; err != nil {
		return nil, 0, 0, err
	}

	// Filter dibangun ulang tiap query supaya kondisi tidak bocor antar query.
	applyFilter := func(tx *gorm.DB) *gorm.DB {
		if filter == "" {
			return tx
		}
		like := "%" + filter + "%"
		return tx.Where("code LIKE ? OR name LIKE ? OR description LIKE ?", like, like, like)
	}

	if err := applyFilter(s.db.Model(&models.Units{})).Count(&filtered).Error; err != nil {
		return nil, 0, 0, err
	}

	var data []models.Units
	err := applyFilter(s.db.Model(&models.Units{})).
		Order("created_at DESC").
		Offset(start).
		Limit(length).
		Find(&data).Error
	if err != nil {
		return nil, 0, 0, err
	}

	return data, total, filtered, nil
}

// FindByReference mengambil satu data; ErrNotFound jika tidak ada.
func (s *UnitsService) FindByReference(reference string) (*models.Units, error) {
	var unit models.Units
	if err := s.db.Where("reference = ?", reference).First(&unit).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &unit, nil
}

// Create menyimpan unit baru; ErrDuplicateCode jika code sudah dipakai.
func (s *UnitsService) Create(in UnitInput) (*models.Units, error) {

	exists, err := s.codeExists(in.Code, "")
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrDuplicateCode
	}

	unit := models.Units{
		ID:          GenerateUUID(),
		Reference:   GenerateShortID("UN", 5),
		Code:        in.Code,
		Name:        in.Name,
		Description: in.Description,
		IsActive:    boolOrDefault(in.IsActive, true),
	}

	if err := s.db.Create(&unit).Error; err != nil {
		return nil, err
	}
	return &unit, nil
}

// Update mengubah unit berdasarkan id; ErrNotFound / ErrDuplicateCode bila perlu.
func (s *UnitsService) Update(id string, in UnitInput) error {

	var unit models.Units
	if err := s.db.Where("id = ?", id).First(&unit).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return err
	}

	exists, err := s.codeExists(in.Code, id)
	if err != nil {
		return err
	}
	if exists {
		return ErrDuplicateCode
	}

	unit.Code = in.Code
	unit.Name = in.Name
	unit.Description = in.Description
	unit.IsActive = boolOrDefault(in.IsActive, true)

	return s.db.Save(&unit).Error
}

// Delete menghapus unit berdasarkan reference; ErrNotFound jika tidak ada.
func (s *UnitsService) Delete(reference string) error {
	result := s.db.Where("reference = ?", reference).Delete(&models.Units{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

// codeExists mengecek duplikat code (exceptID diabaikan saat update).
func (s *UnitsService) codeExists(code, exceptID string) (bool, error) {
	var count int64
	tx := s.db.Model(&models.Units{}).Where("code = ?", code)
	if exceptID != "" {
		tx = tx.Where("id <> ?", exceptID)
	}
	if err := tx.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func boolOrDefault(p *bool, def bool) bool {
	if p == nil {
		return def
	}
	return *p
}


// a := []int{1, 2, 3} 
// b := a[:2] 
// b = append(b, 99) 
// fmt.Println(a)

// jawaban : 
// // mencetak 1 2 99

// var m map[string]int 
// fmt.Println(m["a"]) 
// m["a"] = 1

// Jawaban : 
// // Mencetak 0, lalu panic saat assignment

// func main() { 
// 	ch := make(chan int) ch <- 1 
// 	fmt.Println(<-ch) 
// }
// // 210

// nilai f() == nil? type MyErr struct{} 
// func (e *MyErr) Error() string { 
// 	return "x" 
// } 
// func f() error { 
// 	var p *MyErr = nil return p 
// }

// // false 

// Setelah memanggil ctx, cancel := context.WithTimeout(parent, 2*time.Second), apa yang sebaiknya selalu dilakukan?
// // Memanggil cancel(), biasanya lewat defer cancel()






// a = [3]int{1, 2, 3} 
// b = a b[0] = 99 
// fmt.Println(a[0], b[0])

// // jawaban 
// // 199

// s := "halo🙂" 
// fmt.Println(len(s), len([]rune(s)))
// // fmt.Println(8, 1)

// // jawaban 85


// var ErrNotFound = errors.New("not found") 
// err := fmt.Errorf("get user: %w", ErrNotFound) 
// fmt.Println(err == ErrNotFound, errors.Is(err, ErrNotFound))

// // jawaban : false, true


// type P struct{ N int } 
// m := map[string]P{
// 		"a": {1}
// 	} m["a"].N = 2

// not asign ke field

// var fs []func() 
// for i := 0; i < 3; i++ { 
// 		fs = append(fs, func() { 
// 				fmt.Print(i) 
// 			}) 
// 	} 
// for _, f := range fs { 
// 		f() 
// 	}













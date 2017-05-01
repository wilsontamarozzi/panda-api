package repositories

import (
	"github.com/wilsontamarozzi/panda-api/models"
	"github.com/wilsontamarozzi/panda-api/helpers"
	"strconv"
	"net/url"
	"github.com/wilsontamarozzi/panda-api/database"
	"log"
)

type TaskCategoryRepositoryInterface interface {
	GetAll(q url.Values) models.TaskCategories
	Get(id string) models.TaskCategory
	Delete(id string) error
	Create(tc *models.TaskCategory) error
	Update(tc *models.TaskCategory) error
	CountRows() int
}

type taskCategoryRepository struct{}

func NewTaskCategoryRepository() *taskCategoryRepository {
	return new(taskCategoryRepository)
}

func (repository taskCategoryRepository) GetAll(q url.Values) models.TaskCategories {
	db := database.GetInstance()

	currentPage, _ := strconv.Atoi(q.Get("page"))
	itemPerPage, _ := strconv.Atoi(q.Get("per_page"))
	pagination := helpers.MakePagination(repository.CountRows(), currentPage, itemPerPage)

	if q.Get("description") != "" {
		db = db.Where("description iLIKE ?", "%" + q.Get("description") + "%")
	}

	var taskCategories models.TaskCategories
	taskCategories.Meta.Pagination = pagination

	db.Limit(pagination.ItemPerPage).
		Offset(pagination.StartIndex).
		Order("description desc").
		Find(&taskCategories)

	return taskCategories
}

func (repository taskCategoryRepository) Get(id string) models.TaskCategory {
	db := database.GetInstance()

	var taskCategory models.TaskCategory
	db.Where("uuid = ?", id).
		First(&taskCategory)

	return taskCategory
}

func (repository taskCategoryRepository) Delete(id string) error {
	db := database.GetInstance()

	err := db.Where("uuid = ?", id).Delete(&models.TaskCategory{}).Error
	if err != nil {
		log.Print(err.Error())
	}

	return err;
}

func (repository taskCategoryRepository) Create(tc *models.TaskCategory) error {
	db := database.GetInstance()

	record := models.TaskCategory{ Description : tc.Description }

	err := db.Create(&record).Error
	if err != nil {
		log.Print(err.Error())
	}

	*(tc) = record
	return err
}

func (repository taskCategoryRepository) Update(tc *models.TaskCategory) error {
	db := database.GetInstance()

	record := models.TaskCategory{ Description : tc.Description }

	err := db.Model(&models.TaskCategory{}).
		Where("uuid = ?", tc.UUID).
		Updates(&record).Error

	if err != nil {
		log.Print(err.Error())
	}

	return err
}

func (repository taskCategoryRepository) CountRows() int {
	db := database.GetInstance()
	var count int
	db.Model(&models.TaskCategory{}).Count(&count)

	return count
}
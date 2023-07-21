package repository

import (
	"github.com/zxnlx/common"
	"github.com/zxnlx/pod/domain/model"
	"gorm.io/gorm"
)

type IPodRepository interface {
	InitTable() error
	FindPodById(id int64) (*model.Pod, error)
	CreatePod(pod *model.Pod) (int64, error)
	DeletePodById(id int64) error
	Update(pod *model.Pod) error
	FindAll() ([]model.Pod, error)
}

func NewPodRepository(db *gorm.DB) IPodRepository {
	return &PodRepository{
		db: db,
	}
}

type PodRepository struct {
	db *gorm.DB
}

func (p *PodRepository) InitTable() error {
	common.Info("Init table 11")
	return p.db.AutoMigrate(&model.Pod{}, &model.PodEnv{}, &model.PodPort{})
}

func (p *PodRepository) FindPodById(id int64) (*model.Pod, error) {
	pod := &model.Pod{}

	err := p.db.First(pod, id).Error
	return pod, err
}

func (p *PodRepository) CreatePod(pod *model.Pod) (int64, error) {
	err := p.db.Create(pod).Error
	return pod.ID, err
}

func (p *PodRepository) DeletePodById(id int64) error {
	tx := p.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if tx.Error != nil {
		return tx.Error
	}

	err := tx.Where("id = ?", id).Delete(&model.Pod{}).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Where("pod_id = ?", id).Delete(&model.PodEnv{}).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Where("pod_id = ?", id).Delete(&model.PodPort{}).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (p *PodRepository) Update(pod *model.Pod) error {
	return p.db.Model(pod).Updates(pod).Error
}

func (p *PodRepository) FindAll() ([]model.Pod, error) {
	pods := make([]model.Pod, 0)
	err := p.db.Find(&pods).Error
	return pods, err
}

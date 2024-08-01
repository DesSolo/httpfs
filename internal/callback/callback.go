package callback

import (
	"bytes"
	"context"
	"fmt"
	"httpfs/internal/entities"
	"io"
)

// CallBack интерфейс предварительной обработки загрузки файла
type PreCallBack interface {
	// предварительная обработка
	Handle(context.Context, io.Reader) error
}

// PreCallBackFunc функция предварительной обработки
type PreCallBackFunc func(context.Context, io.Reader) error

func (f PreCallBackFunc) Handle(ctx context.Context, r io.Reader) error {
	return f(ctx, r)
}

// CallBack интерфейс обработки после успешной загрузки файла
type PostCallBack interface {
	Handle(context.Context, entities.Hash) error
}

// PostCallBackFunc функция обработки
type PostCallBackFunc func(context.Context, entities.Hash) error

func (f PostCallBackFunc) Handle(ctx context.Context, h entities.Hash) error {
	return f(ctx, h)
}

// CallBack обработчик
type CallBack struct {
	// предварительные обработчики
	pre []PreCallBack
	// обработчики после успешной загрузки файла
	post []PostCallBack
}

// New создание обработчика
func New() *CallBack {
	return &CallBack{
		pre:  make([]PreCallBack, 0),
		post: make([]PostCallBack, 0),
	}
}

// RegisterPre регистрация предварительного обработчика
func (c *CallBack) RegisterPre(cb PreCallBack) {
	c.pre = append(c.pre, cb)
}

// RegisterPost регистрация обработчика после успешной загрузки файла
func (c *CallBack) RegisterPost(cb PostCallBack) {
	c.post = append(c.post, cb)
}

// Pre выполнение предварительной обработки
func (c *CallBack) Pre(ctx context.Context, r io.Reader) error {
	// таже самая проблема с вычитыванием в память и 
	data, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("io read err: %w", err)
	}

	for _, cb := range c.pre {
		if err := cb.Handle(ctx, bytes.NewBuffer(data)); err != nil {
			return err
		}
	}
	return nil
}

// Post выполнение обработки
func (c *CallBack) Post(ctx context.Context, h entities.Hash) error {
	for _, cb := range c.post {
		if err := cb.Handle(ctx, h); err != nil {
			return err
		}
	}
	return nil
}

/*
 *     Drop-in replacement for Go's sync featuring generics and channels.
 *     Copyright (C) 2025  Dviih
 *
 *     This program is free software: you can redistribute it and/or modify
 *     it under the terms of the GNU Affero General Public License as published
 *     by the Free Software Foundation, either version 3 of the License, or
 *     (at your option) any later version.
 *
 *     This program is distributed in the hope that it will be useful,
 *     but WITHOUT ANY WARRANTY; without even the implied warranty of
 *     MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *     GNU Affero General Public License for more details.
 *
 *     You should have received a copy of the GNU Affero General Public License
 *     along with this program.  If not, see <https://www.gnu.org/licenses/>.
 *
 */

package sync

import (
	"github.com/Dviih/sync/channel"
	"sync"
)

type Pool[T interface{}] struct {
	pool sync.Pool
	New  func() T
}

func (pool *Pool[T]) Get() T {
	t, ok := pool.pool.Get().(T)
	if ok {
		return t
	}

	return pool.new()
}

func (pool *Pool[T]) Put(t T) {
	pool.pool.Put(t)
}

func (pool *Pool[T]) new() T {
	if pool.New == nil {
		return Zero[T]()
	}

	return pool.New()
}

type PoolChan[T interface{}] struct {
	Pool[T]
	channel.Manager[T]
}

func (pc *PoolChan[T]) Sender() chan<- T {
	if pc.Manager.Handler == nil {
		pc.Manager.Handler = func(t T) {
			pc.Pool.Put(t)
		}
	}

	return pc.Manager.Sender()
}

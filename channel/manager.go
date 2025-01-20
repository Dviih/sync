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

package channel

import (
	"reflect"
	"sync"
)

type Manager[T interface{}] struct {
	receivers sync.Map
	dead      sync.Map

	Handler func(T)
	Dead    chan uintptr
}

func (manager *Manager[T]) Sender() chan<- T {
	c := make(chan T)

	go func() {
		for {
			select {
			case t, ok := <-c:
				if !ok {
					return
				}

				go manager.handle(t)
			}
		}
	}()

	return c
}

func (manager *Manager[T]) Receiver(v ...int) <-chan T {
	var c chan T

	switch len(v) {
	case 0:
		c = make(chan T)
	case 1:
		c = make(chan T, v[0])
	default:
		panic("invalid receiver length")
	}

	manager.receivers.Store(c, true)
	return c
}

func (manager *Manager[T]) handle(t T) {
	if manager.Handler != nil {
		manager.Handler(t)
	}

	var wg sync.WaitGroup

	manager.receivers.Range(func(key, value any) bool {
		wg.Add(1)

		go func(receiver chan T) {
			defer wg.Done()

			select {
			case receiver <- t:
				manager.dead.Delete(receiver)
			default:
				if _, ok := manager.dead.Load(receiver); ok {
					if manager.Dead != nil {
						manager.Dead <- reflect.ValueOf(receiver).Pointer()
					}

					manager.receivers.Delete(receiver)
					manager.dead.Delete(receiver)

					return
				}

				manager.dead.Store(receiver, true)
			}
		}(key.(chan T))

		return true
	})
}

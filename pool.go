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

import "github.com/Dviih/sync/channel"

type Pool[T interface{}] struct {
	New func() T
	c   chan T
	m   Mutex
}

func (pool *Pool[T]) Get() T {
	defer pool.m.Unlock()
	pool.m.Lock()

	if len(pool.c) == 0 {
		return pool.new()
	}

	v, ok := <-pool.c
	if !ok {
		return pool.new()
	}

	return v
}

func (pool *Pool[T]) Put(t T) {
	defer pool.m.Unlock()
	pool.m.Lock()

	select {
	case pool.c <- t:
		return
	default:
		if len(pool.c) <= cap(pool.c) {
			if pool.c == nil {
				pool.c = make(chan T, 2)
				pool.c <- t
				return
			}

			panic("resize ain't the issue here")
		}

		c := make(chan T, 2*cap(pool.c))
		for {
			v, ok := <-pool.c
			if !ok {
				break
			}

			c <- v
		}

		c <- t
		pool.c = c
		return

	}
}

func (pool *Pool[T]) new() T {
	if pool.New == nil {
		return Zero[T]()
	}

	return pool.New()
}


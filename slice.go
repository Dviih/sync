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

type Slice[T interface{}] struct {
	data []T
	m    Mutex
}

func (slice *Slice[T]) Index(i int) T {
	defer slice.m.Unlock()

	slice.m.Lock()
	return slice.data[i]
}

func (slice *Slice[T]) Append(v ...T) {
	defer slice.m.Unlock()

	slice.m.Lock()
	slice.data = append(slice.data, v...)
}

func (slice *Slice[T]) Delete(i int) {
	defer slice.m.Unlock()

	slice.m.Lock()
	slice.data = append(slice.data[:i], slice.data[i+1:]...)
}

func (slice *Slice[T]) Len() int {
	defer slice.m.Unlock()

	slice.m.Lock()
	return len(slice.data)
}

func (slice *Slice[T]) Cap() int {
	defer slice.m.Unlock()

	slice.m.Lock()
	return cap(slice.data)
}

func (slice *Slice[T]) Range(fn func(int, T) bool) {
	defer slice.m.Unlock()
	slice.m.Lock()

	for i, t := range slice.data {
		if !fn(i, t) {
			break
		}
	}
}

func (slice *Slice[T]) Slice() []T {
	defer slice.m.Unlock()

	slice.m.Lock()
	return slice.data[:]
}


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

//go:build tinygo

package sync

import "reflect"

func (maps *Map[K, V]) Swap(key K, value V) (V, bool) {
	p, ok := maps.Load(key)
	if !ok {
		p = Zero[V]()
	}

	maps.Store(key, value)
	return p, ok
}

func (maps *Map[K, V]) CompareAndSwap(key K, old, new V) bool {
	c, ok := maps.Load(key)
	if !ok {
		return false
	}

	if reflect.DeepEqual(old, c) {
		return false
	}

	maps.Store(key, new)
	return true
}

func (maps *Map[K, V]) CompareAndDelete(key K, value V) bool {
	c, ok := maps.Load(key)
	if !ok {
		return false
	}

	if reflect.DeepEqual(value, c) {
		return false
	}

	maps.Delete(key)
	return true
}

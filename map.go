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

type Map[K comparable, V interface{}] struct {
	m sync.Map
}

func (maps *Map[K, V]) Load(key K) (V, bool) {
	v, ok := maps.m.Load(key)
	if !ok {
		return Zero[V](), false
	}

	return v.(V), true
}

func (maps *Map[K, V]) Store(key K, value V) {
	maps.m.Store(key, value)
}

func (maps *Map[K, V]) Delete(key K) {
	maps.m.Delete(key)
}

func (maps *Map[K, V]) Clear() {
	maps.m.Clear()
}

func (maps *Map[K, V]) LoadOrStore(key K, value V) (V, bool) {
	v, ok := maps.m.LoadOrStore(key, value)
	if !ok {
		return Zero[V](), false
	}

	return v.(V), true
}

func (maps *Map[K, V]) LoadAndDelete(key K) (V, bool) {
	v, ok := maps.m.LoadAndDelete(key)
	if !ok {
		return Zero[V](), false
	}

	return v.(V), true
}

func (maps *Map[K, V]) Range(fn func(K, V) bool) {
	maps.m.Range(func(k, v interface{}) bool {
		return fn(k.(K), v.(V))
	})
}

func (maps *Map[K, V]) Map() map[K]V {
	m := make(map[K]V)

	maps.Range(func(k K, v V) bool {
		m[k] = v
		return true
	})

	return m
}

type MapChan[K comparable, V interface{}] struct {
	Map[K, V]
	channel.Manager[*channel.KV[K, V]]
}

func (mc *MapChan[K, V]) Sender() chan<- *channel.KV[K, V] {
	if mc.Manager.Handler == nil {
		mc.Manager.Handler = func(kv *channel.KV[K, V]) {
			mc.Map.Store(kv.Key, kv.Value)
		}
	}

	return mc.Manager.Sender()
}

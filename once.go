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
	"reflect"
	"sync/atomic"
)

type Once[T interface{}] struct {
	done   atomic.Bool
	m      Mutex
	result []interface{}
}

func (once *Once[T]) Do(fn T) interface{} {
	if once.done.Load() {
		switch len(once.result) {
		case 0:
			return nil
		case 1:
			return once.result[0]
		default:
			return once.result[:]
		}
	}

	defer once.m.Unlock()
	once.m.Lock()

	v := reflect.ValueOf(fn)

	if v.Type().NumIn() > 0 {
		panic("once: Do func must not have input arguments")
	}

	out := v.Call(nil)
	for i := 0; i < len(out); i++ {
		once.result = append(once.result, out[i].Interface())
	}

	once.done.Store(true)

	switch len(once.result) {
	case 0:
		return nil
	case 1:
		return once.result[0]
	default:
		return once.result[:]
	}
}

func OnceFunc(fn func()) func() {
	once := &Once[func()]{}

	return func() {
		once.Do(fn)
	}
}

func OnceValue[T interface{}](fn func() T) func() T {
	once := &Once[func() T]{}

	return func() T {
		return once.Do(fn).(T)
	}
}

func OnceValues[T1, T2 interface{}](fn func() (T1, T2)) func() (T1, T2) {
	once := &Once[func() (T1, T2)]{}

	return func() (T1, T2) {
		i := once.Do(fn)

		return i.([]interface{})[0].(T1), i.([]interface{})[1].(T2)
	}
}

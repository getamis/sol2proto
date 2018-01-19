// Copyright 2017 AMIS Technologies
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package rand

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Generator", func() {
	Describe("with SHA1", func() {
		Context("with Hex encoder", func() {
			generator := New(Sha1Hash(), HexEncoder())

			It("generate key", func() {
				rawKey := generator.Key()
				key := generator.KeyEncoded()

				Expect(len(rawKey)).ShouldNot(Equal(0))
				Expect(len(key)).ShouldNot(Equal(0))
			})
		})

		Context("with Base64 encoder", func() {
			generator := New(Sha1Hash(), Base64Encoder())

			It("generate key", func() {
				rawKey := generator.Key()
				key := generator.KeyEncoded()

				Expect(len(rawKey)).ShouldNot(Equal(0))
				Expect(len(key)).ShouldNot(Equal(0))
			})
		})

		Context("with UUID encoder", func() {
			generator := New(Sha1Hash(), UUIDEncoder())

			It("generate key", func() {
				rawKey := generator.Key()
				key := generator.KeyEncoded()

				Expect(len(rawKey)).ShouldNot(Equal(0))
				Expect(len(key)).ShouldNot(Equal(0))
			})
		})
	})

	Describe("with SHA256", func() {
		Context("with Hex encoder", func() {
			generator := New(Sha256Hash(), HexEncoder())

			It("generate key", func() {
				rawKey := generator.Key()
				key := generator.KeyEncoded()

				Expect(len(rawKey)).ShouldNot(Equal(0))
				Expect(len(key)).ShouldNot(Equal(0))
			})
		})

		Context("with Base64 encoder", func() {
			generator := New(Sha256Hash(), Base64Encoder())

			It("generate key", func() {
				rawKey := generator.Key()
				key := generator.KeyEncoded()

				Expect(len(rawKey)).ShouldNot(Equal(0))
				Expect(len(key)).ShouldNot(Equal(0))
			})
		})

		Context("with UUID encoder", func() {
			generator := New(Sha1Hash(), UUIDEncoder())

			It("generate key", func() {
				rawKey := generator.Key()
				key := generator.KeyEncoded()

				Expect(len(rawKey)).ShouldNot(Equal(0))
				Expect(len(key)).ShouldNot(Equal(0))
			})
		})
	})

	Describe("with SHA512", func() {
		Context("with Hex encoder", func() {
			generator := New(Sha512Hash(), HexEncoder())

			It("generate key", func() {
				rawKey := generator.Key()
				key := generator.KeyEncoded()

				Expect(len(rawKey)).ShouldNot(Equal(0))
				Expect(len(key)).ShouldNot(Equal(0))
			})
		})

		Context("with Base64 encoder", func() {
			generator := New(Sha512Hash(), Base64Encoder())

			It("generate key", func() {
				rawKey := generator.Key()
				key := generator.KeyEncoded()

				Expect(len(rawKey)).ShouldNot(Equal(0))
				Expect(len(key)).ShouldNot(Equal(0))
			})
		})

		Context("with UUID encoder", func() {
			generator := New(Sha1Hash(), UUIDEncoder())

			It("generate key", func() {
				rawKey := generator.Key()
				key := generator.KeyEncoded()

				Expect(len(rawKey)).ShouldNot(Equal(0))
				Expect(len(key)).ShouldNot(Equal(0))
			})
		})
	})
})

func TestCryptoGeneratorSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Crypto Generator Suite")
}

package x25519

import (
	"crypto/rand"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"testing"

	"github.com/cloudflare/circl/internal/test"
)

type katVector struct {
	TcID    int
	Public  Key
	Private Key
	Shared  Key
}

type timesVector struct {
	T uint32
	W Key
}

func strToKey(s string) (k Key) {
	for j := 0; j < Size; j++ {
		a, _ := strconv.ParseUint(s[2*j:2*j+2], 16, 8)
		k[j] = byte(a)
	}
	return k
}

// Indicates wether long tests should be run
var runLongTest bool

func TestMain(m *testing.M) {
	flag.BoolVar(&runLongTest, "long", false, "runs longer tests")
	flag.Parse()
	os.Exit(m.Run())
}

func TestRFC7748Kat(t *testing.T) {
	const nameFile = "testdata/rfc7748_kat_test.json"
	var kat []struct {
		Public  string `json:"input"`
		Shared  string `json:"output"`
		Private string `json:"scalar"`
	}

	jsonFile, err := os.Open(nameFile)
	if err != nil {
		t.Fatalf("File %v can not be opened. Error: %v", nameFile, err)
	}
	defer jsonFile.Close()
	input, _ := ioutil.ReadAll(jsonFile)

	err = json.Unmarshal(input, &kat)
	if err != nil {
		t.Fatalf("File %v can not be loaded. Error: %v", nameFile, err)
	}

	vec := make([]katVector, len(kat))
	for i := range kat {
		vec[i].Public = strToKey(kat[i].Public)
		vec[i].Shared = strToKey(kat[i].Shared)
		vec[i].Private = strToKey(kat[i].Private)
	}

	var got, want Key
	for _, v := range vec {
		got.Shared(&v.Private, &v.Public)
		want = v.Shared
		if got != want {
			test.ReportError(t, got, want, v)
		}
	}
}

func TestRFC7748Times(t *testing.T) {
	const nameFile = "testdata/rfc7748_times_test.json"
	jsonFile, err := os.Open(nameFile)
	if err != nil {
		t.Fatalf("File %v can not be opened. Error: %v", nameFile, err)
	}
	defer jsonFile.Close()
	input, _ := ioutil.ReadAll(jsonFile)

	var kat []struct {
		Times uint32 `json:"times"`
		Key   string `json:"key"`
	}
	err = json.Unmarshal(input, &kat)
	if err != nil {
		t.Fatalf("File %v can not be loaded. Error: %v", nameFile, err)
	}
	vec := make([]timesVector, len(kat))
	for i := range kat {
		vec[i].T = kat[i].Times
		vec[i].W = strToKey(kat[i].Key)
	}
	var u, r Key
	for _, v := range vec {
		if !runLongTest && v.T == uint32(1000000) {
			t.Log("Skipped one long test, add -long flag to run longer tests")
			continue
		}
		u.SetGenerator()
		k := u
		for i := uint32(0); i < v.T; i++ {
			r.Shared(&k, &u)
			u = k
			k = r
		}
		got, want := k, v.W

		if got != want {
			test.ReportError(t, got, want, v.T)
		}
	}
}

func TestBase(t *testing.T) {
	testTimes := 1 << 10
	var secret, got, want, gen Key
	gen.SetGenerator()
	for i := 0; i < testTimes; i++ {
		_, _ = rand.Read(secret[:])

		got.KeyGen(&secret)
		want.Shared(&secret, &gen)

		if got != want {
			test.ReportError(t, got, want, secret)
		}
	}
}

func TestWycheproof(t *testing.T) {
	// Test vectors from Wycheproof v0.4.12
	const nameFile = "testdata/wycheproof_kat.json"
	jsonFile, err := os.Open(nameFile)
	if err != nil {
		t.Fatalf("File %v can not be opened. Error: %v", nameFile, err)
	}
	defer jsonFile.Close()

	input, _ := ioutil.ReadAll(jsonFile)
	var vecRaw []struct {
		TcID    int      `json:"tcId"`
		Comment string   `json:"comment"`
		Curve   string   `json:"curve"`
		Public  string   `json:"public"`
		Private string   `json:"private"`
		Shared  string   `json:"shared"`
		Result  string   `json:"result"`
		Flags   []string `json:"flags"`
	}

	err = json.Unmarshal(input, &vecRaw)
	if err != nil {
		t.Fatalf("File %v can not be loaded. Error: %v", nameFile, err)
	}
	vec := make([]katVector, len(vecRaw))
	for i, v := range vecRaw {
		vec[i].TcID = v.TcID
		vec[i].Public = strToKey(v.Public)
		vec[i].Private = strToKey(v.Private)
		vec[i].Shared = strToKey(v.Shared)
	}
	var got, want Key
	for _, v := range vec {
		got.Shared(&v.Private, &v.Public)
		want = v.Shared
		if got != want {
			test.ReportError(t, got, want, v)
		}
	}
}

// Benchmarks
func BenchmarkX25519(b *testing.B) {
	var x, y, z Key
	_, _ = rand.Read(x[:])
	_, _ = rand.Read(y[:])

	b.Run("KeyGen", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			x.KeyGen(&y)
		}
	})
	b.Run("Shared", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			z.Shared(&x, &y)
			y = x
			x = z
		}
	})
}

func Example_x25519() {
	var AliceSecret, BobSecret,
		AlicePublic, BobPublic,
		AliceShared, BobShared Key

	// Generating Alice's secret and public keys
	_, _ = rand.Read(AliceSecret[:])
	AlicePublic.KeyGen(&AliceSecret)
	// Generating Bob's secret and public keys
	_, _ = rand.Read(BobSecret[:])
	BobPublic.KeyGen(&BobSecret)
	// Deriving Alice's shared key
	AliceShared.Shared(&AliceSecret, &BobPublic)
	// Deriving Bob's shared key
	BobShared.Shared(&BobSecret, &AlicePublic)

	fmt.Println(AliceShared == BobShared)
	// Output: true
}

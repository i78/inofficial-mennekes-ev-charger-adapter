package mennekes

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMennekesMcp(t *testing.T) {

	t.Run("should fetch status from expected charger address", func(t *testing.T) {
		called := false

		testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			res.WriteHeader(200)
			called = true
		}))
		defer testServer.Close()

		subject := NewMennekesMCP(MennekesConfiguration{
			ChargerHost: testServer.URL,
			ChargerPin1: "0000",
		})

		subject.GetCurrentChargerStatus()

		assert.True(t, called)
	})

	t.Run("should decode status as expected", func(t *testing.T) {
		fakeStatus := `{
			"ChgState": "Idle",
			"ChgDuration": 111,
			"ChgNrg": 222,
			"NrgDemand": 333,
			"ActPwr": 444,
			"ActCurr": 555,
			"MaxCurrT1": 666}`

		testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			res.WriteHeader(200)
			res.Write([]byte(fakeStatus))
		}))
		defer testServer.Close()

		subject := NewMennekesMCP(MennekesConfiguration{
			ChargerHost: testServer.URL,
			ChargerPin1: "0000",
		})

		result, _ := subject.GetCurrentChargerStatus()

		assert.Equal(t, 222.0, result.TotalEnergyConsumption)
		assert.Equal(t, 555.0, result.OutputCurrent)
		assert.Equal(t, 444.0, result.PowerOutput)
	})

	t.Run("should survive empty charger telegram and report error", func(t *testing.T) {
		fakeStatus := ""

		testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			res.WriteHeader(200)
			res.Write([]byte(fakeStatus))
		}))
		defer testServer.Close()

		subject := NewMennekesMCP(MennekesConfiguration{
			ChargerHost: testServer.URL,
			ChargerPin1: "0000",
		})

		_, err := subject.GetCurrentChargerStatus()

		assert.NotNil(t, err)
	})
}

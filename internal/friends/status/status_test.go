package status

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"testing"
	"time"
)

const EnableJSONOutput = false

func TestStatus(t *testing.T) {
	d := &Frame{start: time.Now()}
	defer d.Done()

	var (
		testVersion  = "v1"
		testBUID     = "abcd"
		testPlatform = "bnet"
		testProduct  = "wolfenstein"
		testCustom   = "bethesda"
		testConn     = "https://localhost:10001/client"
		testJSON     = `{"buid":%q,"status":%q,"main_game_status":%q,"player_status":%q,"extended_game_status":[%v],"custom_data":%q}`
		globalStatus = &Status{
			BUID: testBUID,
		}
		productStatus = &Status{
			BUID:     testBUID,
			Product:  testProduct,
			Platform: testPlatform,
		}
	)

	t.Logf("presence.api.version: %s", testVersion)

	t.Run("Time", func(t *testing.T) {
		t.Run("Format", func(t *testing.T) {
			for _, sec := range []float64{UnitMinute * 2, UnitHour * 2, UnitDay * 2, UnitWeek * 2, UnitMonth * 2, UnitYear, UnitYear * 2} {
				if f, s := formatUnit(sec); f == 0 || s == "" {
					t.Errorf("got %v (%v); want %v (%s)", f, s, sec, "UOM")
				}
			}
		})
	})
	/* ------------------------------------------------------------------------ */
	t.Run("Key", func(t *testing.T) {
		t.Run("Global", func(t *testing.T) {
			// key should be global and not contain idle suffix
			t.Run("Default", func(t *testing.T) {
				key := globalStatus.Key(true, false)
				if strings.Contains(key, SxLastActivity) {
					t.Errorf("got %s; want %s", key, strings.TrimSuffix(key, SxLastActivity))
				}
			})
			/* -------------------------------------------------------------------- */
			// key should be global and contain idle suffix
			t.Run("Idle", func(t *testing.T) {
				key := globalStatus.Key(true, true)
				if !strings.Contains(key, SxLastActivity) {
					t.Errorf("got %s; want %s", key, key+SxLastActivity)
				}
			})
			/* -------------------------------------------------------------------- */
		})
		t.Run("Product", func(t *testing.T) {
			// check product key not idle
			t.Run("Default", func(t *testing.T) {
				productKey := productStatus.Key(false, false)
				if !strings.Contains(productKey, testProduct) {
					t.Errorf("got %s; want %s", productStatus.Product, testProduct)
				}
				if !strings.Contains(productKey, testPlatform) {
					t.Errorf("got %s; want %s", productStatus.Platform, testPlatform)
				}
				if strings.Contains(productKey, SxLastActivity) {
					t.Errorf("got %s; want %s", productKey, strings.TrimSuffix(productKey, SxLastActivity))
				}
			})
			/* -------------------------------------------------------------------- */
			// check product key for idle suffix
			t.Run("Idle", func(t *testing.T) {
				key := productStatus.Key(false, true)
				if !strings.Contains(key, SxLastActivity) {
					t.Errorf("got %s; want %s", key, key+SxLastActivity)
				}
			})
			/* -------------------------------------------------------------------- */
		})
	})
	/* ------------------------------------------------------------------------ */
	t.Run("Get", func(t *testing.T) {
		t.Run("Default", func(t *testing.T) {
			if globalStatus.Enum != Offline {
				t.Errorf("got %d (%s); want %d (%s)", globalStatus.Enum, globalStatus.Enum, Offline, Offline)
			}
			if _, err := json.Marshal(&globalStatus); err != nil {
				t.Errorf("got %v; want %v", err, nil)
			}
		})
	})
	/* ------------------------------------------------------------------------ */
	t.Run("Set", func(t *testing.T) {
		t.Run("Global", func(t *testing.T) {
			for _, to := range []Kind{Offline, Online, Offline, Online, Idle, Online, AppearOffline, DND, Offline} {
				t.Run(to.String(), func(t *testing.T) {
					if to == Online {
						globalStatus.Connection = &testConn
						globalStatus.Custom = &testCustom
					}
					if err := d.Eval(globalStatus.Set(to), globalStatus.Enum, to); err != nil {
						t.Error(err)
					}
					d.step++
				})
			}
		})
		/* ---------------------------------------------------------------------- */
		t.Run("Product", func(t *testing.T) {
			for _, to := range []Kind{Offline, Online, Offline, Online, Idle, AppearOffline, DND, Offline} {
				t.Run(to.String(), func(t *testing.T) {
					if to == Online {
						productStatus.Connection = &testConn
						productStatus.Custom = &testCustom
					}
					if err := d.Eval(productStatus.Set(to), productStatus.Enum, to); err != nil {
						t.Error(err)
					}
					d.step++
				})
			}
		})
	})
	/* ------------------------------------------------------------------------ */
	t.Run("JSON", func(t *testing.T) {
		t.Run("Good", func(t *testing.T) {
			t.Run("Body", func(t *testing.T) {
				data := fmt.Sprintf(testJSON, productStatus.BUID, "online", "", "", "", "")
				temp := &Status{}
				err := json.Unmarshal([]byte(data), &temp)
				if err != nil {
					t.Errorf("want %v, got %v", nil, err)
				}
			})
			/* -------------------------------------------------------------------- */
			t.Run("Status", func(t *testing.T) {
				data := fmt.Sprintf(testJSON, productStatus.BUID, "online", "", "", "", "")
				temp := &Status{}
				if err := json.Unmarshal([]byte(data), &temp); err != nil {
					t.Errorf("want %v, got %v", nil, err.Error())
				}
			})
			/* -------------------------------------------------------------------- */
			t.Run("Player", func(t *testing.T) {
				data := fmt.Sprintf(testJSON, productStatus.BUID, "online", "", "online", "", "")
				temp := &Status{}
				if err := json.Unmarshal([]byte(data), &temp); err != nil {
					t.Errorf("want %v, got %v", nil, err.Error())
				}
			})
			/* -------------------------------------------------------------------- */
			t.Run("Custom", func(t *testing.T) {
				amount := MaxCustomDataBytes - 1
				var b strings.Builder
				for i := 0; i < amount; i++ {
					fmt.Fprintf(&b, "%s", "a")
				}
				bad := b.String()
				pattern := `{"buid": %q,"status":"online","custom_data":%q}`
				data := fmt.Sprintf(pattern, productStatus.BUID, bad)
				temp := &Status{}
				err := json.Unmarshal([]byte(data), temp)
				if err != nil {
					t.Errorf("want %v, got %v", err.Error(), nil)
				}
				if temp.Custom != nil {
					if len(*temp.Custom) != amount {
						t.Errorf("want %v, got %v", amount, len(*temp.Custom))
					}
				}
			})
			/* -------------------------------------------------------------------- */
		})
		/* ---------------------------------------------------------------------- */
		t.Run("Bad", func(t *testing.T) {
			t.Run("Body", func(t *testing.T) {
				temp := &Status{}
				if err := json.Unmarshal([]byte("1"), &temp); err == nil {
					t.Errorf("want %v, got %v", "error unmarshaling bad json", nil)
				}
			})
			/* -------------------------------------------------------------------- */
			t.Run("Enum", func(t *testing.T) {
				b, _ := json.Marshal(&Status{Enum: 10})
				temp := &Status{}
				if err := json.Unmarshal(b, &temp); err == nil {
					t.Errorf("want %v, got %v", "error unmarshaling bad json", nil)
				}
			})
			/* -------------------------------------------------------------------- */
			t.Run("BUID", func(t *testing.T) {
				data := fmt.Sprintf(testJSON, "", "", "", "", "", "")
				temp := &Status{}
				if err := json.Unmarshal([]byte(data), &temp); err == nil {
					t.Errorf("want %v, got %v", "error unmarshaling bad buid", nil)
				}
			})
			/* -------------------------------------------------------------------- */
			t.Run("Status", func(t *testing.T) {
				data := fmt.Sprintf(testJSON, productStatus.BUID, "golf", "", "", "", "")
				temp := &Status{}
				if err := json.Unmarshal([]byte(data), &temp); err == nil {
					t.Errorf("want %v, got %v", "error unmarshaling bad buid", nil)
				}
			})
			/* -------------------------------------------------------------------- */
			t.Run("Game", func(t *testing.T) {
				amount := MaxGameStatusDataLength + 1
				var b strings.Builder
				for i := 0; i < amount; i++ {
					fmt.Fprintf(&b, "%s", "a")
				}
				bad := b.String()
				pattern := `{"buid": %q, "status":"online", "main_game_status": %q}`
				data := fmt.Sprintf(pattern, productStatus.BUID, bad)
				temp := productStatus
				err := json.Unmarshal([]byte(data), temp)
				if err == nil {
					t.Errorf("want %v, got %v", "error unmarshaling bad json", nil)
				}
				if err != ErrMaxGameStatusDataLength {
					t.Errorf("want %s, got %s", ErrMaxGameStatusDataLength, err)
				}
				if temp.Game != nil {
					if len(*temp.Game) > MaxGameStatusDataLength {
						t.Errorf("want err %v, got %v", MaxGameStatusDataLength, len(*temp.Game))
					}
					if *temp.Game == bad {
						t.Errorf("want err %s, got %v", ErrMaxGameStatusIDLength.Error(), temp.Game)
					}
				}
			})
			/* -------------------------------------------------------------------- */
			t.Run("Extended", func(t *testing.T) {
				amount := MaxExtendedGameStatuses + 1
				bad := []ExtData{}
				for i := 0; i < amount; i++ {
					tx := ExtData{ID: fmt.Sprintf("attr_%v", i), Arg: []interface{}{1, 2}}
					bad = append(bad, tx)
				}
				b, _ := json.Marshal(bad)
				pattern := `{"buid": %q,"status":"online","extended_game_status":%s}`
				data := fmt.Sprintf(pattern, productStatus.BUID, b)
				temp := productStatus
				err := json.Unmarshal([]byte(data), temp)
				if err == nil {
					t.Errorf("want %v, got %v", "error unmarshaling bad json", nil)
				}
				if err != ErrMaxExtended {
					t.Errorf("want %s, got %s", ErrMaxExtended, err)
				}
				if temp.Extended != nil {
					if len(*temp.Extended) > MaxExtendedGameStatuses {
						t.Errorf("want err %v, got %v", MaxExtendedGameStatuses, len(*temp.Extended))
					}
				}
			})
			/* -------------------------------------------------------------------- */
			t.Run("Custom", func(t *testing.T) {
				amount := MaxCustomDataBytes + 100
				var b strings.Builder
				for i := 0; i < amount; i++ {
					fmt.Fprintf(&b, "%s", "a")
				}
				bad := b.String()
				pattern := `{"buid": %q,"status":"online","custom_data":%q}`
				data := fmt.Sprintf(pattern, productStatus.BUID, bad)
				temp := &Status{}
				err := json.Unmarshal([]byte(data), temp)
				if err == nil {
					t.Errorf("want %v, got %v", "error unmarshaling bad json", nil)
				}
				if err != ErrMaxCustomBytes {
					t.Errorf("want %s, got %s", ErrMaxCustomBytes, err)
				}
				if temp.Custom != nil {
					if len(*temp.Custom) > MaxCustomDataBytes {
						t.Errorf("want err %v, got %v", MaxCustomDataBytes, len(*temp.Custom))
					}
				}
			})
			/* -------------------------------------------------------------------- */
		})
		/* ---------------------------------------------------------------------- */
	})
}

type Frame struct {
	step int
	temp []byte
	err  error

	start time.Time
	stop  time.Time
}

func (f *Frame) Eval(s *Status, from, to Kind) error {
	if s.Enum != to {
		return fmt.Errorf("got %d (%s); want %d (%s)", s.Enum, s.Enum, to, to)
	}
	switch s.Enum {
	case Offline:
		return f.checkOffline(s)
	case AppearOffline:
		return f.checkAppearOffline(s)
	case Online:
		return f.checkOnline(s)
	case Idle:
		return f.checkIdle(s)
	case DND:
		return f.checkDND(s)
	default:
		return fmt.Errorf("invalid enum: %d=%s", s.Enum, s.Enum)
	}
}

func (f *Frame) checkOnline(s *Status) error {
	b, err := json.Marshal(&s)
	if err != nil {
		return err
	}
	temp := &Status{}
	if err := json.Unmarshal(b, &temp); err != nil {
		return err
	}
	if temp.Enum != Online {
		return fmt.Errorf("online: decoding err: got %d (%s); want %d (%s)", temp.Enum, temp.Enum, s.Enum, s.Enum)
	}
	if temp.Custom == nil {
		return fmt.Errorf("online: data err: missing custom data; got %v; want %v", nil, s.Custom)
	}
	return nil
}

func (f *Frame) checkIdle(s *Status) error {
	if s.Enum != Idle {
		return fmt.Errorf("idle: set err: got %d (%s); want %d (%s)", s.Enum, s.Enum, Idle, Idle)
	}
	b, err := json.Marshal(&s)
	if err != nil {
		return err
	}
	temp := &Status{}
	if err := json.Unmarshal(b, &temp); err != nil {
		return err
	}
	if temp.Enum != Idle {
		return fmt.Errorf("idle: decoding err: got %d (%s); want %d (%s)", temp.Enum, temp.Enum, s.Enum, s.Enum)
	}
	if temp.Custom == nil {
		return fmt.Errorf("idle: data err: missing custom data; got %v; want %v", nil, s.Custom)
	}
	return nil
}

func (f *Frame) checkDND(s *Status) error {
	b, err := json.Marshal(&s)
	if err != nil {
		return err
	}
	temp := &Status{}
	if err := json.Unmarshal(b, &temp); err != nil {
		return err
	}
	if temp.Enum != DND {
		return fmt.Errorf("dnd: decoding err: got %d (%s); want %d (%s)", temp.Enum, temp.Enum, s.Enum, s.Enum)
	}
	if temp.Custom == nil {
		return fmt.Errorf("idle: data err: missing custom data; got %v; want %v", nil, s.Custom)
	}
	return nil
}

func (f *Frame) checkAppearOffline(s *Status) error {
	b, err := json.Marshal(&s)
	if err != nil {
		return err
	}
	temp := &Status{}
	if err := json.Unmarshal(b, &temp); err != nil {
		return err
	}
	switch temp.Enum {
	case Offline:
		return nil
	case AppearOffline:
		return fmt.Errorf("appear-offline: data err: got %d (%s); want %d (%s)", temp.Enum, temp.Enum, Offline, Offline)
	default:
		return fmt.Errorf("appear-offline: decoding err: got %d (%s); want %d (%s)", temp.Enum, temp.Enum, s.Enum, s.Enum)
	}
}

func (f *Frame) checkOffline(s *Status) error {
	b, err := json.Marshal(&s)
	if err != nil {
		return err
	}
	temp := &Status{}
	if err := json.Unmarshal(b, &temp); err != nil {
		return err
	}
	if temp.Enum != Offline {
		return fmt.Errorf("offline: decoding err: got %d (%s); want %d (%s)", temp.Enum, temp.Enum, s.Enum, s.Enum)
	}
	if temp.Custom != nil {
		return fmt.Errorf("offline: data err: got %v; want %v", temp.Custom, nil)
	}
	return nil
}

func (f *Frame) Done() {
	f.stop = time.Now()
	if EnableJSONOutput {
		delta := f.stop.Sub(f.start)
		data := map[string]interface{}{
			"ok":       f.err == nil,
			"duration": delta.String(),
			"start":    f.start,
			"stop":     f.stop,
			"step":     f.step,
		}
		if f.err != nil {
			data["err"] = f.err
		}
		b, _ := json.Marshal(data)
		log.Printf("RESULT=%s", b)
	}
}

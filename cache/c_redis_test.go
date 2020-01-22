package cache

import (
	"log"
)

type Test struct {
	Nome      string `json:"nome"`
	Sobrenome string `json:"sobrenome"`
}

func initializarConf() *ICacheClient {
	arg := make(map[string]interface{})
	arg["DB"] = 1
	confi := OptionsCacheClient{
		Host:     "127.0.0.1:6379",
		Password: "",
		Driver:   "redis",
		Args:     arg,
	}
	cache, err := confi.ConfiguraCache()
	if err != nil {
		log.Println(err)
	}
	return cache
}

// func TestCacheRedisCon(t *testing.T) {
// 	arg := make(map[string]interface{})
// 	arg["DB"] = 1
// 	confi := OptionsCacheClient{
// 		Host:     "127.0.0.1:6379",
// 		Password: "",
// 		Driver:   "redis",
// 		Args:     arg,
// 	}
// 	_, err := confi.ConfiguraCache()
//
// 	if err != nil {
// 		t.Error("Expected nothing, got ", nil)
// 	}
//
// }

// func TestGetValues(t *testing.T) {
// 	cache := *initializarConf()
// 	test := &Test{}
//
// 	error := cache.GetValue("vicente", test)
// 	// time.Sleep(time.Second * 30)
//
// 	log.Println(test.Sobrenome)
// 	if error != nil {
// 		t.Error("Expected nothing, got ", error)
// 	}
//
// }
//
// func TestSetValues(t *testing.T) {
// 	cache := *initializarConf()
// 	test := &Test{}
// 	test.Nome = "ronaldo"
// 	test.Sobrenome = "klebernilton"
//
// 	error := cache.SetValue("vicente", test, 0)
// 	// time.Sleep(time.Second * 30)
//
// 	if error != nil {
//
// 		t.Error("Expected nothing, got ", error)
// 	}
// }
//
// func TestDelValues(t *testing.T) {
// 	cache := *initializarConf()
// 	error := cache.DelValue("ronaldo")
// 	// time.Sleep(time.Second * 30)
//
// 	if error != nil {
//
// 		t.Error("Expected nothing, got ", error)
// 	}
// }
//
// func TestAddtValues(t *testing.T) {
// 	cache := *initializarConf()
// 	test := &Test{}
// 	test.Nome = "ronaldo"
// 	test.Sobrenome = "klebernilton"
// 	test1 := &Test{}
// 	test1.Nome = "vicente"
// 	test1.Sobrenome = "klebernilton"
// 	tests := make([]*Test, 2)
// 	tests[0] = test
// 	tests[1] = test1
// 	error := cache.AddValues("vicenteList8", &tests)
//
// 	// time.Sleep(time.Second * 30)
//
// 	if error != nil {
// 		t.Error("Expected nothing, got ", error)
//
// 	}
// }
//
// func TestGetListValues(t *testing.T) {
// 	cache := *initializarConf()
// 	test := &Test{}
//
// 	tests := make([]Test, 1)
// 	tests = append(tests, *test)
//
// 	error := cache.GetListValues("vicenteList8", &tests)
//
// 	// time.Sleep(time.Second * 30)
// 	log.Println(tests)
// 	if error != nil {
// 		t.Error("Expected nothing, got ", error)
// 	}
// }

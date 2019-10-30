package mylib

type WalletMap struct {
	Data map[string]string
}

func New() *WalletMap { return new(WalletMap).Init() }

func (wmap *WalletMap) Init() *WalletMap {

	wmap.Data = make(map[string]string)

	return wmap
}

func (wmap *WalletMap) Put(key string, value string) {
	wmap.Data[key] = value
}

func (wmap *WalletMap) PutEntries(map1 map[string]string) {
	for key, value := range map1 {
		wmap.Data[key] = value
	}
}

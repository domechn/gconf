//Date: 2018/8/22 下午3:19
//
//Description:
package gconf

import (
	"sync"
	"strings"
	"errors"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"bytes"
	"log"
	"gconf/goyml2"
)

func ReadYmlReader(path string) (cnf map[string]interface{}, err error) {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}

	return parseYML(buf)
}

func parseYML(buf []byte) (cnf map[string]interface{}, err error) {
	if len(buf) < 3 {
		return
	}

	if string(buf[0:1]) == "{" {
		log.Println("Look like a Json, try json umarshal")
		err = json.Unmarshal(buf, &cnf)
		if err == nil {
			log.Println("It is Json Map")
			return
		}
	}

	data, err := goyml2.Read(bytes.NewBuffer(buf))
	if err != nil {
		log.Println("Goyaml2 ERR>", string(buf), err)
		return
	}

	if data == nil {
		log.Println("Goyaml2 output nil? Pls report bug\n" + string(buf))
		return
	}
	cnf, ok := data.(map[string]interface{})
	if !ok {
		log.Println("Not a Map? >> ", string(buf), data)
		cnf = nil
	}
	cnf = expandValueEnvForMap(cnf)
	return
}

type yamlConfig struct {
}

func (yc *yamlConfig) parse(filepath string) (Configer, error) {
	cnf, err := ReadYmlReader(filepath)
	if err != nil {
		return nil,err
	}
	y := &yamlConfigContainer{
		data: cnf,
	}
	return y , nil
}

func (yc *yamlConfig) parseData(data []byte) (Configer, error) {
	cnf, err := parseYML(data)
	if err != nil {
		return nil, err
	}

	return &yamlConfigContainer{
		data: cnf,
	}, nil
}

type yamlConfigContainer struct {
	data map[string]interface{}
	sync.RWMutex
}

func (ycr *yamlConfigContainer) String(key string) string {
	if v, err := ycr.getData(key); err == nil {
		if vv, ok := v.(string); ok {
			return vv
		}
	}
	return ""
}

func (ycr *yamlConfigContainer) Strings(key string) []string {
	v := ycr.String(key)
	if v == "" {
		return nil
	}
	return strings.Split(v, ";")
}

func (ycr *yamlConfigContainer) Int(key string) (int, error) {
	if v, err := ycr.getData(key); err != nil {
		return 0, err
	} else if vv, ok := v.(int); ok {
		return vv, nil
	} else if vv, ok := v.(int64); ok {
		return int(vv), nil
	}
	return 0, errors.New("not int value")
}

func (ycr *yamlConfigContainer) Int64(key string) (int64, error) {
	if v, err := ycr.getData(key); err != nil {
		return 0, err
	} else if vv, ok := v.(int64); ok {
		return vv, nil
	}
	return 0, errors.New("not bool value")
}

func (ycr *yamlConfigContainer) Bool(key string) (bool, error) {
	v, err := ycr.getData(key)
	if err != nil {
		return false, err
	}
	return parseBool(v)
}

func (ycr *yamlConfigContainer) Float(key string) (float64, error) {
	if v, err := ycr.getData(key); err != nil {
		return 0.0, err
	} else if vv, ok := v.(float64); ok {
		return vv, nil
	} else if vv, ok := v.(int); ok {
		return float64(vv), nil
	} else if vv, ok := v.(int64); ok {
		return float64(vv), nil
	}
	return 0.0, errors.New("not float64 value")
}

func (ycr *yamlConfigContainer) DefaultString(key string, defaultval string) string {
	v := ycr.String(key)
	if v == "" {
		return defaultval
	}
	return v
}

func (ycr *yamlConfigContainer) DefaultStrings(key string, defaultval []string) []string {
	v := ycr.Strings(key)
	if v == nil {
		return defaultval
	}
	return v
}

func (ycr *yamlConfigContainer) DefaultInt(key string, defaultval int) int {
	v, err := ycr.Int(key)
	if err != nil {
		return defaultval
	}
	return v
}

func (ycr *yamlConfigContainer) DefaultInt64(key string, defaultval int64) int64 {
	v, err := ycr.Int64(key)
	if err != nil {
		return defaultval
	}
	return v
}

func (ycr *yamlConfigContainer) DefaultBool(key string, defaultval bool) bool {
	v, err := ycr.Bool(key)
	if err != nil {
		return defaultval
	}
	return v
}

func (ycr *yamlConfigContainer) DefaultFloat(key string, defaultval float64) float64 {
	v, err := ycr.Float(key)
	if err != nil {
		return defaultval
	}
	return v
}

func (ycr *yamlConfigContainer) Interface(key string) (interface{}, error) {
	return ycr.getData(key)
}

func (ycr *yamlConfigContainer) getData(key string)(interface{}, error){
	if key == "" {
		return nil, errors.New("key is empty")
	}
	ycr.RLock()
	defer ycr.RUnlock()

	keys := strings.Split(key, ".")
	tmpData := ycr.data
	for idx, k := range keys {
		if v, ok := tmpData[k]; ok {
			switch v.(type) {
			case map[string]interface{}:
				{
					tmpData = v.(map[string]interface{})
					if idx == len(keys) - 1 {
						return tmpData, nil
					}
				}
			default:
				{
					return v, nil
				}

			}
		}
	}
	return nil, fmt.Errorf("not exist key %q", key)
}

func newYamlConfigContainer(filePath string) (Configer, error) {
	y := yamlConfig{}
	return y.parse(filePath)
}

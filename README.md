# gconf

将beego中的配置文件读取功能提取出来

c.json
```
{
    "ABC":"a",
    "CCC":12,
    "B":"abc;asc;as"
    "DD":{
        "DDD":23
    }
}
```

```
func init(){
    err := gconf.Register("conffile","/abc/c.json||./c.json||abc/c.json")
    if err != nil {
        panic(err)
    }
}


func getConf(){
    conffile,err := gconf.GetConfiger("conffile")
    if err != nil {
        panic(err)
    }
    conffile.String("ABC")  => "a"
    conffile.Int("CCC")   => 12
    conffile.Strings("B") => []string{"abc","asc","as"}
    conffile.DefaultString("a","bc") => "bc"
    conffile.Int("DD.DDD") => 23
}
```


c.yaml

```
ABC : 12
AC :
    BC : 23
BCD :
    - "a"
    - "b"
    - "c"
```

```
    func init(){
        err := gconf.Register("cyaml",c.yaml")
        if err != nil {
                panic(err)
            }
    }
    
    func getConfig(){
        cyaml, err := gconf.GetConfiger("cyaml")
        if err != nil {
            painc(err)
        }
        cyaml.Int("ABC")  => 12
        cyaml.Int("AC.BC") => 23
        cyaml.Strings("BCD") => []string{"a","b","c"}
    }
```
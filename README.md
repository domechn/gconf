# gconf

将beego中的配置文件读取功能提取出来

c.json
```
{
    "ABC":"a",
    "CCC":12,
    "B":"abc;asc;as"
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
}
```

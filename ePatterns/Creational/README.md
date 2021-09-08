# 创造模式

|模式|描述|支持情况|
|:-:|:-:|:-:|
|Abstract Factory 抽象工厂|提供用于创建已释放对象系列的接口|F|
|Builder 建造者|使用简单对象构建复杂对象|T|
|Factory Method 工厂方法|将对象的实例化推迟到专门用于创建实例的函数|T|
|Object Pool 对象池|实例化和维护一组相同类型的对象实例|T|
|Singleton 单身|将类型的实例化限制为一个对象|T|

## Builder 建造者

``` go
package car 

type Speed float64

const (
    BlueColor Color = "blue"
    GreenColor      = "green"
    ReadColor       = "red"
)

type Wheels string

const (
    SportsWheels Wheels = "sports"
    SteelWheels         = "steel"
)
```

## Factory Method 工厂方法 
## Object Pool 对象池
## Singleton 单身

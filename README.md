## 简介

这是一个 vue + go 实现的校园交易系统

示例地址如下 **<H2>c.214233.xyz</H2>**

## 一键部署
### docker
> docker run -d -p 8080:8080 twoonefour1/campustrader

## 项目目录
```
CampusTrader/
├── frontend
├── internal
```
- frontend vue 前端
- internal, go后端

无数据库抽象层，controller + service 两层结构加个DAO，service层直接操作db对象, 

## 截止当前进度

目前已经实现：
- 用户登录/注册及接口鉴权
- 购买/发布商品
- 流式分页查询 商品+好评率显示 （业务内两次分查询实现，而非视图或复杂查询）
- 最热门分类 （复杂查询）
- 模糊查询 (甚至有搜索建议)
- 业务层商品下架或交易完成log记录（非数据库触发器）
- 分类查看

未实现：
- 用户管理，收藏夹，聊天记录

为什么没实现？

**因为光是实现上面，代码量就要爆了，哥们写了2000多行了，github提交战绩可查**

## 快速使用
```
git clone https://github.com/twoonefour/campustrader.git
cd campustrader
cd frontend && npm build && cd .. && go build -o cmd/main.go campustrader && ./campustrader
```
需要填写数据库变量 `DATABASE_DSN` 为go dsn 格式
## 实现难点

### 难点1 数据库设计

#### 表间关联

其中`user`主键很显然和`order` `product`都有外键约束，保证2nf，避免user删掉，order和product脏数据的问题

`product`中对应user表中卖家id的外键，对应商品类别category id的外键

`order`又和`product`有约束，理由同上，购买一个产品的订单，需要对应product的id

`review`又必须有两个外键，一个卖家user外键，一个卖家user外键

差不多就这些，这种关系函数表达我不是很会写，因为我菜

#### 表间约束
product表内的限制
```
Seller User `gorm:"foreignKey:SellerId;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
```
目前是只有这个限制

级联更新，为什么是限制删除呢？

~~我也不懂，跟我的GPT说去吧~~ 

保证当userid修改的时候，关联的订单一并修改，当user被删除的时候，要是出现了很多依赖数据，需要人工介入，防止数据库数据崩坏

### 难点2
与其说是难点，不如说是我不能接受

**视图+触发器**

为什么我看到题目第一反应是：为什么要写在数据库层？

一般肯定是业务层拦截，直接调用db增加一条log啊

要是你新来了个打工的实习哥们，你忘记跟他说，噢，这里有什么什么trigger，不用手动插一个日志到log表里面

他反手写一个 `db.Model(&model.review{}).Create(&log)`, 表里数据直接乘2，要是项目组长一问谁干得，你不炸了

**总之，业务层就是要显式插入log**，而且触发器性能也很差

当然玩归玩闹归闹，bro还是写了的，只是没往我go代码里加

#### 触发器sql
商品卖出后触发逻辑，插入log
```sql
create trigger after_product_sold
after UPDATE on orders
for each row
begin
    if new.status = 'completed' then
        INSERT INTO product_sold_logs(product_id, buyer_id, seller_id, price)
        VALUES (new.product_id, new.buyer_id, new.seller_id, new.amount);
    end if
end;
```

~~商品下架触发逻辑，插入log~~

数据库层一点都不好写啊，相比起来如果我在业务层就可以传 下架原因，操作人等等参数了，如果只有记录下架时间+下架商品id那就还好

以下为下架时间的trigger，只需要传product_id,mysql会自动生成时间戳
```sql
create trigger after_product_drop
after delete on products
for each row
begin
    INSERT INTO product_drop_log(product_id)
    VALUES (new.product_id);
end;
```

#### 视图sql
```sql
create or replace view v_product as
        select
            products.*, users.nickname,
            users.phone,
            IFNULL(t.avg_rating, 0),
            IFNULL(t.review_count, 0)
        from products
        join users on products.seller_id = users.id
        LEFT JOIN (
            SELECT
                target_user_id,
                AVG(rating) AS avg_rating,
                COUNT(*) AS review_count
            FROM reviews
            GROUP BY target_user_id
        ) t ON t.target_user_id = users.id;
    select * from v_product;
```

### 难点3
哥们一个后端，完全不会前端啊

> 要求6、前端界面：开发简单的图形界面，实现核心数据的增删改查和关键业务查询功能。

叽里咕噜说什么呢，跟我的ai说去吧


### 难点4
存储过程

分类查询
```
CREATE PROCEDURE sp_search_and_count_by_category(
    IN p_keyword VARCHAR(100), 
    IN p_category_id BIGINT UNSIGNED,
    OUT p_total_count INT
)
BEGIN
    -- 获取总数
    SELECT COUNT(*) INTO p_total_count FROM products 
    WHERE name LIKE CONCAT('%', p_keyword, '%') AND category_id = p_category_id AND status = 'available';
    -- 获取结果
    SELECT * FROM products 
    WHERE name LIKE CONCAT('%', p_keyword, '%') AND category_id = p_category_id AND status = 'available'
    ORDER BY created_at DESC;
END;
```

更新卖出产品

```
CREATE PROCEDURE sp_complete_order(IN p_order_id BIGINT UNSIGNED)
BEGIN
    DECLARE EXIT HANDLER FOR SQLEXCEPTION ROLLBACK;
    START TRANSACTION;
        UPDATE orders SET status = 'completed', completed_at = NOW() WHERE id = p_order_id;
        UPDATE products SET status = 'sold' WHERE id = (SELECT product_id FROM orders WHERE id = p_order_id);
    COMMIT;
END;
```

存储过程两条：
```sql
	CREATE PROCEDURE sp_search_and_count_by_category(
		IN p_category_id BIGINT UNSIGNED
	)
	BEGIN
		SELECT * FROM products
		where 
			category_id=p_category_id 
			AND status='available';
	END;
```

```
	CREATE PROCEDURE sp_complete_order(IN p_order_id BIGINT UNSIGNED)
	BEGIN
		UPDATE orders SET status = 'completed', updated_at = NOW() WHERE id = p_order_id;
		UPDATE products SET status = 'sold' WHERE id = (SELECT product_id FROM orders WHERE id = p_order_id);
	END;
```

**在我用这个存储过程的时候bug就来了，因为我是浮标分页，用存储过程原生执行sql，就无法分页，必须要把分页逻辑写进存储过程**
<img width="1279" height="626" alt="image" src="https://github.com/user-attachments/assets/512167b6-231e-4097-8073-18a19c3fe018" />

这里如果调这个接口，无论怎么翻页，都只会返回第一页，因为lastId没有传给存储过程，必须要传，你要传，你就要改存储过程

那你如果再来个需求，我要查询某个用户的某个分类的商品，那你是不是又要改一次存储过程？

哥们，那我一行直接判断一下有没有用户id，然后加一句`.Where('seller_id = ?', user_id)`不比这个强？

狗都不用


## 稿子思路

### 3nf设计是冗余/非冗余传统香烟和电子烟

第一范式，第二范式自不用说，必然是满足的，打破了没有任何好处，本仓库表都是一个表一个主键且不可拆分

问题在于3nf（非主键字段不能相互依赖）

试想数据库的存储逻辑，若执行这样一个语句会发生什么呢

```sql
select email from user where username = '丰川祥子';
```

答案是，没有主键，没有索引，直接触发全表扫描

为什么？mysql底层回想一下就能想明白为什么，想想b+树存什么/怎么存的，你就会得出结论：**隐式回表多次，mysql优化器会直接判断全表比回表多次快很多，触发全表扫**

要解决这个问题该怎么办呢？

可以建索引，username直接建索引(username, email)即可，这里获取数据就是从 `username -> id -> (id, username, email)` 变成了索引后 `(username, email)` 一次就能拿到

那如果我要查 `某个用户评价如何`，那索引就不管用了

在3nf范式下得查两个表: order review，（review得到orderid,orderid取seller_id, seller_id 可能你还要回去取用户名称）

那如果我打破了3nf，会怎么样呢，我需要看某个用户评价，直接 `select avg(rating) from order where target_user_id = 123 group by target_user_id;` 

这里的查询过程从两个跨表查询 order_id -> seller_id，变成了一个单表查询，这其中的速度提升可想而知

并且我是每个商品都要展示用户评价，相比于每个商品都跨表join O(m * n)的时间复杂度，去除3nf的时间复杂度为 O(n), 这种提升是显而易见的

还有一种可能是，3nf的快照问题，假设 有个场景需要统计产品销量，某个用户在下单后，删号了，统计销量就会受到影响，因为某个订单中含有这个用户的订单，但实际上是需要统计这一条数据的，这种情况下打破3nf是有意义的

那什么是最有意义的3nf呢？比如最需要数据安全的

例如，某个供应商品列表如下

<img width="627" height="163" alt="image" src="https://github.com/user-attachments/assets/b2296d15-2619-4e69-a4b8-3414ae650dc7" />

这里间接依赖了非主键供应商id，会出现三大问题，修改插入删除，在这种情况下，3范式则是必须的，因为要**维护数据的唯一真实来源**

### 触发器/视图/存储过程，业务与底层数据的思考

**前言，触发器/视图/存储过程，我认为已经是被时代淘汰了，没必要强行用**

触发器，**数据改动日志，数据校验**

视图，**权限控制，隐藏底层n个join的复杂实现**

存储过程，**把很多sql糅杂在一起实现一个特定的逻辑函数**，或者直接实现一个输入输出，让业务层一行调用

如果非要用，不妨问问自己，为什么一定要在数据库层做呢？上面这些能否可以在业务层做？

**答案是肯定的**

触发器与视图，不光维护麻烦，调用者还需要去思考该怎么用

难道不符合我的逻辑，那我又要改掉触发器/视图/存储过程？这显然是不可能的

像这个例子，我刚开始考虑的时候，只觉得我要搜某个类别，直接传给触发器一行调用。咋一看觉得简单好用，那如果我要搜某个用户发布的某个类别的产品呢？

```sql
	CREATE PROCEDURE sp_search_and_count_by_category(
		IN p_category_id BIGINT UNSIGNED
	)
	BEGIN
		SELECT * FROM products
		where 
			category_id=p_category_id 
			AND status='available';
	END;
```

**你是不是要连进数据库改存储过程了？还要给他重新取个名，因为你不能直接改这个，也许别人在用呢**

触发器和视图都是同样的道理，说白了就是存储过程符合抽象设计概念，触发器不符合显式设计概念，而视图，可以用索引代替，就算为了强行满足3NF去使用视图，**某些性能方面还是不如破坏3nf+索引的场景**

**底层数据就应该是底层数据，而不是带上业务逻辑**

### 索引优化

**前言，where查询多的地方，给他上索引，能优化覆盖则覆盖，不能则考虑不从索引做优化，这才是正确的，不要死磕数据库层**

只分析product，因为这个实现只有product和categories相关的，查询多的场景，涉及数据库底层存储逻辑

#### product
```
type Product struct {
	Id          uint64         `gorm:"column:id;type:BIGINT UNSIGNED;primaryKey;" json:"id"`
	Name        string         `gorm:"column:name;type:VARCHAR(100);comment:商品名称;not null;index:idx_name_ft,class:FULLTEXT;" json:"name"`     // 商品名称
	Description string         `gorm:"column:description;type:TEXT;comment:商品描述;not null;" json:"description"`                                // 商品描述
	Price       float64        `gorm:"column:price;type:DECIMAL(10, 2);comment:价格;not null;" json:"price"`                                    // 价格
	CategoryId  uint64         `gorm:"column:category_id;type:BIGINT UNSIGNED;comment:分类ID;not null;index:idx_cat_status" json:"category_id"` // 分类ID
	SellerId    uint64         `gorm:"column:seller_id;type:BIGINT UNSIGNED;comment:卖家ID;not null;" json:"seller_id"`                         // 卖家ID
	Category    Category       `gorm:"foreignKey:CategoryId" json:"category,omitempty"`
	Seller      User           `gorm:"foreignKey:SellerId;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"seller"`
	Status      string         `gorm:"column:status;type:ENUM('available', 'sold', 'removed');comment:状态;default:available;index:idx_cat_status;index:idx_status" json:"status"` // 状态
	Condition   string         `gorm:"column:condition;type:ENUM('new', 'like_new', 'good', 'fair', 'poor');comment:新旧程度;not null;" json:"condition"`                            // 新旧程度
	ImageUrl    string         `gorm:"column:image_url;type:VARCHAR(255);comment:主图URL;" json:"image_url"`                                                                       // 主图URL
	CreatedAt   time.Time      `gorm:"column:created_at;type:TIMESTAMP;default:CURRENT_TIMESTAMP;" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"column:updated_at;type:TIMESTAMP;default:CURRENT_TIMESTAMP;" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
```

**以下两个sql压力可能会比较大**

一个是模糊查询，如 `like %商品名%`，这种一定会触发全表扫，数据量少还好，数据量一多就容易超时卡死，直接给他上全文索引，或者你上elasticsearch直接解决问题

另一个是根据类别查询，由于非主键查询一定会回表两次，这种就比较复杂，你必须要想明白，*前端要展示的东西是什么，把这些东西列出来，给他单独索引*

这样问题还没解决，**根据最左匹配，你只能查 where category_id = xxx and status = xxx， 这种就不会触发全表扫**

由于我还有个复杂查询，连order表找完成订单，去统计最热分类top3，这种情况就太复杂了，你不能考虑索引一次直接覆盖所有你要的数据，在数据库层两个表会瓶颈，你要么再来一个表统计，要么用sql8.0+窗口函数

但无论怎么样，你category_id是肯定要加索引的，然后order表因为要根据category_id分组，你也要加索引

总之上面这个问题，`(category_id)`，是肯定没错的，由于统计后需要展示种类名称，你再连个category表

**还有在主页面，你要展示的一定是上架的商品，即status = 'available'**, 虽说也会回表，但在这个status也是可以索引加速的，直接给他上索引

**总结**

name全文索引，status单独索引，category_id单独索引

### 锁与事务
本设计没有加任何显式锁，由于go没有乐观锁，我直接靠数据库层的乐观锁+数据库事务保证数据一致

在订单生成中，我的实现如下
```
	err := db.Transaction(func(tx *gorm.DB) error {
		// 乐观锁
		res := tx.Model(&model.Product{}).Where("id=? AND status='available'", itemID).Update("status", "sold")
		if res.RowsAffected == 0 {
			return errors.New("已售出")
		}
		return tx.Create(&model.Order{
			ProductId: item.Id,
			BuyerId:   buyerID,
			Amount:    item.Price,
			SellerId:  item.SellerId,
		}).Error
	})
```

go代码可能很多人没看过，简单说就是一个乐观锁，用status状态原子更新，最后提交事务订单，在数据库默认**RR**（rc以上）的事务状态下，这是可行的，**因为会命中主索引id，产生行锁，防止修改**

**两个操作必须同时完成，所以必须用事务包裹，由于并发问题，可能别人已经购买过这个订单了，可以用乐观锁去尝试修改订单，让他返回错误**

高阶实现你可以加redis lua原子解锁，速度会快很多，但这只是一个小作业，我就不上了


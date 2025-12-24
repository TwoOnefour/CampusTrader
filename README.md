## 简介

这是一个 vue + go 实现的校园交易系统

示例地址如下 **<H2>c.214233.xyz</H2>**


## 截止当前进度

目前已经实现：
- 用户登录/注册及接口鉴权
- 购买/发布商品
- 商品+好评率显示 （业务内两次分查询实现，而非视图或复杂查询）
- 最热门分类 （复杂查询）
- 模糊查询
- 业务层商品下架或交易完成log记录（非数据库触发器）

未实现：
- 用户管理，收藏夹，聊天记录
- 分类浏览

为什么没实现？

**因为光是实现上面，代码量就要爆了，哥们写了2000多行了，github提交战绩可查**

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

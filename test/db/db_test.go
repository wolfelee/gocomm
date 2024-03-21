package test

import (
	"fmt"
	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/wolfelee/gocomm/pkg/cache"
	"math/rand"
	"strconv"
	"testing"
	"time"
	"xorm.io/xorm"
)

var (
	cacheUser1RealIdPrefix = "cache#user1#realId#"
)

type (
	User1Model struct {
		cache.DataBase
		table string
	}

	User1 struct {
		RealId     string    `xorm:"real_id pk" json:"real_id"`
		Number     string    `xorm:"number" json:"number"`     // 学号
		Name       string    `xorm:"name" json:"name"`         // 用户名称
		Password   string    `xorm:"password" json:"password"` // 用户密码
		Gender     string    `xorm:"gender" json:"gender"`     // 男｜女｜未公开
		CreateTime time.Time `xorm:"create_time" json:"create_time"`
		UpdateTime time.Time `xorm:"update_time" json:"update_time"`
	}
)

func (*User1) TableName() string {
	return "user1"
}

func (m *User1Model) Insert(data User1) error {
	return m.Exec(func(session *xorm.Session) error {
		_, err := session.Insert(data)
		return err
	})
}

func (m *User1Model) FindOne(realId string) (User1, bool, error) {
	user1RealIdKey := fmt.Sprintf("%s%v", cacheUser1RealIdPrefix, realId)
	var resp User1

	has, err := m.QueryRow(user1RealIdKey, &resp, func(session *xorm.Session, v interface{}) (bool, error) {
		return session.Where("`real_id` = ?", realId).Get(&resp)
	})
	fmt.Println("--->", err)
	if err != nil {
		return resp, false, err
	}
	return resp, has, nil
}

func (m *User1Model) Update(data User1) error {
	user1RealIdKey := fmt.Sprintf("%s%v", cacheUser1RealIdPrefix, data.RealId)
	return m.Exec(func(session *xorm.Session) error {
		_, err := session.ID(data.RealId).Update(&data)
		return err
	}, user1RealIdKey)
}

func (m *User1Model) Delete(realId string) error {
	user1RealIdKey := fmt.Sprintf("%s%v", cacheUser1RealIdPrefix, realId)
	return m.Exec(func(session *xorm.Session) error {
		data := new(User1)
		_, err := session.ID(realId).Delete(data)
		return err
	}, user1RealIdKey)
}

func NewUser1Model(x *xorm.Session, r *redis.Client) *User1Model {
	return &User1Model{
		DataBase: cache.NewDataBaseTest(x, r),
		table:    "user1",
	}
}

func initClient() (redisdb *redis.Client, err error) {
	redisdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // post端口
		Password: "zhy1996",        // 密码
		DB:       15,               // 使用redis的库
	})
	_, err = redisdb.Ping().Result()
	if err != nil {
		return
	}
	return
}

var eng *xorm.Engine
var redisDb *redis.Client

func InitModel() {
	var err error
	eng, err = xorm.NewEngine("mysql", "root:zhy1996@tcp(localhost:3306)/test")
	if err != nil {
		//t.Fatal(err)
		fmt.Println(err)
		return
	}

	redisDb, err = initClient()
	if err != nil {
		return
	}
	return

}

func TestMain(t *testing.M) {
	InitModel()
	t.Run()
}

func NewModel() *User1Model {
	return &User1Model{
		DataBase: cache.NewDataBaseTest(eng.NewSession(), redisDb),
		table:    "user1",
	}
}

func TestInsert(t *testing.T) {
	for i := 14100; i < 14200; i++ {
		go func(a int) {
			model := NewModel()
			defer model.CloseSession()
			u := User1{
				RealId:     strconv.Itoa(a),
				Number:     strconv.Itoa(a*a + 10086),
				Name:       fmt.Sprintf("%s-%d", "zhaoahiyu", a),
				Password:   fmt.Sprintf("%s|%d", "pass", a),
				Gender:     "?",
				CreateTime: time.Now(),
				UpdateTime: time.Now(),
			}
			//t.Log(model)
			//_ = u
			err := model.Insert(u)
			if err != nil {
				t.Error(err)
				return
			}
		}(i)
	}
	time.Sleep(time.Second)
}

func GetRandomString(l int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyz"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

func BenchmarkInsert(b *testing.B) {
	for i := 0; i < b.N; i++ {
		model := NewModel()
		defer model.CloseSession()
		u := User1{
			RealId:     GetRandomString(10),
			Number:     GetRandomString(8),
			Name:       GetRandomString(6),
			Password:   GetRandomString(16),
			Gender:     GetRandomString(1),
			CreateTime: time.Now(),
			UpdateTime: time.Now(),
		}
		//t.Log(model)
		//_ = u
		err := model.Insert(u)
		if err != nil {
			b.Error(err)
			return
		}
	}
}

func TestDelete(t *testing.T) {
	model := NewModel()
	defer model.CloseSession()
	err := model.Delete("zak8ubifvp")
	if err != nil {
		t.Fatal(err)
	}
}

func TestUpdate(t *testing.T) {
	model := NewModel()
	defer model.CloseSession()
	u := User1{
		RealId:     "zlf6luodvf",
		Number:     GetRandomString(8),
		Name:       GetRandomString(6),
		Password:   GetRandomString(16),
		Gender:     GetRandomString(1),
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}
	//t.Log(model)
	//_ = u
	err := model.Update(u)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestSelect(t *testing.T) {
	model := NewModel()
	defer model.CloseSession()

	tests := []string{"zzmfmofyzo", "zzc6k48jce", "zz7vkwsgus", "zz6tc2fuvk", "zy8ogmu1k9", "zy82g74ujp", "zy3r5dm4jv",
		"zy2k9wq558", "zxqqoou3el", "zxiaup1hu6", "zx710c07em", "zx6puhu9bw", "zwf1sm2cmy", "zvf3u4frs0", "zucxcmi0eb",
		"zt66saoi3z", "zstuqbl4k5", "zsmi2holjh", "zslk7dqp88", "zsdnbu0vlx", "zs7gspi54o", "zryg89bjty", "zqx762x0av",
		"zpij921wgb", "zowrv8emp3", "zo7yjpuq29", "znvhpcupyf", "znhdu69be9", "znfupf0d2k", "zncmtm4cw4", "zmyxcm8fw0",
		"zmsry00drw", "zli69qb92g", "zlfncockb8", "ziygx9ncwl", "zier57jygz", "zhetlkhzqv", "zfw9qsqmd6",
		"zfubolqovl", "zekdq8sme9", "zea9mgtshe", "zdlgp4okva", "zdk7tgcawg", "zdip35fhrf", "zd5me3cb3j", "zcc6qnsjel",
		"zc86q6y5nm", "zblozamz1o", "zb9sc79c89"}
	for _, test := range tests {
		t.Run(test, func(t *testing.T) {
			u, has, err := model.FindOne(test)
			if err != nil {
				t.Fatal(err)
			}
			if !has {
				t.Fatal("not have ")
			}
			t.Log(u)
		})
	}

	time.Sleep(time.Second)
}

package database

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/wendisx/puzzle/pkg/clog"
)

const (
	_test_default_dsn     = "root:ff5740f5-5083-4f71-9ac7-c13144c87b78@tcp(127.0.0.1:53306)/test?charset=utf8mb4"
	_test_ctx_tomeout     = 10 * time.Second
	_test_user_name_basic = "test_user_"
	_test_user_max_limit  = 10
	_test_user_id_basic   = 1
)

var (
	_test_mysql_db MysqlDB
	_test_sql_list = []string{
		`
		insert into user_basic(extern_id,user_name,user_password)
		values 
		(uuid_to_bin(uuid()),?,sha2(?, 256))
		`,
		`
		insert into user_basic(extern_id,user_name,user_password)
		values 
		(uuid_to_bin(uuid()),:user_name,sha2(user_password, 256))
		`,
		`
		insert into user_detail(user_id,nickname,phone,email,avatar)
		values 
		(?,?,?,?,?)
		`,
		`
		insert into user_detail(user_id,nickname,phone,email,avatar)
		values 
		(:user_id,:nickname,:phone,:email,:avatar)
		`,
		`
		update user_basic
		set user_password = sha(?, 256)
		where user_name = ? and deleted = 0
		`,
		`
		update user_basic
		set user_password = sha(:user_password, 256)
		where user_name = :user_name and deleted = 0
		`,
		`
		update user_detail
		set nickname = ?
		where user_id = ?
		`,
		`
		update user_basic
		set deleted = 1
		where user_name = ? and deleted = 0
		`,
		`
		delete from user_basic where deleted = 1 cascade
		`,
		`
		select user_name, user_password 
		from user_basic
		where user_name = ? and deleted = 0
		`,
		`
		select user_id, nickname, phone, email, avatar
		from user_detail
		where user_id = :user_id
		`,
		`
		select user_name, user_password from user_basic
		where user_name like ? and deleted = 0
		`,
		`
		select user_id, nickname phone, email, avatar
		from user_detail
		where user_id > :min_user_id
		`,
		`
		select user_id, nickname, phone, email, avatar
		from user_detail
		order by user_id desc
		limit ?,?
		`,
	}
)

type (
	UserBasic struct {
		MysqlMeta
		UserName     string `db:"user_name"`
		UserPassword string `db:"user_password"`
	}
	UserDetail struct {
		MysqlMeta
		UserId   uint64 `db:"user_id"`
		Nickname string `db:"nickname"`
		Phone    string `db:"phone"`
		Email    string `db:"email"`
		Avatar   string `db:"avatar"`
	}
)

// pass
func Test_connect_db(t *testing.T) {
	_test_mysql_db = InitMysql(_test_default_dsn)
}

// pass
func Test_insert_one(t *testing.T) {
	_test_mysql_db = InitMysql(_test_default_dsn)
	ctx, cancle := context.WithTimeout(context.Background(), _test_ctx_tomeout)
	defer cancle()
	// for user_basic
	for i := 0; i < _test_user_max_limit; i += 1 {
		newUser := UserBasic{
			UserName:     fmt.Sprintf("%s%d", _test_user_name_basic, i),
			UserPassword: fmt.Sprintf("%s%d", _test_user_name_basic, i),
		}
		var err error
		if i%2 == 0 {
			err = InsertWithPlace(ctx, _test_mysql_db, _test_sql_list[0], newUser.UserName, newUser.UserPassword)
		} else {
			err = InsertWithName(ctx, _test_mysql_db, _test_sql_list[1], newUser)
		}
		if err != nil {
			panic(fmt.Sprintf("%s", err.Error()))
		}
	}
	// for user_detail
	for i := 0; i < _test_user_max_limit; i += 1 {
		newDetail := UserDetail{
			UserId:   uint64(_test_user_id_basic + i),
			Nickname: fmt.Sprintf("test_nickname_%d", i),
			Phone:    fmt.Sprintf("phone_%d", i),
			Email:    fmt.Sprintf("email_%d", i),
			Avatar:   fmt.Sprintf("avatar_%d", i),
		}
		var err error
		if i%2 == 0 {
			err = InsertWithName(ctx, _test_mysql_db, _test_sql_list[3], newDetail)
		} else {
			err = InsertWithPlace(ctx, _test_mysql_db, _test_sql_list[2], newDetail.UserId, newDetail.Nickname, newDetail.Phone, newDetail.Email, newDetail.Avatar)
		}
		if err != nil {
			panic(fmt.Sprintf("%s", err.Error()))
		}
	}
}

/*
	update 和 delete 直接下放到实际的程序执行时测试, 因为逻辑和insert一致, 几乎不可能出现预期外的结果
*/
// pass? -- yes just like insert happen
func Test_update_one(t *testing.T) {
}

// pass? -- yes just like insert happen
func Test_delete_one(t *testing.T) {
	_test_mysql_db = InitMysql(_test_default_dsn)
	ctx, cancle := context.WithTimeout(context.Background(), _test_ctx_tomeout)
	defer cancle()
	_test_delall_basic_sql := `
	delete from user_basic
	`
	_test_delall_detail_sql := ` 
	delete from user_detail
	`
	var err error
	err = DeleteWithPlace(ctx, _test_mysql_db, _test_delall_detail_sql)
	err = DeleteWithPlace(ctx, _test_mysql_db, _test_delall_basic_sql)
	if err != nil {
		panic(err.Error())
	}
}

/*
实际上 select 也几乎不可能出现预期外的结果, 但是由于存在可能的model类型问题, 还是需要显式实现最小可执行测试进行部分结果显示测试.
1. one pass
2. more pass
3. page pass
*/
func Test_select_one(t *testing.T) {
	_test_mysql_db = InitMysql(_test_default_dsn)
	ctx, cancle := context.WithTimeout(context.Background(), _test_ctx_tomeout)
	defer cancle()
	// one test
	userBasic := UserBasic{}
	userDetail := UserDetail{}
	var err error
	userBasic, err = QueryWithPlace[UserBasic](ctx, _test_mysql_db, _test_sql_list[9], "zhangsan")
	userDetail, err = QueryWithName[UserDetail](ctx, _test_mysql_db, _test_sql_list[10], map[string]any{
		"user_id": 2,
	})
	if err != nil {
		panic(fmt.Sprintf("%s", err.Error()))
	}
	t.Logf("<=========================>\n")
	t.Logf("%+v\n", userBasic)
	t.Logf("%+v\n", userDetail)
	// more test
	userBasicList, err := QListWithPlace[UserBasic](ctx, _test_mysql_db, _test_sql_list[11], "test_user_%")
	userDetailList, err := QListWithName[UserDetail](ctx, _test_mysql_db, _test_sql_list[12], map[string]any{
		"min_user_id": 2,
	})
	if err != nil {
		panic(fmt.Sprintf("%s", err.Error()))
	}
	t.Logf("<=========================>\n")
	t.Logf("%+v\n", userBasicList)
	t.Logf("%+v\n", userDetailList)
	// page
	_test_current_page := 2
	_test_pagesize := 3
	page, err := QPageWithPlace[UserDetail](ctx, _test_mysql_db, _test_sql_list[13], (_test_current_page-1)*_test_pagesize, _test_pagesize)
	page.CurrentPage = _test_current_page
	page.PageSize = _test_pagesize
	page.Total = len(page.Items)
	t.Logf("<=========================>\n")
	t.Logf("%+v\n", page)
}

// test sqlite3 integration []
func Test_sqlite3_integration(t *testing.T) {
	dsn := `file:../../demo/sqlite/test.db?cache=shared&timeout=30`
	db := InitSqlite(dsn)
	var err error
	// test insert [passed]
	sqlStr := `
	insert into namespace(id,extern_id,namespace_id,name,visible,created_at,updated_at,deleted)
	values
	(23, 1023, 111, 'promotion_engine', 1, '2024-01-12 09:00:00', '2024-01-12 09:00:00', 0),
	(24, 1024, 111, 'coupon_manager', 0, '2024-01-12 10:15:00', '2024-01-12 10:15:00', 0),
	(25, 1025, 112, 'refund_process', 1, '2024-01-13 08:30:00', '2024-01-13 08:30:00', 0),
	(26, 1026, 112, 'loyalty_program', 1, '2024-01-13 14:20:00', '2024-01-13 14:20:00', 0),
	(27, 1027, 113, 'tax_calculator', 1, '2024-01-14 11:00:00', '2024-01-14 11:00:00', 0);
	`
	err = InsertWithPlace(t.Context(), db,
		sqlStr,
	)
	if err != nil {
		clog.Panic(err.Error())
	}

	// test update [passed]
	sqlStrs := []string{
		`
		UPDATE namespace 
		SET visible = 1, updated_at = CURRENT_TIMESTAMP 
		WHERE id = 24;
		`,
		`
		UPDATE namespace 
		SET updated_at = CURRENT_TIMESTAMP 
		WHERE namespace_id = 112;
		`,
		`
		UPDATE namespace 
		SET deleted = 1, updated_at = CURRENT_TIMESTAMP 
		WHERE id = 27;
		`,
	}
	for i := range sqlStrs {
		err = UpdateWithPlace(t.Context(), db, sqlStrs[i])
		if err != nil {
			clog.Panic(err.Error())
		}
	}

	// test delete [passed]
	sqlStrs = []string{
		`DELETE FROM namespace WHERE deleted = 1`,
		`DELETE FROM namespace WHERE extern_id = 1023`,
	}
	for i := range sqlStrs {
		err = DeleteWithPlace(t.Context(), db, sqlStrs[i])
		if err != nil {
			clog.Panic(err.Error())
		}
	}

	// test select [passed]
	type (
		NameSpace struct {
			Id         uint64    `db:"id"`
			ExternId   string    `db:"extern_id"`
			CreatedAt  time.Time `db:"created_at"`
			UpdatedAt  time.Time `db:"updated_at"`
			Deleted    bool      `db:"deleted"`
			Name       string    `db:"name"`
			NamespceId int       `db:"namespace_id"`
			Visible    int       `db:"visible"`
		}
	)
	sqlStrs = []string{
		`
			SELECT id, name, namespace_id, created_at 
			FROM namespace 
			WHERE deleted = 0 AND visible = 1;
		`,
		`
			SELECT id, name, updated_at 
			FROM namespace 
			WHERE created_at BETWEEN '2024-01-12 00:00:00' AND '2024-01-14 23:59:59' 
			AND deleted = 0;
			`,
		`
			SELECT t1.id, t1.name, t1.extern_id 
			FROM namespace t1
			WHERE t1.extern_id IN (1024, 1025) 
			AND t1.visible = 1;
		`,
	}
	for i := range sqlStrs {
		switch i {
		case 0:
			fmt.Printf("Query all functions that have not been deleted and are still visible.\n")
			list, err := QListWithPlace[NameSpace](t.Context(), db, sqlStrs[i])
			if err != nil {
				clog.Panic(err.Error())
			} else {
				clog.Info(fmt.Sprintf("%+v", list))
			}
		case 1:
			fmt.Printf("The function to query creation within a specified time period.")
			list, err := QListWithPlace[NameSpace](t.Context(), db, sqlStrs[i])
			if err != nil {
				clog.Panic(err.Error())
			} else {
				clog.Info(fmt.Sprintf("%+v", list))
			}
		case 2:
			fmt.Printf("Related queries")
			list, err := QListWithPlace[NameSpace](t.Context(), db, sqlStrs[i])
			if err != nil {
				clog.Panic(err.Error())
			} else {
				clog.Info(fmt.Sprintf("%+v", list))
			}
		}
	}
	clog.Info("test sqlite3 integration passed.")
}
